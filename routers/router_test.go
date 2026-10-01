package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestRoutes(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		target     string
		wantStatus int
	}{
		{
			name:       "clear all caches",
			method:     http.MethodDelete,
			target:     "/api/cache/all",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalidate list without params",
			method:     http.MethodDelete,
			target:     "/api/cache/events",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "autocomplete with short input",
			method:     http.MethodGet,
			target:     "/api/locations/autocomplete?input=a",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.target, nil)

			beego.BeeApp.Handlers.ServeHTTP(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					w.Code,
					tc.wantStatus,
				)
			}
		})
	}
}
