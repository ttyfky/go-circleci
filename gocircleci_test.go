package circleci_test

import (
	"fmt"
	"testing"

	"github.com/ttyfky/go-circleci"
)

func TestNewClient(t *testing.T) {
	token := "test_token"
	testClient := circleci.NewClient(token)
	expected := "https://circleci.com"
	if testClient.BaseURL.String() != expected {
		t.Errorf("NewClient BaseURL = %v, expected %v", testClient.BaseURL.String(), expected)
	}
}

func TestProjectSlug(t *testing.T) {
	projectType := "gh"
	org := "ttyfky"
	repo := "go-circleci"
	expected := fmt.Sprintf("%s/%s/%s", projectType, org, repo)
	actual := circleci.ProjectSlug(projectType, org, repo)
	if expected != actual {
		t.Errorf("Invalid ProjectSlug. Expected: %s, Actual:%s", expected, actual)
	}
}
