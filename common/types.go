package common

type Mock struct {
	Path string
	Method HttpMethod
	Name   string
	Body   string
	Header [] Header
	Code int
}

type Header struct {
	Key string
	Value string
}


type PostmanHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}


type HttpMethod string

const(
	GET  HttpMethod = "GET"
	POST HttpMethod = "POST"
	PUT HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
	OPTIONS HttpMethod = "OPTIONS"
	HEAD HttpMethod = "HEAD"
)