package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

type FavoriteModel struct {
	DB *sql.DB
}

type SnippetViewData struct {
		Snippet    *Snippet
		IsFavorite bool
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	s := &Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *FavoriteModel) Insert(snippetID int) error {
	stmt := `INSERT INTO favorites (snippet_id) VALUES (?)`
	_, err := m.DB.Exec(stmt, snippetID)
	return err
}

func (m *FavoriteModel) Delete(snippetID int) error {
	stmt := `DELETE FROM favorites WHERE snippet_id = ?`
	_, err := m.DB.Exec(stmt, snippetID)
	return err
}

func (m *FavoriteModel) IsFavorite(snippetID int) (bool, error) {
	stmt := `SELECT COUNT(*) FROM favorites WHERE snippet_id = ?`
	var count int
	err := m.DB.QueryRow(stmt, snippetID).Scan(&count)
	if err != nil {
			return false, err
	}
	return count > 0, nil
}

func (m *FavoriteModel) GetAll() ([]*Snippet, error) {
		stmt := `SELECT snippets.id, snippets.title, snippets.content, snippets.created, snippets.expires 
						 FROM snippets INNER JOIN favorites ON snippets.id = favorites.snippet_id`

		rows, err := m.DB.Query(stmt)
		if err != nil {
				return nil, err
		}
		defer rows.Close()

		var snippets []*Snippet

		for rows.Next() {
				s := &Snippet{}
				err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
				if err != nil {
						return nil, err
				}
				snippets = append(snippets, s)
		}

		return snippets, nil
}
