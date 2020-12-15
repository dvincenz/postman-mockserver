package postman

import (
	. "github.com/dvincenz/postman-mockserver/common"
)

type postmanCollection struct {
	Collection Collection `json:"collection"`
}

type Collection struct {
	Info struct {
		PostmanID string `json:"_postman_id"`
		Name      string `json:"name"`
		Schema    string `json:"schema"`
	} `json:"info"`
	Item []item `json:"item"`
	Auth struct {
		Type   string `json:"type"`
		Bearer []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"bearer"`
	} `json:"auth"`
	Event []struct {
		Listen string `json:"listen"`
		Script struct {
			ID   string   `json:"id"`
			Type string   `json:"type"`
			Exec []string `json:"exec"`
		} `json:"script"`
	} `json:"event"`
	Variable []struct {
		ID    string `json:"id"`
		Key   string `json:"key"`
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"variable"`
}

type item struct {
	Name                    string `json:"name"`
	ProtocolProfileBehavior struct {
		DisableBodyPruning bool `json:"disableBodyPruning"`
	} `json:"protocolProfileBehavior"`
	Request   request   `json:"request"`
	Response  []respone `json:"response"`
	Item      []item    `json:"item"`
	PostmanID string    `json:"_postman_id"`
	Event     []struct {
		Listen string `json:"listen"`
		Script struct {
			ID   string   `json:"id"`
			Type string   `json:"type"`
			Exec []string `json:"exec"`
		} `json:"script"`
	} `json:"event,omitempty"`
}

type request struct {
	Method string          `json:"method"`
	Header []PostmanHeader `json:"header"`
	URL    struct {
		Raw      string   `json:"raw"`
		Protocol string   `json:"protocol"`
		Host     []string `json:"host"`
		Port     string   `json:"port"`
		Path     []string `json:"path"`
	} `json:"url"`
}

type respone struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	OriginalRequest struct {
		Method string          `json:"method"`
		Header []PostmanHeader `json:"header"`
		URL    struct {
			Raw  string   `json:"raw"`
			Path []string `json:"path"`
		} `json:"url"`
	} `json:"originalRequest"`
	Status                 string          `json:"status"`
	Code                   int             `json:"code"`
	PostmanPreviewlanguage string          `json:"_postman_previewlanguage"`
	Header                 []PostmanHeader `json:"header"`
	Cookie                 []string        `json:"cookie"`
	ResponseTime           int             `json:"responseTime"`
	Body                   string          `json:"body"`
}

type Collections struct {
	Collections []CollectionOverview `json:"collections"`
}

type CollectionOverview struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
	UID   string `json:"uid"`
}
