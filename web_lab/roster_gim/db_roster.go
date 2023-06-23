package main

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

var _ Roster = (*DbRoster)(nil)

type DbRoster struct {
	conn *gorm.DB
}

func NewDbRoster(conn *gorm.DB) *DbRoster {
	return &DbRoster{conn: conn}
}

func (roster *DbRoster) All() ([]Person, error) {
	persons := make([]Person, 0)
	if tx := roster.conn.Find(&persons); tx.Error != nil {
		log.Println("数据库读取失败, ", tx.Error)
		return nil, tx.Error
	}
	return persons, nil
}

func (roster *DbRoster) Registry(person Person) error {
	if tx := roster.conn.Create(&person); tx != nil {
		return tx.Error
	}
	return nil
}

func (roster *DbRoster) Get(name string) (person Person, err error) {
	persons := make([]Person, 0)
	if tx := roster.conn.Where("name=?", name).Find(&persons); tx.Error != nil {
		log.Println("数据库读取失败, ", tx.Error)
		return person, tx.Error
	}
	if len(persons) == 0 {
		return person, errors.New("查无此人")
	}
	return persons[0], nil
}
