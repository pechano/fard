package main

import (
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"github.com/gorilla/mux"
)

func sorenHandler(w http.ResponseWriter, r *http.Request) {
	speaker.Clear()
	f, err := os.Open(filepath.Join("data", "special", "glidefedt.mp3"))
	check(err)
	streamer, format, err := mp3.Decode(f)
	if format.SampleRate == 44100 {
		fmt.Println("running søren.exe")
	}
	check(err)
	speaker.Play(streamer)
}
func hlhandler(w http.ResponseWriter, r *http.Request) {
	HLchannel <- "freeman"
}
func hlsfx(c chan string, list []string) {

	for {
		fmt.Println(<-c)
		sound := list[rand.Intn(len(list))]
		f, err := os.Open(filepath.Join(sound))
		check(err)
		streamer, format, err := wav.Decode(f)
		check(err)
		resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)
		speaker.Play(resampled)
	}
}
func tomhandler(w http.ResponseWriter, r *http.Request) {
	tomchannel <- "tom scot moment"
}
func tomsfx(c chan string, list []string) {

	for {
		fmt.Println(<-c)
		sound := list[rand.Intn(len(list))]
		f, err := os.Open(filepath.Join(sound))
		check(err)
		streamer, format, err := mp3.Decode(f)
		check(err)
		resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)
		speaker.Play(resampled)
	}
}
func badlandshandler(w http.ResponseWriter, r *http.Request) {
	badlandschannel <- "chug"
}
func badlandshit(c chan string, list []string) {

	for {
		fmt.Println(<-c)
		sound := list[rand.Intn(len(list))]
		f, err := os.Open(filepath.Join(sound))
		check(err)
		streamer, format, err := mp3.Decode(f)
		check(err)
		resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)
		speaker.Play(resampled)
	}
}

func duckhandler(w http.ResponseWriter, r *http.Request) {
	duckchannel <- "quack?"
}

func quacker(c chan string) {
	f, err := os.Open(filepath.Join("data", "special", "quack.mp3"))
	check(err)
	streamer, format, err := mp3.Decode(f)
	check(err)
	quackbuffer := beep.NewBuffer(format)
	quackbuffer.Append(streamer)
	streamer.Close()

	f, err = os.Open(filepath.Join("data", "special", "quack2.mp3"))
	check(err)
	streamer, format, err = mp3.Decode(f)
	check(err)
	quackbuffer2 := beep.NewBuffer(format)
	quackbuffer2.Append(streamer)
	streamer.Close()

	f, err = os.Open(filepath.Join("data", "special", "honk.mp3"))
	check(err)
	streamer, format, err = mp3.Decode(f)
	check(err)
	honkbuffer := beep.NewBuffer(format)
	honkbuffer.Append(streamer)
	streamer.Close()
	for {
		fmt.Println(<-c)
		random := rand.Intn(10)
		switch random {
		case 9:
			quack := quackbuffer.Streamer(0, quackbuffer.Len())
			speaker.Play(quack)
		case 8:
			honk := honkbuffer.Streamer(0, honkbuffer.Len())
			speaker.Play(honk)
		default:
			quack2 := quackbuffer2.Streamer(0, quackbuffer2.Len())
			speaker.Play(quack2)
		}
	}
}

func dingHandler(w http.ResponseWriter, r *http.Request) {
	DingChannel <- "dong"
}

func dingding(c chan string) {
	f, err := os.Open(filepath.Join("data", "special", "ding2.mp3"))
	check(err)
	streamer, format, err := mp3.Decode(f)
	check(err)
	dingbuffer := beep.NewBuffer(format)
	dingbuffer.Append(streamer)
	streamer.Close()
	for {
		fmt.Println(<-c)
		ding := dingbuffer.Streamer(0, dingbuffer.Len())
		speaker.Play(ding)
	}
}

func ListFiles(ext string, path string) ([]string, error) {
	var fileList []string
	// Read all the file recursively
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	err := filepath.Walk(path, func(file string, f os.FileInfo, err error) error {
		if filepath.Ext(file) == ext {
			fileList = append(fileList, file)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileList, nil
}

func ScanForHL() {

	err := filepath.WalkDir(filepath.Join("data", "random", "HL"), walk)
	check(err)
}

func walk(s string, d fs.DirEntry, err error) error {
	var list []string
	if err != nil {
		return err
	}
	if !d.IsDir() {
		list = append(list, s)
	}
	return nil
}

func doghandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["distance"]
	ID, err := strconv.Atoi(key)

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}
	dogchannel <- ID
	fmt.Println("barking from " + fmt.Sprint(ID) + " meters")
}

func distantDog(d chan int) {

	for {
		distance := <-d
		switch distance {
		case 0:
			f, err := os.Open(filepath.Join("data", "special", "0m.mp3"))
			check(err)
			streamer, format, err := mp3.Decode(f)
			check(err)
			resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)
			speaker.Play(resampled)
		case 250:
			f, err := os.Open(filepath.Join("data", "special", "250m.mp3"))
			check(err)
			streamer, format, err := mp3.Decode(f)
			check(err)
			resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)
			speaker.Play(resampled)
		case 500:
			f, err := os.Open(filepath.Join("data", "special", "500m.mp3"))
			check(err)
			streamer, format, err := mp3.Decode(f)
			check(err)
			resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)
			speaker.Play(resampled)
		case 750:
			f, err := os.Open(filepath.Join("data", "special", "750m.mp3"))
			check(err)
			streamer, format, err := mp3.Decode(f)
			check(err)
			resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)
			speaker.Play(resampled)
		case 1000:
			f, err := os.Open(filepath.Join("data", "special", "1000m.mp3"))
			check(err)
			streamer, format, err := mp3.Decode(f)
			check(err)
			resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)
			speaker.Play(resampled)
		default:
			continue
		}

	}
}
