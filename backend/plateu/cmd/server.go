package cmd

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fidellr/jastip_way/backend/plateu"
	"github.com/fidellr/jastip_way/backend/plateu/image"
	"github.com/fidellr/jastip_way/backend/plateu/utils"

	_httpDelivery "github.com/fidellr/jastip_way/backend/plateu/internal/delivery/http"
	_mongoRepository "github.com/fidellr/jastip_way/backend/plateu/internal/delivery/mongo"
)

var plateuServerCMD = &cobra.Command{
	Use:   "http",
	Short: "Start http server for plateu",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		echoInstance := echo.New()
		echoInstance.Server.ReadTimeout = time.Duration(viper.GetInt("http.server_read_timeout")) * time.Second
		echoInstance.Server.WriteTimeout = time.Duration(viper.GetInt("http.server_write_timeout")) * time.Second

		echoInstance.GET("/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong plateu")
		})

		initPlateuApplication(echoInstance)

		address := viper.GetString("server.address")
		if err := echoInstance.Start(address); err != nil {
			logrus.Fatalln(err.Error())
		}

		logrus.Infof("Start listening on : %s", address)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCMD.AddCommand(plateuServerCMD)
	plateuServerCMD.PersistentFlags().String("config", "", "Set this flag to use a configuration file")
}

func initConfig() {
	configFile := ""
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if plateuServerCMD.Flags().Lookup("config") != nil {
		configFile = "config.json"
		viper.BindPFlag("config", plateuServerCMD.Flags().Lookup("config"))
		viper.SetConfigType("json")
	}

	if config := viper.GetString("config"); config != "" {
		configFile = config
	}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalln(err.Error())
	}

	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Warn("Plateu is running in Debug Mode")
	}
}

func initPlateuApplication(e *echo.Echo) {
	contextTimeout := time.Duration(viper.GetInt("context.timeout")) * time.Second
	mongoDSN := viper.GetString("mongo.dsn")
	masterSession, err := mgo.Dial(mongoDSN)
	if err != nil {
		logrus.Fatalln(err.Error())
	}

	masterSession.SetSafe(&mgo.Safe{})

	mongoDatabase := viper.GetString("mongo.database")
	if mongoDatabase == "" {
		logrus.Fatalln(errors.New("Please provide a mongo database name"))
	}

	validator := plateu.NewValidator()
	imageRepo := _mongoRepository.NewImageMongo(
		_mongoRepository.ImageSession(masterSession),
		_mongoRepository.ImageDBName(mongoDatabase),
	)
	imageService := image.NewService(
		image.Repository(imageRepo),
		image.Timeout(contextTimeout),
		image.Validator(validator),
	)

	e.HTTPErrorHandler = utils.HandleUncaughtHTTPError
	_httpDelivery.NewImageHandler(e, _httpDelivery.ImageService(imageService))
}
