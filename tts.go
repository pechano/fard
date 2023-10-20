package main

import (
	"fmt"
	"net/http"

	"github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
)

type ttsOptions struct {
  words string
  language string

}
func tts(options ttsOptions){

speech := htgotts.Speech{Folder: "audio", Language: options.language, Handler: &handlers.Native{}}
 speech.Speak(options.words)
}

func getOptions  (w http.ResponseWriter, r *http.Request) {

var options ttsOptions	
  		options.words = r.FormValue("query")
  options.language = r.FormValue("lang")
tts(options)
  fmt.Println("Endpoint hit: TTS")

		const homeButton = `<a href=../>Go home</a><br>`
		const goAgain = `<a href=../pages/ttsmeme.html>Go again</a><br>`

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "xdlmao amirite guis????!\n %s\n %s", homeButton, goAgain)
	}
