package httpproxy

import (
	"net/url"
	"testing"
)

func TestCovertHTTPURL(t *testing.T) {

	backendHost := "localhost:3000"
	var reqURL, _ = url.Parse("https://storage.my-proxy.com/xyz/?q=123")
	res := CovertHTTPURL(reqURL, backendHost, "http")
	expectedResult := "http://localhost:3000/xyz/?q=123"
	if res != expectedResult {
		t.Error(res, "should be", expectedResult)
	}

	backendHost = "backend.test.com"
	reqURL, _ = url.Parse("https://user@test.my-proxy.com/xyz/?q=123")
	res = CovertHTTPURL(reqURL, backendHost, "ws")
	expectedResult = "ws://user@backend.test.com/xyz/?q=123"
	if res != expectedResult {
		t.Error(res, "should be", expectedResult)
	}
}
