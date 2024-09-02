package main

import (
  "net/http"

  "github.com/julienschmidt/httprouter"
  "github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

  router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    app.notFound(w)
  })

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/", http.StripPrefix("/static", fileServer))

  router.HandlerFunc(http.MethodGet, "/", app.home)
  router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
  router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)
  //Funcionalidade de exportar snippet
  router.HandlerFunc(http.MethodGet, "/snippet/export", app.exportSnippet)

  router.HandlerFunc(http.MethodGet, "/snippet/favorite", app.favoriteSnippet)
  router.HandlerFunc(http.MethodGet, "/snippet/unfavorite", app.unfavoriteSnippet)
  router.HandlerFunc(http.MethodGet, "/snippet/favorites", app.listFavorites)


  standard := alice.New(app.recoverPanic, app.logRequest)
  
	return standard.Then(router)
}
