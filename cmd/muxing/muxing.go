package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func WriteSuccessText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

func WriteClientErrorText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	//w.WriteHeader(http.StatusBadRequest)
	//w.Write([]byte(text))
}

func WriteServerErrorText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	//w.WriteHeader(http.StatusInternalServerError)
	//w.Write([]byte(text))
}

func HelloName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["PARAM"]

	text := fmt.Sprintf("Hello, %s!", param)
	WriteSuccessText(w, text)
}

func ReadData(w http.ResponseWriter, r *http.Request) {
	var text string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		text = fmt.Sprintf("Error occured: %s", err)
		WriteServerErrorText(w, text)
		return
	}

	text = fmt.Sprintf("I got message:\n%s", string(body))
	WriteSuccessText(w, text)
}

func SumHeaders(w http.ResponseWriter, r *http.Request) {
	var text string

	var aAsString []string = r.Header["a"]
	if len(aAsString) != 1 {
		text = fmt.Sprintf("Header \"%s\" is not specified", "a")
		WriteClientErrorText(w, text)
		return
	}

	a, err := strconv.Atoi(aAsString[0])
	if err != nil {
		text = fmt.Sprintf("Error occured with header %s: %s", "a", err)
		WriteClientErrorText(w, text)
	}

	var bAsString []string = r.Header["b"]
	if len(bAsString) != 1 {
		text = fmt.Sprintf("Header \"%s\" is not specified", "b")
		WriteClientErrorText(w, text)
		return
	}

	b, err := strconv.Atoi(bAsString[0])
	if err != nil {
		text = fmt.Sprintf("Error occured with header %s: %s", "b", err)
		WriteClientErrorText(w, text)
	}

	text = strconv.Itoa(a + b)
	WriteSuccessText(w, text)
}

// Start /** Starts the web server listener on given host and port.
func Start(host string, port int) {
	router := mux.NewRouter()

	router.HandleFunc("/name/{PARAM}", HelloName).Methods(http.MethodGet)
	router.HandleFunc("/data", ReadData).Methods(http.MethodPost)
	router.HandleFunc("/headers", SumHeaders).Methods(http.MethodPost)
	router.HandleFunc("/bad", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}).Methods(http.MethodGet)
	log.Println(fmt.Printf("Starting API server on %s:%d\n", host, port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router); err != nil {
		log.Fatal(err)
	}
}

//main /** starts program, gets HOST:PORT param and calls Start func.
func main() {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8081
	}
	Start(host, port)
}
