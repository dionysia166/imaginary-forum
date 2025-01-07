package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", app.dynamic(app.home))

	mux.Handle("GET /account/create", app.dynamic(app.accountCreate))
	mux.Handle("POST /account/create", app.dynamic(app.accountCreatePost))
	mux.Handle("GET /account/view/{id}", app.protected(app.accountView))

	mux.Handle("GET /account/login", app.dynamic(app.accountLogin))
	mux.Handle("POST /account/login", app.dynamic(app.accountLoginPost))
	mux.Handle("POST /account/logout", app.protected(app.accountLogoutPost))

	mux.Handle("GET /thread/create", app.protected(app.threadCreate))
	mux.Handle("POST /thread/create", app.protected(app.threadCreatePost))
	mux.Handle("GET /thread/view/{id}", app.dynamic(app.threadView))

	mux.Handle("GET /thread/view/{id}/message/create", app.protected(app.messageCreate))
	mux.Handle("POST /thread/view/{id}/message/create", app.protected(app.messageCreatePost))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	return app.logRequest(commonHeaders(mux))
}

func (app *application) protected(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(handler)))
}

func (app *application) dynamic(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return app.sessionManager.LoadAndSave(http.HandlerFunc(handler))
}
