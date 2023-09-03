package network

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Cookie struct {
    Name  string
    Value string

    Path       string    
    Domain     string    
    Expires    time.Time 
    RawExpires string   

    // MaxAge=0 means no 'Max-Age' attribute specified.
    // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
    // MaxAge>0 means Max-Age attribute present and given in seconds
    MaxAge   int 
    Secure   bool
    HttpOnly bool
    Raw      string
    Unparsed []string
}

func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func FourOhFour(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, s-maxage=604800")
	fmt.Fprint(w, "<h1>404</h1>");
	return
}

func SetCookie(w http.ResponseWriter, r *http.Request, name string, value string) {
    // Initialize a new cookie containing the string "Hello world!" and some
    // non-default attributes.
    cookie := http.Cookie{
        Name:     name,
        Value:    value,
        Path:     "/",
        MaxAge:   3600,
        HttpOnly: true,
        // Secure:   true,
        SameSite: http.SameSiteLaxMode,
    }

    // Use the http.SetCookie() function to send the cookie to the client.
    // Behind the scenes this adds a `Set-Cookie` header to the response
    // containing the necessary cookie data.
    http.SetCookie(w, &cookie)
}

func GetCookie(w http.ResponseWriter, r *http.Request, name string) string {
    // Retrieve the cookie from the request using its name (which in our case is
    // "exampleCookie"). If no matching cookie is found, this will return a
    // http.ErrNoCookie error. We check for this, and return a 400 Bad Request
    // response to the client.
    cookie, err := r.Cookie(name)
    
		if err != nil {
        switch {
        	case errors.Is(err, http.ErrNoCookie):
						return ""
            // http.Error(w, "cookie not found", http.StatusBadRequest)
        	default:
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
    }

    // Echo out the cookie value in the response body.
		return cookie.Value
}

func Authorization(w http.ResponseWriter, r *http.Request) {
	// Handle form submission here
	SetCookie(w, r, "ally-admin-session", "1234")
	Redirect(w, r, "/ally-admin")
}