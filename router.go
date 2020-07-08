package majix

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Tiller interface {
	Param(string) string
}

type Handler func(http.ResponseWriter, *http.Request, Tiller)

type Middleware func(Handler) Handler

type Route struct {
	handler   interface{}
	method    string
	meddlers  []Middleware
	urlParams map[string]string
}

func (r *Route) Middleware(handler ...Middleware) {
	r.meddlers = handler
}

type RouteManagerInterface interface {
	Any(pattern string, handler interface{})
	Get(pattern string, handler interface{}) *Route
	Post(pattern string, handler interface{}) *Route
	Put(pattern string, handler interface{}) *Route
	Patch(pattern string, handler interface{}) *Route
	Options(pattern string, handler interface{}) *Route
	Delete(pattern string, handler interface{}) *Route
	Start(address string)
}

type Util struct {
	route *Route
}

func (u *Util) Param(key string) string {
	return u.route.urlParams[key]
}

type RouteManager struct {
	mapper    map[*regexp.Regexp]*Route
	duplicate map[string]bool
	Util      *Tiller
}

func (r *RouteManager) Error() {

}

func (r *RouteManager) add(method string, pattern string, handler interface{}) *Route {
	_, ok := r.duplicate[pattern]
	if ok {
		panic(pattern + "duplicate route")
	}

	reg, _ := regexp.Compile(pattern)
	route := &Route{
		handler: Handler(handler.(func(http.ResponseWriter, *http.Request, Tiller))),
		method:  method,
	}
	r.mapper[reg] = route
	r.duplicate[pattern] = true

	return route

}

func (r *RouteManager) Start(address string) {
	fmt.Println("Server is running on : ", address)
	http.ListenAndServe(address, r)
}

func (r *RouteManager) Any(pattern string, handler interface{}) {
	r.add("", pattern, handler)
}

func (r *RouteManager) Get(pattern string, handler interface{}) *Route {
	return r.add(http.MethodGet, pattern, handler)
}

func (r *RouteManager) Post(pattern string, handler interface{}) *Route {
	return r.add(http.MethodPost, pattern, handler)
}

func (r *RouteManager) Put(pattern string, handler interface{}) *Route {
	return r.add(http.MethodPut, pattern, handler)
}

func (r *RouteManager) Patch(pattern string, handler interface{}) *Route {
	return r.add(http.MethodPatch, pattern, handler)
}

func (r *RouteManager) Options(pattern string, handler interface{}) *Route {
	return r.add(http.MethodOptions, pattern, handler)
}

func (r *RouteManager) Delete(pattern string, handler interface{}) *Route {
	return r.add(http.MethodDelete, pattern, handler)
}

func (r *RouteManager) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	method := getMethod(req)

	for pattern, route := range r.mapper {

		if matches := pattern.FindStringSubmatch(req.URL.Path); matches != nil {
			if err := checkMethod(w, method, route.method); err != nil {
				return
			}

			urlParams := make(map[string]string)

			for i := 1; i < len(matches); i++ {
				urlParams[pattern.SubexpNames()[i]] = matches[i]
			}
			route.urlParams = urlParams
			//r.App.Util.route = route
			util := &Util{
				route: route,
			}

			revivedFunction := route.handler.(Handler)
			if route.meddlers != nil {
				len := len(route.meddlers) - 1
				revivedFunction = route.meddlers[len](revivedFunction)
				for i := len - 1; i >= 0; i-- {
					revivedFunction = route.meddlers[i](revivedFunction)
				}
			}
			revivedFunction(w, req, util)
			return

		}
	}
	fmt.Fprintf(w, "404 Not Found")
	return
}

func checkMethod(w http.ResponseWriter, method string, routeMethod string) error {
	if routeMethod != "" && strings.ToUpper(method) != routeMethod {
		fmt.Fprintf(w, "Method Not Allowed")
		return errors.New("method not allowed")
	}
	return nil
}

func getMethod(req *http.Request) string {
	var method = req.Method
	if req.Method != http.MethodGet {
		req.ParseForm()
		if req.Form.Get("_method") != "" {
			method = req.Form.Get("_method")
		}
	}
	return method
}

func NewRouter() *RouteManager {
	return &RouteManager{
		mapper:    make(map[*regexp.Regexp]*Route),
		duplicate: make(map[string]bool),
	}
}
