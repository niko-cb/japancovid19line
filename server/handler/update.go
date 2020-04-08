package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/niko-cb/covid19datascraper/server/scrape"
	"github.com/niko-cb/covid19datascraper/server/utils"
	"net/http"
)

const (
	UpdateDataAPIBasePath   = "/update"
	updateDataAPIPathPrefix = APIPathPrefix + UpdateDataAPIBasePath
	updateDataAPIPath       = "/"
	updateADataPIFullPath   = updateDataAPIPathPrefix + updateDataAPIPath
)

func UpdateData(r chi.Router) {
	r.Get(updateDataAPIPath, updateData)
}

// updateData is a handler to update all at once (scrape and dialogflow)
func updateData(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(r)
	jpd := scrape.Scrape()
	dp, err := utils.NewDialogflowSession(
		"japancovid19",
		"./dialogflowcredentials.json",
		"ja",
		"Japan/Tokyo")
	if err != nil {
		panic(err)
	}

	err = dp.CreateOrRecreateIntents()
	if err != nil {
		panic(err)
	}
	if err := render.Render(w, r, newPrefectureDataRes(jpd)); err != nil {
		_ = render.Render(w, r, utils.ErrRender(ctx, err))
	}
	render.Status(r, http.StatusOK)
}
