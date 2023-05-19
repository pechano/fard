package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)


func 	Wgetmemehandler (w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		wgetmemepage(w,r)
	case "POST":

		wgetmeme(w, r)
}}
func wgetmemepage(w http.ResponseWriter, r *http.Request) {
	Wgettemplate := template.Must(template.ParseFiles("./pages/wgetmeme.html"))
	err :=	Wgettemplate.Execute(w, nil)
	check(err)}


func wgetmeme(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)
	imgFile, imgHandler, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer imgFile.Close()
	fmt.Printf("Uploaded File: %+v\n", imgHandler.Filename)
	fmt.Printf("File Size: %+v\n", imgHandler.Size)
	fmt.Printf("MIME Header: %+v\n", imgHandler.Header)
	soundfile := r.FormValue("link")

	GetExternalMP3(soundfile)
	soundfile = filepath.Base(soundfile)

	var newmeme Meme
	newmeme.Img = imgHandler.Filename
	newmeme.SoundFile = soundfile
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

	const homeButton = `<a href=../>Go home</a>`


	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "Successfully Uploaded File\n %s", homeButton)
}

func GetExternalMP3(link string)(){
	soundfile := exec.Command("wget", "-P", "./data/snd/",link)
	err := soundfile.Run()
	check(err)
}
