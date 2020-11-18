package controllers

import (
	"net/http"

	"github.com/jpr98/apis_pf_back/models"
	"github.com/labstack/echo/v4"
)

// Projects represents a projects controller
type Projects struct {
	projectStore models.ProjectStore
}

// NewProjectsController creates a new projects controlelr with a store
func NewProjectsController(ps models.ProjectStore) Projects {
	return Projects{projectStore: ps}
}

// Create handles creating a new project
func (p *Projects) Create(c echo.Context) error {
	project := new(models.Project)
	if err := c.Bind(project); err != nil {
		c.Logger().Error("Can't bind body to JSON")
		return c.String(http.StatusBadRequest, "Can't bind body to json")
	}

	userID := getTokenStringClaimByKey(c, "id")
	createdProject, err := p.projectStore.Create(*project, userID)
	if err != nil {
		c.Logger().Errorf("Can't create project", err)
		return c.String(http.StatusInternalServerError, "Can't create project")
	}

	return c.JSON(http.StatusCreated, createdProject)
}

// GetByID handles looking for a project with a given id
func (p *Projects) GetByID(c echo.Context) error {
	id := c.Param("id")
	project, err := p.projectStore.GetByID(id)
	if err != nil {
		c.Logger().Errorf("Can't find project", err)
		return c.String(http.StatusNotFound, "Can't find project")
	}
	return c.JSON(http.StatusFound, project)
}

type titleSearch struct {
	Title string `json:"title"`
}

// SearchByTitle handles looking for a project with a given title
func (p *Projects) SearchByTitle(c echo.Context) error {
	ts := new(titleSearch)
	if err := c.Bind(&ts); err != nil {
		return c.String(http.StatusBadRequest, "Can't bind request body")
	}

	projects, err := p.projectStore.GetByTitle(ts.Title)
	if err != nil {
		return c.String(http.StatusNotFound, "No projects matching search")
	}

	return c.JSON(http.StatusFound, projects)
}
