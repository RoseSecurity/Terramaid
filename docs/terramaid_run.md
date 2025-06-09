## terramaid run

Generate Mermaid diagrams from Terraform configurations

```
terramaid run [flags]
```

### Options

```
  -c, --chart-type string      Specify the type of Mermaid chart to generate (env: TERRAMAID_CHART_TYPE) (default "flowchart")
  -r, --direction string       Specify the direction of the diagram (env: TERRAMAID_DIRECTION) (default "TD")
  -h, --help                   help for run
  -o, --output string          Output file for Mermaid diagram (env: TERRAMAID_OUTPUT) (default "Terramaid.md")
  -s, --subgraph-name string   Specify the subgraph name of the diagram (env: TERRAMAID_SUBGRAPH_NAME) (default "Terraform")
  -b, --tf-binary string       Path to Terraform binary (env: TERRAMAID_TF_BINARY)
  -p, --tf-plan string         Path to Terraform plan file (env: TERRAMAID_TF_PLAN)
  -t, --timeout duration       Timeout for the entire run (e.g. 5m) (env: TERRAMAID_TIMEOUT)
  -v, --verbose                Enable verbose output (env: TERRAMAID_VERBOSE)
  -w, --working-dir string     Working directory for Terraform (env: TERRAMAID_WORKING_DIR) (default ".")
```

### SEE ALSO

* [terramaid](terramaid.md)	 - A utility for generating Mermaid diagrams from Terraform configurations

