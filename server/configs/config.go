package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)
type Config struct{
	App AppConfig
	Database DatabaseConfig
	JWT JWTConfig
	R2 R2Config
}

type AppConfig struct{
	Env string `mapstructure:"APP_ENV"`
	Port string `mapstructure:"APP_PORT"`
}

type DatabaseConfig struct{
	DB_USER string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_NAME string `mapstructure:"DB_NAME"`
	DB_PORT string `mapstructure:"DB_PORT"`
	DB_url string `mapstructure:"DB_URL"`

}

type JWTConfig struct{
	Secret string `mapstructure:"JWT_SECRET"`
	AccessTokenExpire time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRE"`
	RefreshTokenExpire time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRE"`
}

type R2Config struct {
	AccountID       string `mapstructure:"R2_ACCOUNT_ID"`
	AccessKeyID     string `mapstructure:"R2_ACCESS_KEY_ID"`
	SecretAccessKey string `mapstructure:"R2_ACCESS_KEY_SECRET"`
	BucketName      string `mapstructure:"R2_BUCKET_NAME"`
	PublicURL       string `mapstructure:"R2_PUBLIC_URL"`
}

// LoadConfig khởi tạo cấu hình từ các biến môi trường

func LoadConfig()(*Config, error){
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")
	if err := viper.ReadInConfig(); err != nil{
		log.Println("Không tìm thấy file .env, đang đọc trực tiếp từ Environment Variables của hệ thống")
	}
	var config Config
	if err := viper.Unmarshal(&config.App); err != nil{
		return nil, err
	}
	if err := viper.Unmarshal(&config.Database); err != nil{
		return nil, err
	}

	if err := viper.Unmarshal(&config.JWT); err != nil{
		return nil, err
	}
	if err := viper.Unmarshal(&config.R2); err != nil{
		return nil, err
	}
	return &config, nil
}