package mysql_gorm_test

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
)

type Person struct {
	Id     int
	Name   string
	Age    int
	Tall   float32
	Weight float32
}

func (*Person) TableName() string {
	return "person"
}

var personMike = Person{
	Name:   "Mike",
	Age:    32,
	Tall:   1.73,
	Weight: 71.5,
}

func connectGorm(t *testing.T) *gorm.DB {
	conn, err := gorm.Open(mysql.Open("reader:123456@tcp(mysql.123sou.cn:3306)/lab"))
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func TestGormFind(t *testing.T) {
	conn := connectGorm(t)

	queryPerson(t, conn)
}

func queryPerson(t *testing.T, conn *gorm.DB) {
	var persons []Person

	if resp := conn.Where("Name = 'Tom'").Find(&persons); resp.Error != nil {
		t.Error(resp.Error)
	}

	for _, person := range persons {
		log.Printf("person: %+v\n", person)
	}
}

func TestGormCreate(t *testing.T) {
	conn := connectGorm(t)
	person := personMike

	if resp := conn.Create(&person); resp.Error != nil {
		t.Error(resp.Error)
	}
}

func TestGormUpdate(t *testing.T) {
	conn := connectGorm(t)
	person := Person{Id: 1, Name: "Mary", Age: 0}

	if resp := conn.Select("Name", "Age").Updates(&person); resp.Error != nil {
		t.Error(resp.Error)
	}
}

func TestGormDelete(t *testing.T) {
	conn := connectGorm(t)
	person := Person{Id: 3}

	if resp := conn.Delete(&person); resp.Error != nil {
		t.Error(resp.Error)
	}
}
