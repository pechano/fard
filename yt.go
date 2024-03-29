package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func YTmemehandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		YTmemepage(w, r)
	case "POST":
		YTmeme(w, r)
	}
}
func YTmemepage(w http.ResponseWriter, r *http.Request) {
	YTtemplate := template.Must(template.ParseFiles("./pages/ytmeme.html"))
	err := YTtemplate.Execute(w, nil)
	check(err)
}

func YTmeme(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
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
	soundfile := r.FormValue("YTsound")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "Downloading File...\n")
	fmt.Println("downloading " + soundfile)
	GetYTsnd(soundfile)

	fmt.Fprintf(w, "Creating meme...\n")
	var newmeme Meme
	newmeme.Img = imgHandler.Filename
	newmeme.SoundFile = soundfile + ".mp3"
	newmeme.Title = r.FormValue("title")

	jsonName := r.FormValue("memename") + ".json"
	jsonNamePath := filepath.Join("data", jsonName)
	fmt.Println("New meme submitted: ", newmeme)
	file, _ := json.MarshalIndent(newmeme, "", " ")
	_ = os.WriteFile(jsonNamePath, file, 0644)

	// Create file
	dst, err := os.Create(filepath.Join("data", "img", imgHandler.Filename))
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
	Memebufferchannel <- newmeme
	const homeButton = `<a href=../>Go home</a>`

	fmt.Fprintf(w, "Successfully created new meme from youtube\n %s", homeButton)
}
func GetYTsnd(ytid string) {

	filename := "temp" + RandomString()
	soundfile := exec.Command("yt-dlp", "-x", "--format=mp4", "-o", "data/snd/"+filename, ytid)
	err := soundfile.Run()
	check(err)
	soundfile2 := exec.Command("ffmpeg", "-i", "data/snd/"+filename+".m4a", "-ar", "44100", "data/snd/"+ytid+".mp3")
	err = soundfile2.Run()
	fmt.Println(err)
	check(err)
	os.Remove("data/snd/" + filename + ".m4a")

}
