package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var roster Roster

func main() {
	roster = NewClassRoster()
	m := http.NewServeMux()
	m.Handle("/registry", http.HandlerFunc(registry))
	m.Handle("/query", http.HandlerFunc(query))
	m.Handle("/all", http.HandlerFunc(all))
	if err := http.ListenAndServe(":8080", m); err != nil {
		log.Fatal(err)
	}
}

func registry(writer http.ResponseWriter, request *http.Request) {
	if !strings.EqualFold(request.Method, "post") {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if request.Body == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	data, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("读取数据失败，%v", err)))
		return
	}
	var person Person
	if err := json.Unmarshal(data, &person); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("解析数据失败，%v", err)))
		return
	}
	if err := roster.Registry(person); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(fmt.Sprintf("注册失败，%v", err)))
		return
	}
	writer.Write([]byte("Success"))
}

func query(writer http.ResponseWriter, request *http.Request) {
	if !strings.EqualFold(request.Method, "get") {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	name := request.URL.Query().Get("name")
	if name == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("name参数未设置"))
		return
	}
	person, err := roster.Get(name)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("获取信息失败，%v", err)))
		return
	}
	data, _ := json.Marshal(person)
	if _, err := writer.Write(data); err != nil {
		log.Println("Warning: 发送数据失败，", err)
	}
}
func all(writer http.ResponseWriter, request *http.Request) {
	if !strings.EqualFold(request.Method, "get") {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	persons, _ := roster.All()
	data, _ := json.Marshal(persons)
	writer.Header().Set("Content-Type", "text/json; charset=utf-8")
	if _, err := writer.Write(data); err != nil {
		log.Println("Warning: 发送数据失败，", err)
	}
}
