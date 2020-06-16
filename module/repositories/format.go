package repositories

import (
	"github.com/codingXiang/go-harbor-client/client"
	"time"
)

// VulnerabilityItem is an item in the vulnerability result returned by vulnerability details API.
type VulnerabilityItem struct {
	ID          string `json:"id"`
	Severity    int64  `json:"severity"`
	Pkg         string `json:"package"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Fixed       string `json:"fixedVersion,omitempty"`
}

type RepoResp struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	ProjectID    int64     `json:"project_id"`
	Description  string    `json:"description"`
	PullCount    int64     `json:"pull_count"`
	StarCount    int64     `json:"star_count"`
	TagsCount    int64     `json:"tags_count"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

// RepoRecord holds the record of an repository in DB, all the infors are from the registry notification event.
type RepoRecord struct {
	RepositoryID int64     `json:"id"`
	Name         string    `json:"name"`
	ProjectID    int64     `json:"project_id"`
	Description  string    `json:"description"`
	PullCount    int64     `json:"pull_count"`
	StarCount    int64     `json:"star_count"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type cfg struct {
	Labels map[string]string `json:"labels"`
}

//ComponentsOverview has the total number and a list of components number of different serverity level.
type ComponentsOverview struct {
	Total   int                        `json:"total"`
	Summary []*ComponentsOverviewEntry `json:"summary"`
}

//ComponentsOverviewEntry ...
type ComponentsOverviewEntry struct {
	Sev   int `json:"severity"`
	Count int `json:"count"`
}

//ImgScanOverview mapped to a record of image scan overview.
type ImgScanOverview struct {
	ID              int64               `json:"-"`
	Digest          string              `json:"image_digest"`
	Status          string              `json:"scan_status"`
	JobID           int64               `json:"job_id"`
	Sev             int                 `json:"severity"`
	CompOverviewStr string              `json:"-"`
	CompOverview    *ComponentsOverview `json:"components,omitempty"`
	DetailsKey      string              `json:"details_key"`
	CreationTime    time.Time           `json:"creation_time,omitempty"`
	UpdateTime      time.Time           `json:"update_time,omitempty"`
}

type tagDetail struct {
	Digest        string    `json:"digest"`
	Name          string    `json:"name"`
	Size          int64     `json:"size"`
	Architecture  string    `json:"architecture"`
	OS            string    `json:"os"`
	DockerVersion string    `json:"docker_version"`
	Author        string    `json:"author"`
	Created       time.Time `json:"created"`
	Config        *cfg      `json:"config"`
}

type Signature struct {
	Tag    string            `json:"tag"`
	Hashes map[string][]byte `json:"hashes"`
}

type TagResp struct {
	tagDetail
	Signature    *Signature       `json:"signature"`
	ScanOverview *ImgScanOverview `json:"scan_overview,omitempty"`
}

type ListRepositoriesOption struct {
	client.ListOptions
	ProjectId int64  `url:"project_id,omitempty" json:"project_id,omitempty"`
	Q         string `url:"q,omitempty" json:"q,omitempty"`
	Sort      string `url:"sort,omitempty" json:"sort,omitempty"`
}

type ManifestResp struct {
	Manifest interface{} `json:"manifest"`
	Config   interface{} `json:"config,omitempty" `
}
