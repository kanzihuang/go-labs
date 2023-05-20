package syntax_labs

import (
	"testing"
)

func TestStringLen(t *testing.T) {
	s := "Hello, 万里"
	if got, want := len(s), 13; got != want {
		t.Errorf("len(%q) got: %d, want: %d", s, got, want)
	}

	num := 0
	for range s {
		num++
	}
	if got, want := num, 9; got != want {
		t.Errorf("count range(%q) got: %d, want: %d", s, got, want)
	}
}

// 修改 []byte(string)，不会影响原 string，原 string 保持不变。
func TestStringModify(t *testing.T) {
	s := "Hello"
	data := []byte(s)
	data[0] = 'h'
	if got, want := data[0], byte('h'); got != want {
		t.Errorf("got: %c, want: %c", got, want)
	}
	if got, want := s[0], byte('H'); got != want {
		t.Errorf("got: %c, want: %c", got, want)
	}
	if got, want := s, "Hello"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}
