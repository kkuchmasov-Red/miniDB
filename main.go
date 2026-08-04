package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	INSERT    = "insert"
	SELECT    = "select"
	EXIT      = "exit"
	UPDATE    = "update"
	DELETE    = "delete"
	EXIT_FAST = "0"
	TEST      = "test"
)

type DB struct {
	NameStorage string
	Row         []Row
}

type Row struct {
	ID   int
	Name string
}

type Storage struct {
	Name string
}

func main() {
	//главная 1
	scanner := bufio.NewScanner(os.Stdin)
	var db DB
	db.NameStorage = "text.txt"
	for {
		comandName, args := parserCommand(scanner)

		switch strings.ToLower(comandName) {
		case INSERT:
			db.Insert(args)
		case UPDATE:
			db.Update(args)
		case DELETE:
			db.Delete(args)
		case SELECT:
			db.SelectAll()
		case TEST:
			db.test()
		case EXIT, EXIT_FAST:
			os.Exit(0)
		default:
			fmt.Println("Команда не известна " + comandName)
		}
	}
}

func (db *DB) test() {
	//db_val := *db

	var err error
	_, err = os.Stat(db.NameStorage)

	if err != nil {
		_, err = os.Create(db.NameStorage)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.OpenFile(db.NameStorage, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	file.WriteString("asd\n")

	if err != nil {
		fmt.Println(err)
		return
	}
}

func parserCommand(scanner *bufio.Scanner) (string, []string) {
	var comandName string
	var args []string
	if !scanner.Scan() {
		return comandName, args
	}

	line := scanner.Text()
	stringSplit := strings.Fields(line)
	if len(stringSplit) == 0 {
		return comandName, args
	}

	comandName = stringSplit[0]
	if len(stringSplit) > 1 {
		args = stringSplit[1:]
	}
	return comandName, args
}

func CheckArgs(args []string, needCount int) error {
	if len(args) != needCount {
		return errors.New("Неверное количество аргументов")
	}
	return nil
}

func IDInArgs(args []string, position int) (int, error) {
	return strconv.Atoi(args[position])
}

func (db *DB) Insert(args []string) {
	db_val := *db
	err := CheckArgs(args, 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	id, err := IDInArgs(args, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	if db.Exists(id) {
		fmt.Println("id уже существует")
	} else {
		newRow := Row{
			ID:   id,
			Name: args[1],
		}
		db_val.Row = append(db_val.Row, newRow)
	}
}

func (db *DB) Update(args []string) {
	db_val := *db
	err := CheckArgs(args, 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	id, err := IDInArgs(args, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := range db_val.Row {
		if (db_val.Row)[i].ID == id {
			(db_val.Row)[i].Name = args[1]
			return
		}
	}
	fmt.Println("id не существует")
}

func (db *DB) Delete(args []string) {
	db_val := *db
	err := CheckArgs(args, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	id, err := IDInArgs(args, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := range db_val.Row {
		if (db_val.Row)[i].ID == id {
			db_val.Row = append((db_val.Row)[:i], (db_val.Row)[i+1:]...)
			return
		}
	}

	fmt.Println("id не существует")
}
func (db *DB) Exists(id int) bool {
	db_val := *db
	for _, row := range db_val.Row {
		if row.ID == id {
			return true
		}
	}
	return false
}

func (db *DB) SelectAll() {
	db_val := *db
	for _, row := range db_val.Row {
		fmt.Printf("id %d; name %s;\n", row.ID, row.Name)
	}
}
