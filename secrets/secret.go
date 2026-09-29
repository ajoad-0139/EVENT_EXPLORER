package secrets

import (
	"errors"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/spf13/viper"
)

// secrets model
type Secrets struct {
	GoogleKey       string `mapstructure:"GOOGLE_API_KEY"`
	TicketmasterKey string `mapstructure:"TICKETMASTER_API_KEY"`
}

// base urls
type BaseUrls struct {
	GoogleBaseURL       string
	TicketmasterBaseURL string
}

// defining secrets and baseUrl variables
var secrets Secrets
var baseUrl BaseUrls

func MustLoad() error {
	// creating new viper instance
	v := viper.New()

	//introducing the secret file
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	_ = v.ReadInConfig()

	// bind the secret keys
	for _, key := range []string{"GOOGLE_API_KEY", "TICKETMASTER_API_KEY"} {
		err := v.BindEnv(key)
		if err != nil {
			return err
		}
	}

	// loading secrets to memory
	err := v.Unmarshal(&secrets)
	if err != nil {
		return err
	}

	//whitespace trimming
	secrets.GoogleKey = strings.TrimSpace(secrets.GoogleKey)
	secrets.TicketmasterKey = strings.TrimSpace(secrets.TicketmasterKey)

	//loading base urls to memory
	googleUrl, err := beego.AppConfig.String("googlebaseurl")
	if err != nil {
		return err
	}
	baseUrl.GoogleBaseURL = googleUrl

	ticketMasterUrl, err := beego.AppConfig.String("ticketmasterbaseurl")
	if err != nil {
		return err
	}
	baseUrl.TicketmasterBaseURL = ticketMasterUrl

	//check if any secrets and urls are missing
	var missing []string
	if secrets.GoogleKey == "" {
		missing = append(missing, "GOOGLE_API_KEY")
	}
	if secrets.TicketmasterKey == "" {
		missing = append(missing, "TICKETMASTER_API_KEY")
	}
	if baseUrl.GoogleBaseURL == "" {
		missing = append(missing, "GOOGLE_BASE_URL (app.conf)")
	}
	if baseUrl.TicketmasterBaseURL == "" {
		missing = append(missing, "TICKET_MASTER_BASE_URL (app.conf)")
	}
	if len(missing) > 0 {
		return errors.New("missing env vars: " + strings.Join(missing, ", "))
	}
	return nil
}

// access secrets through methods only
func GetSecrets() Secrets {
	return secrets
}

func GetBaseUrls() BaseUrls {
	return baseUrl
}
