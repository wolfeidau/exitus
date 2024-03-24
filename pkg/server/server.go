package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/auth"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/store"
)

const (
	// DefaultCustomerID TODO NOT THIS.
	DefaultCustomerID = "a4a777ff-fd47-42ab-84b4-1cca19a51f8f"

	// DefaultReporter TODO NOT THIS.
	DefaultReporter = "34a20135-1c9b-4c4d-b590-7771207ed847"

	// DefaultAuthor TODO NOT THIS.
	DefaultAuthor = "3fbecb27-1f23-4ed0-91e4-68f97a1f0364"
)

// Server represents all server handlers.
type Server struct {
	cfg    *conf.Config
	stores *store.Stores
}

// NewServer new api server.
func NewServer(cfg *conf.Config, stores *store.Stores) (*Server, error) {
	return &Server{cfg: cfg, stores: stores}, nil
}

// Customers Get a list of customers. (GET /customers).
func (sv *Server) Customers(ctx echo.Context, params api.CustomersParams) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	query, limit, offset := listArgs(params.Q, params.Limit, params.Offset)
	log.Info().Str("query", query).Int("offset", offset).Int("limit", limit).Msg("ProjectsListOptions")

	opt := store.NewCustomersListOptions(query, offset, limit)

	resCusts, err := sv.stores.Customers.List(ctx.Request().Context(), opt)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &api.CustomersPage{Customers: resCusts})
}

// NewCustomer Create a customer. (POST /customers).
func (sv *Server) NewCustomer(ctx echo.Context) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

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

