package httpproxy

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Httpproxy config
type Httpproxy struct {
	BackendHost string
}

// HTTPProxy proxy http request to backend url
func (p *Httpproxy) HTTPProxy(r *http.Request) (res *http.Response, err error) {
	url := CovertHTTPURL(r.URL, p.BackendHost, "http")
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		res = &http.Response{
			StatusCode: 502,
			Status:     "Create Request Error",
			Body:       ioutil.NopCloser(strings.NewReader(err.Error())),
		}
		return res, err
	}
	req.Header = *CreateProxyRequestHeader(r)
	req.Host = r.Host
	httpClient := &http.Client{}
	res, err = httpClient.Do(req)
	if err != nil {
		res = &http.Response{
			StatusCode: 503,
			Status:     "Inner HTTP Error",
			Body:       ioutil.NopCloser(strings.NewReader(err.Error())),
		}
	}
	return res, err
}
