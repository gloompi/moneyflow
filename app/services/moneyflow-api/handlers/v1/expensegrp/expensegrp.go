// Package expensegrp maintains the group of handlers for expense access.
package expensegrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gloompi/ultimate-service/business/core/expense"
	"github.com/gloompi/ultimate-service/business/web/auth"
	v1Web "github.com/gloompi/ultimate-service/business/web/v1"
	"github.com/gloompi/ultimate-service/foundation/web"
)

// Handlers manages the set of expense endpoints.
type Handlers struct {
	Expense expense.Core
	Auth    *auth.Auth
}

// Create adds a new expense to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var ne expense.NewExpense
	if err := web.Decode(r, &ne); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	exp, err := h.Expense.Create(ctx, ne, v.Now)
	if err != nil {
		return fmt.Errorf("expense[%+v]: %w", &exp, err)
	}

	return web.Respond(ctx, w, exp, http.StatusCreated)
}

// Update updates an expense in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var upd expense.UpdateExpense
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	id := web.Param(r, "id")

	exp, err := h.Expense.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, expense.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, expense.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying expense[%s]: %w", id, err)
		}
	}

	if !claims.AuthorizedByRole(auth.RoleAdmin) && exp.UserID != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Expense.Update(ctx, id, upd, v.Now); err != nil {
		switch {
		case errors.Is(err, expense.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, expense.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] Expense[%+v]: %w", id, &upd, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes an expense from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	id := web.Param(r, "id")

	exp, err := h.Expense.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, expense.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, expense.ErrNotFound):
			// Don't send StatusNotFound here since the call to Delete
			// below won't if this expense is not found. We only know
			// this because we are doing the Query for the UserID.
			return v1Web.NewRequestError(err, http.StatusNoContent)
		default:
			return fmt.Errorf("querying expense[%s]: %w", id, err)
		}
	}

	// If you are not an admin and looking to delete an expense you don't own.
	if !claims.AuthorizedByRole(auth.RoleAdmin) && exp.UserID != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Expense.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, expense.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Query returns a list of expenses with paging.
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

	expenses, err := h.Expense.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for expenses: %w", err)
	}

	return web.Respond(ctx, w, expenses, http.StatusOK)
}

// QueryByID returns a expense by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	exp, err := h.Expense.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, expense.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, expense.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, exp, http.StatusOK)
}