// GetCustomer (GET /customers/{id}).
func (sv *Server) GetCustomer(ctx echo.Context, id string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	resCust, err := sv.stores.Customers.GetByID(ctx.Request().Context(), id)
	if err != nil {
		if _, ok := err.(*store.CustomerNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusOK, resCust)
}

// UpdateCustomer Update a customer. (PUT /customers/{id}).
func (sv *Server) UpdateCustomer(ctx echo.Context, id string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

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

// Projects Get a list of projects. (GET /projects).
func (sv *Server) Projects(ctx echo.Context, params api.ProjectsParams) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	query, limit, offset := listArgs(params.Q, params.Limit, params.Offset)
	log.Info().Str("query", query).Int("offset", offset).Int("limit", limit).Msg("ProjectsListOptions")

	opt := store.NewProjectsListOptions(query, offset, limit)

	resProjs, err := sv.stores.Projects.List(ctx.Request().Context(), opt, DefaultCustomerID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &api.ProjectsPage{Projects: resProjs})
}

// NewProject Create a project. (POST /projects).
func (sv *Server) NewProject(ctx echo.Context) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	newProj := new(api.NewProject)
	if err := ctx.Bind(newProj); err != nil {
		return err
	}

	resProj, err := sv.stores.Projects.Create(ctx.Request().Context(), newProj, DefaultCustomerID)
	if err != nil {
		if err == store.ErrProjectNameAlreadyExists {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusCreated, resProj)
}

// GetProject (GET /projects/{id}).
func (sv *Server) GetProject(ctx echo.Context, id string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	resProj, err := sv.stores.Projects.GetByID(ctx.Request().Context(), id, DefaultCustomerID)
	if err != nil {
		if _, ok := err.(*store.ProjectNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusOK, resProj)
}

// UpdateProject Update a project. (PUT /projects/{id}).
func (sv *Server) UpdateProject(ctx echo.Context, id string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

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

// Issues Get a list of issues. (GET /projects/{project_id}/issues).
func (sv *Server) Issues(ctx echo.Context, projectId string, params api.IssuesParams) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	query, limit, offset := listArgs(params.Q, params.Limit, params.Offset)
	log.Info().Str("query", query).Int("offset", offset).Int("limit", limit).Msg("IssuesListOptions")

	opt := store.NewIssueListOptions(query, offset, limit)

	resIssues, err := sv.stores.Issues.List(ctx.Request().Context(), opt, projectId, DefaultCustomerID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &api.IssuesPage{Issues: resIssues})
}

// NewIssue Create a issue. (POST /projects/{project_id}/issues).
func (sv *Server) NewIssue(ctx echo.Context, projectId string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	newIssue := new(api.NewIssue)
	if err := ctx.Bind(newIssue); err != nil {
		return err
	}

	resIssue, err := sv.stores.Issues.Create(ctx.Request().Context(), newIssue, projectId, DefaultCustomerID, DefaultReporter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, resIssue)
}

// UpdateIssue (PUT /projects/{project_id}/issues/{id}).
func (sv *Server) UpdateIssue(ctx echo.Context, projectId string, id string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	upIssue := new(api.UpdatedIssue)
	if err := ctx.Bind(upIssue); err != nil {
		return err
	}

	resIssue, err := sv.stores.Issues.Update(ctx.Request().Context(), upIssue, id, projectId, DefaultCustomerID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, resIssue)
}

// GetIssue (GET /projects/{project_id}/issues/{id}).
func (sv *Server) GetIssue(ctx echo.Context, projectId string, id string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	resIssue, err := sv.stores.Issues.GetByID(ctx.Request().Context(), id, projectId, DefaultCustomerID)
	if err != nil {
		if _, ok := err.(*store.IssueNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusOK, resIssue)
}

// Comments Get a list of Comments. (GET /projects/{project_id}/issues/{issue_id}/comments).
func (sv *Server) Comments(ctx echo.Context, projectId string, issueId string, params api.CommentsParams) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	query, limit, offset := listArgs(params.Q, params.Limit, params.Offset)
	log.Info().Str("query", query).Int("offset", offset).Int("limit", limit).Msg("CommentsListOptions")

	opt := store.NewCommentListOptions(query, offset, limit)

	resComments, err := sv.stores.Comments.List(ctx.Request().Context(), opt, issueId, projectId, DefaultCustomerID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &api.CommentsPage{Comments: resComments})
}

// NewComment Create a comment on a issue. (POST /projects/{project_id}/issues/{issue_id}/comments).
func (sv *Server) NewComment(ctx echo.Context, projectId string, issueId string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	newComment := new(api.NewComment)
	if err := ctx.Bind(newComment); err != nil {
		return err
	}

	resComment, err := sv.stores.Comments.Create(ctx.Request().Context(), newComment, issueId, projectId, DefaultCustomerID, DefaultAuthor)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, resComment)
}

// UpdateComment (PUT /projects/{project_id}/issues/{issue_id}/comments/{id}).
func (sv *Server) UpdateComment(ctx echo.Context, projectId string, issueId string, id string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	upComment := new(api.UpdatedComment)
	if err := ctx.Bind(upComment); err != nil {
		return err
	}

	resComment, err := sv.stores.Comments.Update(ctx.Request().Context(), upComment, id, issueId, projectId, DefaultCustomerID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, resComment)
}

// GetComment (GET /projects/{project_id}/issues/{issue_id}/comments/{id}).
func (sv *Server) GetComment(ctx echo.Context, projectId string, issueId string, id string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	resComment, err := sv.stores.Comments.GetByID(ctx.Request().Context(), id, issueId, projectId, DefaultCustomerID)
	if err != nil {
		if _, ok := err.(*store.CommentNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusOK, resComment)
}

// Users Get a list of users. (GET /users).
func (sv *Server) Users(ctx echo.Context, params api.UsersParams) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	return ctx.JSON(http.StatusNotImplemented, "not implemented yet")
}

// GetUser (GET /users/{id}).
func (sv *Server) GetUser(ctx echo.Context, id string) error {
	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	return ctx.JSON(http.StatusNotImplemented, "not implemented yet")
}

func userHasAccess(ctx echo.Context) bool {
	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to load user from context")
		return false
	}

	scopes, err := auth.LoadOperationScopesFromContext(ctx, auth.OpenIDScopes)
	if err != nil {
		log.Error().Err(err).Msg("failed to load scopes from context")
		return false
	}

	log.Info().Strs("Scopes", scopes).Object("User", &user).Msg("Scopes check")

	return user.HasScope(scopes)
}
