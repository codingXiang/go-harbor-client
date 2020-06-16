package projects

import (
	client2 "github.com/codingXiang/go-harbor-client/client"
	"github.com/codingXiang/go-harbor-client/module/user"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

const (
	root          = "api.projects.root"
	base          = "api.projects.base"
	metadatasRoot = "api.projects.metadatas.root"
	metadatasBase = "api.projects.metadatas.base"
	logsRoot      = "api.projects.logs.root"
	membersRoot   = "api.projects.members.root"
	membersBase   = "api.projects.members.base"
)

// ProjectsService handles communication with the user related methods of
// the Harbor API.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L45
type Service interface {
	//列出所有專案
	List(opt *ListProjectsOptions) ([]Project, *gorequest.Response, []error)
	//取得特定專案
	Get(id int64) (Project, *gorequest.Response, []error)
	//建立專案
	Create(p *ProjectRequest) (*gorequest.Response, []error)
	//更新專案
	Update(id int64, p Project) (*gorequest.Response, []error)
	//刪除專案
	Delete(id int64) (*gorequest.Response, []error)
	//檢查專案
	Check(name string) (*gorequest.Response, []error)
	//取得 log
	GetLog(id int64, options ListLogOptions) ([]AccessLog, *gorequest.Response, []error)
	//透過 id 取得 metadata
	GetMetadataById(id int64) (map[string]string, *gorequest.Response, []error)
	//加入 metadata
	AddMetadata(id int64, medadata map[string]string) (*gorequest.Response, []error)
	//透過名稱取得 metadata
	GetMetadata(id int64, name string) (map[string]string, *gorequest.Response, []error)
	//更新 metadata
	UpdateMetadata(id int64, name string) (*gorequest.Response, []error)
	//刪除 metadata
	DeleteMetadata(id int64, name string) (*gorequest.Response, []error)
	//取得成員
	GetMembers(id int64) ([]user.User, *gorequest.Response, []error)
	//加入成員
	AddMember(id int64, member MemberRequest) (*gorequest.Response, []error)
	//更新成員角色
	UpdateMemberRole(id int, uid int, role MemberRequest) (*gorequest.Response, []error)
	//取得成員角色
	GetMemberRole(id int, uid int) (Role, *gorequest.Response, []error)
	//刪除成員
	DeleteMember(id int, uid int) (*gorequest.Response, []error)
}

type ProjectsService struct {
	client client2.ClientInterface
}

func NewProjectService(client client2.ClientInterface) Service {
	return &ProjectsService{client: client}
}

func (s *ProjectsService) getConfigString(key string) string {
	return s.client.GetConfig().GetString(key)
}

// List projects
//
// This endpoint returns all projects created by Harbor,
// and can be filtered by project name.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L46
func (s *ProjectsService) List(opt *ListProjectsOptions) ([]Project, *gorequest.Response, []error) {
	var projects []Project
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, s.getConfigString(root)).
		Query(*opt).
		EndStruct(&projects)
	return projects, &resp, errs
}

// Check if the project name user provided already exists.
//
// This endpoint is used to check if the project name user provided already exist.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L100
func (s *ProjectsService) Check(projectName string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.HEAD, s.getConfigString(root)).
		Query(fmt.Sprintf("project_name=%s", projectName)).
		End()
	return &resp, errs
}

// Create a new project.
//
// This endpoint is for user to create a new project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L122
func (s *ProjectsService) Create(p *ProjectRequest) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.POST, s.getConfigString(root)).
		Send(*p).
		End()
	return &resp, errs
}

// Return specific project detail information.
//
// This endpoint returns specific project information by project ID.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L149
func (s *ProjectsService) Get(pid int64) (Project, *gorequest.Response, []error) {
	var project Project
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf(s.getConfigString(base), pid)).
		EndStruct(&project)
	return project, &resp, errs
}

// Update properties for a selected project.
//
// This endpoint is aimed to update the properties of a project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L171
func (s *ProjectsService) Update(pid int64, p Project) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.PUT, fmt.Sprintf(s.getConfigString(base), pid)).
		Send(p).
		End()
	return &resp, errs
}

