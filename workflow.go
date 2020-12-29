package circleci

import "time"

const workflowBasePath = "/workflow"

// WorkflowService is an interface for Workflow API.
type WorkflowService interface {
	Get(id string) (*Workflow, error)
	Approve(id, approvalReqID string) (*Message, error)
	Cancel(id string) (*Message, error)
	GetJobs(id string) (*WorkflowJobs, error)
	Rerun(id string, jobIDs []string, fromFailed bool) (*Message, error)
}

// WorkflowOp handles communication with the project related methods in the CircleCI API v2.
type WorkflowOp struct {
	client *Client
}

var _ WorkflowService = (*WorkflowOp)(nil)

// Workflow represents information workflow.
type Workflow struct {
	PipelineID     string    `json:"pipeline_id,omitempty"`
	CanceledBy     string    `json:"canceled_by,omitempty"`
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	ProjectSlug    string    `json:"project_slug,omitempty"`
	ErroredBy      string    `json:"errored_by,omitempty"`
	Tag            string    `json:"tag,omitempty"`
	Status         string    `json:"status,omitempty"`
	StartedBy      string    `json:"started_by,omitempty"`
	PipelineNumber int       `json:"pipeline_number,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	StoppedAt      time.Time `json:"stopped_at,omitempty"`
}

// WorkflowJobs is jobs belongs to a Workflow.
type WorkflowJobs struct {
	Items []struct {
		CanceledBy        string      `json:"canceled_by,omitempty"`
		Dependencies      []string    `json:"dependencies,omitempty"`
		JobNumber         int         `json:"job_number,omitempty"`
		ID                string      `json:"id,omitempty"`
		StartedAt         time.Time   `json:"started_at,omitempty"`
		Name              string      `json:"name,omitempty"`
		ApprovedBy        string      `json:"approved_by,omitempty"`
		ProjectSlug       string      `json:"project_slug,omitempty"`
		Status            interface{} `json:"status,omitempty"`
		Type              string      `json:"type,omitempty"`
		StoppedAt         time.Time   `json:"stopped_at,omitempty"`
		ApprovalRequestID string      `json:"approval_request_id,omitempty"`
	} `json:"items,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
}

// RerunJob is a payload to send when rerunning jobs in Workflow.
// Multiple job ids can be given in Jobs.
type RerunJob struct {
	Jobs       []string `json:"jobs,omitempty"`
	FromFailed bool     `json:"from_failed,omitempty"`
}

// Get gets detail of workflow.
func (ps *WorkflowOp) Get(id string) (*Workflow, error) {
	w := &Workflow{}
	path := workflowBasePath + "/" + id
	err := ps.client.Get(path, w, nil)
	if err != nil {
		return nil, err
	}
	return w, nil
}

// Approve approves pending workflow job.
func (ps *WorkflowOp) Approve(id, approvalReqID string) (*Message, error) {
	m := &Message{}
	path := workflowBasePath + "/" + id + "/approve/" + approvalReqID
	err := ps.client.Post(path, nil, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Cancel cancels given workflow.
func (ps *WorkflowOp) Cancel(id string) (*Message, error) {
	m := &Message{}
	path := workflowBasePath + "/" + id + "/cancel"
	err := ps.client.Post(path, nil, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetJobs get jobs in the workflow.
func (ps *WorkflowOp) GetJobs(id string) (*WorkflowJobs, error) {
	wj := &WorkflowJobs{}
	path := workflowBasePath + "/" + id + "/job"
	err := ps.client.Post(path, nil, wj)
	if err != nil {
		return nil, err
	}
	return wj, nil
}

// Rerun cancels given workflow
func (ps *WorkflowOp) Rerun(id string, jobIDs []string, fromFailed bool) (*Message, error) {
	m := &Message{}
	path := workflowBasePath + "/" + id + "/rerun"
	err := ps.client.Post(path, &RerunJob{
		Jobs:       jobIDs,
		FromFailed: fromFailed,
	}, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
