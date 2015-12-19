package gzip

import (
	"net/http"

	negroniGzip "github.com/phyber/negroni-gzip/gzip"
	"github.com/pilwon/gogo/adapter/negroni"
	"golang.org/x/net/context"
)

// Gzip is a middleware handler that writes gzip-compressed reponsse
type Gzip struct {
	compressionLevel int
}

// New returns a new Gzip instance
func New(compressionLevel int) *Gzip {
	return &Gzip{
		compressionLevel: compressionLevel,
	}
}

func (l *Gzip) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return negroni.Middleware(negroniGzip.Gzip(l.compressionLevel)).ServeHTTP(c, w, r)
}
