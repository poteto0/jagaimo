package http

type IHttpClient interface {
	GET(host string, port uint16, path string) (string, error)
}

type HttpClient struct{}

func NewHttpClient() IHttpClient {
	return &HttpClient{}
}
