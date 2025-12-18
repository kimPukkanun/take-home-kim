package thaibahttext

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestToThaiBahtText_Examples(t *testing.T) {
	cases := []struct {
		in   decimal.Decimal
		want string
	}{
		{decimal.NewFromFloat(1234), "หนึ่งพันสองร้อยสามสิบสี่บาทถ้วน"},
		{decimal.NewFromFloat(33333.75), "สามหมื่นสามพันสามร้อยสามสิบสามบาทเจ็ดสิบห้าสตางค์"},
	}

	for _, tc := range cases {
		got, err := ToThaiBahtText(tc.in)
		if err != nil {
			t.Fatalf("unexpected err for %s: %v", tc.in.String(), err)
		}
		if got != tc.want {
			t.Fatalf("for %s: got %q, want %q", tc.in.String(), got, tc.want)
		}
	}
}

func TestToThaiBahtText_RulesAndEdgeCases(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"zero", "0", "ศูนย์บาทถ้วน"},
		{"one", "1", "หนึ่งบาทถ้วน"},
		{"eleven", "11", "สิบเอ็ดบาทถ้วน"},
		{"twentyone", "21", "ยี่สิบเอ็ดบาทถ้วน"},
		{"oneHundredOne", "101", "หนึ่งร้อยเอ็ดบาทถ้วน"},
		{"oneThousandOne", "1001", "หนึ่งพันเอ็ดบาทถ้วน"},
		{"million", "1000000", "หนึ่งล้านบาทถ้วน"},
		{"oneSatang", "0.01", "ศูนย์บาทหนึ่งสตางค์"},
		{"tenSatang", "0.10", "ศูนย์บาทสิบสตางค์"},
		{"twentyFiveSatang", "0.25", "ศูนย์บาทยี่สิบห้าสตางค์"},
		{"negative", "-12.50", "ลบสิบสองบาทห้าสิบสตางค์"},
		{"roundUp", "1.005", "หนึ่งบาทหนึ่งสตางค์"},
		{"moreThanMillion", "10000000.999", "สิบล้านหนึ่งบาทถ้วน"},
		{"largeNumber", "12000000000", "หนึ่งหมื่นสองพันล้านบาทถ้วน"},
	}

	for _, tc := range cases {
		in, err := decimal.NewFromString(tc.in)
		if err != nil {
			t.Fatalf("bad input %q: %v", tc.in, err)
		}
		got, err := ToThaiBahtText(in)
		if err != nil {
			t.Fatalf("%s: unexpected err: %v", tc.name, err)
		}
		if got != tc.want {
			t.Fatalf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}
