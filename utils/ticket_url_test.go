package utils

import "testing"

func TestIsApprovedTicketURL(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "approved ticketmaster url",
			url:  "https://ticketmaster.com/event/123",
			want: true,
		},
		{
			name: "approved ticketmaster canada url",
			url:  "https://ticketmaster.ca/event/123",
			want: true,
		},
		{
			name: "approved ticketmaster uk url",
			url:  "https://ticketmaster.co.uk/event/123",
			want: true,
		},
		{
			name: "approved livenation url",
			url:  "https://livenation.com/event/123",
			want: true,
		},
		{
			name: "approved ticketweb url",
			url:  "https://ticketweb.com/event/123",
			want: true,
		},
		{
			name: "approved subdomain",
			url:  "https://www.ticketmaster.com/event/123",
			want: true,
		},
		{
			name: "approved nested subdomain",
			url:  "https://tickets.events.ticketmaster.com/event/123",
			want: true,
		},
		{
			name: "http is not allowed",
			url:  "http://ticketmaster.com/event/123",
			want: false,
		},
		{
			name: "unapproved domain",
			url:  "https://example.com/event/123",
			want: false,
		},
		{
			name: "lookalike domain",
			url:  "https://ticketmaster.com.example.com/event/123",
			want: false,
		},
		{
			name: "lookalike domain with suffix",
			url:  "https://faketicketmaster.com/event/123",
			want: false,
		},
		{
			name: "empty url",
			url:  "",
			want: false,
		},
		{
			name: "invalid url",
			url:  "not-a-url",
			want: false,
		},
		{
			name: "ftp url",
			url:  "ftp://ticketmaster.com/event/123",
			want: false,
		},
		{
			name: "url with username",
			url:  "https://user@ticketmaster.com/event/123",
			want: false,
		},
		{
			name: "url with username and password",
			url:  "https://user:password@ticketmaster.com/event/123",
			want: false,
		},
		{
			name: "uppercase approved domain",
			url:  "https://TICKETMASTER.COM/event/123",
			want: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := IsApprovedTicketURL(tc.url)

			if res != tc.want {
				t.Errorf("approved-ticket-url = %v, wanted-approved-ticket-url %v", res, tc.want)
			}
		})
	}
}
