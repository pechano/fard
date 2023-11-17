package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type status struct {
	Memes int `json:"memes"`
	Loops int `json:"loops"`
}

type Meme struct {
	SoundFile string `json:"file"`
	Img       string `json:"img"`
	Title     string `json:"title"`
	buffer    *beep.Buffer
	ID        int `json:"ID"`
}

type Memecollection struct {
	UID      int
	Memes    []Meme
	Lock     bool
	channel  chan Meme
	fardrate beep.SampleRate
}

var Memebufferchannel chan Meme
var loopbufferchannel chan loop

var collection Memecollection
var LoopCollection LoopsData

func main() {

	teststring := "halli hallo og lige et øjeblik"
	teststring = TruncateTitle(teststring)
	teststring = teststring + RandomString()
	fmt.Println(teststring)

	err := godotenv.Load()
	check(err)
	port := os.Getenv("port")
	if port == "" {
		port = "10000"
	}

	conn, error := net.Dial("udp", "8.8.8.8:80")
	if error != nil {
		fmt.Println(error)
	}
	defer conn.Close()
	ipAddress := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println("Hosting fardserver at:", ipAddress.IP, ":10000")

	//set up handleFuncs for server and restart thereof
	var status status

	collection.channel = make(chan Meme)
	Memebufferchannel = make(chan Meme)
	loopbufferchannel = make(chan loop)
	LoopCollection.channel = make(chan loop)

	f1, err := os.Open(filepath.Join("data", "snd", "fard.mp3"))
	if err != nil {
		log.Fatal(err)
	}

	_, format, err := mp3.Decode(f1)
	if err != nil {
		log.Fatal(err)
	}
	fardrate := format.SampleRate

	speaker.Init(fardrate, fardrate.N(time.Second/10))
	collection.fardrate = fardrate

	fmt.Println("Looking for new memes")
	discoverMemes()
	go collection.Manager()

	go scanForMemes(collection)
	go bufferman(Memebufferchannel, collection)

	go LoopBufferman(loopbufferchannel, LoopCollection)
	go LoopCollection.Manager()

	go scanForLoops(&LoopCollection)

	//set up the main router and the handler for "/shutdown", which will restart the server.
	myRouter := mux.NewRouter().StrictSlash(true)

	myServer := http.Server{Addr: ":" + port, Handler: myRouter}
	myRouter.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK")) // Write response body
		if err := myServer.Close(); err != nil {
			log.Fatal(err)
		}
	})

	//read files to create the meme collection

	fard := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["id"]
		ID, err := strconv.Atoi(key)
		if err != nil {
			fmt.Println("Error during conversion")
			return
		}
		if ID < len(collection.Memes) {
			fart := collection.Memes[ID].buffer.Streamer(0, collection.Memes[ID].buffer.Len())
			speaker.Play(fart)

			fmt.Printf("Endpoint Hit: %s \n", collection.Memes[ID].Title)
		} else {
			fmt.Println("Out of range request made")
		}
	}

	var PoolTemp string
	PoolTemp = "11.1"
	getTemp := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["temp"]
		*&PoolTemp = key
		fmt.Println(PoolTemp) //for testing purposes, can be removed
	}
	myRouter.HandleFunc("/logger/{temp}", getTemp)

	pooltemplate := template.Must(template.ParseFiles("./pages/temp.html"))
	myRouter.HandleFunc("/temp", func(w http.ResponseWriter, r *http.Request) {
		err := pooltemplate.Execute(w, PoolTemp)
		check(err)
	})

	fileserver := http.FileServer(http.Dir("./data"))
	pageserver := http.FileServer(http.Dir("./pages"))
	loopserver := http.FileServer(http.Dir("./loops"))

	myRouter.PathPrefix("/data").Handler(http.StripPrefix("/data", fileserver))
	myRouter.PathPrefix("/pages").Handler(http.StripPrefix("/pages", pageserver))
	myRouter.PathPrefix("/loops").Handler(http.StripPrefix("/loops", loopserver))

	myRouter.HandleFunc("/fard/{id}", fard)
	myRouter.HandleFunc("/soren/", sorenHandler)

	myRouter.HandleFunc("/ding/", dingHandler)

	myRouter.HandleFunc("/tts", getOptions)

	myRouter.HandleFunc("/loop/{id}", loopHandler)
	myRouter.HandleFunc("/stoploop", loopstopper)
	myRouter.HandleFunc("/getloops", LoopCollection.sendloops)

	statushandler := func(w http.ResponseWriter, r *http.Request) {

		status.Memes = len(collection.Memes)
		status.Loops = len(LoopCollection.Loops)
		json.NewEncoder(w).Encode(status)
	}

	myRouter.HandleFunc("/status", statushandler)

	refreshhandler := func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(collection.Memes)
	}

	myRouter.HandleFunc("/refreshmemes", refreshhandler)

	filterhandler := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		fmt.Println(vars)
		key := vars["term"]
		term := key
		check(err)
		FilteredMemes := filterMemesFromJS(collection.Memes, term)
		json.NewEncoder(w).Encode(FilteredMemes)
	}

	myRouter.HandleFunc("/filtermemes/{term}", filterhandler)
	//Fill in the main page template
	tmpl := template.Must(template.ParseFiles("./pages/index.html"))
	myRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, collection.Memes)
		check(err)
	})

	myRouter.HandleFunc("/upload", uploadHandler)

	myRouter.HandleFunc("/uploadmeme", uploadmemeHandler)

	myRouter.HandleFunc("/uploadloop", uploadloopHandler)

	myRouter.HandleFunc("/YTmeme", YTmemehandler)

	myRouter.HandleFunc("/wgetmeme", Wgetmemehandler)

	myRouter.HandleFunc("/memebuilder", func(w http.ResponseWriter, r *http.Request) {
		details := Meme{
			Title:     r.FormValue("title"),
			SoundFile: r.FormValue("file"),
			Img:       r.FormValue("img"),
		}
		fmt.Println(details)
	})

	myRouter.HandleFunc("/newmeme", func(w http.ResponseWriter, r *http.Request) {

		var newmeme = template.Must(template.ParseFiles("./pages/newmeme.html"))
		err := newmeme.Execute(w, nil)
		check(err)
	})

	if err := myServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Printf("Finished")

}

