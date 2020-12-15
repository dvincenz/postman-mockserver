package common


func Map(vs []PostmanHeader, f func(PostmanHeader) Header) []Header {
	vsm := make([]Header, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func TruncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}
