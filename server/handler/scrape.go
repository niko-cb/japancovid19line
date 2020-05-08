package handler

import (
	"errors"
	"net/http"

	"github.com/niko-cb/covid19datascraper/server/context"
	e "github.com/niko-cb/covid19datascraper/server/error"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/niko-cb/covid19datascraper/server/model"
	"github.com/niko-cb/covid19datascraper/server/scrape"
)

const (
	ScrapeDataAPIBasePath = "/scrape"
	scrapeDataAPIPath     = "/"
)

func Scrape(r chi.Router) {
	r.Get(scrapeDataAPIPath, scrapeData)
}

// scrapeData is a handler to retrieve all prefectures' covid19 data
func scrapeData(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)

	ctx := context.New(r)
	jpd := scrape.Do(ctx)
	if err := render.Render(w, r, newPrefectureDataRes(jpd)); err != nil {
		_ = render.Render(w, r, e.ErrRender(ctx, err))
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
		return errors.New("no prefecture data")
	}
	return nil
}
