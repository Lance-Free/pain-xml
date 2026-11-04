package reference

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		raw    string
		want   CreditorReference
		wantOK bool
	}{
		{name: "dutch", raw: "5000056789012345", want: CreditorReference{Value: "5000056789012345", Issuer: IssuerCurrence}, wantOK: true},
		{name: "dutch with spaces", raw: "5000 0567 8901 2345", want: CreditorReference{Value: "5000056789012345", Issuer: IssuerCurrence}, wantOK: true},
		{name: "dutch wrong check digit", raw: "4000056789012345"},
		{name: "dutch too short", raw: "000056789012345"},
		{name: "dutch with letter", raw: "500005678901234A"},
		{name: "rf", raw: "RF18539007547034", want: CreditorReference{Value: "RF18539007547034", Issuer: IssuerISO}, wantOK: true},
		{name: "rf with spaces and lower case", raw: "rf18 5390 0754 7034", want: CreditorReference{Value: "RF18539007547034", Issuer: IssuerISO}, wantOK: true},
		{name: "rf wrong check digits", raw: "RF19539007547034"},
		{name: "rf too long", raw: "RF18539007547034539007547034"},
		{name: "rf with symbol", raw: "RF18-539007547034"},
		{name: "invoice number", raw: "INV-2026-0042"},
		{name: "empty", raw: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := Parse(test.raw)
			if ok != test.wantOK || got != test.want {
				t.Fatalf("Parse(%q) = %+v, %t; want %+v, %t", test.raw, got, ok, test.want, test.wantOK)
			}
		})
	}
}
