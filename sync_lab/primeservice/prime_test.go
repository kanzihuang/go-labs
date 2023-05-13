package primeservice

import (
	"reflect"
	"testing"
)

func TestService_IsPrime(t *testing.T) {
	tests := []struct {
		name   string
		number int
		want   bool
	}{
		{name: "prime", number: 100000007, want: true},
		{name: "not prime", number: 10013, want: false},
		{name: "prime", number: 2, want: true},
		{name: "prime", number: 3, want: true},
		{name: "prime", number: 5, want: true},
		{name: "not prime", number: 25, want: false},
		{name: "prime", number: 73, want: true},
		{name: "not prime", number: -1, want: false},
		{name: "not prime", number: 0, want: false},
		{name: "not prime", number: 1, want: false},
		{name: "not prime", number: 4, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			number, want := tt.number, tt.want
			t.Parallel()
			svc := NewService()
			if got := svc.IsPrime(number); got != want {
				t.Errorf("IsPrime(%v) got %v, want %v", number, got, want)
			}
		})
	}
}

func TestService_GetPrimes(t *testing.T) {
	tests := []struct {
		name  string
		left  int
		right int
		want  []Prime
	}{
		{name: "primes", left: 2, right: 15,
			want: []Prime{2, 3, 5, 7, 11, 13}},
		{name: "primes", left: 20, right: 25,
			want: []Prime{23}},
		{name: "primes", left: 100000, right: 100100,
			want: []Prime{100003, 100019, 100043, 100049, 100057, 100069}},
		{name: "large primes", left: 100000000, right: 100000100,
			want: []Prime{100000007, 100000037, 100000039, 100000049, 100000073, 100000081}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, right, want := tt.left, tt.right, tt.want
			t.Parallel()
			se := NewService()
			got := se.GetPrimes(left, right)
			if reflect.DeepEqual(got, want) == false {
				t.Errorf("GetPrimes(%v, %v) = %v, want %v", left, right, got, want)
			}
		})
	}

}
