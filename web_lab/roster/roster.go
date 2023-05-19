package main

import (
	"errors"
	"sync"
)

var _ Roster = (*ClassRoster)(nil)

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

type ClassRoster struct {
	persons sync.Map
}

func NewClassRoster() *ClassRoster {
	return &ClassRoster{}
}

func (c *ClassRoster) Registry(person Person) error {
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

func (c *ClassRoster) Get(name string) (person Person, err error) {
	v, ok := c.persons.Load(name)
	if ok {
		return *v.(*Person), nil
	} else {
		return person, errors.New("查无此人")
	}
}
func (c *ClassRoster) All() (persons []Person, err error) {
	persons = make([]Person, 0)
	c.persons.Range(func(key, value any) bool {
		persons = append(persons, *value.(*Person))
		return true
	})
	return persons, nil
}
