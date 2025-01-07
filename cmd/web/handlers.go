package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"forum/cmd/internal/models"
	"forum/cmd/internal/validator"
)

// home displays the 10 latest threads.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	threads, err := app.threads.Latests()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Threads = threads

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

// accountCreateForm holds the data for the account creation form.
type createUserForm struct {
	Username string
	Email    string
	Password string
	validator.Validator
}

// accountCreate displays the account creation form.
func (app *application) accountCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = createUserForm{}
	app.render(w, r, http.StatusOK, "account-create.tmpl", data)
}

// accountCreatePost creates an account and redirects to view of account
func (app *application) accountCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := createUserForm{
		Username: r.PostForm.Get("username"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be blank.")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank.")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank.")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 10 characters long.")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field is not a valid email address.")
	form.CheckField(validator.UpperCase(form.Password), "password", "This field must contain at least one uppercase letter.")
	form.CheckField(validator.ContainsNumber(form.Password), "password", "This field must contain at least one number.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "account-create.tmpl", data)
		return
	}

	id, err := app.users.InsertUser(form.Username, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			data := app.newTemplateData(r)
			form.AddFieldError("email", "Address is already in use")
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "account-create.tmpl", data)
			return
		}
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Account created successfully!")
	http.Redirect(w, r, fmt.Sprintf("/account/view/%d", id), http.StatusSeeOther)
}

// accountView displays the account details for a specific user.
func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || accountID < 1 {
		http.NotFound(w, r)
		return
	}

	user, err := app.users.GetUser(accountID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// should each account path be protected? we did not see this in class but found method in book ch11
	userSessionID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if accountID != userSessionID {
		http.Redirect(w, r, "/account/login", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.User = user

	app.render(w, r, http.StatusOK, "account-view.tmpl", data)
}

// accountLoginForm holds the data for the account login form.
type accountLoginForm struct {
	Username string
	Email    string
	Password string
	validator.Validator
}

// accountLogin displays the login form.
func (app *application) accountLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = accountLoginForm{}
	app.render(w, r, http.StatusOK, "account-login.tmpl", data)
}

// accountLoginPost logs in the user.
func (app *application) accountLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := accountLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank.")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field is not a valid email address.")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "account-login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "account-login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/thread/create", http.StatusSeeOther)
}

// userLogoutPost logs out the user.
func (app application) accountLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// createThreadForm holds the data for the thread creation form.
type createThreadForm struct {
	Title string
	validator.Validator
}

// threadCreate displays the thread creation form.
func (app *application) threadCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = createThreadForm{}
	app.render(w, r, http.StatusOK, "thread-create.tmpl", data)
}

// threadCreatePost creates a thread.
func (app *application) threadCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := createThreadForm{
		Title: r.PostForm.Get("title"),
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank.")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters).")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "thread-create.tmpl", data)
		return
	}

	userSessionID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userSessionID == 0 {
		http.Redirect(w, r, "/account/login", http.StatusSeeOther)
		return
	}

	threadID, err := app.threads.Insert(form.Title, userSessionID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Thread created successfully!")
	http.Redirect(w, r, fmt.Sprintf("/thread/view/%d", threadID), http.StatusSeeOther)
}

// threadView displays a single thread.
func (app *application) threadView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	thread, err := app.threads.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Thread = thread

	app.render(w, r, http.StatusOK, "thread-view.tmpl", data)
}

// createMessageForm holds the data for the message creation form.
type createMessageForm struct {
	Message string
	validator.Validator
}

// messageCreate displays the message creation form for a specific thread.
func (app *application) messageCreate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	thread, err := app.threads.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Thread = thread
	data.Form = createMessageForm{}

	app.render(w, r, http.StatusOK, "message-create.tmpl", data)
}

// messageCreatePost creates a message and redirects to updated thead.
func (app *application) messageCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	threadID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || threadID < 1 {
		http.NotFound(w, r)
		return
	}

	thread, err := app.threads.Get(threadID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	form := createMessageForm{
		Message: r.PostForm.Get("message"),
	}

	form.CheckField(validator.NotBlank(form.Message), "message", "This field cannot be blank.")
	form.CheckField(validator.MaxChars(form.Message, 1000), "message", "This field cannot be more than 1000 characters).")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Thread = thread
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "message-create.tmpl", data)
		return
	}

	userSessionID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userSessionID == 0 {
		http.Redirect(w, r, "/account/login", http.StatusSeeOther)
		return
	}

	_, err = app.messages.InsertMessage(form.Message, threadID, userSessionID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Message created successfully!")
	http.Redirect(w, r, fmt.Sprintf("/thread/view/%d", threadID), http.StatusSeeOther)
}
