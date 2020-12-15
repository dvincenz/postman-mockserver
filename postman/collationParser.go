package postman

import (
	"encoding/json"
	. "github.com/dvincenz/postman-mockserver/common"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"strings"

)

func readPostmanFile(path string) map[string]Mock {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Err(err).Msgf("ops, failed to open file " + path)
		return make(map[string]Mock)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return parsePostmanCollectionMock(byteValue)
}

func parsePostmanCollectionMock(payload []byte)map[string]Mock{
	var collection postmanCollection
	json.Unmarshal(payload, &collection)
	if collection.Collection.Info.PostmanID == "" {
		log.Info().Msg("No postman collection object in json found - search for items in json root")
		json.Unmarshal(payload, &collection.Collection)
	}
	mocks := make(map[string]Mock)
	for i := 0; i < len(collection.Collection.Item); i++ {
		mocks = appendMap(mocks, getAllRequest(collection.Collection.Item[i], 0))
	}
	return mocks
}


func getAllRequest(item item, level int) map[string]Mock{
	log.Trace().Msg("Mock: " + strings.Repeat(" ", level)  + item.Name)
	mocks := make(map[string]Mock)
	mocks = appendMap(mocks, getMocks(item.Response))

	for n := 0; n < len(item.Item); n++ {
		mocks = appendMap(mocks, getAllRequest(item.Item[n], level + 1))
	}
	return mocks
}

func getMocks(responses []respone) map[string]Mock{
	mocks := make(map[string]Mock)
	for i := 0; i < len(responses); i++ {
		mock := Mock{
			Method: HttpMethod(responses[i].OriginalRequest.Method),
			Code: responses[i].Code,
			Name: responses[i].Name,
			Body: responses[i].Body,
			Header: Map(responses[i].Header, parseHeaders),
		}

		mocks[strings.ToLower(responses[i].OriginalRequest.Method + responses[i].OriginalRequest.URL.Raw)] = mock
	}
	return mocks
}

func parseHeaders(postmanHeaders PostmanHeader) Header {
	return Header {
		Key: postmanHeaders.Key,
		Value: postmanHeaders.Value,
	}
}


func appendMap(m1 map[string]Mock, m2 map[string]Mock) map[string]Mock{
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

