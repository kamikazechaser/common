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

	ResponseProvider struct {
		logProvider Logger
	}
)

func NewResponseProvider(logProvder Logger) *ResponseProvider {
	return &ResponseProvider{logProvider: logProvder}
}

func (rp *ResponseProvider) Reply(w http.ResponseWriter, req *http.Request, err error) {
	if err == nil {
		return
	}
	logProv := rp.logProvider.Init(req.Context())

	httpErr, ok := err.(*httpError)
	if !ok {
		logProv.Info("internal server errror", "err", err)
		httputil.JSON(w, http.StatusInternalServerError, httpError{
			OK:          false,
			Description: "Internal server error",
		})
		return
	}

	httputil.JSON(w, httpErr.StatusCode, httpErr)
}
