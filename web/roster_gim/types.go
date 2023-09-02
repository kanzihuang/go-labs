package main

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

type Roster interface {
	Registry(person Person) error
	Get(name string) (Person, error)
	All() ([]Person, error)
}
