// This package load API configuration
package config

import (
	"encoding/json"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// Config class contains general project configuration
type Config struct {
	Env          string `validate:"required,oneof='localhost' 'dev' 'rct' 'prod' 'develop' 'recette' 'production'"`
	Level        string `validate:"required,oneof='trace' 'debug' 'info' 'warning' 'warn' 'error' 'fatal' 'panic'"`
	rootPath     string `validate:"required"`
	ReportCaller bool
	Port         int `validate:"gte=1,lte=65535"`
	Database     struct {
		Adapter         string        `json:"adapter" validate:"required"`
		Host            string        `json:"host" validate:"required"`
		Username        string        `json:"username" validate:"required"`
		Password        string        `json:"password" validate:"required"`
		Port            string        `json:"port" validate:"gte=1,lte=65535"`
		Name            string        `json:"name" validate:"required"`
		Charset         string        `json:"charset" validate:"required"`
		MaxOpenConns    int           `json:"maxOpenConns" validate:"gte=1"`
		MaxIdleConns    int           `json:"maxIdleConns" validate:"gte=1"`
		MaxConnLifeTime time.Duration `json:"maxConnLifeTime" validate:"gte=1"`
	} `json:"database"`
}

// UnmarshalJSON : lower case on string fields level and Env, to avoid case sensitive
func (c *Config) UnmarshalJSON(data []byte) error {

	type TmpConfig Config //avoid infinite loop stack exceed.

	var tmpConf TmpConfig //init struct

	err := json.Unmarshal(data, &tmpConf)
	if err != nil {
		return err
	}

	*c = Config(tmpConf)

	(*c).Env = strings.ToLower((*c).Env)
	(*c).Level = strings.ToLower((*c).Level)

	return nil
}

// GetRootPath return workspace path
func (c *Config) GetRootPath() string {

	return c.rootPath
}

// InitRootPath init workspace path
func (c *Config) InitRootPath(val string) {
	if c.rootPath == "" {
		c.rootPath = val
	}
}

// New constructor Config
func New(fileName string, validator validator.Validate) (c *Config, err error) {

	c = &Config{}

	// Define logrus text formatter
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",

		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// s := strings.Split(f.Function, ".")
			// funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return "", " [" + filename + ":" + strconv.Itoa(f.Line) + "]"
		},
		//PadLevelText:    true,
	})

	configFile, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(c)
	if err != nil {
		return
	}

	// Set Rootpath
	rootPath, err := os.Getwd()
	if err != nil {
		return
	}
	c.InitRootPath(rootPath)

	// Set log error Level
	level, err := logrus.ParseLevel(c.Level)
	if err == nil {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Set reportCaller logrus
	logrus.SetReportCaller(c.ReportCaller)

	// Set Gin mode
	if c.Env != "localhost" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
		var f *os.File
		f, err := os.OpenFile("gin.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			logrus.Fatal(err)
		}
		gin.DefaultWriter = f
	}

	//validate fields from config
	if err = validator.Struct(c); err != nil {
		return
	}

	return
}
