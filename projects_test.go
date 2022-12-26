package elestio

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupProjectTestCase(t *testing.T) *Client {
	t.Log("Setup service test case")

	email, apiKey := os.Getenv("ELESTIO_INTEGRATION_EMAIL"), os.Getenv("ELESTIO_INTEGRATION_API_KEY")

	c, err := NewClient(email, apiKey)
	if err != nil {
		t.Fatalf("client authentication error: %s", err)
	}

	return c
}

func TestProjectHandler_Get(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projectID := "1851"

	project, err := c.Project.Get(projectID)
	require.NoError(t, err, "expected no error when getting project")
	require.NotNil(t, project, "expected non-nil project")
	// require.Equal(t, projectID, project.ID, "expected project ID to be 1")
}

func TestProjectHandler_GetList(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	projects, err := c.Project.GetList()
	require.NoError(t, err, "expected no error when getting projects")
	require.NotNil(t, projects, "expected non-nil projects")
	// require.Equal(t, 1, len(projects.Projects), "expected 1 project")
	// require.Equal(t, 1, projects.Projects[0].ID, "expected project ID to be 1")
}

func TestProjectHandler_Create(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	req := CreateProjectRequest{
		Name:            "test-project",
		Description:     "test project",
		TechnicalEmails: "adamkrim.dev@gmail.com",
	}

	project, err := c.Project.Create(req)
	require.NoError(t, err, "expected no error when creating project")
	require.NotNil(t, project, "expected non-nil project")
	require.Equal(t, req.Name, project.Name, "expected project name to be test-project")
}

func TestProjectHandler_Update(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	args := UpdateProjectRequest{
		Name:            "test-project-updated",
		Description:     "test project updated",
		TechnicalEmails: "adamkrim.dev+updated@gmail.com",
	}

	project, err := c.Project.Update("2003", args)
	require.NoError(t, err, "expected no error when updating project")
	require.NotNil(t, project, "expected non-nil project")
	require.Equal(t, args.Name, project.Name, "expected project name to be test-project-updated")
}

func TestProjectHandler_Delete(t *testing.T) {
	t.Skip("Skipping test")
	c := setupServiceTestCase(t)

	err := c.Project.Delete("2003")
	require.NoError(t, err, "expected no error when deleting project")

	project, err := c.Project.Get("2003")
	require.Error(t, err, "expected error when getting project")
	require.Nil(t, project, "expected nil project")
}
