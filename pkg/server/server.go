package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/store"
)

// Server represents all server handlers.
type Server struct {
	cfg    *conf.Config
	stores *store.Stores
}

// NewServer new api server
func NewServer(cfg *conf.Config, stores *store.Stores) (*Server, error) {
	return &Server{cfg: cfg, stores: stores}, nil
}

// Projects Get a list of projects.// (GET /projects)
func (sv *Server) Projects(ctx echo.Context, params api.ProjectsParams) error {
	return nil
}

// NewProject Create a project.// (POST /projects)
func (sv *Server) NewProject(ctx echo.Context) error {

	newProj := new(api.NewProject)
	if err := ctx.Bind(newProj); err != nil {
		return err
	}

	resProj, err := sv.stores.Projects.Create(ctx.Request().Context(), newProj, "a4a777ff-fd47-42ab-84b4-1cca19a51f8f")
	if err != nil {
		if err == store.ErrProjectNameAlreadyExists {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusCreated, resProj)

	return nil
}

// GetProject (GET /projects/{id})
func (sv *Server) GetProject(ctx echo.Context, id string) error {
	return nil
}

// UpdateProject Update a project.// (PUT /projects/{id})
func (sv *Server) UpdateProject(ctx echo.Context, id string) error {
	return nil
}

// Issues Get a list of issues.// (GET /projects/{project_id}/issues)
func (sv *Server) Issues(ctx echo.Context, projectId string, params api.IssuesParams) error {
	return nil
}

// NewIssue Create a issue.// (POST /projects/{project_id}/issues)
func (sv *Server) NewIssue(ctx echo.Context, projectId string) error {
	return nil
}

// GetIssue (GET /projects/{project_id}/issues/{id})
func (sv *Server) GetIssue(ctx echo.Context, projectId string, id string) error {
	return nil
}

// Comments Get a list of Comments.// (GET /projects/{project_id}/issues/{issue_id}/comments)
func (sv *Server) Comments(ctx echo.Context, projectId string, issueId string, params api.CommentsParams) error {
	return nil
}

// NewComment Create a comment on a issue.// (POST /projects/{project_id}/issues/{issue_id}/comments)
func (sv *Server) NewComment(ctx echo.Context, projectId string, issueId string) error {
	return nil
}

// GetComment (GET /projects/{project_id}/issues/{issue_id}/comments/{id})
func (sv *Server) GetComment(ctx echo.Context, projectId string, issueId string, id string) error {
	return nil
}

// Users Get a list of users.// (GET /users)
func (sv *Server) Users(ctx echo.Context, params api.UsersParams) error {
	return nil
}

// GetUser (GET /users/{id})
func (sv *Server) GetUser(ctx echo.Context, id string) error {
	return nil
}
