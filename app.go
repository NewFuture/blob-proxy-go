// httpForward
// {SUBFOLDER}.{CONTAINER}.{STORAGE}.mydomain.com/xyz/ ==>  {STORAGE}.blob.core.windows.net/{CONTAINER}/{SUBFOLDER}/xyz/index.html
// {CONTAINER}.{STORAGE}.mydomain.com/xyz/ ==>  {STORAGE}.blob.core.windows.net/{CONTAINER}/xyz/index.html
// {STORAGE}.mydomain.com/xyz/ ==>  {STORAGE}.blob.core.windows.net/xyz/index.html
package main

import (
	"github.com/NewFuture/blob-proxy-go/blobproxy"
	"github.com/NewFuture/blob-proxy-go/httpproxy"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getEnv(key string, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return
}
func getEnvInt(key string, defaultValue int) int {
	rawvalue := os.Getenv(key)
	if rawvalue == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(rawvalue)
	if err != nil {
		log.Print(err.Error())
		return defaultValue
	}
	return value
}

var serverPort = getEnvInt("PORT", 80)
var httpProxyFrontend = getEnv("HTTP_PROXY_FRONT_DOMAIN", "localhost")
var forceHTTPS, _ = strconv.ParseBool(getEnv("FORCE_HTTPS", "true"))
var blobProxy = &blobproxy.Blobproxy{
	BasicDomainNum:  getEnvInt("BASIC_DOMAIN_NUM", 0),
	BlobSuffix:      getEnv("BLOB_SUFFIX", "blob.core.windows.net"),
	DefaultDocument: getEnv("DEFAULT_DOCUMENT", "index.html"),
	SplitKey:        getEnv("SPLIT_KEY", "--"),
	ForceHTTPS:      forceHTTPS,
}
var httpProxy = &httpproxy.Httpproxy{
	BackendHost: getEnv("HTTP_PROXY_BACKEND", "localhost:3000"),
}

func doProxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Hostname() == "" {
		r.URL.Host = r.Host
	}

	if r.Header.Get("Upgrade") == "websocket" {
		log.Print("websocket", r.URL)
		httpProxy.WebsocketProxy(w, r)
		return
	}

	var res *http.Response
	var err error
	if strings.HasSuffix(r.URL.Hostname(), httpProxyFrontend) {
		log.Print("http", r.URL)
		res, err = httpProxy.Proxy(r)
	} else {
		log.Print("blob", r.URL)
		res, err = blobProxy.Proxy(r)
	}
	defer res.Body.Close()
	if err != nil {
		log.Println(r.URL, res.Status, err.Error())
	}
	// 复制应答
	for k, v := range res.Header {
		w.Header().Set(k, v[0])
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func main() {
	pattern := "/"
	http.HandleFunc(pattern, doProxy)
	strPort := strconv.Itoa(serverPort)
	log.Print("listenning on :", strPort, pattern)
	err := http.ListenAndServe(":"+strPort, nil)
	if err != nil {
		log.Fatal("fail start on ", strPort, pattern, err.Error())
	}
}
