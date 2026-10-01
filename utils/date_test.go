package utils

import "testing"

func TestFormatEventDate(t *testing.T) {
	cases := []struct {
		name string
		date string
		want string
	}{
		{
			name: "valid date",
			date: "2026-10-01",
			want: "Thu, 01 Oct 2026",
		},
		{
			name: "another valid date",
			date: "2026-01-15",
			want: "Thu, 15 Jan 2026",
		},
		{
			name: "invalid date",
			date: "2026/10/01",
			want: "2026/10/01",
		},
		{
			name: "empty date",
			date: "",
			want: "",
		},
		{
			name: "invalid text",
			date: "not-a-date",
			want: "not-a-date",
		},
		{
			name: "invalid day",
			date: "2026-10-32",
			want: "2026-10-32",
		},
		{
			name: "invalid month",
			date: "2026-13-01",
			want: "2026-13-01",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := FormatEventDate(tc.date)

			if res != tc.want {
				t.Errorf("formatted-event-date = %s, wanted-event-date %s", res, tc.want)
			}
		})
	}
}

func TestFormatEventTime(t *testing.T) {
	cases := []struct {
		name string
		time string
		want string
	}{
		{
			name: "morning time",
			time: "09:30:00",
			want: "9:30 AM",
		},
		{
			name: "afternoon time",
			time: "15:45:00",
			want: "3:45 PM",
		},
		{
			name: "midnight",
			time: "00:00:00",
			want: "12:00 AM",
		},
		{
			name: "noon",
			time: "12:00:00",
			want: "12:00 PM",
		},
		{
			name: "evening time",
			time: "20:15:30",
			want: "8:15 PM",
		},
		{
			name: "invalid time",
			time: "20:15",
			want: "Time to be announced",
		},
		{
			name: "empty time",
			time: "",
			want: "Time to be announced",
		},
		{
			name: "invalid text",
			time: "not-a-time",
			want: "Time to be announced",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := FormatEventTime(tc.time)

			if res != tc.want {
				t.Errorf("formatted-event-time = %s, wanted-event-time %s", res, tc.want)
			}
		})
	}
}
