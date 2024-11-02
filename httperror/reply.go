package httperror

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/kamikazechaser/common/httputil"
)

type (
	Logger interface {
		Init(context.Context) *slog.Logger
	}

	ReplyProvider struct {
		logProvider Logger
	}
)

func NewReplyProvider(logProvder Logger) *ReplyProvider {
	return &ReplyProvider{logProvider: logProvder}
}

func (rp *ReplyProvider) ReplyError(w http.ResponseWriter, req *http.Request, err error) {
	if err == nil {
		return
	}
	logProv := rp.logProvider.Init(req.Context())

	// Support common errors here e.g.
	// Timeouts, context cancellations, JSON parsing errors, validations etc.

	httpErr, ok := err.(*httpError)
	if !ok {
		logProv.Error("internal server errror", "err", err)
		httputil.JSON(w, http.StatusInternalServerError, httpError{
			OK:          false,
			Description: "Internal server error",
		})
		return
	}

	httputil.JSON(w, httpErr.StatusCode, httpErr)
}
