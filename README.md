# socks4
A simple socks4 / socks4a dialer for go with authentication support


## Installation

```bash
go get -u github.com/SteffenLoges/socks4
```

## Usage
```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/SteffenLoges/socks4"
)

func main() {

	protocol := socks4.SOCKS4     // socks4.SOCKS4 or socks4.SOCKS4A
	proxyAddr := "127.0.0.1:1080" // proxy address with port
	userID := ""                  // userID string for authentication, leave blank if no authentication is required 
	timeout := time.Second * 5    // use 0 to dial without a timeout

	tr := &http.Transport{
		Dial: socks4.Dialer(protocol, proxyAddr, userID, timeout),
	}

	c := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	resp, err := c.Get("https://en.wikipedia.org/wiki/SOCKS")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

}
```

### or with fasthttp

```go
package main

import (
	"fmt"
	"time"

	"github.com/SteffenLoges/socks4"
	"github.com/valyala/fasthttp"
)

func main() {

	protocol := socks4.SOCKS4     // socks4.SOCKS4 or socks4.SOCKS4A
	proxyAddr := "127.0.0.1:1080" // proxy address with port
	userID := ""                  // userID string for authentication, leave blank if no authentication is required 
	timeout := time.Second * 5

	client := &fasthttp.Client{
		Dial: socks4.FasthttpDialer(protocol, proxyAddr, userID, timeout),
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("https://en.wikipedia.org/wiki/SOCKS")

	err := client.DoTimeout(req, resp, timeout)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

}
```
