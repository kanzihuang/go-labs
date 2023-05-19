package mysql_test

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

type Person struct {
	Name   string
	Age    int
	Tall   float32
	Weight float32
}

var personMike = Person{
	Name:   "Mike",
	Age:    32,
	Tall:   1.73,
	Weight: 71.5,
}

func TestMysqlPing(t *testing.T) {
	conn := getDbConnection(t)
	defer conn.Close()
	if err := conn.Ping(); err != nil {
		t.Fatal(err)
	}
}

func getDbConnection(t *testing.T) *sql.DB {
	conn, err := sql.Open("mysql", "reader:123456@tcp(mysql.123sou.cn:3306)/lab")
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func TestMysqlSelect(t *testing.T) {
	conn := getDbConnection(t)
	defer conn.Close()

	queryPerson(t, conn)
}

func queryPerson(t *testing.T, conn *sql.DB) {
	rows, err := conn.Query("select name, age, tall, weight from person where name <> ?", personMike.Name)
	if err != nil {
		t.Fatal(err)
	}
	var person Person
	for rows.Next() {
		rows.Scan(&person.Name, &person.Age, &person.Tall, &person.Weight)
		if person != personMike {
			log.Printf("person: %+v\n", person)
		}
	}
}

func TestMysqlInsert(t *testing.T) {
	conn := getDbConnection(t)
	defer conn.Close()

	_, err := conn.Exec("insert person(name, age, tall, weight) values(?, ?, ?, ?)",
		personMike.Name, personMike.Age, personMike.Tall, personMike.Weight)
	if err != nil {
		t.Fatal(err)
	}
	queryPerson(t, conn)
}
