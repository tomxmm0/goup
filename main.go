package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func GenerateName(extension string) string {
	name := fmt.Sprintf("%d.%s", rand.Int(), extension)

	if _, err := os.Stat(name); os.IsNotExist(err) {
		return name
	} else {
		return GenerateName(extension)
	}
}

func UploadHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		file, handler, err := req.FormFile("file")

		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}

		defer file.Close()

		data := make([]byte, handler.Size)
		file.Read(data)

		split := strings.Split(handler.Filename, ".")
		name := fmt.Sprintf("files/%s", GenerateName(split[len(split)-1]))

		ioutil.WriteFile(name, data, 0644)

		http.Redirect(w, req, fmt.Sprintf("/%s", name), http.StatusPermanentRedirect)
	} else {
		http.Redirect(w, req, "/home/", http.StatusPermanentRedirect)
	}
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	index, _ := ioutil.ReadFile("public/index.html")
	w.Write(index)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	static_fs := http.FileServer(http.Dir("public/static"))
	files_fs := http.FileServer(http.Dir("files"))

	http.Handle("/", static_fs)
	http.Handle("/files/", http.StripPrefix("/files/", files_fs))
	http.HandleFunc("/api/upload/", UploadHandler)
	http.HandleFunc("/home/", HomeHandler)

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