// This part reads the meme collection and presents them in menu form as a numbered list.
// The user inputs a number with ENTER and the corresponding sound will play.
func showMenu(Memes []Meme) {
	var choice int
	for {
		fmt.Printf("Pick an option and press [ENTER] to <<<FARD>>>! \n")
		for i := 0; i < len(Memes); i++ {
			fmt.Printf("%d: %s \n", i+1, Memes[i].Title)
		}
		fmt.Scan(&choice)
		choice = choice - 1
		if choice >= len(Memes) {
			fmt.Println("Please pick a valid option")
		} else if choice < 0 {
			fmt.Println("Please pick a valid option")
		} else {
			fart := Memes[choice].buffer.Streamer(0, Memes[choice].buffer.Len())
			speaker.Play(fart)
		}
	}

}

// en funktion som kigger efter nye .zip-filer og flytter deres indhold til de rigtige mapper.
// lige pt. bestående af en masser dårlig boilerplate
func discoverMemes() {

	files, err := os.ReadDir(".")
	check(err)
	var cleanlist []string
	for _, e := range files {
		if strings.HasSuffix(e.Name(), ".zip") {
			cleanlist = append(cleanlist, e.Name())
		}
	}
	for _, archive := range cleanlist {

		boosterPack, err := zip.OpenReader(archive)
		if err != nil {
			log.Print(err.Error())
		}
		defer boosterPack.Close()

		for _, f := range boosterPack.File {
			//

			fileExt := filepath.Ext(f.Name)
			switch fileExt {
			case ".json":
				dataPath := filepath.Join("data", f.Name)
				outFile, err := os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					if err != nil {
						log.Println(err.Error())
					}
				}
				fileInArchive, err := f.Open()
				if err != nil {
					if err != nil {
						log.Println(err.Error())
					}
				}
				if _, err := io.Copy(outFile, fileInArchive); err != nil {
					if err != nil {
						log.Println(err.Error())
					}
				}

				outFile.Close()
				fileInArchive.Close()
				os.Remove(f.Name)

			case ".jpg", ".jpeg", ".gif", ".bmp":
				dataPath := filepath.Join("data", "img", f.Name)
				outFile, err := os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())

				if err != nil {
					if err != nil {
						log.Println(err.Error())
					}
				}
				fileInArchive, err := f.Open()
				if err != nil {
					if err != nil {
						log.Println(err.Error())
					}
				}
				if _, err := io.Copy(outFile, fileInArchive); err != nil {
					if err != nil {
						log.Println(err.Error())
					}
				}

				outFile.Close()
				fileInArchive.Close()
				os.Remove(f.Name)
			case ".mp3":

				dataPath := filepath.Join("data", "snd", f.Name)
				outFile, err := os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())

				if err != nil {
					if err != nil {
						log.Println(err.Error())
					}
				}
				fileInArchive, err := f.Open()
				if err != nil {
					if err != nil {
						log.Println(err.Error())
					}
				}
				if _, err := io.Copy(outFile, fileInArchive); err != nil {
					if err != nil {
						log.Println(err.Error())
					}
				}

				outFile.Close()
				fileInArchive.Close()
				os.Remove(f.Name)
			}

		}

		os.Remove(archive)
	}

	fmt.Println("Files extracted")
}

