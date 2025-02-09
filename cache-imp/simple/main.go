package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	ID       int
	UserName string
}
type Server struct {
	db    map[int]*User
	cache map[int]*User
	dbHit int
}

func NewServer() *Server {
	db := make(map[int]*User)
	cache := make(map[int]*User)
	for i := 0; i < 100; i++ {
		db[i+1] = &User{
			ID:       i + 1,
			UserName: fmt.Sprintf("user_%d", i+1),
		}
	}
	return &Server{
		db:    db,
		cache: cache,
	}
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	user, ok := s.tryFromCache(id)
	if ok {
		json.NewEncoder(w).Encode(user)
		return
	}

	user, ok = s.db[id]
	if !ok {
		panic("user not found")
	}
	s.dbHit++

	// hydrate the cache
	s.cache[id] = user

	json.NewEncoder(w).Encode(user)
}

func (s *Server) tryFromCache(id int) (*User, bool) {
	user, ok := s.cache[id]
	return user, ok
}

func main() {
}
