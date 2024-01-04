package main

import (
	"log"
	"os"
	"time"

	. "github.com/confkai/confkai/confkai"
)

type MyConfig struct {
	Environment       func() string
	DatabaseName      func() string
	SlowMessage       func() string
	SlowMessageCached func() string
}

var (
	environment = "environment"
	_           = os.Setenv("my_env", "staging")
)

var config = MyConfig{
	Environment: RegisterTag(environment, Value(os.Getenv("my_env"))).Must(),
	DatabaseName: FirstOf(
		Tag(environment, "dev", Value("my_dev_db")),
		Tag(environment, "staging", Value("my_staging_db")),
		Tag(environment, "prod", Value("my_prod_db")),
	).Must(),
	SlowMessage: FuncValue(func() (string, error) {
		time.Sleep(3 * time.Second)
		return "hello world", nil
	}).Must(),
	SlowMessageCached: Cached(FuncValue(func() (string, error) {
		time.Sleep(3 * time.Second)
		return "hello universe", nil
	})).Must(),
}

func main() {

	log.Println(config.Environment())
	// 2024/01/03 20:18:19 staging
	log.Println(config.DatabaseName())
	// 2024/01/03 20:18:19 my_staging_db
	log.Println(config.SlowMessage())
	// 2024/01/03 20:18:22 hello world
	log.Println(config.SlowMessageCached())
	// 2024/01/03 20:18:25 hello universe
	log.Println(config.SlowMessageCached())
	// 2024/01/03 20:18:25 hello universe

}
