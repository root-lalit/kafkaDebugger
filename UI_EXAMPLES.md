# Kafka Debugger - UI Examples

This document shows examples of what the Kafka Debugger terminal UI looks like in different views.

## Main Menu

```
┌─────────────────────────────────────────────────────────────┐
│  Kafka Debugger - Main Menu                                 │
│                                                              │
│  > Consumer Groups                                           │
│    View and manage consumer groups                           │
│                                                              │
│    Topics                                                    │
│    Browse topics and partitions                              │
│                                                              │
│    Quit                                                      │
│    Exit the application                                      │
│                                                              │
└─────────────────────────────────────────────────────────────┘

✓ Connected to Kafka

[enter] select • [esc/q] back/quit • [↑/↓] navigate
```

## Consumer Groups List

```
┌─────────────────────────────────────────────────────────────┐
│  Consumer Groups                                             │
│                                                              │
│  > test-consumer-group                                       │
│    Consumer Group                                            │
│                                                              │
│    user-analytics-group                                      │
│    Consumer Group                                            │
│                                                              │
│    order-processing-group                                    │
│    Consumer Group                                            │
│                                                              │
└─────────────────────────────────────────────────────────────┘

✓ Loaded 3 consumer groups

[enter] select • [esc/q] back/quit • [↑/↓] navigate
```

## Consumer Group Details with Lag

```
Consumer Group Details

┌──────────────────────────────────────────────────────────────────────────────────────────┐
│ Topic          Partition  Offset      Log End     Lag         Member ID                  │
├──────────────────────────────────────────────────────────────────────────────────────────┤
│ test-topic     0          10          50          40          consumer-1-abc123          │
│ test-topic     1          25          48          23          consumer-2-def456          │
│ test-topic     2          30          30          0           consumer-3-ghi789          │
│ user-events    0          5           30          25          consumer-1-abc123          │
│ user-events    1          8           28          20          consumer-2-def456          │
└──────────────────────────────────────────────────────────────────────────────────────────┘

✓ Group: test-consumer-group | State: Stable | Protocol: consumer

[enter] select • [esc/q] back/quit • [↑/↓] navigate
```

## Topics List

```
┌─────────────────────────────────────────────────────────────┐
│  Topics                                                      │
│                                                              │
│  > test-topic                                                │
│    Partitions: 3, Replicas: 1                                │
│                                                              │
│    user-events                                               │
│    Partitions: 2, Replicas: 1                                │
│                                                              │
│    order-events                                              │
│    Partitions: 4, Replicas: 1                                │
│                                                              │
│    __consumer_offsets                                        │
│    Partitions: 50, Replicas: 1                               │
│                                                              │
└─────────────────────────────────────────────────────────────┘

✓ Loaded 4 topics

[enter] select • [esc/q] back/quit • [↑/↓] navigate
```

## Messages View

```
Messages

┌────────────────────────────────────────────────────────────────────────────────────────────────┐
│ Offset      Partition  Key           Value                              Timestamp              │
├────────────────────────────────────────────────────────────────────────────────────────────────┤
│ 31          0          key-31        Message 31 with sample data        2024-02-18 15:30:12    │
│ 32          0          key-32        Message 32 with sample data        2024-02-18 15:30:13    │
│ 33          0          key-33        Message 33 with sample data        2024-02-18 15:30:14    │
│ 34          0          key-34        Message 34 with sample data        2024-02-18 15:30:15    │
│ 35          0          key-35        Message 35 with sample data        2024-02-18 15:30:16    │
│ 36          0          key-36        Message 36 with sample data        2024-02-18 15:30:17    │
│ 37          0          key-37        Message 37 with sample data        2024-02-18 15:30:18    │
│ 38          0          key-38        Message 38 with sample data        2024-02-18 15:30:19    │
│ 39          0          key-39        Message 39 with sample data        2024-02-18 15:30:20    │
│ 40          0          key-40        Message 40 with sample data        2024-02-18 15:30:21    │
└────────────────────────────────────────────────────────────────────────────────────────────────┘

✓ Showing messages from topic: test-topic, partition: 0

[enter] select • [esc/q] back/quit • [↑/↓] navigate
```

## Error Example

```
┌─────────────────────────────────────────────────────────────┐
│  Consumer Groups                                             │
│                                                              │
│  (No consumer groups available)                              │
│                                                              │
└─────────────────────────────────────────────────────────────┘

❌ Error: failed to connect to broker: dial tcp 127.0.0.1:9092: connect: connection refused

[enter] select • [esc/q] back/quit • [↑/↓] navigate
```

## Key Features Visible in UI

### 1. Consumer Group Lag Monitoring
The most important feature for debugging - the Consumer Group Details view shows:
- **Current Offset**: Where the consumer is currently reading
- **Log End Offset**: The latest message available
- **Lag**: How many messages behind (Log End - Current)
- Color coding: High lag is highlighted for quick identification

### 2. Real-time Updates
- Status messages show loading states
- Errors are displayed prominently
- Connection status is always visible

### 3. Navigation
- Keyboard-only interface
- Consistent controls across all views
- Breadcrumb-style navigation (press Esc to go back)

### 4. Data Presentation
- Tables for structured data (partitions, messages)
- Lists for browsing (groups, topics)
- Truncated values for readability
- Clear column headers

### 5. Helpful Feedback
- Status bar shows operation results
- Error messages explain what went wrong
- Help text at bottom shows available commands
- Empty states guide users when no data is available

## Use Cases Demonstrated

### Debugging Consumer Lag
1. View consumer groups → Select group → Check lag column
2. High lag indicates consumers can't keep up with producers
3. Zero lag means consumers are caught up

### Inspecting Messages
1. View topics → Select topic → See recent messages
2. Verify messages are arriving correctly
3. Check message content, keys, and timestamps

### Monitoring Partition Distribution
1. View consumer group details
2. See which consumer (member ID) owns which partition
3. Verify load is distributed evenly

### Checking Topic Configuration
1. View topics list
2. See partition counts and replication factors
3. Identify misconfigured topics
