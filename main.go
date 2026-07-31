package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	INSERT = "insert"
	SELECT = "select"
	EXIT   = "exit"
)

type DB []Row

type Row struct {
	ID   int
	Name string
}

var db DB

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		var comandName string
		var args []string

		determineCommand(scanner, &comandName, &args)

		switch strings.ToLower(comandName) {
		case INSERT:
			insert(&args)
		case SELECT:
			selectAll()
		case EXIT:
			os.Exit(0)
		default:
			fmt.Println("Команда не известна " + comandName)
		}
	}
}

func determineCommand(scanner *bufio.Scanner, comandName *string, args *[]string) {
	vScanner := *scanner
	if !vScanner.Scan() {
		return
	}

	line := vScanner.Text()
	stringSplit := strings.Split(line, " ")
	if len(stringSplit) == 0 {
		return
	}

	*comandName = stringSplit[0]
	if len(stringSplit) > 1 {
		*args = stringSplit[1:]
	}
}

func insert(args *[]string) {
	vArgs := *args
	if len(vArgs) == 2 {
		id, err := strconv.Atoi(vArgs[0])
		if err != nil {
			fmt.Println("id не равен числу")
			return
		}
		newRow := Row{
			ID:   id,
			Name: vArgs[1],
		}
		db = append(db, newRow)
	} else {
		fmt.Println("Неверное число параметров")
	}

}

func selectAll() {
	for i := range db {
		fmt.Printf("id %d; name %s;\n", db[i].ID, db[i].Name)
	}
}
