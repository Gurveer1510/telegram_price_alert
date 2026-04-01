package config

import "github.com/spf13/viper"

type Config struct {
	ZerodhaApiKey string
	DBHost        string
	DBName        string
	DBUser        string
	DBPass        string
	SSL           string
	ChanBind      string
	BotToken      string
	KiteUser      string
	KitePassword  string
	KiteSecret    string
	TotpSecret    string
}

func GetConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		ZerodhaApiKey: viper.GetString("ZERODHA_API_KEY"),
		DBHost:        viper.GetString("DATABASE_HOST"),
		DBName:        viper.GetString("DATABASE_NAME"),
		DBUser:        viper.GetString("DATABASE_USER"),
		DBPass:        viper.GetString("DATABASE_PASSWORD"),
		SSL:           viper.GetString("SSL"),
		ChanBind:      viper.GetString("CHANNEL_BINDING"),
		BotToken:      viper.GetString("BOT_TOKEN"),
		KiteUser:      viper.GetString("KITE_USER"),
		KitePassword:  viper.GetString("KITE_PASSWORD"),
		TotpSecret:    viper.GetString("TOTP_SECRET"),
		KiteSecret:    viper.GetString("KITE_SECRET"),
	}

	return config, nil
}
