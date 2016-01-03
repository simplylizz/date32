// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clock

import (
	"testing"
	"time"
)

func TestClockHoursMinutesSeconds(t *testing.T) {
	cases := []struct {
		in          Clock
		h, m, s, ms int
	}{
		{New(0, 0, 0, 0), 0, 0, 0, 0},
		{New(1, 2, 3, 4), 1, 2, 3, 4},
		{New(23, 59, 59, 999), 23, 59, 59, 999},
		{New(0, 0, 0, -1), 23, 59, 59, 999},
		{NewAt(time.Date(2015, 12, 4, 18, 50, 42, 173444111, time.UTC)), 18, 50, 42, 173},
	}
	for i, x := range cases {
		h := x.in.Hours()
		m := x.in.Minutes()
		s := x.in.Seconds()
		ms := x.in.Millisec()
		if h != x.h || m != x.m || s != x.s || ms != x.ms {
			t.Errorf("%d: got %02d:%02d:%02d.%03d, want %v (%d)", i, h, m, s, ms, x.in, x.in)
		}
	}
}

func TestClockSinceMidnight(t *testing.T) {
	cases := []struct {
		in Clock
		d  time.Duration
	}{
		{New(1, 2, 3, 4), time.Hour + 2*time.Minute + 3*time.Second + 4},
		{New(23, 59, 59, 999), 23*time.Hour + 59*time.Minute + 59*time.Second + 999},
	}
	for i, x := range cases {
		d := x.in.DurationSinceMidnight()
		c2 := SinceMidnight(d)
		if c2 != x.in {
			t.Errorf("%d: got %v, want %v (%d)", i, c2, x.in, x.in)
		}
	}
}

func TestClockIsInOneDay(t *testing.T) {
	cases := []struct {
		in   Clock
		want bool
	}{
		{New(0, 0, 0, 0), true},
		{New(24, 0, 0, 0), true},
		{New(-24, 0, 0, 0), false},
		{New(48, 0, 0, 0), false},
		{New(0, 0, 0, 1), true},
		{New(2, 0, 0, 1), true},
		{New(-1, 0, 0, 0), false},
		{New(0, 0, 0, -1), false},
	}
	for _, x := range cases {
		got := x.in.IsInOneDay()
		if got != x.want {
			t.Errorf("%v got %v, want %v", x.in, x.in.IsInOneDay(), x.want)
		}
	}
}

func TestClockAdd(t *testing.T) {
	cases := []struct {
		h, m, s, ms int
		in, want    Clock
	}{
		{0, 0, 0, 0, 2 * ClockHour, New(2, 0, 0, 0)},
		{0, 0, 0, 1, 2 * ClockHour, New(2, 0, 0, 1)},
		{0, 0, 0, -1, 2 * ClockHour, New(1, 59, 59, 999)},
		{0, 0, 1, 0, 2 * ClockHour, New(2, 0, 1, 0)},
		{0, 0, -1, 0, 2 * ClockHour, New(1, 59, 59, 0)},
		{0, 1, 0, 0, 2 * ClockHour, New(2, 1, 0, 0)},
		{0, -1, 0, 0, 2 * ClockHour, New(1, 59, 0, 0)},
		{1, 0, 0, 0, 2 * ClockHour, New(3, 0, 0, 0)},
		{-1, 0, 0, 0, 2 * ClockHour, New(1, 0, 0, 0)},
		{-2, 0, 0, 0, 2 * ClockHour, New(0, 0, 0, 0)},
		{-2, 0, -1, -1, 2 * ClockHour, New(0, 0, -1, -1)},
	}
	for i, x := range cases {
		got := x.in.Add(x.h, x.m, x.s, x.ms)
		if got != x.want {
			t.Errorf("%d: %d %d %d.%d: got %v, want %v", i, x.h, x.m, x.s, x.ms, got, x.want)
		}
	}
}

func TestClockIsMidnight(t *testing.T) {
	cases := []struct {
		in   Clock
		want bool
	}{
		{New(0, 0, 0, 0), true},
		{ClockDay, true},
		{24 * ClockHour, true},
		{New(24, 0, 0, 0), true},
		{New(-24, 0, 0, 0), true},
		{New(-48, 0, 0, 0), true},
		{New(48, 0, 0, 0), true},
		{New(0, 0, 0, 1), false},
		{New(2, 0, 0, 1), false},
		{New(-1, 0, 0, 0), false},
		{New(0, 0, 0, -1), false},
	}
	for i, x := range cases {
		got := x.in.IsMidnight()
		if got != x.want {
			t.Errorf("%d: %v got %v, want %v, %d", i, x.in, x.in.IsMidnight(), x.want, x.in.Mod24())
		}
	}
}

