package secrets

import (
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestMustLoad(t *testing.T) {
	cases := []struct {
		name          string
		googleKey     string
		ticketKey     string
		googleURL     string
		ticketURL     string
		wantGoogleKey string
		wantErr       string
	}{
		{
			name:          "all values present",
			googleKey:     "  google-key  ",
			ticketKey:     "ticket-key",
			googleURL:     "http://google/",
			ticketURL:     "http://ticket/",
			wantGoogleKey: "google-key",
			wantErr:       "",
		},
		{
			name:      "missing google key",
			ticketKey: "ticket-key",
			googleURL: "http://google/",
			ticketURL: "http://ticket/",
			wantErr:   "missing env vars: GOOGLE_API_KEY",
		},
		{
			name:      "missing ticketmaster key",
			googleKey: "google-key",
			googleURL: "http://google/",
			ticketURL: "http://ticket/",
			wantErr:   "missing env vars: TICKETMASTER_API_KEY",
		},
		{
			name:      "missing google url",
			googleKey: "google-key",
			ticketKey: "ticket-key",
			ticketURL: "http://ticket/",
			wantErr:   "missing env vars: GOOGLE_BASE_URL (app.conf)",
		},
		{
			name:      "missing ticketmaster url",
			googleKey: "google-key",
			ticketKey: "ticket-key",
			googleURL: "http://google/",
			wantErr:   "missing env vars: TICKET_MASTER_BASE_URL (app.conf)",
		},
		{
			name:    "everything missing",
			wantErr: "missing env vars: GOOGLE_API_KEY, TICKETMASTER_API_KEY, GOOGLE_BASE_URL (app.conf), TICKET_MASTER_BASE_URL (app.conf)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// MustLoad does not overwrite missing keys, so start every case clean
			secrets = Secrets{}
			baseUrl = BaseUrls{}

			t.Setenv("GOOGLE_API_KEY", tc.googleKey)
			t.Setenv("TICKETMASTER_API_KEY", tc.ticketKey)
			_ = beego.AppConfig.Set("googlebaseurl", tc.googleURL)
			_ = beego.AppConfig.Set("ticketmasterbaseurl", tc.ticketURL)

			err := MustLoad()

			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != tc.wantErr {
				t.Errorf(
					"error-message = %q, wanted-error-message %q",
					errMsg,
					tc.wantErr,
				)
			}

			if tc.wantErr != "" {
				return
			}

			if GetSecrets().GoogleKey != tc.wantGoogleKey {
				t.Errorf(
					"google-key = %q, wanted-google-key %q",
					GetSecrets().GoogleKey,
					tc.wantGoogleKey,
				)
			}

			if GetBaseUrls().GoogleBaseURL != tc.googleURL {
				t.Errorf(
					"google-base-url = %q, wanted-google-base-url %q",
					GetBaseUrls().GoogleBaseURL,
					tc.googleURL,
				)
			}
		})
	}
}
