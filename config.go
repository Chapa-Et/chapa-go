package chapa

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"math/big"
)

func InitConfig() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(context.Background(), fmt.Sprintf("Failed to read config: %v", err))
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: file %v", e.Name)
	})
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(length int) string {
	var str string
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			log.Printf("error while generating string: %v", err)
		}
		str += string(alphabet[n.Int64()])
	}
	return str
}