func TestClockMod(t *testing.T) {
	cases := []struct {
		h, want Clock
	}{
		{0, 0},
		{1 * ClockHour, 1 * ClockHour},
		{2 * ClockHour, 2 * ClockHour},
		{23 * ClockHour, 23 * ClockHour},
		{24 * ClockHour, 0},
		{-24 * ClockHour, 0},
		{-48 * ClockHour, 0},
		{25 * ClockHour, ClockHour},
		{49 * ClockHour, ClockHour},
		{-1 * ClockHour, 23 * ClockHour},
		{-23 * ClockHour, ClockHour},
		{New(0, 0, 0, 1), 1},
		{New(0, 0, 1, 0), ClockSecond},
		{New(0, 0, 0, -1), New(23, 59, 59, 999)},
	}
	for i, x := range cases {
		clock := x.h
		got := clock.Mod24()
		if got != x.want {
			t.Errorf("%d: %dh: got %#v, want %#v", i, x.h, got, x.want)
		}
	}
}

func TestClockDays(t *testing.T) {
	cases := []struct {
		h, days int
	}{
		{0, 0},
		{1, 0},
		{23, 0},
		{24, 1},
		{25, 1},
		{48, 2},
		{49, 2},
		{-1, -1},
		{-23, -1},
		{-24, -2},
	}
	for i, x := range cases {
		clock := Clock(x.h) * ClockHour
		if clock.Days() != x.days {
			t.Errorf("%d: %dh: got %v, want %v", i, x.h, clock.Days(), x.days)
		}
	}
}

func TestClockString(t *testing.T) {
	cases := []struct {
		h, m, s, ms           Clock
		hh, hhmm, hhmmss, str string
	}{
		{0, 0, 0, 0, "00", "00:00", "00:00:00", "00:00:00.000"},
		{0, 0, 0, 1, "00", "00:00", "00:00:00", "00:00:00.001"},
		{0, 0, 1, 0, "00", "00:00", "00:00:01", "00:00:01.000"},
		{0, 1, 0, 0, "00", "00:01", "00:01:00", "00:01:00.000"},
		{1, 0, 0, 0, "01", "01:00", "01:00:00", "01:00:00.000"},
		{1, 2, 3, 4, "01", "01:02", "01:02:03", "01:02:03.004"},
		{-1, -1, -1, -1, "22", "22:58", "22:58:58", "22:58:58.999"},
	}
	for _, x := range cases {
		d := Clock(x.h*ClockHour + x.m*ClockMinute + x.s*ClockSecond + x.ms)
		if d.Hh() != x.hh {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.Hh(), x.hh, d)
		}
		if d.HhMm() != x.hhmm {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.HhMm(), x.hhmm, d)
		}
		if d.HhMmSs() != x.hhmmss {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.HhMmSs(), x.hhmmss, d)
		}
		if d.String() != x.str {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.String(), x.str, d)
		}
	}
}

func TestClockParse(t *testing.T) {
	cases := []struct {
		str  string
		want Clock
	}{
		{"00", New(0, 0, 0, 0)},
		{"01", New(1, 0, 0, 0)},
		{"23", New(23, 0, 0, 0)},
		{"00:00", New(0, 0, 0, 0)},
		{"00:01", New(0, 1, 0, 0)},
		{"01:00", New(1, 0, 0, 0)},
		{"01:02", New(1, 2, 0, 0)},
		{"23:59", New(23, 59, 0, 0)},
		{"00:00:00", New(0, 0, 0, 0)},
		{"00:00:01", New(0, 0, 1, 0)},
		{"00:01:00", New(0, 1, 0, 0)},
		{"01:00:00", New(1, 0, 0, 0)},
		{"01:02:03", New(1, 2, 3, 0)},
		{"23:59:59", New(23, 59, 59, 0)},
		{"00:00:00.000", New(0, 0, 0, 0)},
		{"00:00:00.001", New(0, 0, 0, 1)},
		{"00:00:01.000", New(0, 0, 1, 0)},
		{"00:01:00.000", New(0, 1, 0, 0)},
		{"01:00:00.000", New(1, 0, 0, 0)},
		{"01:02:03.004", New(1, 2, 3, 4)},
		{"23:59:59.999", New(23, 59, 59, 999)},
	}
	for _, x := range cases {
		str := MustParse(x.str)
		if str != x.want {
			t.Errorf("%s, got %v, want %v", x.str, str, x.want)
		}
	}
}
