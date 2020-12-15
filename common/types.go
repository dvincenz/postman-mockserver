package common

type Mock struct {
	RawPath string
	Method  HttpMethod
	Name    string
	Body    string
	Header  []Header
	Code    int
	Path    []string
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

func (m HttpMethod) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case OPTIONS:
		return "OPTIONS"
	case HEAD:
		return "HEAD"
	default:
		return ""
	}
}