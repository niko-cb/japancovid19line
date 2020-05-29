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
	DialogflowAPIBasePath           = "/dialogflow"
	dialogflowCreateDataAPIPath     = "/createData"
	dialogflowCreateSymptomsAPIPath = "/createSymptoms"
	dialogflowDeleteAPIPath         = "/delete"
)

func Dialogflow(r chi.Router) {
	r.Get(dialogflowCreateDataAPIPath, createDataIntents)
	r.Get(dialogflowCreateSymptomsAPIPath, createSymptomIntents)
	r.Delete(dialogflowDeleteAPIPath, deleteIntents)
}

func createDataIntents(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	dp := dialogflow.NewSession()
	if err := dp.CreateDataIntents(); err != nil {
		log.Errorf(ctx, err.Error())
		_ = render.Render(w, r, e.ErrRender(ctx, err))
	}
	render.Status(r, http.StatusOK)
}

func createSymptomIntents(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	dp := dialogflow.NewSession()
	if err := dp.CreateSymptomIntents(); err != nil {
		log.Errorf(ctx, err.Error())
		_ = render.Render(w, r, e.ErrRender(ctx, err))
	}
	render.Status(r, http.StatusOK)
}

func deleteIntents(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	dp := dialogflow.NewSession()
	if err := dp.DeleteIntents(); err != nil {
		log.Errorf(ctx, err.Error())
		_ = render.Render(w, r, e.ErrRender(ctx, err))
	}
	render.Status(r, http.StatusOK)
}
