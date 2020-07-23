package client

import (
	"github.com/codingXiang/configer"
	"github.com/codingXiang/go-logger"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"net/url"
	"strings"
)

const (
	libraryVersion = "1.4.0"
	userAgent      = "go-harbor/" + libraryVersion
)

type ClientInterface interface {
	GetClient() *gorequest.SuperAgent
	GetConfig() *viper.Viper
	GetUserAgent() string
	NewRequest(method string, subPath string) *gorequest.SuperAgent
}

type Client struct {
	// HTTP client used to communicate with the API.
	client *gorequest.SuperAgent
	// Base URL for API requests. Defaults to the public GitLab API, but can be
	// set to a domain endpoint to use with a self hosted GitLab server. baseURL
	// should always be specified with a trailing slash.
	baseURL *url.URL
	config  *viper.Viper
	// User agent used when communicating with the GitLab API.
	userAgent string
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty" json:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PageSize int `url:"page_size,omitempty" json:"page_size,omitempty"`
}

func NewClient(config configer.CoreInterface) ClientInterface {
	return newClient(config)
}

// SetBaseURL sets the base URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) setBaseURL(urlStr string) error {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	var err error
	c.baseURL, err = url.Parse(urlStr)
	return err
}

// GetClient 取得 client
func (c *Client) GetClient() *gorequest.SuperAgent {
	return c.client
}

// GetConfig 取得設定檔
func (c *Client) GetConfig() *viper.Viper {
	return c.config
}

// GetUserAgent 取得 User Agent
func (c *Client) GetUserAgent() string {
	return c.userAgent
}

func newClient(config configer.CoreInterface) *Client {
	if data, err := config.ReadConfig(nil); err == nil {
		harborClient := gorequest.New()
		c := &Client{client: harborClient, userAgent: userAgent, config: data}
		var (
			baseURL  = data.GetString("ingress.protocol") + "://" + data.GetString("ingress.domain")
			username = data.GetString("management.user.name")
			password = data.GetString("management.user.password")
		)
		// 設定基礎驗證
		harborClient.SetBasicAuth(username, password)
		// 設定 harbor 位置
		if err := c.setBaseURL(baseURL); err == nil {
			logger.Log.Info("Harbor Server 設定完成，位置為", baseURL)
		} else {
			logger.Log.Error(err.Error())
		}
		return c
	} else {
		logger.Log.Error("設定組態檔發生錯誤", err.Error())
		return nil
	}

}

// NewRequest creates an API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, subPath string) *gorequest.SuperAgent {
	var (
		api = c.config.GetString("api.root")
	)
	u := c.baseURL.String() + api + subPath
	h := c.client.Set("Accept", "application/json")
	if c.userAgent != "" {
		h.Set("User-Agent", c.userAgent)
	}
	logger.Log.Debug("發起 Request", "["+method+"]", u)
	switch method {
	case gorequest.PUT:
		return c.client.Put(u).Set("Content-Type", "application/json")
	case gorequest.POST:
		return c.client.Post(u).Set("Content-Type", "application/json")
	case gorequest.GET:
		return c.client.Get(u)
	case gorequest.HEAD:
		return c.client.Head(u)
	case gorequest.DELETE:
		return c.client.Delete(u)
	case gorequest.PATCH:
		return c.client.Patch(u)
	case gorequest.OPTIONS:
		return c.client.Options(u)
	default:
		return c.client.Get(u)
	}
}

type SearchRepository struct {
	// The ID of the project that the repository belongs to
	ProjectId int32 `json:"project_id,omitempty"`
	// The name of the project that the repository belongs to
	ProjectName string `json:"project_name,omitempty"`
	// The flag to indicate the publicity of the project that the repository belongs to
	ProjectPublic bool `json:"project_public,omitempty"`
	// The name of the repository
	RepositoryName string `json:"repository_name,omitempty"`
	PullCount      int32  `json:"pull_count,omitempty"`
	TagsCount      int32  `json:"tags_count,omitempty"`
}

type StatisticMap struct {
	// The count of the private projects which the user is a member of.
	PrivateProjectCount int `json:"private_project_count,omitempty"`
	// The count of the private repositories belonging to the projects which the user is a member of.
	PrivateRepoCount int `json:"private_repo_count,omitempty"`
	// The count of the public projects.
	PublicProjectCount int `json:"public_project_count,omitempty"`
	// The count of the public repositories belonging to the public projects which the user is a member of.
	PublicRepoCount int `json:"public_repo_count,omitempty"`
	// The count of the total projects, only be seen when the is admin.
	TotalProjectCount int `json:"total_project_count,omitempty"`
	// The count of the total repositories, only be seen when the user is admin.
	TotalRepoCount int `json:"total_repo_count,omitempty"`
}

// Get projects number and repositories number relevant to the user
//
//This endpoint is aimed to statistic all of the projects number
// and repositories number relevant to the logined user,
// also the public projects number and repositories number.
// If the user is admin,
// he can also get total projects number and total repositories number.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L631
func (c *Client) GetStatistics() (StatisticMap, *gorequest.Response, []error) {
	var statistics StatisticMap
	resp, _, errs := c.NewRequest(gorequest.GET, "statistics").
		EndStruct(&statistics)
	return statistics, &resp, errs
}
