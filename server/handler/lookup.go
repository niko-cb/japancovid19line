package handler

import (
	"cloud.google.com/go/datastore"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/niko-cb/covid19datascraper/server/model"
	"github.com/niko-cb/covid19datascraper/server/utils"
	"log"
	"net/http"
)

const (
	LookupAPIBasePath       = "/lookup"
	lookupAPIPathPrefix     = APIPathPrefix + LookupAPIBasePath
	lookupPrefectureAPIPath = "/prefecture"
	lookupAPIFullPath       = lookupAPIPathPrefix + lookupPrefectureAPIPath
)

func LookupPrefecture(r chi.Router) {
	r.Get(lookupPrefectureAPIPath, lookupPrefecture)
}

func lookupPrefecture(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(r)
	dsClient, err := utils.NewDSClient()
	if err != nil {
		return
	}
	prefectureName := r.URL.Query().Get("p")

	//pm := PrefectureMap()
	//log.Println(pm[prefectureName])

	kind := utils.DatastoreKind()

	var pData []*model.PrefectureData
	q := datastore.NewQuery(kind).Filter("Prefecture =", prefectureName).Limit(1)
	if _, err := dsClient.GetAll(ctx, q, &pData); err != nil {
		log.Printf("errored on get: %v", err)
	}

	render.Status(r, http.StatusOK)
	if err := render.Render(w, r, newSinglePrefectureDataRes(pData)); err != nil {
		_ = render.Render(w, r, utils.ErrRender(ctx, err))
	}
}

func newSinglePrefectureDataRes(jpd []*model.PrefectureData) *singlePrefectureData {
	return &singlePrefectureData{PrefectureDataSlice: jpd}
}

// request

type singlePrefectureData struct {
	model.PrefectureDataSlice
}

// render

func (res *singlePrefectureData) Render(w http.ResponseWriter, r *http.Request) error {
	if res.PrefectureDataSlice == nil {
		return errors.New("missing required Prefecture fields")
	}
	return nil
}
