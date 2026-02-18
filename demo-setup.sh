#!/bin/bash

# Demo script to set up sample Kafka data for testing the Kafka Debugger
# This script requires a running Kafka instance

KAFKA_BROKER=${KAFKA_BROKERS:-"localhost:9092"}

echo "=== Kafka Debugger Demo Setup ==="
echo "Using Kafka broker: $KAFKA_BROKER"
echo ""

# Check if kafka-topics command is available
if ! command -v kafka-topics.sh &> /dev/null; then
    echo "Warning: kafka-topics.sh not found in PATH"
    echo "This script requires Kafka CLI tools to be installed"
    echo ""
    echo "You can:"
    echo "1. Install Kafka and add bin/ to your PATH"
    echo "2. Use Docker: docker-compose up -d (from this directory)"
    echo "3. Use the Kafka container directly:"
    echo "   docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --list"
    exit 1
fi

echo "Step 1: Creating test topics..."
kafka-topics.sh --bootstrap-server $KAFKA_BROKER \
    --create --topic test-topic \
    --partitions 3 \
    --replication-factor 1 \
    --if-not-exists

kafka-topics.sh --bootstrap-server $KAFKA_BROKER \
    --create --topic user-events \
    --partitions 2 \
    --replication-factor 1 \
    --if-not-exists

kafka-topics.sh --bootstrap-server $KAFKA_BROKER \
    --create --topic order-events \
    --partitions 4 \
    --replication-factor 1 \
    --if-not-exists

echo ""
echo "Step 2: Producing sample messages..."

# Produce messages to test-topic
for i in {1..50}; do
    echo "key-$i:Message $i with some sample data - timestamp $(date +%s)" | \
    kafka-console-producer.sh \
        --bootstrap-server $KAFKA_BROKER \
        --topic test-topic \
        --property "parse.key=true" \
        --property "key.separator=:"
done

# Produce messages to user-events
for i in {1..30}; do
    echo "user-$((i % 10)):User action event - login, view, or purchase - event_id=$i" | \
    kafka-console-producer.sh \
        --bootstrap-server $KAFKA_BROKER \
        --topic user-events \
        --property "parse.key=true" \
        --property "key.separator=:"
done

# Produce messages to order-events  
for i in {1..40}; do
    echo "order-$i:Order placed - customer_id=$((i % 5)) amount=$((RANDOM % 1000 + 10))" | \
    kafka-console-producer.sh \
        --bootstrap-server $KAFKA_BROKER \
        --topic order-events \
        --property "parse.key=true" \
        --property "key.separator=:"
done

echo ""
echo "Step 3: Creating test consumer groups..."

# Create a consumer group by consuming some messages
timeout 3 kafka-console-consumer.sh \
    --bootstrap-server $KAFKA_BROKER \
    --topic test-topic \
    --group test-consumer-group \
    --from-beginning \
    --max-messages 10 &> /dev/null || true

timeout 3 kafka-console-consumer.sh \
    --bootstrap-server $KAFKA_BROKER \
    --topic user-events \
    --group user-analytics-group \
    --from-beginning \
    --max-messages 5 &> /dev/null || true

timeout 3 kafka-console-consumer.sh \
    --bootstrap-server $KAFKA_BROKER \
    --topic order-events \
    --group order-processing-group \
    --from-beginning \
    --max-messages 15 &> /dev/null || true

echo ""
echo "=== Demo Setup Complete! ==="
echo ""
echo "Sample data created:"
echo "  - 3 topics with messages (test-topic, user-events, order-events)"
echo "  - 3 consumer groups (test-consumer-group, user-analytics-group, order-processing-group)"
echo "  - Consumer groups have some lag for demonstration"
echo ""
echo "You can now run the Kafka Debugger:"
echo "  ./kafkaDebugger"
echo ""
echo "Or with custom broker:"
echo "  KAFKA_BROKERS=\"$KAFKA_BROKER\" ./kafkaDebugger"
echo ""
