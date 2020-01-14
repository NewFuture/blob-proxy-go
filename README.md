# blob-proxy-go

## blob forward rules

### single level subdomian

`{STORAGE}--{COTANINER}--{SUBFOLDER}.my-proxy.com/xyz` ==> `{STORAGE}.blob.core.windows.net/{CONTAINER}/{SUBFOLDER}/xyz/index.html`

### muti-level subdomain
* `{SUBFOLDER}.{CONTAINER}.{STORAGE}.mydomain.com/xyz/` ==> `{STORAGE}.blob.core.windows.net/{CONTAINER}/{SUBFOLDER}/xyz/index.html`
* `{CONTAINER}.{STORAGE}.mydomain.com/xyz/` ==> `{STORAGE}.blob.core.windows.net/{CONTAINER}/xyz/index.html`
* `{STORAGE}.mydomain.com/xyz/` ==>  `{STORAGE}.blob.core.windows.net/xyz/index.html`

## configuration

| Name | Default | comment |
| :--- | :----: | :---- |
| `BLOB_SUFFIX` | `blob.core.windows.net` | suffix domain for blob |
| `DEFAULT_DOCUMENT` | `index.html` | append default file when list folder  |

| `FORCE_HTTPS` | `true` | using https connect to upstream |

| `SPLIT_KEY` | `--` | using one single level subdaomin as {STORAGE}--{COTANINER}--{SUBFOLDER}--{SubSUBFOLDER}.my-proxy.com|
| `BASIC_DOMAIN_NUM` | 0 | when set will ignore `split_key`|


# docker support

[blob-proxy-go](https://github.com/NewFuture/blob-proxy-go/packages/102924)

* stable

`docker pull docker.pkg.github.com/newfuture/blob-proxy-go/blob-proxy-go`

* for beta

`docker pull newfuture/blob-proxy-go:beta`
