package main

import (
	"fmt"
	"os"
)

func main() {
	var (
		command     string
		commandArgs []string
	)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Missing experiment name.")
		os.Exit(1)
	}

	command = args[1]
	if len(args) > 2 {
		commandArgs = args[2:]
	}

	switch command {
	case "list":
		handleList(commandArgs)
	case "info":
		handleInfo(commandArgs)
	default:
		handleExperiment(command, commandArgs)
	}
}

// handleList writes to stdout the list of available experiments.
func handleList(cArgs []string) {
	if cArgs != nil {
		fmt.Println("Command 'list' does not take arguments.")
		os.Exit(1)
	}
	fmt.Println("\nExperiments:")
	for _, exp := range experiments {
		fmt.Printf("  - %s\n", exp.Name)
	}
	os.Exit(0)
}

// handleInfo writes to stdout the information related to the experiment
// denoted by cArgs[0].
func handleInfo(cArgs []string) {
	if cArgs == nil {
		fmt.Println("Command 'info' requires experiment names.")
		os.Exit(1)
	}
	expLookup := expMap()
	out := ""
	for _, arg := range cArgs {
		exp, ok := expLookup[arg]
		if !ok {
			fmt.Printf("Experiment '%s' does not exist.\n", arg)
			os.Exit(1)
		}
		out += exp.Description
	}
	fmt.Println()
	fmt.Println(out)
	os.Exit(0)
}

// handleExperiment runs the experiment denoted by c with arguments cArgs.
func handleExperiment(c string, cArgs []string) {
	exp, ok := expMap()[c]
	if !ok {
		fmt.Printf("Experiment '%s' does not exist.\n", c)
		os.Exit(1)
	}

	fmt.Println("Running experiment...")

	if err := exp.Run(cArgs...); err != nil {
		fmt.Printf("Error: %s.", err.Error())
		os.Exit(1)
	}

	fmt.Println("Done running.")
	os.Exit(0)
}
