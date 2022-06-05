// Package testgrp provides group of test handlers.
package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/gloompi/ultimate-service/business/sys/validate"
	"github.com/gloompi/ultimate-service/foundation/web"
	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	}

	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
