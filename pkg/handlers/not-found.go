package handlers

import (
	"net/http"
	"github.com/bfaulk96/scoreboard-api/pkg/models"
	"github.com/bfaulk96/scoreboard-api/pkg/models/responses"
)

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
	aw.Respond(ar, &responses.Message{Message: "Unsupported URL"}, http.StatusNotFound)
}