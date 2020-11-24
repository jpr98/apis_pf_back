package controllers

import (
	"net/http"
	"strconv"
	"strings"

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

type projectSearch struct {
	SearchType string `json:"search_type"`
	Keywords   string `json:"keywords"`
}

// SearchProject handles looking for a project with a given title
func (p *Projects) SearchProject(c echo.Context) error {
	ps := new(projectSearch)
	if err := c.Bind(&ps); err != nil {
		return c.String(http.StatusBadRequest, "Can't bind request body")
	}

	var projects []models.Project
	var err error

	switch ps.SearchType {
	case "title":
		projects, err = p.projectStore.GetByTitle(ps.Keywords)

	case "category":
		projects, err = p.projectStore.GetByCategory(ps.Keywords)

	case "tags":
		tags := strings.Fields(ps.Keywords)
		projects, err = p.projectStore.GetByTags(tags)

	default:
		return c.String(http.StatusBadRequest, "Please provide a valid search type (title, tags, category)")
	}

	if err != nil {
		return c.String(http.StatusNotFound, "No projects matching search")
	}

	return c.JSON(http.StatusFound, projects)
}

// GetByOwner returns all the projects a user owns
func (p *Projects) GetByOwner(c echo.Context) error {
	id := c.Param("userId")
	if id == "" {
		return c.String(http.StatusBadRequest, "No user id")
	}

	projects, err := p.projectStore.GetByOwnerID(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, projects)
}

// GetVotedFor returns the projects voted for by a user
func (p *Projects) GetVotedFor(c echo.Context) error {
	id := c.Param("userId")
	if id == "" {
		return c.String(http.StatusBadRequest, "No user id")
	}

	projects, err := p.projectStore.GetVotedProjects(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, projects)
}

// VoteForProject handles a user voting or unvoting a project
func (p *Projects) VoteForProject(c echo.Context) error {
	id := c.Param("id")
	upvoteStr := c.QueryParam("upvote")
	upvote, err := strconv.ParseBool(upvoteStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid query parameter")
	}

	userID := getTokenStringClaimByKey(c, "id")
	if err := p.projectStore.Vote(id, userID, upvote); err != nil {
		return c.String(http.StatusBadRequest, "Couldn't upvote project")
	}
	return c.JSON(http.StatusOK, "")
}

// Delete deletes a project permanently
func (p *Projects) Delete(c echo.Context) error {
	id := c.Param("id")
	userID := getTokenStringClaimByKey(c, "id")

	project, err := p.projectStore.GetByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if project.Owner.Hex() != userID {
		return c.String(http.StatusUnauthorized, "You must be the project owner to delete it")
	}

	if err := p.projectStore.Delete(id); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "Project deleted")
}

// View increments a projects views
func (p *Projects) View(c echo.Context) error {
	id := c.Param("id")

	if err := p.projectStore.View(id); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusAccepted, "Project views updated")
}

type commentRequest struct {
	Text string `json:"text"`
}

// Comment handles adding a comment to a project
func (p *Projects) Comment(c echo.Context) error {
	id := c.Param("id")
	user := getTokenStringClaimByKey(c, "id")

	cr := new(commentRequest)
	if err := c.Bind(cr); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := p.projectStore.AddComment(id, user, cr.Text); err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusAccepted, "Comment successfully created")
}

type contributionRequest struct {
	Amount float32 `json:"amount"`
}

// Contribute handles adding a contribution to a project
func (p *Projects) Contribute(c echo.Context) error {
	id := c.Param("id")
	user := getTokenStringClaimByKey(c, "id")

	cr := new(contributionRequest)
	if err := c.Bind(cr); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := p.projectStore.AddContribution(id, user, cr.Amount); err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusAccepted, "Contribution successfully added")
}
