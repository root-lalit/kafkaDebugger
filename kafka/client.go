package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

// Client wraps Kafka operations
type Client struct {
	brokers []string
	config  *sarama.Config
	admin   sarama.ClusterAdmin
	client  sarama.Client
}

// ConsumerGroupInfo holds information about a consumer group
type ConsumerGroupInfo struct {
	GroupID     string
	State       string
	Protocol    string
	PartitionInfo []PartitionInfo
}

// PartitionInfo holds partition details including lag
type PartitionInfo struct {
	Topic     string
	Partition int32
	Offset    int64
	LogEndOffset int64
	Lag       int64
	MemberID  string
}

// TopicInfo holds topic information
type TopicInfo struct {
	Name       string
	Partitions int32
	Replicas   int16
}

// Message represents a Kafka message
type Message struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       string
	Value     string
	Timestamp time.Time
}

// NewClient creates a new Kafka client
func NewClient(brokers []string) (*Client, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Return.Errors = true

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster admin: %w", err)
	}

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		admin.Close()
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &Client{
		brokers: brokers,
		config:  config,
		admin:   admin,
		client:  client,
	}, nil
}

// Close closes the Kafka client
func (c *Client) Close() error {
	if c.admin != nil {
		c.admin.Close()
	}
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// ListConsumerGroups returns a list of all consumer groups
func (c *Client) ListConsumerGroups() ([]string, error) {
	groups, err := c.admin.ListConsumerGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to list consumer groups: %w", err)
	}

	var groupIDs []string
	for groupID := range groups {
		groupIDs = append(groupIDs, groupID)
	}
	return groupIDs, nil
}

// DescribeConsumerGroup returns detailed information about a consumer group
func (c *Client) DescribeConsumerGroup(groupID string) (*ConsumerGroupInfo, error) {
	groups, err := c.admin.DescribeConsumerGroups([]string{groupID})
	if err != nil {
		return nil, fmt.Errorf("failed to describe consumer group: %w", err)
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("consumer group not found")
	}

	group := groups[0]
	info := &ConsumerGroupInfo{
		GroupID:  group.GroupId,
		State:    group.State,
		Protocol: group.ProtocolType,
	}

	// Get offset information
	coordinator, err := c.client.Coordinator(groupID)
	if err != nil {
		return info, nil // Return basic info even if we can't get coordinator
	}

	request := &sarama.OffsetFetchRequest{
		Version:       1,
		ConsumerGroup: groupID,
	}

	// Get all topics for this group
	offsets, err := coordinator.FetchOffset(request)
	if err != nil {
		return info, nil
	}

	// Calculate lag for each partition
	for topic, partitions := range offsets.Blocks {
		for partition, block := range partitions {
			logEndOffset, err := c.client.GetOffset(topic, partition, sarama.OffsetNewest)
			if err != nil {
				continue
			}

			var memberID string
			for _, member := range group.Members {
				if assignment, err := member.GetMemberAssignment(); err == nil {
					for _, partAssignment := range assignment.Topics[topic] {
						if partAssignment == partition {
							memberID = member.MemberId
							break
						}
					}
				}
			}

			lag := logEndOffset - block.Offset
			if lag < 0 {
				lag = 0
			}

			info.PartitionInfo = append(info.PartitionInfo, PartitionInfo{
				Topic:        topic,
				Partition:    partition,
				Offset:       block.Offset,
				LogEndOffset: logEndOffset,
				Lag:          lag,
				MemberID:     memberID,
			})
		}
	}

	return info, nil
}

// ListTopics returns a list of all topics
func (c *Client) ListTopics() ([]TopicInfo, error) {
	topics, err := c.admin.ListTopics()
	if err != nil {
		return nil, fmt.Errorf("failed to list topics: %w", err)
	}

	var topicList []TopicInfo
	for name, detail := range topics {
		topicList = append(topicList, TopicInfo{
			Name:       name,
			Partitions: detail.NumPartitions,
			Replicas:   detail.ReplicationFactor,
		})
	}
	return topicList, nil
}

// GetPartitions returns partition information for a topic
func (c *Client) GetPartitions(topic string) ([]int32, error) {
	partitions, err := c.client.Partitions(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to get partitions: %w", err)
	}
	return partitions, nil
}

// GetMessages fetches messages from a topic partition
func (c *Client) GetMessages(topic string, partition int32, offset int64, limit int) ([]Message, error) {
	consumer, err := sarama.NewConsumerFromClient(c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to consume partition: %w", err)
	}
	defer partitionConsumer.Close()

	var messages []Message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < limit; i++ {
		select {
		case msg := <-partitionConsumer.Messages():
			messages = append(messages, Message{
				Topic:     msg.Topic,
				Partition: msg.Partition,
				Offset:    msg.Offset,
				Key:       string(msg.Key),
				Value:     string(msg.Value),
				Timestamp: msg.Timestamp,
			})
		case <-ctx.Done():
			return messages, nil
		}
	}

	return messages, nil
}

// GetLatestMessages fetches the latest N messages from a topic partition
func (c *Client) GetLatestMessages(topic string, partition int32, count int) ([]Message, error) {
	// Get the latest offset
	offset, err := c.client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest offset: %w", err)
	}

	// Calculate the starting offset
	startOffset := offset - int64(count)
	if startOffset < 0 {
		startOffset = 0
	}

	return c.GetMessages(topic, partition, startOffset, count)
}

// GetCurrentOffset returns the current offset for a topic partition
func (c *Client) GetCurrentOffset(topic string, partition int32) (int64, error) {
	offset, err := c.client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return 0, fmt.Errorf("failed to get current offset: %w", err)
	}
	return offset, nil
}
