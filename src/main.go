package main

import (
	"application"
	"fmt"
	"net/http"
	"sort"
	"store"
	"strconv"
)

func main() {
	a := application.NewApplication()
	s := store.NewStore()

	a.Get("/api/todos", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("GET /api/todos")
		todos := s.FindAll()
		list := make([]store.Todo, 0, 0)
		for _, todo := range todos {
			list = append(list, todo)
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].Id < list[j].Id
		})
		application.Json(rw, list)
	})

	a.Post("/api/todos", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("POST /api/todos")
		var t store.Todo
		application.Bind(r, &t)
		s.Create(t)
		application.Json(rw, s.FindAll())
	})

	a.Delete("/api/todos", func(rw http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(application.QueryParam(r, "id"))
		fmt.Println("DELETE /api/todos/", id)
		if err != nil {
			panic(err)
		}
		s.Destroy(id)
		application.Json(rw, s.FindAll())
	})

	a.Put("/api/todos", func(rw http.ResponseWriter, r *http.Request) {
		var t store.Todo
		application.Bind(r, &t)
		fmt.Println("PUT /api/todos/", t.Id)
		s.Update(t)
		application.Json(rw, s.FindAll())
	})

	a.Static("../examples/vanillajs")
	a.Start(":3000")
}
