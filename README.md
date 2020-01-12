# blob-proxy-go

## blob forward rules

* `{SUBFOLDER}.{CONTAINER}.{STORAGE}.mydomain.com/xyz/` ==> `{STORAGE}.blob.core.windows.net/{CONTAINER}/{SUBFOLDER}/xyz/index.html`
* `{CONTAINER}.{STORAGE}.mydomain.com/xyz/` ==> `{STORAGE}.blob.core.windows.net/{CONTAINER}/xyz/index.html`
* `{STORAGE}.mydomain.com/xyz/` ==>  `{STORAGE}.blob.core.windows.net/xyz/index.html`

## configuration

| Name | Default | comment |
| :--- | :----: | ----: |
| `BLOB_SUFFIX` | `blob.core.windows.net` | suffix domain for blob |
| `DEFAULT_DOCUMENT` | `index.html` | append default file when list folder  |
| `BASIC_DOMAIN_NUM` | 2 | the basic domain count to ignore. set 3 for {STORAGE}.preview.my.com|

