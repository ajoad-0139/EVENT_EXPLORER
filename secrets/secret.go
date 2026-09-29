package secrets

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

// secrets model
type Secrets struct {
	GoogleKey       string `mapstructure:"GOOGLE_API_KEY"`
	TicketmasterKey string `mapstructure:"TICKETMASTER_API_KEY"`
}

// defining secrets variable
var secrets Secrets

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

	//check if any secrets is missing
	var missing []string
	if secrets.GoogleKey == "" {
		missing = append(missing, "GOOGLE_API_KEY")
	}
	if secrets.TicketmasterKey == "" {
		missing = append(missing, "TICKETMASTER_API_KEY")
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
