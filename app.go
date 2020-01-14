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
var basicDomainNum, _ = getEnvInt("BASIC_DOMAIN_NUM", 0)
var splitKey = getEnv("SPLIT_KEY", "--")
var forceHTTPS, _ = strconv.ParseBool(getEnv("FORCE_HTTPS", "true"))
var defaultDocument = getEnv("DEFAULT_DOCUMENT", "index.html")

// Host for blob
var blobHost = getEnv("BLOB_SUFFIX", "blob.core.windows.net")

func revesre(a []string) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}
func parseURL(url *url.URL) (strURL string, err error) {
	subDomains := strings.Split(url.Hostname(), ".")
	var storageInfo []string
	if basicDomainNum <= 0 {
		storageInfo = strings.Split(subDomains[0], splitKey)
	} else if len(subDomains) <= basicDomainNum {
		err = http.ErrNotSupported
		return
	} else {
		// 域名倒装
		storageInfo = subDomains[:len(subDomains)-basicDomainNum]
		revesre(storageInfo)
	}

	var sb strings.Builder
	if !forceHTTPS && url.Scheme != "" {
		sb.WriteString(url.Scheme)
	} else {
		sb.WriteString("https://")
	}
	sb.WriteString(storageInfo[0]) // blob
	sb.WriteByte('.')
	storageInfo[0] = blobHost
	sb.WriteString(strings.Join(storageInfo, "/"))
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
	log.Print(url)
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
