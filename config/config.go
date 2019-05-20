package config

import (
	"os"

	"github.com/jinzhu/configor"
	amazonpay "github.com/qor/amazon-pay-sdk-go"
	"github.com/qor/auth/providers/facebook"
	"github.com/qor/auth/providers/github"
	"github.com/qor/auth/providers/google"
	"github.com/qor/auth/providers/twitter"
	"github.com/qor/gomerchant"
	"github.com/qor/location"
	"github.com/qor/mailer"
	"github.com/qor/mailer/logger"
	"github.com/qor/media/oss"
	"github.com/qor/oss/qiniu"
	"github.com/qor/redirect_back"
	"github.com/qor/session/manager"
	"github.com/unrolled/render"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

var Config = struct {
	HTTPS bool `default:"false" env:"HTTPS"`
	Port  uint `default:"7000" env:"PORT"`
	DB    struct {
		Name     string `env:"DBName" default:"qor_example"`
		Adapter  string `env:"DBAdapter" default:"mysql"`
		Host     string `env:"DBHost" default:"localhost"`
		Port     string `env:"DBPort" default:"3306"`
		User     string `env:"DBUser"`
		Password string `env:"DBPassword"`
	}
	S3 struct {
		AccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
		Region          string `env:"AWS_REGION"`
		S3Bucket        string `env:"AWS_BUCKET"`
	}

	Qiniu struct {
		AccessID  string `env:"QINIU_ACCESS_ID"`
		AccessKey string `env:"QINIU_ACCESS_KEY"`
		Bucket    string `env:"QINIU_BUCKET"`
		Region    string `env:"QINIU_REGION"`
		Endpoint  string `env:"QINIU_ENDPOINT"`
	}

	AmazonPay struct {
		MerchantID   string `env:"AmazonPayMerchantID"`
		AccessKey    string `env:"AmazonPayAccessKey"`
		SecretKey    string `env:"AmazonPaySecretKey"`
		ClientID     string `env:"AmazonPayClientID"`
		ClientSecret string `env:"AmazonPayClientSecret"`
		Sandbox      bool   `env:"AmazonPaySandbox"`
		CurrencyCode string `env:"AmazonPayCurrencyCode" default:"JPY"`
	}
	SMTP         SMTPConfig
	Github       github.Config
	Google       google.Config
	Facebook     facebook.Config
	Twitter      twitter.Config
	GoogleAPIKey string `env:"GoogleAPIKey"`
	BaiduAPIKey  string `env:"BaiduAPIKey"`
}{}

var (
	Root           = os.Getenv("GOPATH") + "/src/github.com/dfang/qor-example"
	Mailer         *mailer.Mailer
	Render         = render.New()
	AmazonPay      amazonpay.AmazonPayService
	PaymentGateway gomerchant.PaymentGateway
	RedirectBack   = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
)

func init() {
	if err := configor.Load(&Config, "config/database.yml", "config/smtp.yml", "config/application.yml"); err != nil {
		panic(err)
	}

	location.GoogleAPIKey = Config.GoogleAPIKey
	location.BaiduAPIKey = Config.BaiduAPIKey

	// if Config.S3.AccessKeyID != "" {
	// 	oss.Storage = s3.New(&s3.Config{
	// 		AccessID:  Config.S3.AccessKeyID,
	// 		AccessKey: Config.S3.SecretAccessKey,
	// 		Region:    Config.S3.Region,
	// 		Bucket:    Config.S3.S3Bucket,
	// 	})
	// }

	// oss.Storage = qiniu.New(&qiniu.Config{
	// 	AccessID:  "MkFws9gjO_CScK5pXrahfBEWf9viOD_khTomtL3f",
	// 	AccessKey: "xVGWVTQTKFAlEEOFj6t4RRasJek5995UPlcMvv3M",
	// 	Bucket:    "zdtech",
	// 	Region:    "huadong",
	// 	Endpoint:  "https://up.qiniup.com",
	// })
	oss.Storage = qiniu.New(&qiniu.Config{
		AccessID:  Config.Qiniu.AccessID,
		AccessKey: Config.Qiniu.AccessKey,
		Bucket:    Config.Qiniu.Bucket,
		Region:    Config.Qiniu.Region,
		// Endpoint:  Config.Qiniu.Endpoint,
		Endpoint: "https://up.qiniup.com",
	})

	AmazonPay = amazonpay.New(&amazonpay.Config{
		MerchantID: Config.AmazonPay.MerchantID,
		AccessKey:  Config.AmazonPay.AccessKey,
		SecretKey:  Config.AmazonPay.SecretKey,
		Sandbox:    true,
		Region:     "jp",
	})

	// dialer := gomail.NewDialer(Config.SMTP.Host, Config.SMTP.Port, Config.SMTP.User, Config.SMTP.Password)
	// sender, err := dialer.Dial()

	// Mailer = mailer.New(&mailer.Config{
	// 	Sender: gomailer.New(&gomailer.Config{Sender: sender}),
	// })
	Mailer = mailer.New(&mailer.Config{
		Sender: logger.New(&logger.Config{}),
	})
}
