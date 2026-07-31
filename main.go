package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	INSERT = "insert"
	SELECT = "select"
	EXIT   = "exit"
)

type DB []Row

type Row struct {
	id   int
	name string
}

var db DB

func main() {
	for true {
		var comandName string
		fmt.Print("Введите команду - ")
		fmt.Scan(&comandName)
		switch strings.ToLower(comandName) {
		case INSERT:
			insert()
		case SELECT:
			selectAll()
		case EXIT:
			os.Exit(0)
		}
	}
}

func insert() {
	var newRow Row
	lendb := len(db)
	if lendb != 0 {
		newRow.id = db[lendb-1].id + 1
	}
	fmt.Print("Введите name ")
	fmt.Scan(&newRow.name)
	db = append(db, newRow)
}

func selectAll() {
	for i := range db {
		fmt.Printf("id %d; name %s;\n", db[i].id, db[i].name)
	}
}
