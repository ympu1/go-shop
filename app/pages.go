package main

import (
	"html/template"
)

type Page struct {
	ID          int64
	Name        string
	H1          string
	Title       string
	Description string
	Text        template.HTML
}

func (page *Page) updateToDB() error {
	stmt, err := db.Prepare("UPDATE pages SET title=?, h1=?, description=?, text=? WHERE name = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(page.Title, page.H1, page.Description, page.Text, page.Name)
	if err != nil {
		return err
	}

	return nil
}

func getPageByName(name string) (Page, error) {
	var page Page

	row := db.QueryRow("SELECT title, h1, description, text FROM pages WHERE name = ?", name)
	err := row.Scan(&page.Title, &page.H1, &page.Description, &page.Text)
	if err != nil {
		return page, err
	}

	return page, nil
}