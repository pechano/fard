package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)


	func uploadmeme (w http.ResponseWriter, r *http.Request) {
		// Maximum upload of 10 MB files
		r.ParseMultipartForm(10 << 20)

		// Get handler for filename, size and headers
		imgFile, imgHandler, err := r.FormFile("image")
		sndFile, sndHandler, err := r.FormFile("sound")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}



		defer imgFile.Close()
		defer sndFile.Close()
		fmt.Printf("Uploaded File: %+v\n", imgHandler.Filename)
		fmt.Printf("Uploaded File: %+v\n", sndHandler.Filename)
		fmt.Printf("File Size: %+v\n", imgHandler.Size)
		fmt.Printf("MIME Header: %+v\n", imgHandler.Header)
		var newmeme Meme
		newmeme.Img = imgHandler.Filename
		newmeme.SoundFile = sndHandler.Filename
		newmeme.Title = r.FormValue("title")

		jsonName := r.FormValue("memename")+".json"
		jsonNamePath := filepath.Join("data",jsonName)
		fmt.Println("New meme submitted: ",newmeme)
		file, _ := json.MarshalIndent(newmeme,""," ")
		_ = ioutil.WriteFile(jsonNamePath,file, 0644)

		// Create file
		dst, err := os.Create(filepath.Join("data","img",imgHandler.Filename))
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, imgFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create file
		dst, err = os.Create(filepath.Join("data","snd",sndHandler.Filename))
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, sndFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		const homeButton = `<a href=../>Go home</a>`


		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "Successfully Uploaded File\n %s", homeButton)
	}
func buildMeme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Println(r.Header)
	fmt.Println(r.Body)
	var build Meme
	_ = json.NewDecoder(r.Body).Decode(&build)
	fmt.Println("before encode: ",build)
	json.NewEncoder(w).Encode(build)
	fmt.Println("New meme received: ",build)
}


	 func display (w http.ResponseWriter, page string, data interface{}) {
	var newmeme = template.Must(template.ParseFiles("./pages/newmeme.html"))
		newmeme.ExecuteTemplate(w, page+".html", data)
	}
	 func uploadHandler (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			display(w, "../upload", nil)
		case "POST":
			uploadFile(w, r)
		}
	}

	 func uploadmemeHandler (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			display(w, "../uploadmeme", nil)
		case "POST":
			uploadmeme(w, r)
	}}

func handlePost(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "post\n")
}

	func uploadFile (w http.ResponseWriter, r *http.Request) {
		// Maximum upload of 10 MB files
		r.ParseMultipartForm(10 << 20)

		// Get handler for filename, size and headers
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}

		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		// Create file
		dst, err := os.Create(handler.Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		const homeButton = `<a href=../>Go home</a>`


		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "Successfully Uploaded File\n %s", homeButton)
	}
