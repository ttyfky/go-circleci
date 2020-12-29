package circleci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	UserAgent = "gocircleci/1.0.0"

	queryLimit         = 100 // maximum that CircleCI allows
	defaultHTTPTimeout = 20
	defaultPathPrefix  = "/api/v2/"
)

var (
	defaultBaseURL = &url.URL{Host: "circleci.com", Scheme: "https"}
	defaultLogger  = log.New(os.Stderr, "", log.LstdFlags)
)

// APIError represents an error from CircleCI
type APIError struct {
	HTTPStatusCode int
	Message        string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.HTTPStatusCode, e.Message)
}

// Client is a CircleCI client.
type Client struct {
	// CircleCI API endpoint (defaults to DefaultEndpoint)
	BaseURL    *url.URL
	pathPrefix string
	// HTTPClient to use for connecting to CircleCI (defaults to http.DefaultClient)
	HTTPClient *http.Client
	token      string

	Project  ProjectService
	EnvVar   ProjectEnvVarService
	Workflow WorkflowService
	Job      JobService
	Context  ContextService
}

// NewClient creates new CircleCI client with given API token.
// Optionally HTTP client or some other fields can be also given.
func NewClient(token string, opts ...Option) *Client {
	c := &Client{
		BaseURL: defaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * defaultHTTPTimeout,
		},
		token:      token,
		pathPrefix: defaultPathPrefix,
	}

	// Apply options for client.
	for _, o := range opts {
		o(c)
	}

	c.Project = &ProjectServiceOp{client: c}
	c.EnvVar = &ProjectEnvVarOp{client: c}
	c.Workflow = &WorkflowOp{client: c}
	c.Job = &JobOp{client: c}
	c.Context = &ContextOp{client: c}
	return c
}

// NewRequest creates a new http.Request with given parameters.
func (c *Client) NewRequest(method, path string, body, opts interface{}) (req *http.Request, err error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	if opts != nil {
		optsQuery, err := query.Values(opts)
		if err != nil {
			return nil, err
		}

		for k, values := range u.Query() {
			for _, v := range values {
				optsQuery.Add(k, v)
			}
		}
		u.RawQuery = optsQuery.Encode()
	}

	// A bit of JSON ceremony
	var js []byte = nil

	if body != nil {
		js, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err = http.NewRequest(method, u.String(), bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", UserAgent)
	if c.token != "" {
		req.Header.Add("Circle-Token", c.token)
	} else {
		return nil, fmt.Errorf("API token must not be blank")
	}
	return req, nil
}

// CreateAndDo performs a web request to CircleCI with the given method (GET,
// POST, PUT, DELETE) and relative path.
// If the data argument is non-nil, it will be used as the body of the request
// for POST and PUT requests.
// The options argument is used for specifying request options.
// Any data returned from CircleCI will be marshalled into resource argument.
func (c *Client) CreateAndDo(method, relPath string, data, options, resource interface{}) error {
	if strings.HasPrefix(relPath, "/") {
		// make sure it's a relative path
		relPath = strings.TrimLeft(relPath, "/")
	}
	relPath = path.Join(c.pathPrefix, relPath)

	req, err := c.NewRequest(method, relPath, data, options)
	if err != nil {
		return err
	}

	return c.do(req, resource)
}

// doGetHeaders executes a request, decoding the response into `v` and also returns any response headers.
func (c *Client) do(req *http.Request, v interface{}) error {
	var resp *http.Response
	var err error

	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	// retry scenario, close resp and any continue will retry
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &APIError{HTTPStatusCode: resp.StatusCode,
			Message: fmt.Sprintf("unable to read response body: %s", err),
		}
	}

	if resp.StatusCode >= 300 {
		if len(body) > 0 {
			message := Message{}
			err = json.Unmarshal(body, &message)
			if err != nil {
				return &APIError{
					HTTPStatusCode: resp.StatusCode,
					Message:        fmt.Sprintf("unable to parse API response: %s", err),
				}
			}
			return &APIError{HTTPStatusCode: resp.StatusCode, Message: message.Message}
		}

		return &APIError{HTTPStatusCode: resp.StatusCode}
	}

	if v != nil {
		err = json.Unmarshal(body, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get performs a GET request for the given path and saves the result in the
// given resource.
func (c *Client) Get(path string, resource, options interface{}) error {
	return c.CreateAndDo("GET", path, nil, options, resource)
}

// Post performs a POST request for the given path and saves the result in the
// given resource.
func (c *Client) Post(path string, data, resource interface{}) error {
	return c.CreateAndDo("POST", path, data, nil, resource)
}

// Put performs a PUT request for the given path and saves the result in the
// given resource.
func (c *Client) Put(path string, data, resource interface{}) error {
	return c.CreateAndDo("PUT", path, data, nil, resource)
}

// Delete performs a DELETE request for the given path
func (c *Client) Delete(path string) error {
	return c.CreateAndDo("DELETE", path, nil, nil, nil)
}

// Message represents messages.
type Message struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

// ProjectSlug assemle ProjectSlug of CircleCI.
// projectType: bitbucket, github(gh)
// org: organization name or user nme
// repo: repository name
func ProjectSlug(projectType, org, repo string) string {
	return fmt.Sprintf("%s/%s/%s", projectType, org, repo)
}
