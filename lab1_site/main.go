package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"lab1/models"
)

var order []models.Feature

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/feature/", detailHandler)
	http.HandleFunc("/order", orderHandler)
	http.HandleFunc("/order/delete", deleteHandler) // новый хендлер

	http.ListenAndServe(":8080", nil)
}

// Главная страница (список признаков)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	var featuresToShow []models.Feature

	// Обработка поиска (выполняется всегда)
	query := r.FormValue("search")
	if query != "" {
		query = strings.ToLower(query)
		for _, f := range models.Features {
			if strings.Contains(strings.ToLower(f.Name), query) {
				featuresToShow = append(featuresToShow, f)
			}
		}
	} else {
		featuresToShow = models.Features
	}

	// Обработка добавления в заказ (только если есть id)
	if r.Method == "POST" {
		idStr := r.FormValue("id")
		if idStr != "" {
			id, _ := strconv.Atoi(idStr)
			for _, f := range models.Features {
				if f.ID == id {
					order = append(order, f)
					break
				}
			}
		}
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, featuresToShow)
}

// Страница признака
func detailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	var selected models.Feature
	for _, f := range models.Features {
		if f.ID == id {
			selected = f
		}
	}

	// добавление в заказ
	if r.Method == "POST" {
		order = append(order, selected)
		http.Redirect(w, r, "/order", http.StatusSeeOther)
		return
	}

	tmpl, _ := template.ParseFiles("templates/detail.html")
	tmpl.Execute(w, selected)
}

// Страница заказа
func orderHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/order.html")
	tmpl.Execute(w, order)
}

// Удаление признака из заказа
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	var newOrder []models.Feature
	for _, f := range order {
		if f.ID != id {
			newOrder = append(newOrder, f)
		}
	}
	order = newOrder

	http.Redirect(w, r, "/order", http.StatusSeeOther)
}
