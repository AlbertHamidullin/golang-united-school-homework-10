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
	if len(text) > 0 {
		w.Write([]byte(text))
	}
}

func WriteClientErrorText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	// w.WriteHeader(http.StatusBadRequest)
	// if len(text) > 0 {
	// 	w.Write([]byte(text))
	// }
}

func WriteServerErrorText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	// w.WriteHeader(http.StatusInternalServerError)
	// if len(text) > 0 {
	// 	w.Write([]byte(text))
	// }
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

func ReadHeaderAsInt(r *http.Request, name string) (i int, err error) {
	var s string = r.Header.Get(name)
	if len(s) == 0 {
		err = fmt.Errorf("header \"%s\" is not specified", name)
		return -1, err
	}

	i, err = strconv.Atoi(s)
	if err != nil {
		err = fmt.Errorf("error occured with header %s: %w", name, err)
		return -1, err
	}

	return i, nil
}

func SumHeaders(w http.ResponseWriter, r *http.Request) {
	var text string

	a, err := ReadHeaderAsInt(r, "a")
	if err != nil {
		WriteClientErrorText(w, err.Error())
		return
	}

	b, err := ReadHeaderAsInt(r, "b")
	if err != nil {
		WriteClientErrorText(w, err.Error())
		return
	}

	text = strconv.Itoa(a + b)
	w.Header().Set("a+b", text)
	WriteSuccessText(w, "")
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
