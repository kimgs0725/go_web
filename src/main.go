package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type Route struct {
	method  string
	pattern *regexp.Regexp
	handler http.Handler
}

type Application struct {
	routes []*Route
}

func init() {

}

func (a *Application) Add(method, path string, handler http.Handler) {
	a.routes = append(a.routes, &Route{
		method,
		regexp.MustCompile(path),
		handler,
	})
}

func (a *Application) AddFunc(method, path string, handler func(rw http.ResponseWriter, r *http.Request)) {
	a.Add(method, path, http.HandlerFunc(handler))
}

func (a *Application) Get(path string, handler func(rw http.ResponseWriter, r *http.Request)) {
	a.AddFunc(http.MethodGet, path, handler)
}

func (a *Application) Post(path string, handler func(rw http.ResponseWriter, r *http.Request)) {
	a.AddFunc(http.MethodPost, path, handler)
}

func (a *Application) Delete(path string, handler func(rw http.ResponseWriter, r *http.Request)) {
	a.AddFunc(http.MethodDelete, path, handler)
}

func (a *Application) Put(path string, handler func(rw http.ResponseWriter, r *http.Request)) {
	a.AddFunc(http.MethodPut, path, handler)
}

func (a *Application) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	for _, route := range a.routes {
		matched := route.pattern.MatchString(r.URL.Path) && route.method == r.Method
		if matched {
			route.handler.ServeHTTP(rw, r)
			return
		}
	}
	http.NotFound(rw, r)
}

func (a *Application) Static(root string) {
	fs := http.FileServer(http.Dir(root))
	a.Add(http.MethodGet, "/*", fs)
}

func (a *Application) Start(port string) {
	fmt.Println("server is running http://localhost:3000")
	http.ListenAndServe(port, a)
}

func NewApplication() *Application {
	return &Application{}
}

type Todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Store struct {
	db []Todo
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) FindAll() []Todo {
	return s.db
}

func (s *Store) FindIndex(id int) (int, bool) {
	var foundIdx int
	for i, todo := range s.db {
		if todo.Id == id {
			foundIdx = i
			break
		}
	}
	return foundIdx, foundIdx > -1
}

func (s *Store) Create(t Todo) {
	s.db = append(s.db, t)
}

func (s *Store) Update(t Todo) {
	i, found := s.FindIndex(t.Id)
	if found {
		s.db[i].Completed = t.Completed
	}
}

func (s *Store) Destroy(id int) {
	i, found := s.FindIndex(id)
	if found {
		s.db = append(s.db[:i], s.db[i+1:]...)
	}
}

func Bind(r *http.Request, i interface{}) {
	json.NewDecoder(r.Body).Decode(i)
}

func QueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func Json(rw http.ResponseWriter, i interface{}) {
	enc := json.NewEncoder(rw)
	enc.Encode(&i)
}

func main() {
	a := NewApplication()
	s := NewStore()

	a.Get("/api/todos", func(rw http.ResponseWriter, r *http.Request) {
		Json(rw, s.FindAll())
	})

	a.Post("/api/todos", func(rw http.ResponseWriter, r *http.Request) {
		var t Todo
		Bind(r, &t)
		s.Create(t)
		Json(rw, s.FindAll())
	})

	a.Delete("/api/todos", func(rw http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(QueryParam(r, "id"))
		fmt.Println("id = ", id)
		if err != nil {
			panic(err)
		}
		s.Destroy(id)
		Json(rw, s.FindAll())
	})

	a.Put("/api/todos", func(rw http.ResponseWriter, r *http.Request) {
		var t Todo
		Bind(r, &t)
		s.Update(t)
		Json(rw, s.FindAll())
	})

	a.Static("../examples/vanillajs")
	a.Start(":3000")
}
