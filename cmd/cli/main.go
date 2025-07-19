package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	
	"github.com/flack/chaos-engineering-as-a-platform/pkg/config"
)

var rootCmd = &cobra.Command{
	Use:   "chaos-cli",
	Short: "Chaos Engineering Platform CLI",
	Long:  `Command line interface for the Chaos Engineering Platform`,
}

var createExperimentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new chaos experiment",
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation will be added later
		fmt.Println("Creating new experiment...")
	},
}

var listExperimentsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all chaos experiments",
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation will be added later
		fmt.Println("Listing experiments...")
	},
}

var executeExperimentCmd = &cobra.Command{
	Use:   "execute [experiment-id]",
	Short: "Execute a chaos experiment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation will be added later
		fmt.Printf("Executing experiment %s...\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(createExperimentCmd)
	rootCmd.AddCommand(listExperimentsCmd)
	rootCmd.AddCommand(executeExperimentCmd)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Set global config for CLI
	_ = cfg

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}