package main

import "time"

// Todo : A new type todo
type Todo struct {
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
}

// Todos : A new type todos which is a slice of todo
type Todos []Todo
