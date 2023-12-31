package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func uploadmeme(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	imgFile, imgHandler, err := r.FormFile("image")
	check(err)
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

	newmeme.Title = r.FormValue("title")
	truncatedname := TruncateTitle(newmeme.Title) + "_" + RandomString()
	imageExtension := filepath.Ext(imgHandler.Filename)

	newmeme.Img = truncatedname + imageExtension

	newmeme.SoundFile = truncatedname + ".mp3"

	jsonName := truncatedname + ".json"

	jsonNamePath := filepath.Join("data", jsonName)
	fmt.Println("New meme submitted: ", newmeme)
	file, _ := json.MarshalIndent(newmeme, "", " ")
	_ = os.WriteFile(jsonNamePath, file, 0644)

	// Create file
	dst, err := os.Create(filepath.Join("data", "img", newmeme.Img))
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
	dst, err = os.Create(filepath.Join("data", "snd", newmeme.SoundFile))
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

	Memebufferchannel <- newmeme
	const homeButton = `<a href=../>Go home</a>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "Successfully Uploaded File\n %s", homeButton)
}

func display(w http.ResponseWriter, page string, data interface{}) {
	var newmeme = template.Must(template.ParseFiles("./pages/newmeme.html"))
	newmeme.ExecuteTemplate(w, page+".html", data)
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "../upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}
func uploadloopHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "../uploadloop", nil)
	case "POST":
		uploadNewLoop(w, r)
	}
}
func uploadmemeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "../uploadmeme", nil)
	case "POST":
		uploadmeme(w, r)
	}
}

func handlePost(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "post\n")
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
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

func subscribehandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	response := <-subscribechannel
	fmt.Fprint(w, response)

}
