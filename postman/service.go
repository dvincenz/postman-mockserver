package postman

import (
	"fmt"
	. "github.com/dvincenz/postman-mockserver/common"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//todo may not use global variable
var mocks map[string]Mock


func StartServer() {
	var err error
	mocks, err := GetMocksFromPostman()
	if err != nil{
		log.Error().Msg("error in get postman collection " + err.Error())
		return
	}

	log.Info().Msg("total " + strconv.Itoa(len(mocks)) + " mocks found")
	http.HandleFunc("/updatecollections", reloadCollectionHandler)

	log.Info().Msg("Startup mock server....")
	createServer()
}

func StartServerFromStaticFile(){
	path := viper.GetString("static.path")
	LoadStaticPostmanCollection(path)
	log.Info().Msg("total " + strconv.Itoa(len(mocks)) + " mocks found")
	if viper.GetBool("static.watchFile") {
		go func() {
			createServer()
		}()
		fileWatcher(LoadStaticPostmanCollection, path)
	}else{
		createServer()
	}
}

func LoadStaticPostmanCollection(path string){
	mocks = readPostmanFile(path)
}



func createServer(){
	port := strconv.Itoa(viper.GetInt("port"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		postmanRouter(w, r)
	})
	log.Info().Msg("start to listen http://localhost:"+ port)
	http.ListenAndServe("localhost:" +  port, nil)
}


func reloadCollectionHandler(w http.ResponseWriter, r *http.Request){
	if HttpMethod(r.Method) == POST {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil || bodyBytes == nil{
			mocks, err = GetMocksFromPostman()
			log.Debug().Msg("Reload mocks from Postman")
		} else {
			log.Warn().Msg("Get empty reload command - fetch mocks from postman")
			mocks = parsePostmanCollectionMock(bodyBytes)
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
			//todo make header ignore list
			if !strings.EqualFold(header.Key, "content-length") {
				w.Header().Set(header.Key, header.Value)
			}
			//log.Trace().Msgf("add header: " + header.Key)
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


func fileWatcher (executor func(path string), path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Msgf("error: ", err)
	}
	defer watcher.Close()
	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					executor(path)
					log.Info().Msgf("modified file execute ", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error().Msgf("error reloading event:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal().Msgf("error reloading file, you may need to restart the application", err)
	}
	<-done

}
