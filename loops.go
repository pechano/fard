package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"github.com/gorilla/mux"
)

type loop struct {
	Name         string `json:"Name"`
	Filename     string `json:"Filename"`
	Img          string `json:"Img"`
	Buffer       *beep.Buffer
	ID           int  `json:",omitempty"`
	FramePerfect bool `json:"Perfect"`
}

type LoopsData struct {
	channel chan loop
	Playing bool `json:"Playing"`
	Loops   []loop
}

func (looplist *LoopsData) sendloops(w http.ResponseWriter, r *http.Request) {
	fmt.Println("responding to getloops")
	json.NewEncoder(w).Encode(looplist.Loops)

}

func loopstopper(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Sending stop signal to loop routine")
	speaker.Clear()
}

func scanForLoops(loops *LoopsData) {
	files, err := os.ReadDir(filepath.Join("data", "loops"))
	check(err)
	var cleanlist []string
	for _, e := range files {
		if strings.HasSuffix(e.Name(), ".json") {
			cleanlist = append(cleanlist, e.Name())
		}
	}
	var rawlist []string
	for _, file := range rawlist {
		if strings.HasSuffix(file, ".json") {
			cleanlist = append(cleanlist, file)
		}
	}

	for _, file := range cleanlist {

		var loop loop
		jsonFile, err := os.Open(filepath.Join("data", "loops", file))
		check(err)

		defer jsonFile.Close()

		byteValue, _ := io.ReadAll(jsonFile)

		json.Unmarshal(byteValue, &loop)
if loop.FramePerfect == false {loops.channel <- loop ; return}
		f, err := os.Open(filepath.Join("data", "loops", loop.Filename))
		check(err)
		streamer, format, err := wav.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)

		loop.Buffer = beep.NewBuffer(format)
		loop.Buffer.Append(resampled)
		loops.channel <- loop
	}
}

func loopHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	ID, err := strconv.Atoi(key)
	if err != nil {
		fmt.Println("Error during conversion")
		return
	}
	if ID < len(LoopCollection.Loops) {

		switch LoopCollection.Loops[ID].FramePerfect {
		case true:

			fardloop := LoopCollection.Loops[ID].Buffer.Streamer(0, LoopCollection.Loops[ID].Buffer.Len())
			reloopedFard := beep.Loop(-1, fardloop)
			speaker.Play(reloopedFard)
			fmt.Printf("Endpoint Hit: %s \n", LoopCollection.Loops[ID].Name)

		case false:
			f, err := os.Open(filepath.Join("data", "loops", LoopCollection.Loops[ID].Filename))
			check(err)
			streamer, format, err := mp3.Decode(f)
			check(err)
			if format.SampleRate != 44100 {

				resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)
				speaker.Play(resampled)
				fmt.Printf("Endpoint Hit: %s \n", LoopCollection.Loops[ID].Name)
			} else {
				fmt.Printf("Endpoint Hit: %s \n", LoopCollection.Loops[ID].Name)

				speaker.Play(streamer)
			}
		}

	} else {
		fmt.Println("Out of range request made")
	}
}

func (l *LoopsData) Manager() {
	for {

		newloop := <-l.channel
		nextID := len(l.Loops)
		newloop.ID = nextID
		l.Loops = append(l.Loops, newloop)
		fmt.Println("new loop added: " + newloop.Name)

	}
}

func uploadNewLoop(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files

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
	var newloop loop
	newloop.Img = imgHandler.Filename
	newloop.Filename = sndHandler.Filename
	shortSample := r.FormValue("perfect")
	if shortSample == "yes" {
		newloop.FramePerfect = true
	} else {
		newloop.FramePerfect = false
	}

	if newloop.FramePerfect {
		if filepath.Ext(newloop.Filename) != ".wav" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "Loop not added. File format MUST be .wav with 44100 Hz sample rate\n")
			return
		}
	}

	newloop.Name = r.FormValue("title")
	truncatedname := TruncateTitle(newloop.Name)

	jsonName := truncatedname + RandomString() + ".json"

	jsonNamePath := filepath.Join("data", "loops", jsonName)
	fmt.Println("New loop submitted: ", newloop)
	file, _ := json.MarshalIndent(newloop, "", " ")
	_ = os.WriteFile(jsonNamePath, file, 0644)

	// Create file
	dst, err := os.Create(filepath.Join("data", "loops", imgHandler.Filename))
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
	dst, err = os.Create(filepath.Join("data", "loops", sndHandler.Filename))
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

	loopbufferchannel <- newloop
	const homeButton = `<a href=../>Go home</a>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "Successfully Uploaded File\n %s", homeButton)
}

func LoopBufferman(bufferchannel chan loop, loops LoopsData) {
	for {
		loopNoBuffer := <-bufferchannel

		if !loopNoBuffer.FramePerfect {
			loops.channel <- loopNoBuffer
			return
		}
		fmt.Println("Received new loop data, creating buffer")

		f, err := os.Open(filepath.Join("data", "loops", loopNoBuffer.Filename))
		fmt.Println("we made it to here")
		check(err)
		streamer, format, err := wav.Decode(f)

		check(err)
		resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)

		loopNoBuffer.Buffer = beep.NewBuffer(format)
		loopNoBuffer.Buffer.Append(resampled)
		streamer.Close()
		fmt.Println("New soundbite buffered, sending to Meme manager")
		loops.channel <- loopNoBuffer

	}
}
