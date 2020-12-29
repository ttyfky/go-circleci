package circleci_test

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/ttyfky/go-circleci"
)

func ExampleProjectServiceOp_Get() {
	token := os.Getenv("CIRCLECI_TOKEN")
	client := circleci.NewClient(token)

	project, err := client.Project.Get(projectSlug())
	if err != nil {
		log.Fatal(err)
	}
	printPretty(project)
}

func ExampleProjectEnvVarOp_List() {
	token := os.Getenv("CIRCLECI_TOKEN")
	client := circleci.NewClient(token)

	envVarList, err := client.EnvVar.List(projectSlug())
	if err != nil {
		log.Fatal(err)
	}
	printPretty(envVarList)
}

func ExampleWorkflowOp_Get() {
	token := os.Getenv("CIRCLECI_TOKEN")
	client := circleci.NewClient(token)

	workflowID := "ID"
	workflow, err := client.Workflow.Get(workflowID)
	if err != nil {
		log.Fatal(err)
	}
	printPretty(workflow)
}

func ExampleContextOp_Get() {
	token := os.Getenv("CIRCLECI_TOKEN")
	client := circleci.NewClient(token)

	id := "context_id"
	contextList, err := client.Context.Get(id)
	if err != nil {
		log.Fatal(err)
	}
	printPretty(contextList)
}

func ExampleContextOp_Create() {
	token := os.Getenv("CIRCLECI_TOKEN")
	client := circleci.NewClient(token)

	id := "context_id"

	envVarName := "test_key"
	envVar, err := client.Context.UpsertEnvVar(id, envVarName, "test_value")
	println("Created env Var")
	printPretty(envVar)
	contextList, err := client.Context.ListEnvVar(id)
	if err != nil {
		log.Fatal(err)
	}
	printPretty(contextList)
	err = client.Context.RemoveEnvVar(id, envVarName)
	if err != nil {
		log.Fatal(err)
	}
	println("Deleted env var")
	contextList, err = client.Context.ListEnvVar(id)
	if err != nil {
		log.Fatal(err)
	}
	printPretty(contextList)

}

func projectSlug() string {
	projectType := "gh"
	org := "ttyfky"
	repo := "go-circleci"
	return circleci.ProjectSlug(projectType, org, repo)
}

func printPretty(obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Fatal(err)
	}
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, b, "", "\t")
	println(prettyJSON.String())
}
