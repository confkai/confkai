# confkai

Confkai is a configuration as code library for Go. With this library
you compose your configurations with functions. This no-frills
library includes the basic composition functions to get you started. But,
you can write your own functions or import a separate module to support values
from other providers like GCP Secrets or AWS ParamStore! 

Benefits of using Confkai for configuration:   
1. Changes are tracked by git.    
1. Changes can be audited in every PR.  
1. Roll back Safe deployments.  
1. Lazy Loading by default.  
1. A single source of truth.  
1. Eager loading, Caching, and much more included with the base library.  
1. No dependencies in the base library.  


## Examples

```go
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
	log.Println(config.DatabaseName())
	log.Println(config.SlowMessage())
	log.Println(config.SlowMessageCached())
	log.Println(config.SlowMessageCached())
}
// output: 2024/01/03 20:18:19 staging
// output: 2024/01/03 20:18:19 my_staging_db
// output: 2024/01/03 20:18:22 hello world
// output: 2024/01/03 20:18:25 hello universe
// output: 2024/01/03 20:18:25 hello universe
```