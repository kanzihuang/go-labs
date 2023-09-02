package encoding

import (
	"encoding/json"
	"io"
	"os"
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

func TestJsonLoad(t *testing.T) {
	f, err := os.Open("./persons.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	testJsonLoadFrom(t, f)
}

func testJsonLoadFrom(t *testing.T, file *os.File) {
	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	var persons []Person
	err = json.Unmarshal(data, &persons)
	if err != nil {
		t.Fatal(err)
	}

	for _, person := range persons {
		if person != personMike {
			t.Errorf("got: %+v, want: %+v\n", person, personMike)
		}
	}
}

func TestJsonSave(t *testing.T) {
	var err error
	var f *os.File
	f, err = os.CreateTemp("", "persons.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	persons := []Person{
		personMike, personMike,
	}
	var data []byte
	data, err = json.Marshal(persons)
	f.Write(data)
	f.Seek(0, io.SeekStart)
	testJsonLoadFrom(t, f)
}
