package elestio

import (
	"encoding/json"
	"fmt"
)

type (
	// ProjectHandler is the client handler for project endpoints.
	ProjectHandler struct {
		client *Client
	}

	Project struct {
		ID             json.Number `json:"id"`
		Name           string      `json:"project_name"`
		Description    string      `json:"description"`
		TechnicalEmail string      `json:"technical_emails"`
		NetworkCIDR    string      `json:"networkCIDR"`
		CreationDate   string      `json:"creation_date"`
	}

	CreateProjectRequest struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		TechnicalEmail string `json:"technicalEmails"`
	}

	UpdateProjectRequest struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		TechnicalEmail string `json:"technicalEmails"`
	}
)

// Get is the method to get a project.
func (h *ProjectHandler) Get(projectID string) (*Project, error) {
	projects, err := h.GetList()
	if err != nil {
		return nil, err
	}

	for _, project := range *projects {
		if project.ID.String() == projectID {
			return &project, nil
		}
	}

	return nil, fmt.Errorf("project not found")
}

// GetList is the method to get a list of projects.
func (h *ProjectHandler) GetList() (*[]Project, error) {
	type projetListRequest struct {
		JWT string `json:"jwt"`
	}

	type projectListResponse struct {
		APIResponse
		ProjectList struct {
			Projects []Project `json:"projects"`
		} `json:"data"`
	}

	req := projetListRequest{
		JWT: h.client.jwt,
	}

	bts, err := h.client.sendPostRequest(
		fmt.Sprintf("%s/api/projects/getList", h.client.BaseURL),
		req,
	)
	if err != nil {
		return nil, err
	}

	var res projectListResponse
	if err = checkAPIResponse(bts, &res); err != nil {
		return nil, err
	}

	return &res.ProjectList.Projects, nil
}

// Create creates a new project.
func (h *ProjectHandler) Create(req CreateProjectRequest) (*Project, error) {
	type createProjectFullRequest struct {
		CreateProjectRequest
		JWT string `json:"jwt"`
	}

	type createProjectResponse struct {
		APIResponse
		Project Project `json:"data"`
	}

	fullReq := createProjectFullRequest{req, h.client.jwt}

	bts, err := h.client.sendPostRequest(
		fmt.Sprintf("%s/api/projects/addProject", h.client.BaseURL),
		fullReq,
	)
	if err != nil {
		return nil, err
	}

	var res createProjectResponse
	if err = checkAPIResponse(bts, &res); err != nil {
		return nil, err
	}

	return &res.Project, nil
}

// Update is the method to update a project.
func (h *ProjectHandler) Update(projectID string, req UpdateProjectRequest) (*Project, error) {
	type updateProjectFullRequest struct {
		UpdateProjectRequest
		ProjectID string `json:"projectId"`
		JWT       string `json:"jwt"`
	}

	type updateProjectResponse struct {
		APIResponse
		Project Project `json:"data"`
	}

	fullReq := updateProjectFullRequest{
		UpdateProjectRequest: req,
		ProjectID:            projectID,
		JWT:                  h.client.jwt,
	}

	bts, err := h.client.sendPutRequest(
		fmt.Sprintf("%s/api/projects/editProject", h.client.BaseURL),
		fullReq,
	)
	if err != nil {
		return nil, err
	}

	var res updateProjectResponse
	if err = checkAPIResponse(bts, &res); err != nil {
		return nil, err
	}

	return &res.Project, nil
}

// Delete is the method to delete a project.
func (h *ProjectHandler) Delete(projectID string) error {
	type deleteProjectFullRequest struct {
		ProjectID string `json:"projectId"`
		JWT       string `json:"jwt"`
	}

	type deleteProjectResponse struct {
		APIResponse
	}

	req := deleteProjectFullRequest{
		ProjectID: projectID,
		JWT:       h.client.jwt,
	}

	bts, err := h.client.sendDeleteRequest(
		fmt.Sprintf("%s/api/projects/deleteProject", h.client.BaseURL),
		req,
	)
	if err != nil {
		return err
	}

	var res deleteProjectResponse
	if err = checkAPIResponse(bts, &res); err != nil {
		return err
	}

	return nil
}