// Delete project by projectID.
//
// This endpoint is aimed to delete project by project ID.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L203
func (s *ProjectsService) Delete(pid int64) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf(s.getConfigString(base), pid)).
		End()
	return &resp, errs
}

// Get access logs accompany with a relevant project.
//
// This endpoint let user search access logs filtered by operations and date time ranges.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L230
func (s *ProjectsService) GetLog(pid int64, opt ListLogOptions) ([]AccessLog, *gorequest.Response, []error) {
	var accessLog []AccessLog
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf(s.getConfigString(logsRoot), pid)).
		Query(opt).
		EndStruct(&accessLog)
	return accessLog, &resp, errs
}

// Get project all metadata.
//
// This endpoint returns metadata of the project specified by project ID.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L307
func (s *ProjectsService) GetMetadataById(pid int64) (map[string]string, *gorequest.Response, []error) {
	var metadata map[string]string
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf(s.getConfigString(base), pid)).
		EndStruct(&metadata)
	return metadata, &resp, errs
}

// Add metadata for the project.
//
// This endpoint is aimed to add metadata of a project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L329
func (s *ProjectsService) AddMetadata(pid int64, metadata map[string]string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.POST, fmt.Sprintf(s.getConfigString(metadatasRoot), pid)).
		Send(metadata).
		End()
	return &resp, errs
}

// Get project metadata
//
// This endpoint returns specified metadata of a project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L364
func (s *ProjectsService) GetMetadata(pid int64, specified string) (map[string]string, *gorequest.Response, []error) {
	var metadata map[string]string
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf(s.getConfigString(membersBase), pid, specified)).
		EndStruct(&metadata)
	return metadata, &resp, errs
}

// Update metadata of a project.
//
// This endpoint is aimed to update the metadata of a project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L391
func (s *ProjectsService) UpdateMetadata(pid int64, metadataName string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.PUT, fmt.Sprintf(s.getConfigString(metadatasBase), pid, metadataName)).
		End()
	return &resp, errs
}

// Delete metadata of a project
//
// This endpoint is aimed to delete metadata of a project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L422
func (s *ProjectsService) DeleteMetadata(pid int64, metadataName string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf(s.getConfigString(metadatasBase), pid, metadataName)).
		End()
	return &resp, errs
}

// Return a project's relevant role members.
//
// This endpoint is for user to search a specified project’s relevant role members.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L452
func (s *ProjectsService) GetMembers(pid int64) ([]user.User, *gorequest.Response, []error) {
	var users []user.User
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf(s.getConfigString(membersRoot), pid)).
		EndStruct(&users)
	return users, &resp, errs
}

// Add project role member accompany with relevant project and user.
//
// This endpoint is for user to add project role member accompany with relevant project and user.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L483
func (s *ProjectsService) AddMember(pid int64, member MemberRequest) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.POST, fmt.Sprintf(s.getConfigString(membersRoot), pid)).
		Send(member).
		End()
	return &resp, errs
}

// Role holds the details of a role.
type Role struct {
	RoleID   int    `json:"role_id"`
	RoleCode string `json:"role_code"`
	Name     string `json:"role_name"`
	RoleMask int    `json:"role_mask"`
}

// Return role members accompany with relevant project and user.
//
// This endpoint is for user to get role members accompany with relevant project and user.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L522
func (s *ProjectsService) GetMemberRole(pid, uid int) (Role, *gorequest.Response, []error) {
	var role Role
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf(s.getConfigString(membersBase), pid, uid)).
		EndStruct(&role)
	return role, &resp, errs
}

// Update project role members accompany with relevant project and user.
//
// This endpoint is for user to update current project role members accompany with relevant project and user.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L559
func (s *ProjectsService) UpdateMemberRole(pid, uid int, role MemberRequest) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.PUT, fmt.Sprintf(s.getConfigString(membersBase), pid, uid)).
		Send(role).
		End()
	return &resp, errs
}

// Delete project role members accompany with relevant project and user.
//
// This endpoint is aimed to remove project role members already added to the relevant project and user.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L597
func (s *ProjectsService) DeleteMember(pid, uid int) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf(s.getConfigString(membersBase), pid, uid)).
		End()
	return &resp, errs
}
