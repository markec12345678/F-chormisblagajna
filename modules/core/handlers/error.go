package handlers

import (
	"fmt"
	"net/http"

	"github.com/nutrixpos/pos/common/logger"
)

// respondError logs the full error server-side and sends a safe generic message to the client.
func respondError(w http.ResponseWriter, log logger.ILogger, err error, statusCode int) {
	log.Error(err.Error())
	http.Error(w, http.StatusText(statusCode), statusCode)
}

// respondErrorMsg logs a contextual message and sends a safe generic message to the client.
func respondErrorMsg(w http.ResponseWriter, log logger.ILogger, context string, err error, statusCode int) {
	log.Error(fmt.Sprintf("%s: %s", context, err.Error()))
	http.Error(w, http.StatusText(statusCode), statusCode)
}
