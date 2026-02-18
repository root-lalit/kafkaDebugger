# Kafka Debugger

A powerful Terminal UI (TUI) application for debugging and monitoring Apache Kafka clusters. This tool provides an intuitive interface for common Kafka operations including consumer group management, topic browsing, partition viewing, and message inspection.

## Features

- **Consumer Group Management**
  - List all consumer groups
  - Describe consumer group details
  - View consumer group lag in real-time
  - Monitor partition offsets and log end offsets
  - See member assignments

- **Topic Operations**
  - Browse all topics
  - View topic details (partitions, replicas)
  - Get messages from any partition
  - Fetch latest messages from current offset
  - View partition information

- **Message Inspection**
  - Display message content with key, value, and metadata
  - View messages from specific offsets
  - Get latest N messages from any partition
  - See message timestamps

- **User Interface**
  - Clean, intuitive terminal UI
  - Keyboard navigation
  - Real-time updates
  - Color-coded information
  - Table views for easy data browsing

## Installation

### Prerequisites

- Go 1.20 or higher
- Access to a Kafka cluster

### Build from Source

```bash
# Clone the repository
git clone https://github.com/root-lalit/kafkaDebugger.git
cd kafkaDebugger

# Build the application
go build -o kafkaDebugger .

# Run the application
./kafkaDebugger
```

### Install using Go

```bash
go install github.com/root-lalit/kafkaDebugger@latest
```

## Usage

### Basic Usage

By default, the application connects to `localhost:9092`:

```bash
./kafkaDebugger
```

### Custom Kafka Brokers

Set the `KAFKA_BROKERS` environment variable to connect to different brokers:

```bash
# Single broker
export KAFKA_BROKERS="kafka1.example.com:9092"
./kafkaDebugger

# Multiple brokers
export KAFKA_BROKERS="kafka1.example.com:9092,kafka2.example.com:9092,kafka3.example.com:9092"
./kafkaDebugger

# Or inline
KAFKA_BROKERS="kafka.example.com:9092" ./kafkaDebugger
```

## Navigation

### Keyboard Controls

- **Arrow Keys (↑/↓)**: Navigate through lists and tables
- **Enter**: Select an item or confirm action
- **Esc/q**: Go back to previous view or quit (from main menu)
- **Ctrl+C**: Force quit from anywhere

### Application Flow

1. **Main Menu**: Choose between Consumer Groups, Topics, or Quit
2. **Consumer Groups**:
   - View list of all consumer groups
   - Select a group to see detailed information
   - View lag, offsets, and member assignments
3. **Topics**:
   - Browse all topics
   - Select a topic to view messages
   - See latest messages from partition 0 (default)

## Examples

### Monitoring Consumer Group Lag

1. Start the application
2. Select "Consumer Groups" from the main menu
3. Choose the consumer group you want to monitor
4. View the lag information for each partition
   - **Lag**: Number of messages the consumer is behind
   - **Offset**: Current consumer position
   - **Log End Offset**: Latest message in the partition

### Viewing Messages from a Topic

1. Start the application
2. Select "Topics" from the main menu
3. Choose the topic you want to inspect
4. View the latest 20 messages from the topic
5. Messages are displayed with:
   - Offset
   - Partition
   - Key
   - Value (truncated if too long)
   - Timestamp

### Common Use Cases

#### Debug Consumer Lag Issues

```bash
KAFKA_BROKERS="prod-kafka:9092" ./kafkaDebugger
# Navigate to Consumer Groups → Select your group → Check lag values
```

#### Inspect Recent Messages

```bash
KAFKA_BROKERS="prod-kafka:9092" ./kafkaDebugger
# Navigate to Topics → Select your topic → View messages
```

#### Monitor Multiple Partitions

```bash
# The consumer group details view shows all partitions
# with their individual offsets and lag information
```

## Architecture

The application is built with:

- **Go**: Primary language for performance and concurrency
- **IBM Sarama**: Kafka client library for Go
- **Bubble Tea**: Terminal UI framework
- **Bubbles**: Pre-built TUI components (lists, tables, inputs)
- **Lipgloss**: Styling and layout for terminal UIs

## Project Structure

```
kafkaDebugger/
├── main.go           # Application entry point
├── kafka/
│   └── client.go     # Kafka client wrapper and operations
├── ui/
│   └── ui.go         # Terminal UI implementation
├── go.mod            # Go module dependencies
└── README.md         # This file
```

## Configuration

### Environment Variables

- `KAFKA_BROKERS`: Comma-separated list of Kafka broker addresses (default: `localhost:9092`)

### Future Configuration Options

The following configuration options may be added in future versions:

- SASL authentication
- SSL/TLS configuration
- Custom consumer group protocols
- Message format options
- Refresh intervals

## Troubleshooting

### Connection Issues

If you can't connect to Kafka:

1. Verify the broker addresses are correct
2. Check network connectivity: `telnet kafka-broker 9092`
3. Ensure your firewall allows connections to Kafka ports
4. Check Kafka broker logs for connection errors

### No Consumer Groups Listed

- Verify consumer groups exist: `kafka-consumer-groups.sh --bootstrap-server localhost:9092 --list`
- Ensure the broker addresses are correct
- Check if you have permissions to access consumer group information

### Empty Topics or No Messages

- Verify the topic exists and has messages
- Check if you have permissions to read from the topic
- Try producing a test message to the topic

## Development

### Building

```bash
go build -o kafkaDebugger .
```

### Running Tests

```bash
go test ./...
```

### Dependencies

Update dependencies:

```bash
go get -u ./...
go mod tidy
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is open source and available under the MIT License.

## Acknowledgments

- [IBM Sarama](https://github.com/IBM/sarama) - Kafka client for Go
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## Support

For issues, questions, or contributions, please use the GitHub issue tracker.
