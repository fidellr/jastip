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

	"github.com/fidellr/jastip_way/backend/rover"
	"github.com/fidellr/jastip_way/backend/rover/content"
	delivery "github.com/fidellr/jastip_way/backend/rover/internal/delivery"
	_httpDelivery "github.com/fidellr/jastip_way/backend/rover/internal/delivery/http"
	_mongoRepository "github.com/fidellr/jastip_way/backend/rover/internal/delivery/mongo"
)

var roverServerCMD = &cobra.Command{
	Use:   "http",
	Short: "Start http server for rover",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		echoInstance := echo.New()
		echoInstance.Server.ReadTimeout = time.Duration(viper.GetInt("http.server_read_timeout")) * time.Second
		echoInstance.Server.WriteTimeout = time.Duration(viper.GetInt("http.server_write_timeout")) * time.Second

		echoInstance.GET("/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong content")
		})

		initRoverApplication(echoInstance)

		address := viper.GetString("server.address")
		if err := echoInstance.Start(address); err != nil {
			logrus.Fatalln(err.Error())
		}

		logrus.Infof("Start listening on: %s", address)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCMD.AddCommand(roverServerCMD)
	roverServerCMD.PersistentFlags().String("config", "", "Set this flag to use a configuration file")
}

func initConfig() {
	configFile := ""
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if roverServerCMD.Flags().Lookup("config") != nil {
		configFile = "config.json"
		viper.BindPFlag("config", roverServerCMD.Flags().Lookup("config"))
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
		logrus.Warn("Rover is Running in Debug Mode")
	}
}

func initRoverApplication(e *echo.Echo) {
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

	validator := rover.NewValidator()
	contentRepo := _mongoRepository.NewContentMongo(
		_mongoRepository.ContentSession(masterSession),
		_mongoRepository.ContentDBName(mongoDatabase),
	)
	roverService := content.NewService(
		content.Repository(contentRepo),
		content.Timeout(contextTimeout),
		content.Validator(validator),
	)
	e.HTTPErrorHandler = delivery.HandleUncaughtHTTPError
	_httpDelivery.NewContentHandler(e, _httpDelivery.ContentService(roverService))
}
