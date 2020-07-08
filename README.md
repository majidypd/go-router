# majix/router


Package `majix/router` implements a fast request router and dispatcher for matching incoming requests to
their respective handler additional supporting url parameters.


This project allows you to have a material framework that you can implement your Microservices very easily.


* Requests can be matched based on URL host, path, path prefix, verbs, header and query values, HTTP methods or using custom matchers.
* URL hosts, paths and query values can have variables with an optional regular expression.
* Supports the addition of middlewares to a Router, which are executed in the order they are added if a match is found, including its subrouters

---

* [Install](#install)
* [Examples](#examples)
* [Matching Routes](#matching-routes)
* [Url Parameters](#url-parameters)
* [Middleware](#middleware)
* [Full Example](#full-example)

---

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get github.com/majidypd/go-router
```

## Examples

Let's start registering a couple of URL paths and handlers:

```go
package main
import "net/http"

func main(){
	router := NewRouter()
	router.Get("/product/(?P<id>[0-9]+)", func(w http.ResponseWriter,r *http.Request,u Tiller) {
            // do something
	})
	router.Start(":8080")
}
```
So what exactly happened here?

Here we register one router mapping URL paths to handler. This is equivalent to how `http.HandleFunc()` works: if an incoming request URL matches one of the paths, the corresponding handler is called passing (`http.ResponseWriter`, `*http.Request`, `Tiller`) as parameters.

Paths can have variables. They are defined using the format `(?P<id>[0-9]+)`. If a regular expression pattern is not defined, the matched variable will be anything until the next slash. For example:

## Matching Routes
```go
router := NewRouter()
router.Get("/product/(?P<id>[0-9]+)", func(w http.ResponseWriter,r *http.Request,u Tiller) {})
router.Get("/products", func(w http.ResponseWriter,r *http.Request,u Tiller) {})
```

## Url Parameters
In order to get url parameters you can use `Tiller` helper as `u.Param(key name)` like:
```go
router.Get("/product/(?P<id>[0-9]+)", func(w http.ResponseWriter,r *http.Request,u Tiller) {
       	fmt.Fprint(w,u.Param("id"))
})
```


This package any verbs you need to implement [REST API](https://dev.socrata.com/docs/verbs.html) such as :
 
> Any: Accept all 
```go
Any(http.ResponseWriter, *http.Request, Tiller)
Get(http.ResponseWriter, *http.Request, Tiller) 
Post(http.ResponseWriter, *http.Request, Tiller) 
Put(http.ResponseWriter, *http.Request, Tiller) 
Patch(http.ResponseWriter, *http.Request, Tiller) 
Options(http.ResponseWriter, *http.Request, Tiller) 
Delete(http.ResponseWriter, *http.Request, Tiller) 
```

## Middleware

Thi package supports the addition of middlewares to a Router, which are executed in the order they are added if a match is found,
including its subrouters. Middlewares are (typically) small pieces of code which take one request, do something with it,
and pass it down to another middleware or the final handler. Some common use cases for middleware are request logging,
header manipulation, or ResponseWriter hijacking.

```go
func main(){
	router := NewRouter()
	router.Get("/list/(?P<id>[0-9]+)", func(w http.ResponseWriter,r *http.Request,u Tiller) {
            // do something
        }).Middleware(authentication, authorization)
	router.Start(":8080")
}

// Authentication Middleware
func authentication(f Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request, u Tiller) {
            // This one goes first
	    f(w, r, u)
	}
}

// Authorization Middleware
func authorization(f Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request, u Tiller) {
            // This one goes after fisrt middleware
	    f(w, r, u)
	}
}
```

## Full Example

```go
package main

import (
	"fmt"
	"net/http"
)

func main(){
	router := NewRouter()
	router.Get("/list/(?P<id>[0-9]+)", func(w http.ResponseWriter,r *http.Request,u Tiller) {
		fmt.Fprint(w,u.Param("id"))
	}).Middleware(checkLogin, authorization)
	router.Start(":8080")
}

func checkLogin(f Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request, u Tiller) {
		f(w, r, u)
	}
}

func authorization(f Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request, u Tiller) {
		f(w, r, u)
	}
}
```