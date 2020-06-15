package postman

import (
	"encoding/json"
	. "github.com/dvincenz/postman-mockserver/common"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
)




type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}


func (c *Client) getCollation() (postmanCollation, error) {
	rel := &url.URL{Path: "/collections/6073583-521f6652-56f9-418e-9e9e-bdbf1007cefc"}
	fullUrl := c.BaseURL.ResolveReference(rel)
	request, err := http.NewRequest("GET", fullUrl.String(), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Api-Key", "PMAK-5ed0b0f080008c004d3d1777-9a76967d71de62dfd934d32c6bbad68584")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return postmanCollation{}, err
	}
	defer response.Body.Close()
	var postmanCollation = postmanCollation{}
	json.NewDecoder(response.Body).Decode(&postmanCollation)

	log.Trace().Msgf("collation", postmanCollation)

	return postmanCollation, err
}

func GetMocksFromPostman() (map[string]Mock, error){
	log.Debug().Msg("load collation from postman. This may take some time as u know postman is sometimes sloooow....")
	var client = new(Client)
	client.httpClient = &http.Client{}
	url, err := url.Parse("https://api.getpostman.com/")
	client.BaseURL = url
	if err != nil {
		return map[string]Mock{}, err
	}
	collation, err := client.getCollation()
	if err != nil {
		return map[string]Mock{}, err
	}
	mocks := make(map[string]Mock)
	for i := 0; i < len(collation.Collection.Item); i++ {
		mocks = appendMap(mocks, getAllRequest(collation.Collection.Item[i], 0))
	}
	return mocks, nil
}
