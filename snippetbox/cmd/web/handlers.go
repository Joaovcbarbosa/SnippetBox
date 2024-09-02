package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"snippetbox.otaviolemos.com/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData()
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData()
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fieldErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 character long"
	}

	

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) exportSnippet(w http.ResponseWriter, r *http.Request) {
		// Obter o ID do snippet a partir da URL
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
				http.NotFound(w, r)
				return
		}

		// Obter o snippet do banco de dados
		snippet, err := app.snippets.Get(id)
		if err == models.ErrNoRecord {
				http.NotFound(w, r)
				return
		} else if err != nil {
				app.serverError(w, err)
				return
		}

		// Obter o formato escolhido pelo usuário
		format := r.URL.Query().Get("format")

		// Definir o nome do arquivo e o tipo de conteúdo baseado no formato
		var fileName string
		var contentType string
		var content string

		switch format {
		case "md":
				fileName = fmt.Sprintf("%s.md", snippet.Title)
				contentType = "text/markdown"
				content = fmt.Sprintf("# %s\n\n%s", snippet.Title, snippet.Content)
		case "txt":
				fallthrough // Se o usuário escolher txt ou se o formato não for reconhecido, usa txt como padrão
		default:
				fileName = fmt.Sprintf("%s.txt", snippet.Title)
				contentType = "text/plain"
				content = snippet.Content
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", contentType)

		// Escrever o conteúdo do snippet no response
		w.Write([]byte(content))
}


