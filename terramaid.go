package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

var filePath string
var mermaidDiagram strings.Builder

type ResourceMode string

func init() {
	flag.StringVar(&filePath, "file", "", "Path to the Terraform plan JSON file")
}

// Define the structure of the Terraform plan JSON
type Plan struct {
	PlannedValues struct {
		RootModule struct {
			Resources []StateResource `json:"resources"`
		} `json:"root_module"`
	} `json:"planned_values"`
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type StateResource struct {
	Address         string                 `json:"address,omitempty"`
	Mode            ResourceMode           `json:"mode,omitempty"`
	Type            string                 `json:"type,omitempty"`
	Name            string                 `json:"name,omitempty"`
	Index           interface{}            `json:"index,omitempty"`
	ProviderName    string                 `json:"provider_name,omitempty"`
	SchemaVersion   uint64                 `json:"schema_version,"`
	AttributeValues map[string]interface{} `json:"values,omitempty"`
	SensitiveValues json.RawMessage        `json:"sensitive_values,omitempty"`
	DependsOn       []string               `json:"depends_on,omitempty"`
	Tainted         bool                   `json:"tainted,omitempty"`
	DeposedKey      string                 `json:"deposed_key,omitempty"`
}

type ResourceChange struct {
	Address string `json:"address"`
	Mode    string `json:"mode"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Change  Change `json:"change"`
}

type Change struct {
	Actions []string `json:"actions"`
}

func main() {
	flag.Parse()

	if filePath == "" {
		fmt.Println("Error: file path is required")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var plan Plan
	err = json.NewDecoder(file).Decode(&plan)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Write the initial graph definition
	mermaidDiagram.WriteString("```mermaid\ngraph TD\n")

	// Iterate through resource changes and add nodes and edges
	for _, rc := range plan.ResourceChanges {
		if len(rc.Change.Actions) > 0 && (rc.Change.Actions[0] == "plan" || rc.Change.Actions[0] == "update" || rc.Change.Actions[0] == "create") {
			// Add a node for the resource
			resourceNode := fmt.Sprintf("%s(\"%s\"):::type_%s\n", rc.Name, rc.Address, rc.Type)
			mermaidDiagram.WriteString(resourceNode)

			// Add edges for dependencies
			for _, res := range plan.PlannedValues.RootModule.Resources {
				if res.Address == rc.Address {
					for _, dep := range res.DependsOn {
						edge := fmt.Sprintf("%s --> %s\n", rc.Name, dep)
						mermaidDiagram.WriteString(edge)
					}
					break
				}
			}
		}
	}

	// Close the mermaid diagram definition
	mermaidDiagram.WriteString("```")

	// Write the diagram to a file
	err = os.WriteFile("Terramaid.md", []byte(mermaidDiagram.String()), 0644)
	if err != nil {
		fmt.Println("Error writing diagram file:", err)
		return
	}

	fmt.Println("Mermaid diagram generated successfully")
}
