package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/niko-cb/covid19datascraper/server/utils"
)

const (
	DialogflowAPIBasePath   = "/dialogflow"
	dialogflowCreateAPIPath = "/createIntents"
)

func Dialogflow(r chi.Router) {
	r.Get(dialogflowCreateAPIPath, createDialogflowIntents)
}

func createDialogflowIntents(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(r)
	dp := utils.NewDialogflowSession()
	if err := dp.CreateOrRecreateIntents(); err != nil {
		utils.LogErrorf(ctx, err.Error())
		render.Status(r, http.StatusInternalServerError)
	}
	render.Status(r, http.StatusOK)
}
