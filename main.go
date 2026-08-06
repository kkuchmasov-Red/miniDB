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
		case EXIT, EXIT_FAST:
			err := db.RewriteStorage()
			if err != nil {
				fmt.Println(err)
			}
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

func IDInArgs(args []string, position int) (uint8, error) {
	val, err := strconv.Atoi(args[position])
	return uint8(val), err
}

func (row Row) MarshalBinary() (bytes.Buffer, error) {
	var allBufer bytes.Buffer
	var buf bytes.Buffer

	if _, err := buf.Write([]byte{row.ID}); err != nil {
		return allBufer, err
	}
	if _, err := buf.Write([]byte{uint8(len(row.Name))}); err != nil {
		return allBufer, err
	}
	if _, err := buf.WriteString(row.Name); err != nil {
		return allBufer, err
	}

	if err := binary.Write(&allBufer, binary.LittleEndian, uint16(len(buf.Bytes()))); err != nil {
		return allBufer, err
	}
	if _, err := allBufer.Write(buf.Bytes()); err != nil {
		return allBufer, err
	}
	return allBufer, nil
}

func (row Row) UnmarshalBinary(bin []byte) Row {
	row.ID = bin[0:1][0]
	lengthName := bin[1:2][0]
	row.Name = string(bin[2 : lengthName+2])
	return row
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
		var err error
		lenghtRowByte := make([]byte, 2)
		_, err = io.ReadFull(file, lenghtRowByte)
		if err == io.EOF {
			break
		}
		var lenghtRow uint16
		binary.Read(bytes.NewReader(lenghtRowByte), binary.LittleEndian, &lenghtRow)

		rowBinary := make([]byte, lenghtRow)
		_, err = io.ReadFull(file, rowBinary)

		newRow := Row{}.UnmarshalBinary(rowBinary)
		db.Rows = append(db.Rows, newRow)
	}
}

func (db *DB) RewriteStorage() error {

	file, err := db.Storage(os.O_WRONLY | os.O_TRUNC)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, val := range db.Rows {

		buf, err := val.MarshalBinary()
		if err != nil {
			return err
		}
		_, err = file.Write(buf.Bytes())
	}

	return err
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
