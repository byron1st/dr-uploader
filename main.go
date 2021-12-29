package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/byron1st/dr-uploader/lib"
	"github.com/spf13/viper"
)

func main() {
	c, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := lib.ConnectDB(c.DBUri, c.DBName); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		relation, err := lib.Parse(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if err := relation.Upload(); err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type config struct {
	DBUri  string
	DBName string
}

func readConfig() (config, error) {
	c := config{}

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return c, err
	}

	c.DBUri, err = getValue("DB_URI")
	if err != nil {
		return c, err
	}

	c.DBName, err = getValue("DB_NAME")
	if err != nil {
		return c, err
	}

	return c, nil
}

func getValue(key string) (string, error) {
	if value, ok := viper.Get(key).(string); !ok {
		return "", errors.New(fmt.Sprintf("failed to get %s", key))
	} else {
		return value, nil
	}
}
