package static

import (
	"net/http"
	"path"
	"strings"

	"golang.org/x/net/context"
)

// NoDirListing is a middleware disabling directory listing of http.FileServer
// func NoDirListing(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
// 			http.Error(w, "forbidden", http.StatusForbidden)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	}
// 	return http.HandlerFunc(fn)
// }

// Static is a middleware handler that serves static files in the given directory/filesystem.
type Static struct {
	// Dir is the directory to serve static files from
	Dir http.FileSystem
	// Prefix is the optional prefix used to serve the static directory content
	Prefix string
	// IndexFile defines which file to serve as index if it exists.
	IndexFile string
}

// New returns a new instance of Static
func New(directory http.FileSystem) *Static {
	return &Static{
		Dir:       directory,
		Prefix:    "",
		IndexFile: "index.html",
	}
}

func (s *Static) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	if r.Method != "GET" && r.Method != "HEAD" {
		return c
	}
	file := r.URL.Path
	// if we have a prefix, filter requests by stripping the prefix
	if s.Prefix != "" {
		if !strings.HasPrefix(file, s.Prefix) {
			return c
		}
		file = file[len(s.Prefix):]
		if file != "" && file[0] != '/' {
			return c
		}
	}
	f, err := s.Dir.Open(file)
	if err != nil {
		// discard the error?
		return c
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return c
	}

	// try to serve index file
	if fi.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusFound)
			return nil
		}

		file = path.Join(file, s.IndexFile)
		f, err = s.Dir.Open(file)
		if err != nil {
			return c
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			return c
		}
	}

	http.ServeContent(w, r, file, fi.ModTime(), f)

	return nil
}
