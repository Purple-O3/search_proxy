package viperwrapper

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func DecodeConfig(configPath string, config any) error {
	fileName := filepath.Base(configPath)
	fileNames := strings.Split(fileName, ".")
	if len(fileNames) != 2 {
		return errors.New("fileNames len not equal 2")
	}

	vp := viper.New()
	vp.AddConfigPath(filepath.Dir(configPath))
	vp.SetConfigName(fileNames[0])
	vp.SetConfigType(fileNames[1])
	err := vp.ReadInConfig()
	if err != nil {
		return err
	}
	err = vp.Unmarshal(config)
	if err != nil {
		return err
	}
	return nil
}
