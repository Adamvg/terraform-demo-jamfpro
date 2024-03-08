package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// Define the top-level TerraformPlan struct
type TerraformPlan struct {
	FormatVersion    string              `json:"format_version"`
	TerraformVersion string              `json:"terraform_version"`
	Variables        map[string]Variable `json:"variables"`
	PlannedValues    PlannedValues       `json:"planned_values"`
	ResourceChanges  []ResourceChange    `json:"resource_changes"`
	Configuration    Configuration       `json:"configuration"`
	Timestamp        string              `json:"timestamp"`
	Errored          bool                `json:"errored"`
}

// Variables

type Variable struct {
	Value       interface{} `json:"value"`
	Description string      `json:"description,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

// Planned Values

type PlannedValues struct {
	RootModule RootModule `json:"root_module"`
}

// Define a struct for the RootModule part
type RootModule struct {
	Resources []Resource `json:"resources"`
}

// Define a struct for each Resource
type Resource struct {
	Address string         `json:"address"`
	Type    string         `json:"type"`
	Values  ResourceValues `json:"values"`
}

// Define a struct for the Values part
type ResourceValues struct {
	Name string `json:"name"`
}

// Resource Change

type ResourceChange struct {
	Address  string `json:"address"`
	Mode     string `json:"mode"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Provider string `json:"provider_name"`
	Change   Change `json:"change"`
}

type Change struct {
	Actions         []string               `json:"actions"`
	Before          map[string]interface{} `json:"before"`
	After           map[string]interface{} `json:"after"`
	AfterUnknown    map[string]interface{} `json:"after_unknown"`
	BeforeSensitive bool                   `json:"before_sensitive"`
	AfterSensitive  map[string]interface{} `json:"after_sensitive"`
}

// Configuration

type Configuration struct {
	ProviderConfig map[string]ProviderConfig `json:"provider_config"`
	RootModule     RootModuleConfig          `json:"root_module"`
}

type ProviderConfig struct {
	Name              string              `json:"name"`
	FullName          string              `json:"full_name"`
	VersionConstraint string              `json:"version_constraint"`
	Expressions       ProviderExpressions `json:"expressions"`
}

type ProviderExpressions struct {
	ClientID     Expression `json:"client_id"`
	ClientSecret Expression `json:"client_secret"`
	InstanceName Expression `json:"instance_name"`
	LogLevel     Expression `json:"log_level"`
}

type Expression struct {
	ConstantValue string   `json:"constant_value,omitempty"`
	References    []string `json:"references,omitempty"`
}

type RootModuleConfig struct {
	Resources []ResourceConfig          `json:"resources"`
	Variables map[string]VariableConfig `json:"variables"`
}

type ResourceConfig struct {
	Address           string      `json:"address"`
	Mode              string      `json:"mode"`
	Type              string      `json:"type"`
	Name              string      `json:"name"`
	ProviderConfigKey string      `json:"provider_config_key"`
	Expressions       Expressions `json:"expressions"`
	SchemaVersion     int         `json:"schema_version"`
}

type Expressions struct {
	Name Expression `json:"name"`
}

type VariableConfig struct {
	Default     interface{} `json:"default"`
	Description string      `json:"description,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

func main() {
	tfPlanPath := flag.String("tfplan", "", "Path to the Terraform plan file in JSON format")
	flag.Parse()

	if *tfPlanPath == "" {
		fmt.Println("Usage: -tfplan <path to terraform plan json>")
		return
	}

	planFile, err := os.ReadFile(*tfPlanPath)
	if err != nil {
		fmt.Printf("Error reading plan file: %v\n", err)
		return
	}

	var plan TerraformPlan
	err = json.Unmarshal(planFile, &plan)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	// Define your security-related conditions here
	securityResources := map[string]bool{
		"jamfpro_api_integration":               true,
		"jamfpro_disk_encryption_configuration": true,
		// Add more resources or properties that you consider security-related
	}

	securityChangesDetected := false

	for _, change := range plan.ResourceChanges {
		// Check if the resource type is one of the security related resources
		if _, ok := securityResources[change.Type]; ok {
			// Check the actions for create, update, or delete
			for _, action := range change.Change.Actions {
				if action == "create" || action == "update" || action == "delete" {
					securityChangesDetected = true
					fmt.Printf("Security-related change detected: %s action on %s\n", action, change.Address)
					break // Break out of the inner loop once a security-related change is found
				}
			}
			if securityChangesDetected {
				break // Break out of the outer loop once a security-related change is found
			}
		}
	}

	if securityChangesDetected {
		fmt.Println("Security-related changes detected in the terraform plan. Setting the 'Security' group for the GitHub PR approval.")
		// Set the GitHub Actions environment variable for the approval group
		fmt.Println("::set-output name=approval_group::Security")
	} else {
		fmt.Println("No security-related changes detected.")
	}
}