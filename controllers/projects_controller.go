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

	createdProject, err := p.projectStore.Create(*project)
	if err != nil {
		c.Logger().Errorf("Can't create project", err)
		return c.String(http.StatusInternalServerError, "Can't create project")
	}

	return c.JSON(http.StatusCreated, createdProject)
}
