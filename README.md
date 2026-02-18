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

# Build the application using Makefile
make build

# Or build directly with Go
go build -o kafkaDebugger .

# Run the application
./kafkaDebugger
```

### Install using Go

```bash
go install github.com/root-lalit/kafkaDebugger@latest
```

## Quick Start with Docker

For a quick test environment, you can use the included Docker Compose file:

```bash
# Start Kafka and Zookeeper
docker-compose up -d

# Wait for Kafka to be ready (about 30 seconds)
sleep 30

# Run the demo setup script to create sample data
./demo-setup.sh

# Run the Kafka Debugger
./kafkaDebugger

# When done, stop the services
docker-compose down
```

Or use the Makefile for convenience:

```bash
# Start Kafka environment
make docker-up

# Wait for Kafka to be ready
sleep 30

# Set up demo data
make demo

# Run the application
make run

# When done, stop the services
make docker-down
```

## Makefile Commands

The project includes a Makefile for common tasks:

- `make build` - Build the application
- `make run` - Build and run the application  
- `make clean` - Remove build artifacts
- `make test` - Run tests
- `make tidy` - Tidy Go dependencies
- `make docker-up` - Start Docker Compose Kafka environment
- `make docker-down` - Stop Docker Compose environment
- `make demo` - Set up demo data in Kafka
- `make help` - Show all available commands

## Usage

### Basic Usage

By default, the application connects to `localhost:9092`:

```bash
./kafkaDebugger
```

### Managing Broker Aliases

The tool supports storing broker configurations with aliases for easy access. Configurations are stored in `~/.kafkaDebugger/config.json`.

#### Add a Broker Alias

```bash
# Add a single broker
./kafkaDebugger -add-alias local:localhost:9092

# Add multiple brokers
./kafkaDebugger -add-alias prod:kafka1.example.com:9092,kafka2.example.com:9092,kafka3.example.com:9092
```

#### List All Aliases

```bash
./kafkaDebugger -list-aliases
```

#### Remove an Alias

```bash
./kafkaDebugger -remove-alias prod
```

#### Use an Alias

```bash
./kafkaDebugger -alias local
./kafkaDebugger -alias prod
```

### Custom Kafka Brokers

You can specify brokers in multiple ways (in order of priority):

1. **Command-line flag** (highest priority):
```bash
./kafkaDebugger -brokers kafka1.example.com:9092,kafka2.example.com:9092
```

2. **Alias from config file**:
```bash
./kafkaDebugger -alias prod
```

3. **Environment variable**:
```bash
export KAFKA_BROKERS="kafka1.example.com:9092,kafka2.example.com:9092"
./kafkaDebugger

# Or inline
KAFKA_BROKERS="kafka.example.com:9092" ./kafkaDebugger
```

4. **Default** (lowest priority): `localhost:9092`

### Running from Anywhere

The tool can be run from any directory. The configuration file is stored in your home directory at `~/.kafkaDebugger/config.json`, so broker aliases are available system-wide.

```bash
# Works from any directory
cd /tmp
/path/to/kafkaDebugger -alias local
```

## Navigation

### Keyboard Controls

- **Arrow Keys (↑/↓)**: Navigate through lists and tables
- **Enter**: Select an item or confirm action
- **Esc/q**: Go back to previous view or quit (from main menu)
- **Ctrl+C**: Force quit from anywhere

### Application Flow

1. **Main Menu**: Choose between Consumer Groups, Topics, Reconnect, or Quit
2. **Connection Issues**: If the initial connection to Kafka fails, the main menu will display a "Reconnect to Kafka" option
3. **Consumer Groups**:
   - View list of all consumer groups
   - Select a group to see detailed information
   - View lag, offsets, and member assignments
4. **Topics**:
   - Browse all topics
   - Select a topic to view messages
   - See latest messages from partition 0 (default)

### Handling Connection Errors

If the application cannot connect to Kafka brokers:

1. The main menu will show "Reconnect to Kafka" option
2. Error details will be displayed at the bottom of the screen
3. Select "Reconnect to Kafka" to retry the connection
4. Verify your broker addresses and network connectivity

## Examples

### Using Broker Aliases

```bash
# Add your development environment
./kafkaDebugger -add-alias dev:dev-kafka:9092

# Add your production environment
./kafkaDebugger -add-alias prod:kafka1.prod.com:9092,kafka2.prod.com:9092,kafka3.prod.com:9092

# List all configured environments
./kafkaDebugger -list-aliases

# Connect to development
./kafkaDebugger -alias dev

# Connect to production
./kafkaDebugger -alias prod
```

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
├── main.go           # Application entry point with CLI flags
├── config/
│   └── config.go     # Configuration management for broker aliases
├── kafka/
│   └── client.go     # Kafka client wrapper and operations
├── ui/
│   └── ui.go         # Terminal UI implementation
├── go.mod            # Go module dependencies
└── README.md         # This file
```

## Configuration

### Broker Aliases Configuration File

The application stores broker configurations in `~/.kafkaDebugger/config.json`. This file is automatically created when you add your first alias.

Example configuration:
```json
{
  "defaultAlias": "local",
  "brokers": [
    {
      "name": "local",
      "brokers": ["localhost:9092"]
    },
    {
      "name": "prod",
      "brokers": [
        "kafka1.example.com:9092",
        "kafka2.example.com:9092",
        "kafka3.example.com:9092"
      ]
    }
  ]
}
```

### Command-Line Options

- `-brokers`: Specify broker addresses directly
- `-alias`: Use a named broker configuration
- `-add-alias`: Add or update a broker alias
- `-list-aliases`: List all configured aliases
- `-remove-alias`: Remove a broker alias

### Environment Variables

- `KAFKA_BROKERS`: Comma-separated list of Kafka broker addresses (default: `localhost:9092`)

### Configuration Priority

The application uses the following priority order (highest to lowest):
1. `-brokers` command-line flag
2. `-alias` command-line flag
3. `KAFKA_BROKERS` environment variable
4. `defaultAlias` from config file
5. Default: `localhost:9092`

### Future Configuration Options

The following configuration options may be added in future versions:

- SASL authentication
- SSL/TLS configuration
- Custom consumer group protocols
- Message format options
- Refresh intervals

## Troubleshooting

### "Kafka Client Not Initialized" Error

If you see this error:

1. **Check broker connectivity**: Ensure you can reach the Kafka brokers from your machine
   ```bash
   telnet kafka-broker 9092
   # or
   nc -zv kafka-broker 9092
   ```

2. **Verify broker addresses**: Make sure the broker addresses are correct
   ```bash
   # List your configured aliases
   ./kafkaDebugger -list-aliases
   
   # Try connecting with explicit brokers
   ./kafkaDebugger -brokers localhost:9092
   ```

3. **Use the Reconnect option**: If the initial connection fails, the main menu will show "Reconnect to Kafka". Select this option to retry the connection.

4. **Check Kafka broker status**: Ensure Kafka is running and accepting connections
   ```bash
   # If using Docker
   docker ps | grep kafka
   
   # Check broker logs
   docker logs kafka
   ```

### Connection Issues

If you can't connect to Kafka:

1. Verify the broker addresses are correct
2. Check network connectivity: `telnet kafka-broker 9092`
3. Ensure your firewall allows connections to Kafka ports
4. Check Kafka broker logs for connection errors
5. Try using the reconnect feature from the main menu

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
