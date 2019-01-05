package configs

import (
	"flag"

	"github.com/sirupsen/logrus"
)

// EnvConfig Set default config values from Env variables
type EnvConfig struct {
	Port             int    `env:"PORT" envDefault:"80"`
	LogLevel         int    `env:"LOG_LEVEL" envDefault:"5"`
	LogPlain         bool   `env:"LOG_PLAIN" envDefault:"false"`
	Environment      string `env:"ENVIRONMENT" envDefault:"development"`
	AuthToken        string `env:"AUTH_TOKEN" envDefault:"35216c9e-dea4-458c-babd-325f2ef0eefb"` // UUID v4 token
	AESSecret        string `env:"AES_SECRET" envDefault:"bQeThWmZq4t7w!z%C&F)J@NcRfUjXn2r"`     // 256-bit
	JWTSecret        string `env:"JWT_SECRET" envDefault:"SpjDimBfySs24H5QOErfH95XzN2sXmzVcrLigggWLJA"`
	JWTExpirySeconds int    `env:"JWT_EXPIRY_SECONDS" envDefault:"3600"`
	CompanyAHost     string `env:"COMPANY_A_HOST" envDefault:"http://localhost"`
	CompanyAPort     int    `env:"COMPANY_A_PORT" envDefault:"80"`
	CompanyBHost     string `env:"COMPANY_B_HOST" envDefault:"http://localhost"`
	CompanyBPort     int    `env:"COMPANY_B_PORT" envDefault:"80"`
}

var (
	// Port ...
	Port int
	// Environment ...
	Environment string
	// AuthToken ...
	AuthToken string
	// AESSecret ...
	AESSecret string
	// JWTSecret ...
	JWTSecret string
	// JWTExpirySeconds ...
	JWTExpirySeconds int
	// CompanyAHost ...
	CompanyAHost string
	// CompanyAPort ...
	CompanyAPort int
	// CompanyBHost ...
	CompanyBHost string
	// CompanyBPort ...
	CompanyBPort int
)

// ParseInitVars ...
func ParseInitVars(config EnvConfig) (logLevel int, logPlain bool) {

	flag.IntVar(&Port, "port", config.LogLevel, "Set service port")
	flag.IntVar(&logLevel, "log-level", config.LogLevel, "Sets the logrus log level")
	flag.BoolVar(&logPlain, "log-plain", config.LogPlain, "Tells logrus not to log as json")
	flag.StringVar(&Environment, "environment", config.Environment, "Environment the service is running")
	flag.StringVar(&AuthToken, "auth-token", config.AuthToken, "Auth Token")
	flag.StringVar(&AESSecret, "aes-secret", config.AESSecret, "AES Secret")
	flag.StringVar(&JWTSecret, "jwt_secret", config.JWTSecret, "JWT Secret")
	flag.IntVar(&JWTExpirySeconds, "jwt-expiry-seconds", config.JWTExpirySeconds, "JWT expiry time in seconds")
	flag.StringVar(&CompanyAHost, "company-a-host", config.CompanyAHost, "Company A Host")
	flag.IntVar(&CompanyAPort, "company-a-port", config.CompanyAPort, "Company A Port")
	flag.StringVar(&CompanyBHost, "company-b-host", config.CompanyBHost, "Company B Host")
	flag.IntVar(&CompanyBPort, "company-b-port", config.CompanyBPort, "Company B Port")

	// Make flag parse flags based on those setup above
	flag.Parse()

	if lowest := int(logrus.PanicLevel); logLevel < lowest {
		logLevel = lowest
	}
	if greatest := int(logrus.DebugLevel); logLevel > greatest {
		logLevel = greatest
	}
	return
}
