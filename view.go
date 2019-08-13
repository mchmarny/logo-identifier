package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	ev "github.com/mchmarny/gcputil/env"
)

var (
	templates *template.Template
)

func initHandlers() {
	tmpls, err := template.ParseGlob("template/*.html")
	if err != nil {
		logger.Fatalf("Error while parsing templates: %v", err)
	}
	templates = tmpls
}

func getCurrentUserID(r *http.Request) string {
	c, _ := r.Cookie(userIDCookieName)
	if c != nil {
		return c.Value
	}
	return ""
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	data["version"] = ev.MustGetEnvVar("RELEASE", "NOT SET")

	if err := templates.ExecuteTemplate(w, "index", data); err != nil {
		logger.Printf("Error in index template: %s", err)
	}

}

func errorHandler(w http.ResponseWriter, r *http.Request, err error, code int) {

	logger.Printf("Error: %v", err)
	errMsg := fmt.Sprintf("%+v", err)

	w.WriteHeader(code)
	if err := templates.ExecuteTemplate(w, "error", map[string]interface{}{
		"error":       errMsg,
		"status_code": code,
		"status":      http.StatusText(code),
	}); err != nil {
		logger.Printf("Error in error template: %s", err)
	}

}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	uid := getCurrentUserID(r)
	if uid == "" {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
		return
	}

	logger.Printf("User has ID: %s, getting data...", uid)
	usr, err := getUser(r.Context(), uid)
	if err != nil {
		logger.Printf("Error while getting user data: %v", err)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
		return
	}

	data["name"] = usr.UserName
	data["email"] = usr.Email
	data["pic"] = usr.Picture
	data["version"] = ev.MustGetEnvVar("RELEASE", "NOT SET")

	logger.Printf("Data: %v", data)

	if err := templates.ExecuteTemplate(w, "view", data); err != nil {
		logger.Printf("Error in view template: %s", err)
	}

}

func logoHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	uid := getCurrentUserID(r)
	if uid == "" {
		logger.Println("User not authenticated")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		logger.Println("Nil request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageURL := r.URL.Query().Get("imageUrl")
	logger.Printf("Logo request: %s", imageURL)

	result, err := getLogoFromURL(r.Context(), imageURL)
	if err != nil {
		logger.Printf("Error while quering logo service: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event := &UserQuery{
		QueryID:  makeUUID(),
		Created:  time.Now(),
		UserID:   uid,
		ImageURL: imageURL,
		Result:   result,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		logger.Printf("Error while encoding logo response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = saveQuery(r.Context(), event)
	if err != nil {
		logger.Printf("Error while saving logo event: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
