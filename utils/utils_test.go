package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"event-explorer/models"

	beego "github.com/beego/beego/v2/server/web"
	beegoctx "github.com/beego/beego/v2/server/web/context"
)

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

func newTestController() (*beego.Controller, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	ctx := beegoctx.NewContext()
	ctx.Reset(w, r)

	c := &beego.Controller{}
	c.Init(ctx, "TestController", "Test", nil)
	c.EnableRender = false // no .tpl files needed in unit tests

	return c, w
}

func TestJsonError(t *testing.T) {
	c, w := newTestController()

	JsonError(c, http.StatusBadRequest, "something failed")

	if w.Code != http.StatusBadRequest {
		t.Errorf(
			"response-status = %d, wanted-response-status %d",
			w.Code,
			http.StatusBadRequest,
		)
	}
}

func TestJsonSuccess(t *testing.T) {
	c, w := newTestController()

	JsonSuccess(c, http.StatusCreated, map[string]string{"hello": "world"})

	if w.Code != http.StatusCreated {
		t.Errorf(
			"response-status = %d, wanted-response-status %d",
			w.Code,
			http.StatusCreated,
		)
	}
}

func TestRenderErrorMsg(t *testing.T) {
	c, _ := newTestController()

	RenderErrorMsg(c, http.StatusTeapot, "Title", "Message")

	if c.Ctx.Output.Status != http.StatusTeapot {
		t.Errorf(
			"response-status = %d, wanted-response-status %d",
			c.Ctx.Output.Status,
			http.StatusTeapot,
		)
	}

	if c.TplName != "error.tpl" {
		t.Errorf("tpl-name = %s, wanted-tpl-name error.tpl", c.TplName)
	}

	page, ok := c.Data["error"].(models.ErrorPage)
	if !ok {
		t.Fatalf("error-page has type %T, wanted models.ErrorPage", c.Data["error"])
	}

	if page.Title != "Title" || page.Message != "Message" {
		t.Errorf("error-page = %+v, wanted title Title and message Message", page)
	}
}

func TestRenderError(t *testing.T) {
	cases := []struct {
		name       string
		status     int
		wantStatus int
		wantTitle  string
	}{
		{
			name:       "bad request",
			status:     http.StatusBadRequest,
			wantStatus: http.StatusBadRequest,
			wantTitle:  "Something is missing",
		},
		{
			name:       "not found",
			status:     http.StatusNotFound,
			wantStatus: http.StatusNotFound,
			wantTitle:  "Event not found",
		},
		{
			name:       "too many requests",
			status:     http.StatusTooManyRequests,
			wantStatus: http.StatusTooManyRequests,
			wantTitle:  "Too many requests",
		},
		{
			name:       "gateway timeout",
			status:     http.StatusGatewayTimeout,
			wantStatus: http.StatusGatewayTimeout,
			wantTitle:  "Taking too long",
		},
		{
			name:       "request timeout becomes gateway timeout",
			status:     http.StatusRequestTimeout,
			wantStatus: http.StatusGatewayTimeout,
			wantTitle:  "Taking too long",
		},
		{
			name:       "other errors become bad gateway",
			status:     http.StatusInternalServerError,
			wantStatus: http.StatusBadGateway,
			wantTitle:  "Something went wrong",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := newTestController()

			RenderError(c, tc.status, "upstream failure")

			if c.Ctx.Output.Status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					c.Ctx.Output.Status,
					tc.wantStatus,
				)
			}

			page, ok := c.Data["error"].(models.ErrorPage)
			if !ok {
				t.Fatalf("error-page has type %T, wanted models.ErrorPage", c.Data["error"])
			}

			if page.Title != tc.wantTitle {
				t.Errorf(
					"error-title = %s, wanted-error-title %s",
					page.Title,
					tc.wantTitle,
				)
			}
		})
	}
}

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
