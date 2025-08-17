package core

type Header struct {
	Name  string
	Value string
}

func NewHeader(name, value string) Header {
	return Header{
		Name:  name,
		Value: value,
	}
}

type HttpResponse struct {
	Version    string
	StatusCode uint32
	Reason     string
	Headers    []Header
	Body       string
}

// TODO
func NewHttpResponse(rawResponse string) HttpResponse {
	return HttpResponse{
		Version:    "HTTP/1.1",
		StatusCode: 200,
		Reason:     "OK",
		Headers:    []Header{},
		Body:       rawResponse,
	}
}
