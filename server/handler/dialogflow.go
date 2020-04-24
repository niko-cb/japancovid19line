package handler

import (
	"net/http"

	"github.com/niko-cb/covid19datascraper/server/context"
	"github.com/niko-cb/covid19datascraper/server/dialogflow"
	"github.com/niko-cb/covid19datascraper/server/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const (
	DialogflowAPIBasePath   = "/dialogflow"
	dialogflowCreateAPIPath = "/createIntents"
)

func Dialogflow(r chi.Router) {
	r.Get(dialogflowCreateAPIPath, createDialogflowIntents)
}

func createDialogflowIntents(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	dp := dialogflow.NewSession()
	if err := dp.CreateOrRecreateIntents(); err != nil {
		log.Errorf(ctx, err.Error())
		render.Status(r, http.StatusInternalServerError)
	}
	render.Status(r, http.StatusOK)
}
