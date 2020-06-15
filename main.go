package main

import (
	"fmt"
	. "github.com/dvincenz/postman-mockserver/common"
	. "github.com/dvincenz/postman-mockserver/postman"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//todo may not use global variable
var mocks map[string]Mock


func main() {
	//mocks := postman.ReadPostmanFile ("json/ixCloudPortal.postman_collection.json")
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	log.Info().Msg("Postman mock servers started")
	var err error
	mocks, err = GetMocksFromPostman()
	if err != nil{
		log.Error().Msg("error in get postman collation " + err.Error())
		return

	}
	log.Info().Msg("total " + strconv.Itoa(len(mocks)) + " mocks found")
	http.HandleFunc("/updatecollation", reloadCollationHandler)
	createMockServer()

}


func createMockServer(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		postmanRouter(w, r)
	})
	http.ListenAndServe("localhost:8080", nil)
}

func reloadCollationHandler(w http.ResponseWriter, r *http.Request){
	if HttpMethod(r.Method) == POST {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil || bodyBytes == nil{
			mocks, err = GetMocksFromPostman()
			log.Debug().Msg("Reload mocks from Postman")
		} else {
			log.Warn().Msg("Get empty reload command - fetch mocks from postman")
			mocks = ParsePostmanCollationMock(bodyBytes)
		}
	}

	w.WriteHeader(200)
}

func  postmanRouter(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if  HttpMethod(r.Method) == OPTIONS {
		handleOptionsRequest(&w)
		return
	}
	path := strings.ToLower(r.Method + html.EscapeString(r.URL.Path))
	log.Trace().Msg("requested path: " + path)
	if html.EscapeString(r.URL.RawQuery) != "" {
		path = path + "?" + strings.ToLower(html.EscapeString(r.URL.RawQuery))
	}
	if mock, ok := mocks[path]; ok {
		w.Header().Set("Content-Type", "application/json")
		for _, header := range mock.Header {
			w.Header().Set(header.Key, header.Value)
		}
		if mock.Code > 0 {
			w.WriteHeader(mock.Code)
		}
		fmt.Fprint(w, mock.Body)
		return
	}
	log.Warn().Msg("Requested path not found: " + path)
	w.WriteHeader(404)

}

func handleOptionsRequest(w * http.ResponseWriter){
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Headers", "content-type")
	(*w).Header().Set("Access-Control-Allow-Methods:", "POST,PUT")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).WriteHeader(200)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}