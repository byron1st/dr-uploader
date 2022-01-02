package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/byron1st/dr-uploader/lib"
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
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}

		relation, err := lib.Parse(line)
		if err != nil {
			log.Fatal(err)
		}

		if err := relation.Upload(); err != nil {
			log.Fatal(err)
		}
		count++
	}

	fmt.Printf("Inserted=%d", count)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err := lib.DisconnectDB(); err != nil {
		log.Fatal(err)
	}
}

type config struct {
	DBUri  string
	DBName string
}

func readConfig() (config, error) {
	c := config{}

	dbURIFull := os.Getenv("DB")
	if dbURIFull == "" {
		return c, errors.New("no DB URI set")
	}

	idx := strings.LastIndex(dbURIFull, "/")
	if idx == -1 {
		return c, errors.New("wrong DB URI")
	}
	c.DBUri = dbURIFull[:idx]
	c.DBName = dbURIFull[idx+1:]

	fmt.Printf("DB:%s/%s", c.DBUri, c.DBName)

	return c, nil
}
