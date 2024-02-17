package spotistats

import (
	"flag"
	"github.com/erfgypO/spotistats/lib/data"
	"github.com/joho/godotenv"
	"log"
)

func Run() {
	err := godotenv.Load()

	if err != nil && err.Error() != "open .env: no such file or directory" {
		log.Fatal("Error loading .env file")
	}

	err = data.SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	runApiFlag := flag.Bool("api", false, "Start the api server")
	flag.Parse()

	if *runApiFlag {
		runApi()
	} else {
		runScrapper()
	}
}
