package common


func Map(vs []PostmanHeader, f func(PostmanHeader) Header) []Header {
	vsm := make([]Header, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}