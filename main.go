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

		if err := relation.Upload(c.ProjectID); err != nil {
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
	DBUri     string
	DBName    string
	ProjectID string
}

func readConfig() (config, error) {
	c := config{}

	dbURIFull := os.Getenv("DB")
	if dbURIFull == "" {
		return c, errors.New("no DB URI set")
	}

	projectID := os.Getenv("PROJECT")
	if projectID == "" {
		return c, errors.New("no project ID set")
	}

	idx := strings.LastIndex(dbURIFull, "/")
	if idx == -1 {
		return c, errors.New("wrong DB URI")
	}
	c.DBUri = dbURIFull[:idx]
	c.DBName = dbURIFull[idx+1:]

	c.ProjectID = projectID

	fmt.Printf("DB:%s/%s,PROJ=%s", c.DBUri, c.DBName, c.ProjectID)

	return c, nil
}
