package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

type Plan struct {
	FormatVersion    string                 `json:"format_version"`
	TerraformVersion string                 `json:"terraform_version"`
	PlannedValues    PlannedValues          `json:"planned_values"`
	ResourceChanges  []interface{}          `json:"resource_changes"`
	OutputChanges    map[string]interface{} `json:"output_changes"`
	PriorState       map[string]interface{} `json:"prior_state"`
	Configuration    map[string]interface{} `json:"configuration"`
}

type PlannedValues struct {
	RootModule struct {
		Resources []Resource `json:"resources"`
	} `json:"root_module"`
}

type Resource struct {
	Address      string                 `json:"address"`
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	ProviderName string                 `json:"provider_name"`
	Values       map[string]interface{} `json:"values"`
}

// DependencyGraph represents a graph of dependencies between resources
type DependencyGraph map[string][]string

func main() {
	// Parse command line arguments for Terraform plan file
	var planFile string
	flag.StringVar(&planFile, "planfile", "tfplan.json", "Path to the Terraform plan JSON file")
	flag.Parse()

	// Read the Terraform plan file
	jsonFile, err := os.Open(planFile)
	if err != nil {
		log.Fatalf("Error opening JSON file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Error reading Terraform plan file: %v", err)
	}

	// Unmarshal JSON content into struct
	var plan Plan
	err = json.Unmarshal(byteValue, &plan)
	if err != nil {
		log.Fatalf("Error unmarshaling Terraform plan: %v", err)
	}

	// Construct the dependency graph
	graph := constructDependencyGraph(plan)

	// Generate the diagram
	err = generateDiagram(graph)
	if err != nil {
		log.Fatalf("Error generating Mermaid diagram: %v", err)
	}
}

func constructDependencyGraph(plan Plan) DependencyGraph {
	graph := make(DependencyGraph)
	for _, resource := range plan.PlannedValues.RootModule.Resources {
		address := resource.Address
		dependencies := extractDependencies(resource.Values)

		log.Printf("Resource: %s, Dependencies: %v\n", address, dependencies)

		graph[address] = dependencies
	}
	return graph
}

func extractDependencies(values map[string]interface{}) []string {
	var dependencies []string
	// Regular expression to match Terraform references
	re := regexp.MustCompile(`\${([^}]+)\.([^}]+)\.([^}]+)}`)

	for _, value := range values {
		if strVal, ok := value.(string); ok {
			// Find all references in the string value
			matches := re.FindAllStringSubmatch(strVal, -1)
			for _, match := range matches {
				if len(match) == 4 {
					// Extract the resource address from the match
					resourceAddr := fmt.Sprintf("%s.%s", match[1], match[2])
					dependencies = append(dependencies, resourceAddr)
				}
			}
		}
	}
	return dependencies
}

func generateDiagram(graph DependencyGraph) error {
	terraFile, err := os.Create("./Terramaid.md")
	if err != nil {
		return fmt.Errorf("error creating Terramaid file: %v", err)
	}
	defer terraFile.Close()

	fmt.Fprintln(terraFile, "```mermaid")
	fmt.Fprintln(terraFile, "graph TD;")
	for address, dependencies := range graph {
		for _, dep := range dependencies {
			fmt.Fprintf(terraFile, "%s --> %s;\n", dep, address)
		}
	}
	fmt.Fprintln(terraFile, "```")

	fmt.Println("Mermaid diagram generated successfully.")
	return nil
}
