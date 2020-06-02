package handler

import (
	"net/http"

	"github.com/niko-cb/covid19datascraper/app/models"

	"github.com/niko-cb/covid19datascraper/app/config/context"
	"github.com/niko-cb/covid19datascraper/app/controller/errors"

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
	dp := models.NewSession()
	if err := dp.CreateDataIntents(); err != nil {
		errors.Errorf(ctx, err.Error())
		_ = render.Render(w, r, errors.ErrRender(ctx, err))
	}
	render.Status(r, http.StatusOK)
}

func createSymptomIntents(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	dp := models.NewSession()
	if err := dp.CreateSymptomIntents(); err != nil {
		errors.Errorf(ctx, err.Error())
		_ = render.Render(w, r, errors.ErrRender(ctx, err))
	}
	render.Status(r, http.StatusOK)
}

func deleteIntents(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	dp := models.NewSession()
	if err := dp.DeleteIntents(); err != nil {
		errors.Errorf(ctx, err.Error())
		_ = render.Render(w, r, errors.ErrRender(ctx, err))
	}
	render.Status(r, http.StatusOK)
}
