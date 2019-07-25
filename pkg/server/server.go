package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/store"
)

// DefaultCustomerID TODO NOT THIS
const DefaultCustomerID = "a4a777ff-fd47-42ab-84b4-1cca19a51f8f"

// Server represents all server handlers.
type Server struct {
	cfg    *conf.Config
	stores *store.Stores
}

// NewServer new api server
func NewServer(cfg *conf.Config, stores *store.Stores) (*Server, error) {
	return &Server{cfg: cfg, stores: stores}, nil
}

// Customers Get a list of customers. (GET /customers)
func (sv *Server) Customers(ctx echo.Context, params api.CustomersParams) error {
	query, limit, offset := listArgs(params.Q, params.Limit, params.Offset)
	log.Info().Str("query", query).Int("offset", offset).Int("limit", limit).Msg("ProjectsListOptions")

	opt := store.NewCustomersListOptions(query, offset, limit)

	resCusts, err := sv.stores.Customers.List(ctx.Request().Context(), opt)
	if err != nil {
		log.Error().Err(err).Msg("Projects failed")
		return err
	}

	return ctx.JSON(http.StatusOK, &api.CustomersPage{Customers: resCusts})
}

// NewCustomer Create a customer. (POST /customers)
func (sv *Server) NewCustomer(ctx echo.Context) error {
	newCust := new(api.NewCustomer)
	if err := ctx.Bind(newCust); err != nil {
		return err
	}

	resCust, err := sv.stores.Customers.Create(ctx.Request().Context(), newCust)
	if err != nil {
		if err == store.ErrCustomerNameAlreadyExists {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusCreated, resCust)
}

// GetCustomer (GET /customers/{id})
func (sv *Server) GetCustomer(ctx echo.Context, id string) error {

	resCust, err := sv.stores.Customers.GetByID(ctx.Request().Context(), id)
	if err != nil {
		if _, ok := err.(*store.CustomerNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Error().Err(err).Msg("get customer failed")
		return err
	}

	return ctx.JSON(http.StatusOK, resCust)
}

// UpdateCustomer Update a customer. (PUT /customers/{id})
func (sv *Server) UpdateCustomer(ctx echo.Context, id string) error {
	upCust := new(api.UpdatedCustomer)
	if err := ctx.Bind(upCust); err != nil {
		return err
	}

	resCust, err := sv.stores.Customers.Update(ctx.Request().Context(), upCust, id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, resCust)
}

// Projects Get a list of projects. (GET /projects)
func (sv *Server) Projects(ctx echo.Context, params api.ProjectsParams) error {

	query, limit, offset := listArgs(params.Q, params.Limit, params.Offset)
	log.Info().Str("query", query).Int("offset", offset).Int("limit", limit).Msg("ProjectsListOptions")

	opt := store.NewProjectsListOptions(query, offset, limit)

	resProjs, err := sv.stores.Projects.List(ctx.Request().Context(), opt, DefaultCustomerID)
	if err != nil {
		log.Error().Err(err).Msg("Projects failed")
		return err
	}

	return ctx.JSON(http.StatusOK, &api.ProjectsPage{Projects: resProjs})
}

// NewProject Create a project. (POST /projects)
func (sv *Server) NewProject(ctx echo.Context) error {

	newProj := new(api.NewProject)
	if err := ctx.Bind(newProj); err != nil {
		return err
	}

	resProj, err := sv.stores.Projects.Create(ctx.Request().Context(), newProj, DefaultCustomerID)
	if err != nil {
		if err == store.ErrProjectNameAlreadyExists {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		log.Error().Err(err).Msg("new project failed")
		return err
	}

	return ctx.JSON(http.StatusCreated, resProj)
}

// GetProject (GET /projects/{id})
func (sv *Server) GetProject(ctx echo.Context, id string) error {

	resProj, err := sv.stores.Projects.GetByID(ctx.Request().Context(), id, DefaultCustomerID)
	if err != nil {
		if _, ok := err.(*store.ProjectNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Error().Err(err).Msg("get project failed")
		return err
	}

	return ctx.JSON(http.StatusOK, resProj)
}

// UpdateProject Update a project. (PUT /projects/{id})
func (sv *Server) UpdateProject(ctx echo.Context, id string) error {
	upProj := new(api.UpdatedProject)
	if err := ctx.Bind(upProj); err != nil {
		return err
	}

	resProj, err := sv.stores.Projects.Update(ctx.Request().Context(), upProj, id, DefaultCustomerID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, resProj)
}

// Issues Get a list of issues. (GET /projects/{project_id}/issues)
func (sv *Server) Issues(ctx echo.Context, projectId string, params api.IssuesParams) error {
	return ctx.JSON(http.StatusNotImplemented, "not implemented yet")
}

// NewIssue Create a issue. (POST /projects/{project_id}/issues)
func (sv *Server) NewIssue(ctx echo.Context, projectId string) error {
	return ctx.JSON(http.StatusNotImplemented, "not implemented yet")
}

// GetIssue (GET /projects/{project_id}/issues/{id})
func (sv *Server) GetIssue(ctx echo.Context, projectId string, id string) error {
	return ctx.JSON(http.StatusNotImplemented, "not implemented yet")
}

// Comments Get a list of Comments. (GET /projects/{project_id}/issues/{issue_id}/comments)
func (sv *Server) Comments(ctx echo.Context, projectId string, issueId string, params api.CommentsParams) error {
	return ctx.JSON(http.StatusNotImplemented, "not implemented yet")
}

// NewComment Create a comment on a issue. (POST /projects/{project_id}/issues/{issue_id}/comments)
func (sv *Server) NewComment(ctx echo.Context, projectId string, issueId string) error {
	return ctx.JSON(http.StatusNotImplemented, "not implemented yet")
}

// GetComment (GET /projects/{project_id}/issues/{issue_id}/comments/{id})
func (sv *Server) GetComment(ctx echo.Context, projectId string, issueId string, id string) error {
	return ctx.JSON(http.StatusNotImplemented, "not implemented yet")
}

// Users Get a list of users. (GET /users)
func (sv *Server) Users(ctx echo.Context, params api.UsersParams) error {
	return ctx.JSON(http.StatusNotImplemented, "not implemented yet")
}

// GetUser (GET /users/{id})
func (sv *Server) GetUser(ctx echo.Context, id string) error {
	return ctx.JSON(http.StatusNotImplemented, "not implemented yet")
}
