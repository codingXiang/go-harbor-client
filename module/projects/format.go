package projects

import (
	"github.com/codingXiang/go-harbor-client/client"
	"time"
)

// ProjectMetadata holds the metadata of a project.
type ProjectMetadata struct {
	ID           int64     `json:"id"`
	ProjectID    int64     `json:"project_id"`
	Name         string    `json:"name"`
	Value        string    `json:"value"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
	Deleted      int       `json:"deleted"`
}

// Project holds the details of a project.
type Project struct {
	ProjectID    int64             `json:"project_id"`
	OwnerID      int               `json:"owner_id"`
	Name         string            `json:"name"`
	CreationTime time.Time         `json:"creation_time"`
	UpdateTime   time.Time         `json:"update_time"`
	Deleted      bool              `json:"deleted"`
	OwnerName    string            `json:"owner_name"`
	Togglable    bool              `json:"togglable"`
	Role         int               `json:"current_user_role_id"`
	RepoCount    int64             `json:"repo_count"`
	Metadata     map[string]string `json:"metadata"`
}

// AccessLog holds information about logs which are used to record the actions that user take to the resourses.
type AccessLog struct {
	LogID     int       `json:"log_id"`
	Username  string    `json:"username"`
	ProjectID int64     `json:"project_id"`
	RepoName  string    `json:"repo_name"`
	RepoTag   string    `json:"repo_tag"`
	GUID      string    `json:"guid"`
	Operation string    `json:"operation"`
	OpTime    time.Time `json:"op_time"`
}

// ProjectRequest holds informations that need for creating project API
type ProjectRequest struct {
	Name     string            `url:"name,omitempty" json:"project_name"`
	Public   *int              `url:"public,omitempty" json:"public"` //deprecated, reserved for project creation in replication
	Metadata map[string]string `url:"-" json:"metadata"`
}

type ListProjectsOptions struct {
	client.ListOptions
	Name   string `url:"name,omitempty" json:"name,omitempty"`
	Public bool   `url:"public,omitempty" json:"public,omitempty"`
	Owner  string `url:"owner,omitempty" json:"owner,omitempty"`
}

// LogQueryParam is used to set query conditions when listing
// access logs.
type ListLogOptions struct {
	client.ListOptions
	Username   string     `url:"username,omitempty"`        // the operator's username of the log
	Repository string     `url:"repository,omitempty"`      // repository name
	Tag        string     `url:"tag,omitempty"`             // tag name
	Operations []string   `url:"operation,omitempty"`       // operations
	BeginTime  *time.Time `url:"begin_timestamp,omitempty"` // the time after which the operation is done
	EndTime    *time.Time `url:"end_timestamp,omitempty"`   // the time before which the operation is doen
}

type MemberRequest struct {
	UserName string `json:"username"`
	Roles    []int  `json:"roles"`
}