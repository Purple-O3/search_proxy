package tools

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"

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

func Str2Bytes(s string) []byte {
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)[:len(s):len(s)]
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func TimeCost() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		return time.Since(start)
	}
}
