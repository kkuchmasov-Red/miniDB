package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//НЕ тест
const (
	INSERT    = "insert"
	SELECT    = "select"
	EXIT      = "exit"
	UPDATE    = "update"
	DELETE    = "delete"
	EXIT_FAST = "0"
	TEST      = "test"
)

type DB []Row

type Row struct {
	ID   int
	Name string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var db DB
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
			test()
		case EXIT, EXIT_FAST:
			os.Exit(0)
		default:
			fmt.Println("Команда не известна " + comandName)
		}
	}
}

func test() {

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
		*db = append(*db, newRow)
	}
}

func (db *DB) Update(args []string) {
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

	for i := range *db {
		if (*db)[i].ID == id {
			(*db)[i].Name = args[1]
			return
		}
	}
	fmt.Println("id не существует")
}

func (db *DB) Delete(args []string) {
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

	for i := range *db {
		if (*db)[i].ID == id {
			*db = append((*db)[:i], (*db)[i+1:]...)
			return
		}
	}

	fmt.Println("id не существует")
}
func (db *DB) Exists(id int) bool {
	for _, row := range *db {
		if row.ID == id {
			return true
		}
	}
	return false
}

func (db *DB) SelectAll() {
	for _, row := range *db {
		fmt.Printf("id %d; name %s;\n", row.ID, row.Name)
	}
}
