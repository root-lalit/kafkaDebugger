package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/root-lalit/kafkaDebugger/kafka"
)

// View represents the different views in the application
type View int

const (
	MainMenuView View = iota
	ConsumerGroupsListView
	ConsumerGroupDetailsView
	TopicsListView
	TopicDetailsView
	PartitionView
	MessagesView
)

// Model represents the application state
type Model struct {
	brokers           string
	kafkaClient       *kafka.Client
	currentView       View
	mainMenu          list.Model
	consumerGroupList list.Model
	topicsList        list.Model
	detailsTable      table.Model
	messagesTable     table.Model
	input             textinput.Model
	width             int
	height            int
	err               error
	loading           bool
	selectedGroup     string
	selectedTopic     string
	selectedPartition int32
	statusMsg         string
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type errMsg struct{ err error }
type statusMsg string
type kafkaClientMsg *kafka.Client
type consumerGroupsMsg []string
type topicsMsg []kafka.TopicInfo
type groupDetailsMsg *kafka.ConsumerGroupInfo
type messagesMsg []kafka.Message

func (e errMsg) Error() string { return e.err.Error() }

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))
)

// NewModel creates a new application model
func NewModel(brokers string) Model {
	// Create main menu items
	mainMenuItems := []list.Item{
		item{title: "Consumer Groups", desc: "View and manage consumer groups"},
		item{title: "Topics", desc: "Browse topics and partitions"},
		item{title: "Quit", desc: "Exit the application"},
	}

	mainMenu := list.New(mainMenuItems, list.NewDefaultDelegate(), 0, 0)
	mainMenu.Title = "Kafka Debugger - Main Menu"
	mainMenu.SetShowStatusBar(false)

	// Create text input for various operations
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.CharLimit = 100

	return Model{
		brokers:     brokers,
		currentView: MainMenuView,
		mainMenu:    mainMenu,
		input:       ti,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.connectKafka()
}

// updateMainMenu updates the main menu items based on connection status
func (m *Model) updateMainMenu() {
	var mainMenuItems []list.Item

	if m.kafkaClient == nil {
		// If not connected, add Reconnect option
		mainMenuItems = []list.Item{
			item{title: "Reconnect to Kafka", desc: "Retry connection to Kafka brokers"},
			item{title: "Quit", desc: "Exit the application"},
		}
	} else {
		// If connected, show normal menu
		mainMenuItems = []list.Item{
			item{title: "Consumer Groups", desc: "View and manage consumer groups"},
			item{title: "Topics", desc: "Browse topics and partitions"},
			item{title: "Quit", desc: "Exit the application"},
		}
	}

	m.mainMenu.SetItems(mainMenuItems)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.mainMenu.SetSize(msg.Width-4, msg.Height-4)
		if m.consumerGroupList.Items() != nil {
			m.consumerGroupList.SetSize(msg.Width-4, msg.Height-4)
		}
		if m.topicsList.Items() != nil {
			m.topicsList.SetSize(msg.Width-4, msg.Height-4)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.currentView == MainMenuView {
				if m.kafkaClient != nil {
					m.kafkaClient.Close()
				}
				return m, tea.Quit
			}
			// Return to previous view
			return m.goBack(), nil

		case "esc":
			return m.goBack(), nil

		case "enter":
			return m.handleEnter()
		}

	case errMsg:
		m.err = msg.err
		m.loading = false
		// Update main menu to show reconnect option if client is nil
		if m.currentView == MainMenuView {
			m.updateMainMenu()
		}
		return m, nil

	case statusMsg:
		m.statusMsg = string(msg)
		m.loading = false
		return m, nil

	case kafkaClientMsg:
		m.kafkaClient = msg
		m.statusMsg = "Connected to Kafka"
		m.loading = false
		// Update main menu to show normal options
		if m.currentView == MainMenuView {
			m.updateMainMenu()
		}
		return m, nil

	case consumerGroupsMsg:
		m.loading = false
		items := make([]list.Item, len(msg))
		for i, group := range msg {
			items[i] = item{title: group, desc: "Consumer Group"}
		}
		m.consumerGroupList = list.New(items, list.NewDefaultDelegate(), m.width-4, m.height-4)
		m.consumerGroupList.Title = "Consumer Groups"
		m.currentView = ConsumerGroupsListView
		return m, nil

	case topicsMsg:
		m.loading = false
		items := make([]list.Item, len(msg))
		for i, topic := range msg {
			items[i] = item{
				title: topic.Name,
				desc:  fmt.Sprintf("Partitions: %d, Replicas: %d", topic.Partitions, topic.Replicas),
			}
		}
		m.topicsList = list.New(items, list.NewDefaultDelegate(), m.width-4, m.height-4)
		m.topicsList.Title = "Topics"
		m.currentView = TopicsListView
		return m, nil

	case groupDetailsMsg:
		m.loading = false
		m.currentView = ConsumerGroupDetailsView

		columns := []table.Column{
			{Title: "Topic", Width: 30},
			{Title: "Partition", Width: 10},
			{Title: "Offset", Width: 15},
			{Title: "Log End", Width: 15},
			{Title: "Lag", Width: 15},
			{Title: "Member ID", Width: 40},
		}

		rows := make([]table.Row, len(msg.PartitionInfo))
		for i, p := range msg.PartitionInfo {
			rows[i] = table.Row{
				p.Topic,
				fmt.Sprintf("%d", p.Partition),
				fmt.Sprintf("%d", p.Offset),
				fmt.Sprintf("%d", p.LogEndOffset),
				fmt.Sprintf("%d", p.Lag),
				p.MemberID,
			}
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(m.height-10),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		m.detailsTable = t
		m.statusMsg = fmt.Sprintf("Group: %s | State: %s | Protocol: %s", msg.GroupID, msg.State, msg.Protocol)
		return m, nil

	case messagesMsg:
		m.loading = false
		m.currentView = MessagesView

		columns := []table.Column{
			{Title: "Offset", Width: 15},
			{Title: "Partition", Width: 10},
			{Title: "Key", Width: 25},
			{Title: "Value", Width: 60},
			{Title: "Timestamp", Width: 25},
		}

		rows := make([]table.Row, len(msg))
		for i, msg := range msg {
			value := msg.Value
			if len(value) > 60 {
				value = value[:57] + "..."
			}
			rows[i] = table.Row{
				fmt.Sprintf("%d", msg.Offset),
				fmt.Sprintf("%d", msg.Partition),
				msg.Key,
				value,
				msg.Timestamp.Format("2006-01-02 15:04:05"),
			}
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(m.height-8),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		m.messagesTable = t
		m.statusMsg = fmt.Sprintf("Showing messages from topic: %s, partition: %d", m.selectedTopic, m.selectedPartition)
		return m, nil
	}

	// Update the current view's component
	switch m.currentView {
	case MainMenuView:
		m.mainMenu, cmd = m.mainMenu.Update(msg)
		cmds = append(cmds, cmd)
	case ConsumerGroupsListView:
		m.consumerGroupList, cmd = m.consumerGroupList.Update(msg)
		cmds = append(cmds, cmd)
	case TopicsListView:
		m.topicsList, cmd = m.topicsList.Update(msg)
		cmds = append(cmds, cmd)
	case ConsumerGroupDetailsView:
		m.detailsTable, cmd = m.detailsTable.Update(msg)
		cmds = append(cmds, cmd)
	case MessagesView:
		m.messagesTable, cmd = m.messagesTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content string

	switch m.currentView {
	case MainMenuView:
		content = m.mainMenu.View()
	case ConsumerGroupsListView:
		content = m.consumerGroupList.View()
	case TopicsListView:
		content = m.topicsList.View()
	case ConsumerGroupDetailsView:
		title := titleStyle.Render("Consumer Group Details")
		content = fmt.Sprintf("%s\n\n%s", title, m.detailsTable.View())
	case MessagesView:
		title := titleStyle.Render("Messages")
		content = fmt.Sprintf("%s\n\n%s", title, m.messagesTable.View())
	}

	// Add status bar
	var statusBar string
	if m.loading {
		statusBar = statusStyle.Render("⏳ Loading...")
	} else if m.err != nil {
		statusBar = errorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err))
	} else if m.statusMsg != "" {
		statusBar = statusStyle.Render("✓ " + m.statusMsg)
	}

	help := helpStyle.Render("\n[enter] select • [esc/q] back/quit • [↑/↓] navigate")

	return fmt.Sprintf("%s\n\n%s%s\n", content, statusBar, help)
}

// Helper functions

// isClientConnected checks if the Kafka client is initialized and sets an error if not
func (m *Model) isClientConnected() bool {
	if m.kafkaClient == nil {
		m.err = fmt.Errorf("kafka client not initialized - please reconnect")
		return false
	}
	return true
}

func (m Model) connectKafka() tea.Cmd {
	return func() tea.Msg {
		brokers := strings.Split(m.brokers, ",")
		client, err := kafka.NewClient(brokers)
		if err != nil {
			return errMsg{err}
		}
		return kafkaClientMsg(client)
	}
}

func (m Model) goBack() Model {
	m.err = nil
	m.statusMsg = ""

	switch m.currentView {
	case ConsumerGroupsListView, TopicsListView:
		m.currentView = MainMenuView
	case ConsumerGroupDetailsView:
		m.currentView = ConsumerGroupsListView
	case TopicDetailsView, MessagesView:
		m.currentView = TopicsListView
	}

	return m
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case MainMenuView:
		selected := m.mainMenu.SelectedItem()
		if selected == nil {
			return m, nil
		}

		switch selected.(item).title {
		case "Reconnect to Kafka":
			m.loading = true
			m.err = nil
			m.statusMsg = ""
			return m, m.connectKafka()
		case "Consumer Groups":
			if !m.isClientConnected() {
				return m, nil
			}
			m.loading = true
			return m, m.loadConsumerGroups()
		case "Topics":
			if !m.isClientConnected() {
				return m, nil
			}
			m.loading = true
			return m, m.loadTopics()
		case "Quit":
			if m.kafkaClient != nil {
				m.kafkaClient.Close()
			}
			return m, tea.Quit
		}

	case ConsumerGroupsListView:
		selected := m.consumerGroupList.SelectedItem()
		if selected == nil {
			return m, nil
		}
		m.selectedGroup = selected.(item).title
		m.loading = true
		return m, m.loadGroupDetails(m.selectedGroup)

	case TopicsListView:
		selected := m.topicsList.SelectedItem()
		if selected == nil {
			return m, nil
		}
		m.selectedTopic = selected.(item).title
		m.selectedPartition = 0 // Default to partition 0
		m.loading = true
		return m, m.loadLatestMessages(m.selectedTopic, m.selectedPartition)
	}

	return m, nil
}

