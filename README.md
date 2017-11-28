# tokenauth

TokenAuth is a authorization middleware for [gin](https://github.com/gin-gonic/gin).

## Installation

```
go get github.com/vaijab/gin-tokenauth
```

## Usage

This example uses file based token store, see [below](#filestore) for more
details how filestore works.

```Go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vaijab/gin-tokenauth"
	"github.com/vaijab/gin-tokenauth/filestore"
)

func main() {
	store, err := filestore.New("tokens.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	r.Use(tokenauth.New(store))

	r.GET("/secrets", func(c *gin.Context) {
		c.String(200, "p4ssw0rd\n")
	})

	r.Run()
}
```

```shell
> curl -i http://localhost:8080/secrets -H 'Authorization: Bearer jUyaoAFFZ5Ay3fxXG2boT5'
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Fri, 27 Oct 2017 12:11:16 GMT
Content-Length: 9

p4ssw0rd
```

## Token Stores

Different token stores can be implemented quite easily, the API stays the same.

### filestore

This store is based on a yaml file. During initialization, a file watcher is
attached which ensures that changes to the tokens file are reflected
immediately.

Tokens file does not have to exist at first. It can be created or removed and
filestore will either create tokens or remove them entirely. Only `token` and
`is_disabled` fields are used at the moment.

```yaml
# tokens.yaml
---
tokens:
  - name: foo
    token: 'jUyaoAFFZ5Ay3fxXG2boT5'
    is_disabled: false
    description: 'Token for user foo'
  - name: bar
    token: 'jUyaoAFFZ5Ay3fxXG2boT5'
    is_disabled: true
    description: 'Disabled token for user bar'
```
