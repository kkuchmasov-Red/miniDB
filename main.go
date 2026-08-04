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
	db.SetStorage()

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
		case EXIT, EXIT_FAST:
			os.Exit(0)
		default:
			fmt.Println("Команда не известна " + comandName)
		}
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

func (db *DB) SetStorage() {

	db.NameStorage = "text.txt"
	file, err := db.Storage(os.O_RDONLY | os.O_CREATE)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		args := strings.Split(text, "|")
		if len(args) != 2 {
			fmt.Println("Данные повреждены")
			continue
		}

		id, err := IDInArgs(args, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		newRow := Row{
			ID:   id,
			Name: args[1],
		}
		db.Row = append(db.Row, newRow)

		fmt.Print(text + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func (db DB) Storage(flag int) (*os.File, error) {
	file, err := os.OpenFile(db.NameStorage, flag, 0666)
	if err != nil {
		return file, err
	}
	return file, nil
}

func (db DB) RowToString(row Row) string {
	return fmt.Sprintf("%d|%s", row.ID, row.Name)
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
		db.Row = append(db.Row, newRow)
		db.InsertinStorage(newRow)
	}

}

func (db DB) InsertinStorage(row Row) {
	file, err := db.Storage(os.O_WRONLY | os.O_APPEND | os.O_CREATE)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	file.WriteString(db.RowToString(row) + "\n")
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

	for i := range db.Row {
		if (db.Row)[i].ID == id {
			(db.Row)[i].Name = args[1]
			db.RewriteStorage()
			return
		}
	}
	fmt.Println("id не существует")
}

func (db *DB) RewriteStorage() {
	file, err := db.Storage(os.O_WRONLY)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for i := range db.Row {
		file.WriteString(db.RowToString(db.Row[i]) + "\n")
	}

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

	for i := range db.Row {
		if db.Row[i].ID == id {
			db.Row = append(db.Row[:i], db.Row[i+1:]...)
			db.RewriteStorage()
			return
		}
	}

	fmt.Println("id не существует")
}
func (db *DB) Exists(id int) bool {
	for _, row := range db.Row {
		if row.ID == id {
			return true
		}
	}
	return false
}

func (db *DB) SelectAll() {
	file, err := db.Storage(os.O_RDONLY)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		fmt.Print(text + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	//for _, row := range db.Row {
	//	fmt.Printf("id %d; name %s;\n", row.ID, row.Name)
	//}
}
