package handler

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/niko-cb/covid19datascraper/server/model"
	"github.com/niko-cb/covid19datascraper/server/utils"
	"net/http"

	"github.com/niko-cb/covid19datascraper/server/scrape"
)

const (
	ScrapeDataAPIBasePath   = "/scrape"
	scrapeDataAPIPathPrefix = APIPathPrefix + ScrapeDataAPIBasePath
	scrapeDataAPIPath       = "/"
	scrapeDataAPIFullPath   = scrapeDataAPIPathPrefix + scrapeDataAPIPath
)

func ScrapeData(r chi.Router) {
	r.Get(scrapeDataAPIPath, scrapeData)
}

// scrapeData is a handler to retrieve all prefectures' covid19 data
func scrapeData(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(r)
	jpd := scrape.Scrape()

	render.Status(r, http.StatusOK)
	if err := render.Render(w, r, newPrefectureDataRes(jpd)); err != nil {
		_ = render.Render(w, r, utils.ErrRender(ctx, err))
	}
}

func newPrefectureDataRes(jpd []*model.PrefectureData) *prefectureDataRes {
	return &prefectureDataRes{PrefectureDataSlice: jpd}
}

// request

type prefectureDataRes struct {
	model.PrefectureDataSlice
}

// render

func (res *prefectureDataRes) Render(w http.ResponseWriter, r *http.Request) error {
	if res.PrefectureDataSlice == nil {
		return errors.New("missing required Prefecture fields")
	}
	return nil
}
