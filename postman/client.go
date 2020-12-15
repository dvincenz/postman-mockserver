package postman

import (
	"encoding/json"
	"fmt"
	. "github.com/dvincenz/postman-mockserver/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL   *url.URL
	UserAgent string
	httpClient *http.Client
}

func getCollection(collectionId string) (postmanCollection, error) {
	responseJson, err := getFromPostman("/collections/" + collectionId)
	if err != nil {
		return postmanCollection{}, err
	}
	var postmanCollection = postmanCollection{}
	json.Unmarshal([]byte(responseJson), &postmanCollection)

	return postmanCollection, err
}

func getCollections()(Collections, error) {
	responseJson, err := getFromPostman("/collections")
	if err != nil {
		return Collections{}, err
	}

	var postmanCollections = Collections{}
	json.Unmarshal([]byte(responseJson), &postmanCollections)
	if len(postmanCollections.Collections) == 0 {
		return Collections{}, fmt.Errorf("no collections found")
	}
	return postmanCollections, nil
}

func GetMocksFromPostman() (map[string]Mock, error){
	log.Debug().Msg("load collections from postman...")
	collections, err := getCollections()
	if err != nil {
		return map[string]Mock{}, err
	}
	mocks := make(map[string]Mock)
	for _, collectionOverview  := range collections.Collections {
		if len(viper.GetStringSlice("postman.collections")) != 0 && strings.ToLower(viper.GetStringSlice("postman.collections")[0]) != "all" {
			if !isIdInList(viper.GetStringSlice("postman.collections"), collectionOverview.UID) {
				continue
			}
		}
		//todo make this stuff concurrent
		collection, err :=  getCollection(collectionOverview.UID)
		if err != nil {
			log.Error().Msg("error get mock for collection " + collectionOverview.UID + " this collection would be skipped")
		}
		for i := 0; i < len(collection.Collection.Item); i++ {
			mocks = appendMap(mocks, getAllRequest(collection.Collection.Item[i], 0))
		}
	}
	return mocks, nil
}

func isIdInList(list []string, id string) bool {
	for _, v := range list {
		if strings.ToLower(v) == id {
			return true
		}
	}
	return false
}


func getFromPostman (path string) (string, error){
	return requestPostman(path, "GET", nil)
}

func requestPostman(path string, method string, body io.Reader) (string, error) {
	log.Debug().Msg("send request to postman for " + path + " ...")
	var client = new(Client)
	client.httpClient = &http.Client{}
	fullUrl, err := url.Parse(viper.GetString("postman.url"))
	if err != nil {
		return "", err
	}
	if fullUrl.Host == "" {
		return "", fmt.Errorf("host not available, please check your configuration")
	}
	client.BaseURL = fullUrl
	if viper.GetString("postman.token") == "" {
		return "", fmt.Errorf("postman token is not present in config, please check config file")
	}
	request, err := http.NewRequest(method, fullUrl.String() + path, body)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Api-Key", viper.GetString("postman.token"))
	response, err := client.httpClient.Do(request)
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request to postman failed, " + response.Status)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("request to postman failed, can not read body")
	}
	bodyString := string(bodyBytes)
	log.Trace().Msg("Body get from postman: " + TruncateString(bodyString, 100))
	return bodyString, err
}