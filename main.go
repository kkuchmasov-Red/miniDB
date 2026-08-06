package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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
	Rows        []Row
}

type Row struct {
	ID   uint8
	Name string
}

type Storage struct {
	Name string
}

func main() {
	//главная 1
	scanner := bufio.NewScanner(os.Stdin)

	var db DB
	db.NameStorage = "bd_bin"
	db.LoadStorage()

	for {
		comandName, args := parserCommand(scanner)

		switch strings.ToLower(comandName) {
		case INSERT:
			db.Insert(args)
		case DELETE:
			db.Delete(args)
		case UPDATE:
			db.Update(args)
		case SELECT:
			db.SelectAll()
		case TEST:
			test()
		case EXIT, EXIT_FAST:
			db.RewriteStorage()
			os.Exit(0)
		default:
			fmt.Println("Команда не известна " + comandName)
		}
	}
}

func test() {
	file, err := os.OpenFile("bd_bin", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	type row struct {
		id   uint8
		name string
	}
	var rows []row

	rows = append(rows, row{id: 1, name: "asd"})
	rows = append(rows, row{id: 2, name: "ssss"})
	rows = append(rows, row{id: 3, name: "m,cxcz"})

	var buf bytes.Buffer

	for _, val := range rows {

		binary.Write(&buf, binary.LittleEndian, val.id)
		lenName := uint8(len(val.name))
		binary.Write(&buf, binary.LittleEndian, lenName)
		binary.Write(&buf, binary.LittleEndian, []byte(val.name))
		binary.Write(&buf, binary.LittleEndian, []byte("\n"))
	}
	file.Write(buf.Bytes())

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

func IDInArgs(args []string, position int) (uint8, error) {
	val, err := strconv.Atoi(args[position])
	return uint8(val), err
}

func (db DB) RowToString(row Row) string {
	return fmt.Sprintf("%d|%s", row.ID, row.Name)
}

func (db DB) Storage(flag int) (*os.File, error) {
	file, err := os.OpenFile(db.NameStorage, flag, 0666)
	if err != nil {
		return file, err
	}
	return file, nil
}

func (db *DB) LoadStorage() {
	file, err := db.Storage(os.O_RDONLY | os.O_CREATE)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for {
		id := make([]byte, 1)
		_, err := file.Read(id)

		if err == io.EOF {
			break
		}

		lenName := make([]byte, 1)
		_, err = file.Read(lenName)

		if err == io.EOF {
			break
		}

		name := make([]byte, lenName[0])
		_, err = file.Read(name)

		if err == io.EOF {
			break
		}

		newRow := Row{
			ID:   id[0],
			Name: string(name),
		}
		db.Rows = append(db.Rows, newRow)
	}
}

func (db *DB) RewriteStorage() {
	file, err := db.Storage(os.O_WRONLY | os.O_TRUNC)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var buf bytes.Buffer

	for _, val := range db.Rows {
		binary.Write(&buf, binary.LittleEndian, val.ID)
		binary.Write(&buf, binary.LittleEndian, uint8(len(val.Name)))
		binary.Write(&buf, binary.LittleEndian, []byte(val.Name))
	}
	file.Write(buf.Bytes())
}

func (db *DB) Exists(id uint8) bool {
	for _, row := range db.Rows {
		if row.ID == id {
			return true
		}
	}
	return false
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
		db.Rows = append(db.Rows, newRow)
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

	for i := range db.Rows {
		if (db.Rows)[i].ID == id {
			(db.Rows)[i].Name = args[1]
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

	for i := range db.Rows {
		if db.Rows[i].ID == id {
			db.Rows = append(db.Rows[:i], db.Rows[i+1:]...)
			return
		}
	}

	fmt.Println("id не существует")
}

func (db *DB) SelectAll() {
	for _, row := range db.Rows {
		fmt.Printf("id %d; name %s;\n", row.ID, row.Name)
	}
}
