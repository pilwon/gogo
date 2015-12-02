# gogo


## Usage

```go
package main

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/pilwon/gogo"
	_ "github.com/pilwon/gogo/router/httprouter"
	"golang.org/x/net/context"
)

func main() {
	app := gogo.New()

	app.Use(NewRecoveryMiddleware())
	app.Use(NewGorillaLogger())
	app.Use(NewGzipMiddleware(gzip.DefaultCompression))
	app.Use(NewStaticMiddleware(http.Dir("static")))

	app.Get("/", func(c context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "gogo server\n")
	})

	app.Get("/hello/:name", func(c context.Context, w http.ResponseWriter, r *http.Request) {
		name, _ := gogo.Param(c, "name")
		fmt.Fprintf(w, "Hello, %s!\n", name)
	})

	app.Run(":8080")
}
```


## Middlewares

### Supported Adapters

* [interpose](https://github.com/carbocation/interpose)
* [negroni](https://github.com/codegangsta/negroni)


## Credits

* [negroni](https://github.com/codegangsta/negroni)
* [kami](https://github.com/guregu/kami)
* [express](https://github.com/strongloop/express)
