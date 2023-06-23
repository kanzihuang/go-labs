package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/header", header)
	http.HandleFunc("/form", form)
	http.HandleFunc("/body", body)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func body(writer http.ResponseWriter, request *http.Request) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintln(writer, "read all err: ", err)
		return
	}
	writer.Write(b)
}

func form(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Fprintf(writer, "parse form error: %v\n", err)
		return
	}
	fmt.Fprintf(writer, "form: %v\n", request.Form)
}

func header(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "header: %v\n", request.Header)
}

func home(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "Hi there, I love %s!\n", request.URL.Path[1:])
	fmt.Fprintf(writer, "request.Host: %v\n", request.Host)
	fmt.Fprintf(writer, "request.URL.Host: %v\n", request.URL.Host)
	fmt.Fprintf(writer, "request.URL.Path: %v\n", request.URL.Path)
	fmt.Fprintf(writer, "request.URL.RawPath: %v\n", request.URL.RawPath)
}
