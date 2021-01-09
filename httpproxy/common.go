package httpproxy

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

// CovertHTTPURL covert url to backend url
func CovertHTTPURL(reqURL *url.URL, backendHost string, schema string) (strURL string) {
	var proxyURL = reqURL
	proxyURL.Scheme = schema
	proxyURL.Host = backendHost
	return proxyURL.String()
}

// CreateProxyRequestHeader from origin request
func CreateProxyRequestHeader(r *http.Request) *http.Header {
	requestHeader := &http.Header{}
	//用遍历header实现完整复制
	for k, v := range r.Header {
		requestHeader.Set(k, v[0])
	}
	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := r.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		requestHeader.Set("X-Forwarded-For", clientIP)
		if r.TLS != nil {
			requestHeader.Set("X-Forwarded-Proto", "https")
		} else {
			requestHeader.Set("X-Forwarded-Proto", "http")
		}
	}
	requestHeader.Set("Host", r.Host)
	return requestHeader
}
