package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("POST /form", app.FormData)

	mux.HandleFunc("GET /form", app.getForm)
	mux.HandleFunc("GET /", app.showOrdersPage)

	return app.logRequest(app.cors(mux))

}
