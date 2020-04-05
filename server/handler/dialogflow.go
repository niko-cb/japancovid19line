package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/niko-cb/covid19datascraper/server/utils"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	DialogflowAPIBasePath   = "/dialogflow"
	DialogflowAPIPathPrefix = APIPathPrefix + DialogflowAPIBasePath
	dialogflowAPIPath       = "/"
	dialogflowAPIFullPath   = DialogflowAPIPathPrefix + dialogflowAPIPath

	dialogflowCreateAPIPath     = "/createIntents"
	dialogflowCreateAPIFullPath = DialogflowAPIPathPrefix + dialogflowCreateAPIPath
)

func Dialogflow(r chi.Router) {
	r.Post(dialogflowAPIPath, dialogflowTest)
	r.Get(dialogflowCreateAPIPath, createDialogflowIntents)
}

// TODO
func dialogflowTest(w http.ResponseWriter, r *http.Request) {
	dp, err := utils.NewDialogflowSession(
		"japancovid19",
		"./dialogflowcredentials.json",
		"ja",
		"Japan/Tokyo")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body content", http.StatusInternalServerError)
	}
	type inboundMessage struct {
		Message string
	}

	var m inboundMessage
	err = json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}

	// Use NLP
	response := dp.ProcessNLP(m.Message, "test")
	fmt.Printf("%#v", response)

	render.Status(r, http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createDialogflowIntents(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err.Error())
	}
	render.Status(r, http.StatusOK)

}
