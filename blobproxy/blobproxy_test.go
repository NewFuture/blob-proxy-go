package blobproxy

import (
	"net/url"
	"testing"
)

func TestCovertBlobURL(t *testing.T) {
	bproxy := &Blobproxy{
		BlobSuffix:      "blob.core.windows.net",
		SplitKey:        "--",
		ForceHTTPS:      true,
		DefaultDocument: "index.html",
	}
	var reqURL, _ = url.Parse("http://storage--container--subfolder.my-proxy.com/xyz/")
	res, err := bproxy.CovertBlobURL(reqURL)
	expectedResult := "https://storage.blob.core.windows.net/container/subfolder/xyz/index.html"
	if res != expectedResult {
		t.Error(res, "should be", expectedResult, err)
	}
}
