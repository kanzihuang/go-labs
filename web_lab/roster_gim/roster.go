package main

import (
	"errors"
	"sync"
)

var _ Roster = (*MapRoster)(nil)

type Person struct {
	Name   string
	Age    int
	Tall   float32
	Weight float32
}

type Roster interface {
	Registry(person Person) error
	Get(name string) (Person, error)
	All() ([]Person, error)
}

type MapRoster struct {
	persons sync.Map
}

func NewClassRoster() *MapRoster {
	return &MapRoster{}
}

func (c *MapRoster) Registry(person Person) error {
	if person.Name != "" {
		_, loaded := c.persons.LoadOrStore(person.Name, &person)
		if loaded {
			return errors.New("已注册")
		} else {
			return nil
		}
	} else {
		return errors.New("姓名为空")
	}
}

func (c *MapRoster) Get(name string) (person Person, err error) {
	v, ok := c.persons.Load(name)
	if ok {
		return *v.(*Person), nil
	} else {
		return person, errors.New("查无此人")
	}
}
func (c *MapRoster) All() (persons []Person, err error) {
	persons = make([]Person, 0)
	c.persons.Range(func(key, value any) bool {
		persons = append(persons, *value.(*Person))
		return true
	})
	return persons, nil
}