func (m Model) loadConsumerGroups() tea.Cmd {
	return func() tea.Msg {
		if m.kafkaClient == nil {
			return errMsg{fmt.Errorf("kafka client not initialized")}
		}
		groups, err := m.kafkaClient.ListConsumerGroups()
		if err != nil {
			return errMsg{err}
		}
		return consumerGroupsMsg(groups)
	}
}

func (m Model) loadTopics() tea.Cmd {
	return func() tea.Msg {
		if m.kafkaClient == nil {
			return errMsg{fmt.Errorf("kafka client not initialized")}
		}
		topics, err := m.kafkaClient.ListTopics()
		if err != nil {
			return errMsg{err}
		}
		return topicsMsg(topics)
	}
}

func (m Model) loadGroupDetails(groupID string) tea.Cmd {
	return func() tea.Msg {
		if m.kafkaClient == nil {
			return errMsg{fmt.Errorf("kafka client not initialized")}
		}
		details, err := m.kafkaClient.DescribeConsumerGroup(groupID)
		if err != nil {
			return errMsg{err}
		}
		return groupDetailsMsg(details)
	}
}

func (m Model) loadLatestMessages(topic string, partition int32) tea.Cmd {
	return func() tea.Msg {
		if m.kafkaClient == nil {
			return errMsg{fmt.Errorf("kafka client not initialized")}
		}
		messages, err := m.kafkaClient.GetLatestMessages(topic, partition, 20)
		if err != nil {
			return errMsg{err}
		}
		return messagesMsg(messages)
	}
}
