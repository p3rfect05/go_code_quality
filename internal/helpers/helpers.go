package helpers

import (
	"fmt"
	"github.com/p3rfect05/go_proj/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	traceError := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(traceError)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}
