package handler

import (
	"net/http"

	"github.com/niko-cb/covid19datascraper/app/controller/errors"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/niko-cb/covid19datascraper/app/config/context"
	"github.com/niko-cb/covid19datascraper/app/controller/logic/scrape"
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
	ctx := context.New(r)
	if err := scrape.Do(ctx); err != nil {
		errors.Errorf(ctx, err.Error())
		_ = render.Render(w, r, errors.ErrRender(ctx, err))
		return
	}
	render.Status(r, http.StatusOK)
}
