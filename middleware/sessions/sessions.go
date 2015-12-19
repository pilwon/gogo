package sessions

import (
	"net/http"

	gorillaSessions "github.com/gorilla/sessions"
	"github.com/pilwon/gogo/middleware"
	"golang.org/x/net/context"
)

type key int

const (
	storeKey = 0
)

// New returns session middleware
func New(name string, store gorillaSessions.Store) middleware.Handler {
	return middleware.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		session, _ := store.Get(r, sessionName)
		c = context.WithValue(c, storeKey, store)
		// gogo.Session(c).Values["foo"]
		return c
	})
}

//
func AddFlash(value interface{}, vars ...string) {

}

// Flashes returns a slice of flash messages from the session.
//
// A single variadic argument is accepted, and it is optional:
// it defines the flash key. If not defined "_flash" is used by default.
func Flashes(vars ...string) []interface{} {

}

func Name() string {

}

// Save is a convenience method to save this session.
// It is the same as calling store.Save(request, response, session).
// You should call Save before writing to the response or returning from the handler.
func Save(r *http.Request, w http.ResponseWriter) error {

}

// Store returns the session store used to register the session.
func Store() gorillaSessions.Store {

}
