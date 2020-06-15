package postman

import (
	"encoding/json"
	. "github.com/dvincenz/postman-mockserver/common"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"strings"

)

func ReadPostmanFile(path string) map[string]Mock {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Err(err).Msgf("ops, failed to pen file " + path)
	}
	log.Debug().Msg("Successfully Opened json file, read requests...")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return ParsePostmanCollationMock(byteValue)
}

func ParsePostmanCollationMock(payload []byte)map[string]Mock{
	var collation postmanCollation
	json.Unmarshal(payload, &collation)

	mocks := make(map[string]Mock)
	for i := 0; i < len(collation.Collection.Item); i++ {
		mocks = appendMap(mocks, getAllRequest(collation.Collection.Item[i], 0))
	}
	return mocks
}


func getAllRequest(item item, level int) map[string]Mock{
	log.Trace().Msg(strings.Repeat(" ", level)  + item.Name)
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

