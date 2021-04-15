package store

import "fmt"

type Todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Store struct {
	sequence int
	db       map[int]Todo
}

func NewStore() *Store {
	s := &Store{}
	s.sequence = 1
	s.db = make(map[int]Todo)
	return s
}

func (s *Store) FindAll() map[int]Todo {
	return s.db
}

func (s *Store) Find(id int) (Todo, bool) {
	val, exist := s.db[id]
	return val, exist
}

func (s *Store) Create(t Todo) {
	t.Id = s.sequence
	s.sequence++
	fmt.Println(t)
	s.db[t.Id] = t
}

func (s *Store) Update(t Todo) {
	todo, found := s.Find(t.Id)
	if found {
		todo.Completed = t.Completed
	}
}

func (s *Store) Destroy(id int) {
	_, found := s.Find(id)
	if found {
		delete(s.db, id)
	}
}
