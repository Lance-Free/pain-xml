package iso7064

import "testing"

func TestValid(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{value: "NL91ABNA0417164300", want: true},
		{value: "DE89370400440532013000", want: true},
		{value: "RF18539007547034", want: true},
		{value: "NL91ABNA0417164301"},
		{value: "nl91abna0417164300"},
		{value: "NL91"},
		{value: ""},
	}

	for _, test := range tests {
		if got := Valid(test.value); got != test.want {
			t.Errorf("Valid(%q) = %t; want %t", test.value, got, test.want)
		}
	}
}
