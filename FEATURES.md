# Kafka Debugger - Features Overview

## Application Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                        MAIN MENU                                │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │  Consumer Groups                                         │  │
│  │  └─> View and manage consumer groups                    │  │
│  │                                                           │  │
│  │  Topics                                                   │  │
│  │  └─> Browse topics and partitions                        │  │
│  │                                                           │  │
│  │  Quit                                                     │  │
│  │  └─> Exit the application                                │  │
│  └─────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                           ↓
                    (Select option)
                           ↓
        ┌──────────────────┴──────────────────┐
        │                                     │
        ▼                                     ▼
┌──────────────────┐                  ┌──────────────────┐
│ CONSUMER GROUPS  │                  │     TOPICS       │
│                  │                  │                  │
│ List all groups  │                  │ List all topics  │
│                  │                  │                  │
│ ▼ Select group   │                  │ ▼ Select topic   │
│                  │                  │                  │
│ ┌──────────────┐ │                  │ ┌──────────────┐ │
│ │ GROUP DETAILS││                  │ │   MESSAGES   ││
│ │              ││                  │ │              ││
│ │ • Group Info ││                  │ │ • Latest 20  ││
│ │ • Partitions ││                  │ │   messages   ││
│ │ • Lag        ││                  │ │ • Key/Value  ││
│ │ • Offsets    ││                  │ │ • Timestamp  ││
│ │ • Members    ││                  │ │ • Offset     ││
│ └──────────────┘ │                  │ └──────────────┘ │
└──────────────────┘                  └──────────────────┘
```

## Features by View

### Main Menu
- **Navigation**: Arrow keys to move, Enter to select
- **Options**: Consumer Groups, Topics, or Quit
- **Status Bar**: Shows connection status and errors

### Consumer Groups View
**List View:**
- Shows all consumer groups in the cluster
- Displays group name and description
- Quick navigation with arrow keys

**Details View:**
- Group metadata (ID, state, protocol)
- Partition-level information in a table:
  - Topic name
  - Partition number
  - Current offset
  - Log end offset
  - **Lag** (how far behind the consumer is)
  - Member ID (which consumer owns this partition)
- Real-time lag monitoring
- Sortable columns

### Topics View
**List View:**
- Shows all topics in the cluster
- Displays partition count and replication factor
- Search and filter capabilities

**Messages View:**
- Shows latest 20 messages from the selected topic
- Displays in a clean table format:
  - Offset
  - Partition
  - Key (message key)
  - Value (message content, truncated if long)
  - Timestamp
- Full message details visible when selected

## Key Features

### 🔍 Consumer Group Debugging
- **Lag Monitoring**: See how far behind consumers are
- **Offset Tracking**: View current position vs. latest available
- **Member Assignment**: See which consumer owns which partition
- **Group State**: Check if group is active, empty, or dead

### 📊 Topic Inspection
- **Message Browsing**: View recent messages
- **Partition Information**: See partition distribution
- **Metadata**: Replication and partition counts

### ⚡ Performance
- **Fast Navigation**: Instant switching between views
- **Efficient Data Fetching**: Only loads what you need
- **Real-time Updates**: Fresh data on every view

### 🎨 User Experience
- **Clean Interface**: Minimalist, focused design
- **Color Coding**: Status messages, errors, and highlights
- **Keyboard Shortcuts**: Efficient navigation
- **Responsive**: Adapts to terminal size

## Common Workflows

### 1. Check Consumer Lag
```
1. Launch application
2. Select "Consumer Groups"
3. Find your consumer group
4. Press Enter to view details
5. Check the "Lag" column for each partition
```

### 2. View Recent Messages
```
1. Launch application
2. Select "Topics"
3. Find your topic
4. Press Enter to view messages
5. Browse the latest 20 messages
```

### 3. Debug Processing Issues
```
1. Check consumer group lag (is it growing?)
2. View topic messages (are messages arriving?)
3. Check partition assignments (are all partitions assigned?)
4. Compare offsets (are consumers making progress?)
```

## Keyboard Reference

| Key | Action |
|-----|--------|
| ↑/↓ | Navigate up/down in lists and tables |
| Enter | Select item or confirm action |
| Esc | Go back to previous view |
| q | Quit from main menu or go back |
| Ctrl+C | Force quit from anywhere |

## Status Indicators

- ✓ **Green**: Success/Normal operation
- ⏳ **Blue**: Loading data
- ❌ **Red**: Error occurred

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| KAFKA_BROKERS | Comma-separated list of Kafka brokers | localhost:9092 |

### Examples

```bash
# Local Kafka
KAFKA_BROKERS="localhost:9092" ./kafkaDebugger

# Production cluster
KAFKA_BROKERS="kafka1:9092,kafka2:9092,kafka3:9092" ./kafkaDebugger

# With custom port
KAFKA_BROKERS="kafka.example.com:9093" ./kafkaDebugger
```
