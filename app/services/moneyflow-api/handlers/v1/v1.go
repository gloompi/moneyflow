// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"net/http"

	"github.com/gloompi/ultimate-service/app/services/moneyflow-api/handlers/v1/expensegrp"
	"github.com/gloompi/ultimate-service/app/services/moneyflow-api/handlers/v1/incomegrp"
	"github.com/gloompi/ultimate-service/app/services/moneyflow-api/handlers/v1/usergrp"
	"github.com/gloompi/ultimate-service/business/core/expense"
	"github.com/gloompi/ultimate-service/business/core/income"
	"github.com/gloompi/ultimate-service/business/core/user"
	"github.com/gloompi/ultimate-service/business/web/auth"
	"github.com/gloompi/ultimate-service/business/web/v1/mid"
	"github.com/gloompi/ultimate-service/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *zap.SugaredLogger
	Auth *auth.Auth
	DB   *sqlx.DB
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.Auth)
	admin := mid.Authorize(auth.RoleAdmin)

	// Register user management and authentication endpoints.
	ugh := usergrp.Handlers{
		User: user.NewCore(cfg.Log, cfg.DB),
		Auth: cfg.Auth,
	}
	app.Handle(http.MethodGet, version, "/users/token", ugh.Token)
	app.Handle(http.MethodGet, version, "/users/:page/:rows", ugh.Query, authen, admin)
	app.Handle(http.MethodGet, version, "/users/:id", ugh.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/users", ugh.Create)
	app.Handle(http.MethodPut, version, "/users/:id", ugh.Update, authen)
	app.Handle(http.MethodDelete, version, "/users/:id", ugh.Delete, authen)

	// Register income management endpoints.
	igh := incomegrp.Handlers{
		Income: income.NewCore(cfg.Log, cfg.DB),
	}
	app.Handle(http.MethodGet, version, "/incomes/:page/:rows", igh.Query, authen, admin)
	app.Handle(http.MethodGet, version, "/incomes/:id", igh.QueryByID, authen, admin)
	app.Handle(http.MethodGet, version, "/incomes/user/:user_id", igh.QueryByUserID, authen)
	app.Handle(http.MethodPost, version, "/incomes", igh.Create, authen)
	app.Handle(http.MethodPut, version, "/incomes/:id", igh.Update, authen)
	app.Handle(http.MethodDelete, version, "/incomes/:id", igh.Delete, authen)

	// Register expense management endpoints.
	egh := expensegrp.Handlers{
		Expense: expense.NewCore(cfg.Log, cfg.DB),
	}
	app.Handle(http.MethodGet, version, "/expenses/:page/:rows", egh.Query, authen, admin)
	app.Handle(http.MethodGet, version, "/expenses/:id", egh.QueryByID, authen, admin)
	app.Handle(http.MethodPost, version, "/expenses", egh.Create, authen)
	app.Handle(http.MethodPut, version, "/expenses/:id", egh.Update, authen)
	app.Handle(http.MethodDelete, version, "/expenses/:id", egh.Delete, authen)
}
