package main

import (
	"errors"
	"sync"
)

var _ Roster = (*MapRoster)(nil)

type MapRoster struct {
	persons sync.Map
}

func NewMapRoster() *MapRoster {
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
