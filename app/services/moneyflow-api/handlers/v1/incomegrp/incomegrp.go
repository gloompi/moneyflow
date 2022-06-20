// Package incomegrp maintains the group of handlers for income access.
package incomegrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gloompi/ultimate-service/business/core/income"
	"github.com/gloompi/ultimate-service/business/web/auth"
	v1Web "github.com/gloompi/ultimate-service/business/web/v1"
	"github.com/gloompi/ultimate-service/foundation/web"
)

// Handlers manages the set of income endpoints.
type Handlers struct {
	Income income.Core
	Auth   *auth.Auth
}

// Create adds a new income to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var ni income.NewIncome
	if err := web.Decode(r, &ni); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	inc, err := h.Income.Create(ctx, ni, v.Now)
	if err != nil {
		return fmt.Errorf("income[%+v]: %w", &inc, err)
	}

	return web.Respond(ctx, w, inc, http.StatusCreated)
}

// Update updates an income in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var upd income.UpdateIncome
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	id := web.Param(r, "id")

	inc, err := h.Income.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, income.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, income.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying income[%s]: %w", id, err)
		}
	}

	// If you are not an admin and looking to update a income you don't own.
	if !claims.AuthorizedByRole(auth.RoleAdmin) && claims.AuthorizedByUserId(inc.UserID) {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Income.Update(ctx, id, upd, v.Now); err != nil {
		switch {
		case errors.Is(err, income.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, income.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] Income[%+v]: %w", id, &upd, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes an income from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	id := web.Param(r, "id")

	inc, err := h.Income.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, income.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, income.ErrNotFound):
			// Don't send StatusNotFound here since the call to Delete
			// below won't if this income is not found. We only know
			// this because we are doing the Query for the UserID.
			return v1Web.NewRequestError(err, http.StatusNoContent)
		default:
			return fmt.Errorf("querying income[%s]: %w", id, err)
		}
	}

	// If you are not an admin and looking to delete an income you don't own.
	if !claims.AuthorizedByRole(auth.RoleAdmin) && !claims.AuthorizedByUserId(inc.UserID) {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Income.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, income.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Query returns a list of incomes with paging.
func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page := web.Param(r, "page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid page format, page[%s]", page), http.StatusBadRequest)
	}
	rows := web.Param(r, "rows")
	rowsPerPage, err := strconv.Atoi(rows)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid rows format, rows[%s]", rows), http.StatusBadRequest)
	}

	incomes, err := h.Income.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for incomes: %w", err)
	}

	return web.Respond(ctx, w, incomes, http.StatusOK)
}

// QueryByID returns a income by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	inc, err := h.Income.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, income.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, income.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, inc, http.StatusOK)
}

// QueryByUserID returns a list of incomes for a user.
func (h Handlers) QueryByUserID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	userID := web.Param(r, "user_id")

	// If you are not an admin and looking to delete an income you don't own.
	if !claims.AuthorizedByRole(auth.RoleAdmin) && !claims.AuthorizedByUserId(userID) {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	incomes, err := h.Income.QueryByUserID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, income.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, income.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("userID[%s]: %w", userID, err)
		}
	}

	return web.Respond(ctx, w, incomes, http.StatusOK)
}
