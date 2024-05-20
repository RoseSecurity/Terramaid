package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

// Define structures to parse Terraform plan JSON
type Plan struct {
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type ResourceChange struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Change  Change `json:"change"`
}

type Change struct {
	Actions []string `json:"actions"`
}

func main() {
	// Parse and read the plan file
	planfile := flag.String("planfile", "tfplan.json", "path to the Terraform plan file")
	flag.Parse()

	data, err := os.ReadFile(*planfile)
	if err != nil {
		log.Fatalf("Error reading plan file: %v\n", err)
	}

	var plan Plan
	err = json.Unmarshal(data, &plan)
	if err != nil {
		log.Fatalf("Error parsing plan file: %v\n", err)
	}

	// Write the Mermaid diagram to the output file
	outFile, err := os.Create("Terramaid.md")
	if err != nil {
		log.Fatalf("Error creating output file: %v\n", err)
	} else {
		fmt.Println("Terramaid file created")
		defer outFile.Close()
	}

	fmt.Fprintln(outFile, "```mermaid")
	fmt.Fprintln(outFile, "graph TD;")

	// Create a map to keep track of node names and their assigned variable
	nodeMap := make(map[string]string)
	varNameCounter := 0

	getVarName := func() string {
		varName := fmt.Sprint('A' + varNameCounter)
		varNameCounter++
		return varName
	}

	for _, rc := range plan.ResourceChanges {
		// Assign or retrieve variable names for the nodes
		sourceKey := rc.Type
		targetKey := rc.Address

		if _, exists := nodeMap[sourceKey]; !exists {
			nodeMap[sourceKey] = getVarName()
		}
		if _, exists := nodeMap[targetKey]; !exists {
			nodeMap[targetKey] = getVarName()
		}

		sourceVar := nodeMap[sourceKey]
		targetVar := nodeMap[targetKey]

		source := fmt.Sprintf("%s(%s)", sourceVar, rc.Type)
		target := fmt.Sprintf("%s(%s)", targetVar, rc.Address)

		if contains(rc.Change.Actions, "create") {
			fmt.Fprintf(outFile, "%s -->|created| %s\n", source, target)
		} else if contains(rc.Change.Actions, "update") {
			fmt.Fprintf(outFile, "%s -->|updated| %s\n", source, target)
		} else if contains(rc.Change.Actions, "delete") {
			fmt.Fprintf(outFile, "%s -->|deleted| %s\n", source, target)
		}
	}
	fmt.Fprintln(outFile, "```")
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
