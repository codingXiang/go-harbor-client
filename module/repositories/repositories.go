package repositories

import (
	client2 "github.com/codingXiang/go-harbor-client/client"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

const (
	repo               = "api.repositories."
	root               = repo + "root"
	base               = repo + "base"
	tagRoot            = repo + "tags.root"
	tagBase            = repo + "tags.base"
	tagManifest        = repo + "tags.manifest.root"
	tagManifestVersion = repo + "tags.manifest.version"
	labelRoot          = repo + "label.root"
	labelBase          = repo + "label.base"
	signatures         = repo + "signatures"
	top                = repo + "top"
)

// RepositoriesService handles communication with the user related methods of
// the Harbor API.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L891
type Service interface {
	List(opt *ListRepositoriesOption) ([]RepoRecord, *gorequest.Response, []error)
	Update(name string, d RepositoryDescription) (*gorequest.Response, []error)
	Delete(name string) (*gorequest.Response, []error)
	GetTag(projectName string, repoName string, tag string) (TagResp, *gorequest.Response, []error)
	DeleteTag(projectName string, repoName string, tag string) (*gorequest.Response, []error)
	ListTags(projectName string, repoName string) ([]TagResp, *gorequest.Response, []error)
	GetTagManifests(name string, tag string, version string) (ManifestResp, *gorequest.Response, []error)
	ScanImage(name string, tag string) (*gorequest.Response, []error)
	GetImageDetails(name string, tag string) ([]VulnerabilityItem, *gorequest.Response, []error)
	GetSignature(name string) ([]Signature, *gorequest.Response, []error)
	GetTop(top interface{}) ([]RepoResp, *gorequest.Response, []error)
}

type RepositoriesService struct {
	client client2.ClientInterface
}

func NewRepositoriesService(client client2.ClientInterface) Service {
	return &RepositoriesService{client: client}
}

func (s *RepositoriesService) getConfigString(key string) string {
	return s.client.GetConfig().GetString(key)
}

// Get repositories accompany with relevant project and repo name.
//
// This endpoint let user search repositories accompanying
// with relevant project ID and repo name.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L892
func (s *RepositoriesService) List(opt *ListRepositoriesOption) ([]RepoRecord, *gorequest.Response, []error) {
	var v []RepoRecord
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, s.getConfigString(root)).
		Query(*opt).
		EndStruct(&v)
	return v, &resp, errs
}

// Delete a repository.
//
// This endpoint let user delete a repository with name.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L948
func (s *RepositoriesService) Delete(repoName string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf(s.getConfigString(base), repoName)).
		End()
	return &resp, errs
}

type RepositoryDescription struct {
	Description string `url:"description,omitempty" json:"description,omitempty"`
}

// Update description of the repository.
//
// This endpoint is used to update description of the repository.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L971
func (s *RepositoriesService) Update(repoName string, d RepositoryDescription) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.PUT, fmt.Sprintf(s.getConfigString(base), repoName)).
		Send(d).
		End()
	return &resp, errs
}

// Get the tag of the repository.
//
// This endpoint aims to retrieve the tag of the repository.
// If deployed with Notary, the signature property of
// response represents whether the image is singed or not.
// If the property is null, the image is unsigned.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L999
func (s *RepositoriesService) GetTag(projectName string, repoName, tag string) (TagResp, *gorequest.Response, []error) {
	var v TagResp
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf(s.getConfigString(tagBase), projectName, repoName, tag)).
		EndStruct(&v)
	return v, &resp, errs
}

// Delete a tag in a repository.
//
// This endpoint let user delete tags with repo name and tag.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1025
func (s *RepositoriesService) DeleteTag(projectName string, repoName, tag string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf(s.getConfigString(tagBase), projectName, repoName, tag)).
		End()
	return &resp, errs
}

// Get tags of a relevant repository.
//
// This endpoint aims to retrieve tags from a relevant repository.
// If deployed with Notary, the signature property of
// response represents whether the image is singed or not.
// If the property is null, the image is unsigned.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1054
func (s *RepositoriesService) ListTags(projectName string, repoName string) ([]TagResp, *gorequest.Response, []error) {
	var v []TagResp
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf(s.getConfigString(tagRoot), projectName, repoName)).
		EndStruct(&v)
	return v, &resp, errs
}

// Get manifests of a relevant repository.
//
// This endpoint aims to retreive manifests from a relevant repository.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1079
func (s *RepositoriesService) GetTagManifests(repoName, tag string, version string) (ManifestResp, *gorequest.Response, []error) {
	var v ManifestResp
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, func() string {
			if version == "" {
				return fmt.Sprintf(s.getConfigString(tagManifest), repoName, tag)
			}
			return fmt.Sprintf(s.getConfigString(tagManifestVersion), repoName, tag, version)
		}()).
		EndStruct(&v)
	return v, &resp, errs
}

// Scan the image.
//
// Trigger jobservice to call Clair API to scan the image
// identified by the repo_name and tag.
// Only project admins have permission to scan images under the project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1113
func (s *RepositoriesService) ScanImage(repoName, tag string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.POST, fmt.Sprintf("/repositories/%s/tags/%s/scan", repoName, tag)).
		End()
	return &resp, errs
}

// Get vulnerability details of the image.
//
// Call Clair API to get the vulnerability based on the previous successful scan.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1177
func (s *RepositoriesService) GetImageDetails(repoName, tag string) ([]VulnerabilityItem, *gorequest.Response, []error) {
	var v []VulnerabilityItem
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("repositories/%s/tags/%s/vulnerability/details", repoName, tag)).
		EndStruct(&v)
	return v, &resp, errs
}

// Get signature information of a repository.
//
// This endpoint aims to retrieve signature information of a repository, the data is
// from the nested notary instance of Harbor.
// If the repository does not have any signature information in notary, this API will
// return an empty list with response code 200, instead of 404
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1211
func (s *RepositoriesService) GetSignature(repoName string) ([]Signature, *gorequest.Response, []error) {
	var v []Signature
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("repositories/%s/signatures", repoName)).
		EndStruct(&v)
	return v, &resp, errs
}

// Get public repositories which are accessed most.
//
// This endpoint aims to let users see the most popular public repositories.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1241
func (s *RepositoriesService) GetTop(top interface{}) ([]RepoResp, *gorequest.Response, []error) {
	var v []RepoResp
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, func() string {
			if t, ok := top.(int); ok {
				return fmt.Sprintf("repositories/top?count=%d", t)
			}
			return fmt.Sprintf("repositories/top")
		}()).
		EndStruct(&v)
	return v, &resp, errs
}
