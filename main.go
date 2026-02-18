package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/root-lalit/kafkaDebugger/config"
	"github.com/root-lalit/kafkaDebugger/ui"
)

func main() {
	// Define command-line flags
	brokerFlag := flag.String("brokers", "", "Comma-separated list of Kafka broker addresses")
	aliasFlag := flag.String("alias", "", "Use broker alias from config file")
	addAliasFlag := flag.String("add-alias", "", "Add broker alias (format: name:broker1:9092,broker2:9092)")
	listAliasFlag := flag.Bool("list-aliases", false, "List all broker aliases")
	removeAliasFlag := flag.String("remove-alias", "", "Remove broker alias")

	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		cfg = &config.Config{Brokers: []config.BrokerConfig{}}
	}

	// Handle alias management commands
	if *addAliasFlag != "" {
		handleAddAlias(cfg, *addAliasFlag)
		return
	}

	if *listAliasFlag {
		handleListAliases(cfg)
		return
	}

	if *removeAliasFlag != "" {
		handleRemoveAlias(cfg, *removeAliasFlag)
		return
	}

	// Determine broker addresses
	var brokers string

	// Priority: 1. Command-line flag, 2. Alias, 3. Environment variable, 4. Default from config, 5. localhost
	if *brokerFlag != "" {
		brokers = *brokerFlag
	} else if *aliasFlag != "" {
		brokerCfg, err := cfg.GetBrokerByAlias(*aliasFlag)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		brokers = strings.Join(brokerCfg.Brokers, ",")
	} else if envBrokers := os.Getenv("KAFKA_BROKERS"); envBrokers != "" {
		brokers = envBrokers
	} else if cfg.DefaultAlias != "" {
		brokerCfg, err := cfg.GetBrokerByAlias(cfg.DefaultAlias)
		if err == nil {
			brokers = strings.Join(brokerCfg.Brokers, ",")
		} else {
			brokers = "localhost:9092"
		}
	} else {
		brokers = "localhost:9092"
	}

	// Create the UI model
	m := ui.NewModel(brokers)

	// Start the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func handleAddAlias(cfg *config.Config, aliasSpec string) {
	parts := strings.SplitN(aliasSpec, ":", 2)
	if len(parts) != 2 {
		fmt.Println("Error: Invalid format. Use: name:broker1:9092,broker2:9092")
		os.Exit(1)
	}

	name := parts[0]
	brokerList := strings.Split(parts[1], ",")

	if err := cfg.AddBroker(name, brokerList); err != nil {
		fmt.Printf("Error adding alias: %v\n", err)
		os.Exit(1)
	}

	if err := config.SaveConfig(cfg); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Successfully added/updated alias '%s' with brokers: %s\n", name, strings.Join(brokerList, ", "))

	// Show config file location
	configPath, _ := config.GetConfigPath()
	fmt.Printf("  Config saved to: %s\n", configPath)
}

func handleListAliases(cfg *config.Config) {
	brokers := cfg.ListBrokers()
	if len(brokers) == 0 {
		fmt.Println("No broker aliases configured.")
		fmt.Println("\nTo add an alias, use:")
		fmt.Println("  ./kafkaDebugger -add-alias name:broker1:9092,broker2:9092")
		return
	}

	fmt.Println("Configured broker aliases:")
	for _, broker := range brokers {
		defaultMarker := ""
		if broker.Name == cfg.DefaultAlias {
			defaultMarker = " (default)"
		}
		fmt.Printf("  %s%s:\n", broker.Name, defaultMarker)
		for _, b := range broker.Brokers {
			fmt.Printf("    - %s\n", b)
		}
	}

	// Show config file location
	configPath, _ := config.GetConfigPath()
	fmt.Printf("\nConfig file: %s\n", configPath)
}

func handleRemoveAlias(cfg *config.Config, name string) {
	if err := cfg.RemoveBroker(name); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := config.SaveConfig(cfg); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Successfully removed alias '%s'\n", name)
}