func filterMemesFromJS(memes []Meme, input string) (filteredmemes []Meme) {

	inputSlice := strings.Split(input, " ")
	resultsize := len(inputSlice)
	result := make([][]Meme, resultsize+1)
	temp := memes

	result[0] = temp
	for i := 0; i <= resultsize-1; i++ {

		var temp []Meme
		for _, meme := range result[i] {
			if strings.Contains(meme.Title, inputSlice[i]) {
				temp = append(temp, meme)
			}
		}
		result[i+1] = temp
	}
	return result[resultsize]
}

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}

func savememes() {
	var fard Meme
	file, _ := json.MarshalIndent(fard, "", " ")

	_ = os.WriteFile("test.json", file, 0644)
	b, err := json.Marshal(fard)
	check(err)
	fmt.Println(fard)
	fmt.Println(string(b[:]))
	var fart2 Meme
	err = json.Unmarshal(b, &fart2)
	check(err)
	fmt.Println(fart2)
}

func (m *Memecollection) Manager() {

	for {

		newmeme := <-m.channel
		nextID := len(m.Memes)
		newmeme.ID = nextID
		m.Memes = append(m.Memes, newmeme)
		fmt.Println("new meme added: " + newmeme.Title + " with ID: " + fmt.Sprint(newmeme.ID))

	}
}

//Big function to scan the data folder for .json files, create buffers, and send the Meme struct off to the .Manager function.

func scanForMemes(collection Memecollection) {
	files, err := os.ReadDir(filepath.Join("data"))
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

		var fard Meme
		jsonFile, err := os.Open(filepath.Join("data", file))
		check(err)

		defer jsonFile.Close()

		byteValue, _ := io.ReadAll(jsonFile)

		json.Unmarshal(byteValue, &fard)
		fmt.Println("Preparing: " + fard.Title)

		f, err := os.Open(filepath.Join("data", "snd", fard.SoundFile))

		if err != nil {
			os.Remove(filepath.Join("data", file))
			fmt.Println("Removing faulty json file: " + file)
			continue
		}

		streamer, format, err := mp3.Decode(f)
		check(err)
		if format.SampleRate != collection.fardrate {
			resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)

			fard.buffer = beep.NewBuffer(format)
			fard.buffer.Append(resampled)
			streamer.Close()
			collection.channel <- fard
		} else {

			fard.buffer = beep.NewBuffer(format)
			fard.buffer.Append(streamer)
			streamer.Close()
			collection.channel <- fard

		}

	}

}

func bufferman(bufferchannel chan Meme, collection Memecollection) {
	for {
		memeNoBuffer := <-bufferchannel
		fmt.Println("Received new meme, creating buffer")

		f, err := os.Open(filepath.Join("data", "snd", memeNoBuffer.SoundFile))
		check(err)
		streamer, format, err := mp3.Decode(f)

		check(err)
		resampled := beep.Resample(4, format.SampleRate, collection.fardrate, streamer)

		memeNoBuffer.buffer = beep.NewBuffer(format)
		memeNoBuffer.buffer.Append(resampled)
		streamer.Close()
		fmt.Println("New soundbite buffered, sending to Meme manager")
		collection.channel <- memeNoBuffer

	}
}

func RandomString() (output string) {
	possibles := []rune("abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ")
	runes := []rune("")
	for i := 0; i < 5; i++ {
		random := rand.Intn(len(possibles) - 1)
		runes = append(runes, possibles[random])
	}
	return string(runes)
}

func TruncateTitle(input string) (output string) {
	m1 := regexp.MustCompile("( .*)")
	output = m1.ReplaceAllString(input, "")
	return output
}
