package main

import (
	"log"
	"os"
	"time"

	. "github.com/confkai/confkai/confkai"
)

type MyConfig struct {
	Environment       Valuer[string]
	DatabaseName      Valuer[string]
	SlowMessage       Valuer[string]
	SlowMessageCached Valuer[string]
}

var (
	environment = "environment"
	_           = os.Setenv("my_env", "staging")
)

var config = MyConfig{
	Environment: RegisterTag(environment, Value(os.Getenv("my_env"))),
	DatabaseName: FirstOf(
		Tag(environment, "dev", Value("my_dev_db")),
		Tag(environment, "staging", Value("my_staging_db")),
		Tag(environment, "prod", Value("my_prod_db")),
	),
	SlowMessage: FuncValue(func() (string, error) {
		time.Sleep(3 * time.Second)
		return "hello world", nil
	}),
	SlowMessageCached: Cached(FuncValue(func() (string, error) {
		time.Sleep(3 * time.Second)
		return "hello universe", nil
	})),
}

func main() {
	log.Println(config.Environment.Must())
	log.Println(config.DatabaseName.Must())
	log.Println(config.SlowMessage.Must())
	log.Println(config.SlowMessageCached.Must())
	log.Println(config.SlowMessageCached.Must())
}

// output: 2024/01/03 20:18:19 staging
// output: 2024/01/03 20:18:19 my_staging_db
// output: 2024/01/03 20:18:22 hello world
// output: 2024/01/03 20:18:25 hello universe
// output: 2024/01/03 20:18:25 hello universe
