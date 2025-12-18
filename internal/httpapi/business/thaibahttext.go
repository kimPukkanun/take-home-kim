package thaibahttext

import (
	"strings"

	"github.com/shopspring/decimal"
)

// ToThaiBahtText converts a decimal amount to Thai Baht text.
//
// Rules:
// - Output is Thai text with suffix บาท.
// - If there is no fractional part (satang == 0), append "ถ้วน".
// - If there is a fractional part, append "<satang>สตางค์".
//
// Notes/assumptions:
// - Supports negative values by prefixing "ลบ".
// - Rounds to 2 decimal places using bankers rounding (shopspring/decimal default).
// - Large integers are supported up to 10^30-1 ("แสนล้านล้านล้าน..." style via repeated "ล้าน") (got this one from testing).
func ToThaiBahtText(amount decimal.Decimal) (string, error) {
	sign := ""
	if amount.IsNegative() {
		sign = "ลบ"
		amount = amount.Abs()
	}

	// Normalize to 2 decimal places for satang handling.
	a := amount.Round(2)

	// Split into baht and satang.
	baht := a.Floor()
	satang := a.Sub(baht).Mul(decimal.NewFromInt(100)).Round(0)

	bahtInt := baht.BigInt()
	satangInt := satang.BigInt()

	bahtText := intToThaiText(bahtInt.String())
	if bahtText == "" {
		bahtText = "ศูนย์"
	}

	var b strings.Builder
	b.WriteString(sign)
	b.WriteString(bahtText)
	b.WriteString("บาท")

	if satangInt.Sign() == 0 {
		b.WriteString("ถ้วน")
		return b.String(), nil
	}

	sSatang := satangInt.String()
	if len(sSatang) == 1 {
		sSatang = "0" + sSatang
	}
	sText := intToThaiText(sSatang)
	// satang is 1..99 here, but keep a safe fallback.
	if sText == "" {
		sText = "ศูนย์"
	}
	b.WriteString(sText)
	b.WriteString("สตางค์")
	return b.String(), nil
}

var thaiDigits = []string{"ศูนย์", "หนึ่ง", "สอง", "สาม", "สี่", "ห้า", "หก", "เจ็ด", "แปด", "เก้า"}
var thaiPositions = []string{"", "สิบ", "ร้อย", "พัน", "หมื่น", "แสน"}

// intToThaiText converts a non-negative integer string (base-10, no sign) to Thai reading text.
// It supports arbitrarily large values by grouping in 6-digit chunks separated by "ล้าน".
func intToThaiText(numStr string) string {
	numStr = strings.TrimLeft(numStr, "0")
	if numStr == "" {
		return ""
	}

	// Split to 6-digit groups from the right.
	groups := split6(numStr)

	var out strings.Builder
	for i, g := range groups {
		gText := groupToThaiText(g)
		if gText != "" {
			out.WriteString(gText)
		}

		// Add "ล้าน" between groups except after last.
		// This must be emitted even if the lower group is zero (e.g., 1,000,000 = หนึ่งล้าน).
		if i != len(groups)-1 {
			out.WriteString("ล้าน")
		}
	}
	return out.String()
}

func split6(s string) []string {
	// returns left-to-right groups, each up to 6 digits.
	var rev []string
	for len(s) > 0 {
		start := 0
		if len(s) > 6 {
			start = len(s) - 6
		}
		rev = append(rev, s[start:])
		s = s[:start]
	}
	// reverse to left-to-right
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	return rev
}

func groupToThaiText(group string) string {
	group = strings.TrimLeft(group, "0")
	if group == "" {
		return ""
	}

	// Left pad to 6 digits.
	if len(group) < 6 {
		group = strings.Repeat("0", 6-len(group)) + group
	}

	digits := make([]int, 6)
	for i := 0; i < 6; i++ {
		digits[i] = int(group[i] - '0')
	}

	var out strings.Builder
	for pos := 5; pos >= 0; pos-- {
		d := digits[5-pos] // map to left-to-right? We'll compute index differently.
		_ = d
	}

	// Iterate positions from left (แสน) to right (หน่วย)
	for i := 0; i < 6; i++ {
		d := digits[i]
		place := 5 - i // 5..0 => แสน..หน่วย

		if d == 0 {
			continue
		}

		switch place {
		case 1: // tens
			// "หนึ่งสิบ" -> "สิบ"
			if d == 1 {
				out.WriteString("สิบ")
				continue
			}
			// "สองสิบ" -> "ยี่สิบ"
			if d == 2 {
				out.WriteString("ยี่")
				out.WriteString("สิบ")
				continue
			}
			out.WriteString(thaiDigits[d])
			out.WriteString("สิบ")
			continue
		case 0: // ones
			// If ones is 1 and there is any non-zero digit before it in this group -> "เอ็ด"
			// (e.g., 11 = สิบเอ็ด, 101 = หนึ่งร้อยเอ็ด, 1001 = หนึ่งพันเอ็ด)
			if d == 1 {
				for j := 0; j < i; j++ {
					if digits[j] != 0 {
						out.WriteString("เอ็ด")
						goto nextDigit
					}
				}
			}
			out.WriteString(thaiDigits[d])
		nextDigit:
			continue
		default:
			out.WriteString(thaiDigits[d])
			out.WriteString(thaiPositions[place])
		}
	}

	return out.String()
}
