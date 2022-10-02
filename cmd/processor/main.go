package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/viveknathani/nattukaka/database"
	"github.com/viveknathani/nattukaka/entity"
)

var databaseServer string = ""

// getDatabase will init and return a db
func getDatabase() *database.Database {

	db := &database.Database{}
	err := db.Initialize(databaseServer)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return db
}

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	databaseServer = os.Getenv("DATABASE_URL")
}

func getIPInfo(ip string) ([]byte, error) {

	response, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func processContent(db *database.Database, line []byte) error {

	var jsonObject entity.Log
	err := json.Unmarshal(line, &jsonObject)
	if err != nil {
		return err
	}
	jsonObject.Info, err = getIPInfo(jsonObject.IP)
	if err != nil {
		return err
	}
	db.InsertLog(&jsonObject)
	return nil
}

func runStatementAndGetOutput(statement string) string {

	cmd := `echo "` + statement + `" | psql -d nattukaka`
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}

func main() {

	db := getDatabase()
	defer db.Close()
	file, err := os.Open("/var/logs.txt")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		err = processContent(db, fileScanner.Bytes())
		if err != nil {
			log.Println(err)
		}
	}

	top10VisitedPaths := "select path, count(*) from logs where path not like '/static/%' group by path order by count desc limit 10;"
	top10IPs := `select ip, info -> 'country' as "country", count(*) from logs group by ip, info -> 'country';`
	countryCount := `select info -> 'country' as "country", count(*) from logs group by info -> 'country';`
	indianStateCount := `select info->'regionName' as "state", count(*) from logs where info -> 'country' = '\"India\"' group by info->'regionName';`

	out := ""
	out += runStatementAndGetOutput(top10VisitedPaths)
	out += runStatementAndGetOutput(top10IPs)
	out += runStatementAndGetOutput(countryCount)
	out += runStatementAndGetOutput(indianStateCount)

	fmt.Println(out)
}
