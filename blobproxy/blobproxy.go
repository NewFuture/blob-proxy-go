package blobproxy

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Blobproxy config
type Blobproxy struct {
	BlobSuffix      string
	SplitKey        string
	DefaultDocument string
	BasicDomainNum  int
	ForceHTTPS      bool
}

func revesre(a []string) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

// CovertBlobURL covert url to blob url
func (b *Blobproxy) CovertBlobURL(url *url.URL) (strURL string, err error) {
	subDomains := strings.Split(url.Hostname(), ".")
	var storageInfo []string
	if b.BasicDomainNum <= 0 {
		storageInfo = strings.Split(subDomains[0], b.SplitKey)
	} else if len(subDomains) <= b.BasicDomainNum {
		err = http.ErrNotSupported
		return
	} else {
		// 域名倒装
		storageInfo = subDomains[:len(subDomains)-b.BasicDomainNum]
		revesre(storageInfo)
	}

	var sb strings.Builder
	if !b.ForceHTTPS && url.Scheme != "" {
		sb.WriteString(url.Scheme)
	} else {
		sb.WriteString("https://")
	}
	sb.WriteString(storageInfo[0]) // blob
	sb.WriteByte('.')
	storageInfo[0] = b.BlobSuffix
	sb.WriteString(strings.Join(storageInfo, "/"))
	sb.WriteString(url.Path)
	// 自动补全默认文档
	if strings.HasSuffix(url.Path, "/") {
		sb.WriteString(b.DefaultDocument)
	}

	strURL = sb.String()
	return
}

// Proxy proxy http request to blob
func (b *Blobproxy) Proxy(r *http.Request) (res *http.Response, err error) {
	url, err := b.CovertBlobURL(r.URL)
	if err != nil {
		res = &http.Response{
			StatusCode: 400,
			Status:     "RequestURL convert error",
			Body:       ioutil.NopCloser(strings.NewReader("not suport " + r.URL.String())),
		}
		return res, err
	}
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		res = &http.Response{
			StatusCode: 502,
			Status:     "Create Request Error",
			Body:       ioutil.NopCloser(strings.NewReader(err.Error())),
		}
		return res, err
	}
	//用遍历header实现完整复制
	for k, v := range r.Header {
		if k != "Host" {
			req.Header.Set(k, v[0])
		}
	}
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
