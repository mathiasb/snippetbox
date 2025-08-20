package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mathiasb/snippetbox/pkg/forms"
	"github.com/mathiasb/snippetbox/pkg/models"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Check if this is an HTMX request
	if app.isHTMXRequest(r) {
		// For HTMX requests, render just the snippets partial
		app.render(w, r, "home.snippets.partial.tmpl", &templateData{Snippets: s})
	} else {
		// For regular requests, render the full page
		app.render(w, r, "home.page.tmpl", &templateData{Snippets: s})
	}
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		// Check if this is an HTMX request
		if app.isHTMXRequest(r) {
			// For HTMX requests, render just the form partial with errors
			app.render(w, r, "create.form.partial.tmpl", &templateData{Form: form})
		} else {
			// For regular requests, render the full page with errors
			app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		}
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Check if this is an HTMX request
	if app.isHTMXRequest(r) {
		// For HTMX requests, return a success response (e.g., redirect via HTMX)
		w.Header().Set("HX-Redirect", fmt.Sprintf("/snippet/%d", id))
		w.WriteHeader(http.StatusOK)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		// Check if this is an HTMX request
		if app.isHTMXRequest(r) {
			// For HTMX requests, render just the form partial with errors
			app.render(w, r, "signup.form.partial.tmpl", &templateData{Form: form})
		} else {
			// For regular requests, render the full page with errors
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		}
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			// Check if this is an HTMX request
			if app.isHTMXRequest(r) {
				// For HTMX requests, render just the form partial with errors
				app.render(w, r, "signup.form.partial.tmpl", &templateData{Form: form})
			} else {
				// For regular requests, render the full page with errors
				app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
			}
		} else {
			app.serverError(w, err)
		}
		return
	}
	
	// Check if this is an HTMX request
	if app.isHTMXRequest(r) {
		// For HTMX requests, return a success response (e.g., redirect via HTMX)
		w.Header().Set("HX-Redirect", "/user/login")
		w.WriteHeader(http.StatusOK)
		return
	}
	
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic error
	// message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			// Check if this is an HTMX request
			if app.isHTMXRequest(r) {
				// For HTMX requests, render just the form partial with errors
				app.render(w, r, "login.form.partial.tmpl", &templateData{Form: form})
			} else {
				// For regular requests, render the full page with errors
				app.render(w, r, "login.page.tmpl", &templateData{Form: form})
			}
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Add the ID of the current user to the session, so that they are now 'logged
	// in'.
	app.session.Put(r, "authenticatedUserID", id)

	// Check if this is an HTMX request
	if app.isHTMXRequest(r) {
		// For HTMX requests, return a success response (e.g., redirect via HTMX)
		w.Header().Set("HX-Redirect", "/snippet/create")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.session.Remove(r, "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.session.Put(r, "flash", "You've been logged out successfully!")
	
	// Check if this is an HTMX request
	if app.isHTMXRequest(r) {
		// For HTMX requests, return a success response (e.g., redirect via HTMX)
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

I also need to update the testutils_test.go file to add the missing imports and helper functions:

testutils_test.go
