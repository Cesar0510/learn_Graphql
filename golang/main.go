package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func messageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Go Web Development")
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	h := md5.New()
	token := fmt.Sprintf("%x", h.Sum(nil))
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		fmt.Fprintf(w, token)

	} else {
		fmt.Fprintf(w, "POST")
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./public/"+token+".png", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {

	http.HandleFunc("/menssage", messageHandler)
	http.HandleFunc("/upload", upload)
	server := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening... in the Port")

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	server.ListenAndServe()
}
