package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	AdminSessionKey      string `yaml:"AdminSessionKey"`
	TempUploadPath       string `yaml:"TempUploadPath"`
	ImageUploadPath      string `yaml:"ImageUploadPath"` 
	CookieSecretKey      string `yaml:"CookieSecretKey"`
	CookieSessionName    string `yaml:"CookieSessionName"`
	FillFormMessage      string `yaml:"FillFormMessage"`
	AccessDeneidMessage  string `yaml:"AccessDeneidMessage"`
	SessionErrorMessage  string `yaml:"SessionErrorMessage"`
	GettingFileError     string `yaml:"GettingFileError"`
	DBAccess             string `yaml:"DBAccess"`
	ThumbPostfix         string `yaml:"ThumbPostfix"`
	TelegramURLTemplate  string `yaml:"TelegramURLTemplate"`
	TelegramBotTokken    string `yaml:"TelegramBotTokken"`
	TelegramChatID       string `yaml:"TelegramChatID"`
	OrderMessageTemplate string `yaml:"OrderMessageTemplate"`
	ThumbWidth           int    `yaml:"ThumbWidth"`
}

func (config *Config) fillFromYML() error {
	content, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	return nil
}