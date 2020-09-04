package period

import (
	"fmt"
	"strings"
)

// used for stages in arithmetic
type period64 struct {
	// always positive values
	years, months, days, hours, minutes, seconds int
	// true if the period is negative
	neg   bool
	input string
}

func (period Period) toPeriod64(input string) *period64 {
	if period.IsNegative() {
		return &period64{
			years: int(-period.years), months: int(-period.months), days: int(-period.days),
			hours: int(-period.hours), minutes: int(-period.minutes), seconds: int(-period.seconds),
			neg:   true,
			input: input,
		}
	}
	return &period64{
		years: int(period.years), months: int(period.months), days: int(period.days),
		hours: int(period.hours), minutes: int(period.minutes), seconds: int(period.seconds),
		input: input,
	}
}

func (p64 *period64) toPeriod() (Period, error) {
	var f []string
	if p64.years > 32767 {
		f = append(f, "years")
	}
	if p64.months > 32767 {
		f = append(f, "months")
	}
	if p64.days > 32767 {
		f = append(f, "days")
	}
	if p64.hours > 32767 {
		f = append(f, "hours")
	}
	if p64.minutes > 32767 {
		f = append(f, "minutes")
	}
	if p64.seconds > 32767 {
		f = append(f, "seconds")
	}

	if len(f) > 0 {
		if p64.input == "" {
			p64.input = p64.String()
		}
		return Period{}, fmt.Errorf("%s: integer overflow occurred in %s", p64.input, strings.Join(f, ","))
	}

	if p64.neg {
		return Period{
			int16(-p64.years), int16(-p64.months), int16(-p64.days),
			int16(-p64.hours), int16(-p64.minutes), int16(-p64.seconds),
		}, nil
	}

	return Period{
		int16(p64.years), int16(p64.months), int16(p64.days),
		int16(p64.hours), int16(p64.minutes), int16(p64.seconds),
	}, nil
}

func (p64 *period64) normalise64(precise bool) *period64 {
	return p64.rippleUp(precise).moveFractionToRight()
}

func (p64 *period64) rippleUp(precise bool) *period64 {
	// remember that the fields are all fixed-point 1E1

	p64.minutes = p64.minutes + (p64.seconds/600)*10
	p64.seconds = p64.seconds % 600

	p64.hours = p64.hours + (p64.minutes/600)*10
	p64.minutes = p64.minutes % 600

	// 32670-(32670/60)-(32670/3600) = 32760 - 546 - 9.1 = 32204.9
	if !precise || p64.hours > 32204 {
		p64.days += (p64.hours / 240) * 10
		p64.hours = p64.hours % 240
	}

	if !precise || p64.days > 32760 {
		dE6 := p64.days * oneE6
		p64.months += dE6 / daysPerMonthE6
		p64.days = (dE6 % daysPerMonthE6) / oneE6
	}

	p64.years = p64.years + (p64.months/120)*10
	p64.months = p64.months % 120

	return p64
}

// moveFractionToRight applies the rule that only the smallest field is permitted to have a decimal fraction.
func (p64 *period64) moveFractionToRight() *period64 {
	// remember that the fields are all fixed-point 1E1

	y10 := p64.years % 10
	if y10 != 0 && (p64.months != 0 || p64.days != 0 || p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0) {
		p64.months += y10 * 12
		p64.years = (p64.years / 10) * 10
	}

	m10 := p64.months % 10
	if m10 != 0 && (p64.days != 0 || p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0) {
		p64.days += (m10 * daysPerMonthE6) / oneE6
		p64.months = (p64.months / 10) * 10
	}

	d10 := p64.days % 10
	if d10 != 0 && (p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0) {
		p64.hours += d10 * 24
		p64.days = (p64.days / 10) * 10
	}

	hh10 := p64.hours % 10
	if hh10 != 0 && (p64.minutes != 0 || p64.seconds != 0) {
		p64.minutes += hh10 * 60
		p64.hours = (p64.hours / 10) * 10
	}

	mm10 := p64.minutes % 10
	if mm10 != 0 && p64.seconds != 0 {
		p64.seconds += mm10 * 60
		p64.minutes = (p64.minutes / 10) * 10
	}

	return p64
}
