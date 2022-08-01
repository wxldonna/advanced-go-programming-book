package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		fmt.Printf("FileName=[%s], FormName=[%s]\n", part.FileName(), part.FormName())
		if part.FileName() == "" { // this is FormData
			data, _ := ioutil.ReadAll(part)
			fmt.Printf("FormData=[%s]\n", string(data))
		} else { // This is FileData
			serverPath := "/home/vagrant/advanced-go-programming-book/examples/UploadFIle/server/"
			dst, err := os.Create(serverPath + part.FileName() + ".upload")
			if err != nil {
				panic(err)
			}
			defer dst.Close()
			io.Copy(dst, part)
		}
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	for _, v := range r.Header {
		log.Printf("header is %v", v)
	}

	io.Copy(w, r.Body)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/test", test)
	http.ListenAndServe(":8080", nil)
}
