package handler

import (
	"net/http"

	e "github.com/niko-cb/covid19datascraper/server/error"

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
	render.Status(r, http.StatusOK)
	ctx := context.New(r)
	dp := (&dialogflow.Processor{}).NewSession()
	if err := dp.CreateOrRecreateIntents(); err != nil {
		log.Errorf(ctx, err.Error())
		_ = render.Render(w, r, e.ErrRender(ctx, err))
	}
}
