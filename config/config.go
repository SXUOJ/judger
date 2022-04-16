package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	ShowDetails bool
	UnSafe      bool
)

type Conf struct {
	Web struct {
		Dev      bool   `yaml:"dev"`
		LogLevel string `yaml:"log_level"`
		Listen   string `yaml:"listen"`
	} `yaml:"web"`
}

func LoadConf() *Conf {
	var (
		conf = new(Conf)
	)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := setupLogDir(); err != nil {
		logrus.Fatalln(err)
	}

	if err := setupLogOutput(); err != nil {
		logrus.Fatalln(err)
	}

	yamlFileBytes, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		logrus.Fatalln(err)
	}

	if err = yaml.Unmarshal(yamlFileBytes, conf); err != nil {
		logrus.Fatalln(err)
	}

	err = setupLogLevel(conf.Web.LogLevel)
	if err != nil {
		logrus.Fatalln(err)
	}

	err = setupGinLog()
	if err != nil {
		logrus.Fatalln(err)
	}

	if conf.Web.Dev == false {
		gin.SetMode(gin.ReleaseMode)
	} else {
		ShowDetails = true
		UnSafe = true
	}
	return conf
}

func setupLogDir() error {
	var err error
	if _, err = os.Stat("./logs/"); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir("./logs/", 0755)
	}
	return err
}

func setupLogOutput() error {
	var err error

	logFileName := time.Now().Format("2006-01-02")
	logFile, err := os.OpenFile("./logs/"+logFileName+"-app-all.log", syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0755)
	if err != nil {
		return err
	}
	logOut := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(logOut)

	return err
}

func setupLogLevel(pLevel string) error {
	switch pLevel {
	case "":
		logrus.SetLevel(logrus.DebugLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		return errors.New("unknown log level: " + pLevel)
	}
	return nil
}

func setupGinLog() error {
	var err error
	logErrorFileName := time.Now().Format("2006-01-02")
	logErrorFile, err := os.OpenFile("./logs/"+logErrorFileName+"-gin-error.log", syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
	if err != nil {
		return err
	}
	gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, logErrorFile)
	logInfoFileName := time.Now().Format("2006-01-02")
	logInfoFile, err := os.OpenFile("./logs/"+logInfoFileName+"-gin-info.log", syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
	if err != nil {
		return err
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logInfoFile)
	return err
}
