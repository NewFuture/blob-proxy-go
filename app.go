// httpForward
// {SUBFOLDER}.{CONTAINER}.{STORAGE}.mydomain.com/xyz/ ==>  {STORAGE}.blob.core.windows.net/{CONTAINER}/{SUBFOLDER}/xyz/index.html
// {CONTAINER}.{STORAGE}.mydomain.com/xyz/ ==>  {STORAGE}.blob.core.windows.net/{CONTAINER}/xyz/index.html
// {STORAGE}.mydomain.com/xyz/ ==>  {STORAGE}.blob.core.windows.net/xyz/index.html
package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
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
func getEnvInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
}

var serverPort, _ = getEnvInt("PORT", 80)
var basicDomainNum, _ = getEnvInt("BASIC_DOMAIN_NUM", 2)
var defaultDocument = getEnv("DEFAULT_DOCUMENT", "index.html")

// Host for blob
var blobHost = getEnv("BLOB_SUFFIX", "blob.core.windows.net")

func parseURL(url *url.URL) (strURL string, err error) {
	subDomains := strings.Split(url.Hostname(), ".")
	length := len(subDomains)
	if length <= basicDomainNum {
		err = http.ErrNotSupported
		return
	}
	subDomains = subDomains[:length-basicDomainNum]
	length = len(subDomains)
	blob := subDomains[length-1]
	var sb strings.Builder
	if url.Scheme != "" {
		sb.WriteString(url.Scheme)
	} else {
		sb.WriteString("https://")
	}
	sb.WriteString(blob)
	sb.WriteByte('.')
	sb.WriteString(blobHost)
	//前缀
	if length > 1 {
		sb.WriteByte('/')
		sb.WriteString(subDomains[length-2])
		for _, prefix := range subDomains[:length-2] {
			sb.WriteByte('/')
			sb.WriteString(prefix)
		}
	}
	sb.WriteString(url.Path)
	// 自动补全默认文档
	if strings.HasSuffix(url.Path, "/") {
		sb.WriteString(defaultDocument)
	}

	strURL = sb.String()
	return
}

func httpForward(pattern string, port int) {
	http.HandleFunc(pattern, doGo)
	strPort := strconv.Itoa(port)
	log.Print("listenning on :", strPort, pattern)
	err := http.ListenAndServe(":"+strPort, nil)
	if err != nil {
		log.Fatal("fail start on", strPort, pattern, err.Error())
	}
}

func doGo(w http.ResponseWriter, r *http.Request) {
	if r.URL.Hostname() == "" {
		r.URL.Host = r.Host
	}

	log.Print(r.URL)
	url, err := parseURL(r.URL)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("not suport " + r.URL.String()))
		log.Print(err.Error())
		return
	}
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Print("http.NewRequest ", err.Error())
		return
	}
	//用遍历header实现完整复制
	for k, v := range r.Header {
		if k != "Host" {
			req.Header.Set(k, v[0])
		}
	}
	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte(err.Error()))
		log.Print("cli.Do(req) ", err.Error())
		return
	}
	defer res.Body.Close()

	// 复制应答
	for k, v := range res.Header {
		w.Header().Set(k, v[0])
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func main() {
	httpForward("/", serverPort)
}
