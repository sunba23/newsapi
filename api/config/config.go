package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	ServerHost                string `mapstructure:"SERVER_HOST"`
	ServerReadTimeoutSeconds  int    `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerShutdownWaitSeconds int    `mapstructure:"SERVER_SHUTDOWN_WAIT"`

	LoggingPretty bool   `mapstructure:"LOGGING_PRETTY"`
	LoggingLevel  string `mapstructure:"LOGGING_LEVEL"`

	PostgresConnStr string `mapstructure:"POSTGRES_CONN_STR" validate:"required"`

	GoogleOauthRedirectUrl  string   `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`
	GoogleOauthClientId     string   `mapstructure:"GOOGLE_OAUTH_CLIENT_ID" validate:"required"`
	GoogleOauthClientSecret string   `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET" validate:"required"`
	GoogleOauthScopes       []string `mapstructure:"GOOGLE_OAUTH_SCOPES"`

	SessionSecret string `mapstructure:"SESSION_SECRET" validate:"required"`
}

func InitConfig(dotEnvFilenames ...string) (*Config, error) {
	if err := godotenv.Load(dotEnvFilenames...); err != nil {
		log.Error().Err(err).Send()
	}

	return NewConfig()
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	setDefaults()
	bindEnvVars(cfg)
	err := viper.Unmarshal(
		cfg,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
				mapstructure.TextUnmarshallerHookFunc(),
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall the config: %v", err)
	}

	if err = validateConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func setDefaults() {
	var defaults = map[string]any{
		"SERVER_HOST":               "0.0.0.0:8000",
		"SERVER_READ_TIMEOUT":       15,
		"SERVER_SHUTDOWN_WAIT":      3,
		"LOGGING_PRETTY":            true,
		"LOGGING_LEVEL":             "debug",
		"GOOGLE_OAUTH_REDIRECT_URL": "http://localhost:8000/auth/google/callback",
		"GOOGLE_OAUTH_SCOPES":       []string{"https://www.googleapis.com/auth/userinfo.email"},
	}

	for key, value := range defaults {
		if val, ok := value.(map[string]any); ok {
			for subKey, subValue := range val {
				viper.SetDefault(fmt.Sprintf("%v.%v", key, subKey), subValue)
			}
		} else {
			viper.SetDefault(key, value)
		}
	}
}

func bindEnvVars(cfg *Config) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfgType := reflect.TypeOf(cfg).Elem()
	for i := range cfgType.NumField() {
		field := cfgType.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			continue
		}
		if err := viper.BindEnv(tag); err != nil {
			log.Fatal().Str("env_var", tag).Err(err).Msg("Failed to bind env var")
		}
	}
}

func validateConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("config validation errors:\n%v", err)
	}
	return nil
}
