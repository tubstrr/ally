package network

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tubstrr/ally/database/sessions"
	"github.com/tubstrr/ally/database/users"
	"github.com/tubstrr/ally/utilities/validation"
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

// Utility functions
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

// Cookie functions
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

func DeleteCookie(w http.ResponseWriter, r *http.Request, name string) {
		// Set the cookie value to empty, and set the max age to -1, i.e. delete it
		// immediately.
		cookie := http.Cookie{
				Name:   name,
				Value:  "",
				Path:   "/",
				MaxAge: -1,
		}

		// Use the http.SetCookie() function to send the cookie to the client.
		// Behind the scenes this adds a `Set-Cookie` header to the response
		// containing the necessary cookie data.
		http.SetCookie(w, &cookie)
}

// Form functions
func Authorization(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Get the form data
	username := r.FormValue("username")
	username = validation.ConvertUsername(username)
	password := r.FormValue("password")

	// Check if all the fields are filled out
	if (username == "" || password == "") {
		// Figure out which fields are missing
		missing_fields := ""
		if (username == "") { missing_fields += "username," }
		if (password == "") { missing_fields += "password," }

		// Remove the last comma
		missing_fields = missing_fields[:len(missing_fields)-1]
		
		// Redirect with error
		redirect_url := "/ally-admin/login?error=missing_fields&missing_fields=" + missing_fields
		Redirect(w, r, redirect_url)
		return
	}

	// Check if the user exists
	user := users.GetUserByUsername(username)
	if (user.Username == "") {
		// Redirect with error
		redirect_url := "/ally-admin/login?error=invalid_username"
		Redirect(w, r, redirect_url)
		return
	}

	// Check if the password is correct
	if (!validation.VerifyPassword(user.Password, password)) {
		// Redirect with error
		redirect_url := "/ally-admin/login?error=invalid_password"
		Redirect(w, r, redirect_url)
		return
	}

	// If we get here, the user is valid, so set the session cookie
	// Handle form submission here
	session := uuid.New().String()
	
	// Set the cookie
	sessions.SetSessionToken(user.Id, session)

	SetCookie(w, r, "ally-admin-session", session)
	Redirect(w, r, "/ally-admin")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Handle form submission here
	session := GetCookie(w, r, "ally-admin-session")
	sessions.DeleteSessionToken(session)
	DeleteCookie(w, r, "ally-admin-session")
	Redirect(w, r, "/ally-admin/login")
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Get the form data
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm_password := r.FormValue("confirm_password")
	redirect := r.FormValue("redirect")

	// Convert role to int
	role, err := strconv.Atoi(r.FormValue("role"))
	if (err != nil) {
		role = 2
	}
		

	if (redirect == "") {
		redirect = "/ally-admin"
	}

	// Check if all the fields are filled out
	if (username == "" || email == "" || password == "" || confirm_password == "") {
		// Figure out which fields are missing
		missing_fields := ""
		if (username == "") { missing_fields += "username," }
		if (email == "") { missing_fields += "email," }
		if (password == "") { missing_fields += "password," }
		if (confirm_password == "") { missing_fields += "confirm_password," }
		
		// Remove the last comma
		missing_fields = missing_fields[:len(missing_fields)-1]

		// Redirect with error
		redirect_url := redirect + "?error=missing_fields&missing_fields=" + missing_fields
		Redirect(w, r, redirect_url)
		return
	}
	
	// Check if the passwords match
	if (password != confirm_password) {
		// Redirect with error
		redirect_url := redirect + "?error=passwords_do_not_match"
		Redirect(w, r, redirect_url)
		return
	}

	// Check for valid email
	if (!validation.ValidateEmail(email)) {
		// Redirect with error
		redirect_url := redirect + "?error=invalid_email"
		Redirect(w, r, redirect_url)
	}

	// Convert username
	username = validation.ConvertUsername(username)

	// Convert password
	password = validation.ConvertPassword(password)

	// Check if the users table is empty
	// We know it's safe to create an account if the users table is empty
	if (users.IsUserTableEmpty()) {
		role = 1
		users.CreateUser(username, email, password, role)
		Redirect(w, r, "/ally-admin/forms/auth?username=" + username + "&password=" + r.FormValue("password"))
		return
	}
	
	// Is the username taken?
	if (!users.IsValidUsername(username)) {
		// Redirect with error
		redirect_url := redirect + "?error=username_taken"
		Redirect(w, r, redirect_url)
		return
	}

	// Is the email taken?
	if (!users.IsValidEmail(email)) {
		// Redirect with error
		redirect_url := redirect + "?error=email_taken"
		Redirect(w, r, redirect_url)
		return
	}

	// If we get here, the user is valid
	// But we need to check if the user has permission to create an account
	user := sessions.GetUserFromSession(GetCookie(w, r, "ally-admin-session"))

	if (user.Role <= 2 && user.Role != 0) {
		users.CreateUser(username, email, password, role)
		Redirect(w, r, redirect + "?success=account_created")
	} else {
		Redirect(w, r, redirect + "?error=permission_denied")
	}

}

func RedirectIfUserLoggedIn(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	loggedIn := IsUserLoggedIn(w, r)
	// If the user is logged in, redirect them to the dashboard
	if (loggedIn) {
		Redirect(w, r, "/ally-admin")
	}
}

func IsUserLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	loggedIn := false
	// First get the session cookie
	session := GetCookie(w, r, "ally-admin-session")
	// Then check if the session cookie is empty
	if (session == "") {
		loggedIn = false
	} else {
		// Then check if the session cookie is valid
		loggedIn = sessions.CheckSessionToken(session)
		if (!loggedIn) {
			// If the session cookie is invalid, delete it
			DeleteCookie(w, r, "ally-admin-session")
		}
	}
	return loggedIn
}

func GetUserID(w http.ResponseWriter, r *http.Request) int {
	session := GetCookie(w, r, "ally-admin-session")
	id := sessions.GetUserIDFromSession(session)
	return id
}