package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	
)

type URL struct{
	ID string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortURL string `json:"short_url"`
	
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string{
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("hasher: ", hasher)
	data := hasher.Sum(nil)
	fmt.Println("hasher data: ", data)
	hash := hex.EncodeToString(data)
	fmt.Println("EncodeToString: ", hash)
	fmt.Println("funal string: ", hash[:8])
	return hash[:8]
}

func createURL(originalURL string) string{
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID: id,
		OriginalURL: originalURL,
		ShortURL: shortURL,
	
	}

	return shortURL
}

func getURL(id string) (URL, error){
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}

	return url, nil
}

func RootPageURL(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Hello world")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request){
	var data struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	return
	}

	shortURL_ := createURL(data.URL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL:  shortURL_}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func redirectURLHandler(w http.ResponseWriter, r *http.Request){
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil{
		http.Error(w, "Invalid request", http.StatusNotFound)
	}

	http.Redirect(w,r,url.OriginalURL,http.StatusFound)
}

func main(){

	OriginalURL := "https://github.com/harshal-rembhotkar/"
	generateShortURL(OriginalURL)
	http.HandleFunc("/",RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/",redirectURLHandler)
	fmt.Println("starting server on port 3000...")
	err := http.ListenAndServe(":3000", nil)

	if err != nil{
		fmt.Println("error on starting server:", err)
	}
}