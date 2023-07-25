package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type status struct{
	Memes int `json:"current_memes"`
	NewMemes int `json:"next_memes"`
}

type Meme struct{
	SoundFile string `json:"file"`
	Img string `json:"img"`
	Title string `json:"title"`
	buffer *beep.Buffer
	ID int `json:",omitempty"`
}


type Memecollection struct{
	UID int
	Memes []Meme
	Lock bool
}



func main() {
	err := godotenv.Load()
	check(err)
	port := os.Getenv("port")
if port == "" {port = "10000"}


	conn, error := net.Dial("udp", "8.8.8.8:80")
	if error != nil {
		fmt.Println(error)
	}

	defer conn.Close()
	ipAddress := conn.LocalAddr().( * net.UDPAddr)
	fmt.Println("Hosting fardserver at:",ipAddress.IP,":10000")
	//Load in template related to uploads


	//set up handleFuncs for server and restart thereof
	//and initiate the loop that will allow for restarts of the server once the /shutdown endpoint is hit





	cycle := 0
	var Oldlist []string

	var Oldmemes []Meme
	var status status


	for {
		//set up the main router and the handler for "/shutdown", which will restart the server.
		myRouter := mux.NewRouter().StrictSlash(true)

		myServer := http.Server{Addr: ":"+port, Handler: myRouter}

		myRouter.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))    // Write response body
			if err := myServer.Close(); err != nil {
				log.Fatal(err)
			}
		})


		if cycle > 0 {fmt.Println("Current reloads:",cycle)}

		//read files to create the meme collection
		fmt.Println("Looking for new memes")
		discoverMemes()
		fmt.Println("Gathering Memes")
		memelist := getlist()
		memelist = filterMemes(Oldlist,memelist)

		fmt.Println("Preparing meme database")
		memesNoBuffer := preparememes(memelist, Oldmemes )
		fmt.Println("Buffering memes")
		Memes := getBuffers(memesNoBuffer,cycle)
		cycle++
		fmt.Println("Memes ready:")
		for _, meme := range Memes{
			fmt.Println(meme.Title)
		}
		Oldlist = getlist()

status.Memes = len(Memes)
status.NewMemes = len(Memes)
	
		fard := func (w http.ResponseWriter, r *http.Request)(){
			vars := mux.Vars(r)
			key := vars["id"]
			ID, err := strconv.Atoi(key)
			if err != nil {fmt.Println("Error during conversion")
			return}
			if ID < len(Memes) {
				fart := Memes[ID].buffer.Streamer(0,Memes[ID].buffer.Len() )
				speaker.Play(fart)

			fmt.Printf("Endpoint Hit: %s \n",Memes[ID].Title)
			} else {fmt.Println("Out of range request made")}
		}



var PoolTemp string
		PoolTemp = "11.1"
		getTemp := func (w http.ResponseWriter, r *http.Request)(){
			vars := mux.Vars(r)
			key := vars["temp"]
*&PoolTemp = key
			fmt.Println(PoolTemp)//for testing purposes, can be removed
		}
		myRouter.HandleFunc("/logger/{temp}", getTemp)


		pooltemplate := template.Must(template.ParseFiles("./pages/temp.html"))
		myRouter.HandleFunc("/temp", func(w http.ResponseWriter, r *http.Request) {
			err :=	pooltemplate.Execute(w, PoolTemp)
			check(err)
		})


		fileserver := http.FileServer(http.Dir("./data"))
		pageserver := http.FileServer(http.Dir("./pages"))

		myRouter.PathPrefix("/data").Handler(http.StripPrefix("/data",fileserver))
		myRouter.PathPrefix("/pages").Handler(http.StripPrefix("/pages",pageserver))

		myRouter.HandleFunc("/fard/{id}", fard)
		myRouter.HandleFunc("/tts", getOptions)
		
		
		statushandler := func(w http.ResponseWriter, r *http.Request)(){ 
			json.NewEncoder(w).Encode(status)
		}

myRouter.HandleFunc("/status", statushandler)


		refreshhandler := func(w http.ResponseWriter, r *http.Request)(){ 
			json.NewEncoder(w).Encode(Memes)
		}


		myRouter.HandleFunc("/refreshmemes", refreshhandler)
		
		filterhandler := func(w http.ResponseWriter, r *http.Request)(){ 

			vars := mux.Vars(r)
			fmt.Println(vars)
			key := vars["term"]
			term := key
			check(err)
			FilteredMemes := filterMemesFromJS(Memes, term)
			json.NewEncoder(w).Encode(FilteredMemes)
		}

		//generates two fake loops for testing purposes
		var looplist playerinfo 	
		looplist.fakeloops()


		myRouter.HandleFunc("/loops",looplist.sendloops)

		myRouter.HandleFunc("/filtermemes/{term}", filterhandler)
		//Fill in the main page template
		tmpl := template.Must(template.ParseFiles("./pages/index.html"))
		myRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			err :=	tmpl.Execute(w, Memes)
			check(err)
		})

		reloadhandler := func(w http.ResponseWriter, r *http.Request){
			MemesReloaded := reloadmemes(Memes,Oldlist,1)
			status.NewMemes = len(MemesReloaded)
		}

myRouter.HandleFunc("/reload", reloadhandler)

		myRouter.HandleFunc("/upload", uploadHandler)

		myRouter.HandleFunc("/uploadmeme", uploadmemeHandler)
		myRouter.HandleFunc("/YTmeme", YTmemehandler)

		myRouter.HandleFunc("/wgetmeme", Wgetmemehandler)
		myRouter.HandleFunc("/builder", buildMeme).Methods("POST")
		//initiate the builder template

		myRouter.HandleFunc("/memebuilder", func(w http.ResponseWriter, r *http.Request){
			details := Meme{
				Title: r.FormValue("title"),
				SoundFile: r.FormValue("file"),
				Img: r.FormValue("img"),
			}
			fmt.Println(details)
		})

		myRouter.HandleFunc("/newmeme",func(w http.ResponseWriter, r *http.Request){

	var newmeme = template.Must(template.ParseFiles("./pages/newmeme.html"))
			err := newmeme.Execute(w,nil)
			check(err)
		})


		if err := myServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		Oldmemes = Memes
		log.Printf("Finished")
	}


}
	  func reloadmemes(currentMemes []Meme, Oldlist []string, cycle int) (MemesReloaded []Meme) {
discoverMemes()
		memelist := getlist()
		memelist = filterMemes(Oldlist,memelist)

		fmt.Println("Preparing meme database")
		memesNoBuffer := preparememes(memelist, currentMemes )
		fmt.Println("Buffering memes")
		MemesReloaded = getBuffers(memesNoBuffer,cycle)
return MemesReloaded
	}

//This part reads the meme collection and presents them in menu form as a numbered list.
//The user inputs a number with ENTER and the corresponding sound will play.
func showMenu(Memes []Meme)(){
	var choice int
	for {
		fmt.Printf("Pick an option and press [ENTER] to <<<FARD>>>! \n")
		for i:=0;i<len(Memes);i++{
			fmt.Printf("%d: %s \n",i+1,Memes[i].Title)
		}
		fmt.Scan(&choice)
		choice = choice -1
		if choice >= len(Memes){fmt.Println("Please pick a valid option") } else if
		choice < 0 {fmt.Println("Please pick a valid option") } else 
		{
			fart := Memes[choice].buffer.Streamer(0,Memes[choice].buffer.Len() )
			speaker.Play(fart)
	}}



}

//en funktion som kigger efter nye .zip-filer og flytter deres indhold til de rigtige mapper.
//lige pt. bestående af en masser dårlig boilerplate
func discoverMemes(){

	files,err := os.ReadDir(".")
	check(err)
	var cleanlist []string
	for _, e := range files {
		if strings.HasSuffix(e.Name(), ".zip") {
			cleanlist = append(cleanlist, e.Name())}
	}
	for _, archive := range cleanlist{

		boosterPack,err := zip.OpenReader(archive)
		if err != nil {log.Print(err.Error())}
		defer	boosterPack.Close()

		for _, f := range boosterPack.File {
			//

			fileExt := filepath.Ext(f.Name)
			switch fileExt {
			case ".json":
				dataPath := filepath.Join("data",f.Name)
				outFile, err := os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()) 
				if err != nil {
					if err != nil {log.Println(err.Error())}
				}
				fileInArchive, err := f.Open()
				if err != nil {
					if err != nil {log.Println(err.Error())}
				}
				if _, err := io.Copy(outFile, fileInArchive); err != nil {
					if err != nil {log.Println(err.Error())}
				}

				outFile.Close()
				fileInArchive.Close()	
				os.Remove(f.Name)


			case ".jpg",".jpeg",".gif",".bmp":
				dataPath := filepath.Join("data","img",f.Name)
				outFile, err := os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()) 

				if err != nil {
					if err != nil {log.Println(err.Error())}
				}
				fileInArchive, err := f.Open()
				if err != nil {
					if err != nil {log.Println(err.Error())}
				}
				if _, err := io.Copy(outFile, fileInArchive); err != nil {
					if err != nil {log.Println(err.Error())}
				}

				outFile.Close()
				fileInArchive.Close()	
				os.Remove(f.Name)
			case ".mp3":

				dataPath := filepath.Join("data","snd",f.Name)
				outFile, err := os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()) 


				if err != nil {
					if err != nil {log.Println(err.Error())}
				}
				fileInArchive, err := f.Open()
				if err != nil {
					if err != nil {log.Println(err.Error())}
				}
				if _, err := io.Copy(outFile, fileInArchive); err != nil {
					if err != nil {log.Println(err.Error())}
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



//En funktion som tager en sample-rate og et filnavn og returnerer en streamer der kan spiller igen og igen?
//Sample rate er unødvendigt, da FARD altid være være udgangspunktet.
func filterMemesFromJS(memes []Meme, input string)(filteredmemes []Meme){

	for _, meme := range memes {
if strings.Contains(meme.Title,input){
			filteredmemes = append(filteredmemes, meme)
		}
	}
return filteredmemes
}


func getlist()(cleanlist []string){

	files,err := os.ReadDir(filepath.Join("data"))
	check(err)

	for _, e := range files {
		if strings.HasSuffix(e.Name(), ".json") {
			cleanlist = append(cleanlist, e.Name())}
	}
	var rawlist []string
	for _, file := range rawlist{
		if strings.HasSuffix(file, ".json") {
			cleanlist = append(cleanlist, file)}

	}
	return cleanlist
}


func check(e error) {
	if e != nil {
		log.Println(e)
	}
}

func preparememes(files []string, Oldmemes []Meme)(Memes []Meme){

	ID := len(Oldmemes)
	if len(Oldmemes) == 0 {
		for _, file := range files{

			jsonFile, err := os.Open(filepath.Join("data",file))
			check(err)

			defer jsonFile.Close()

			byteValue, _ := ioutil.ReadAll(jsonFile)
			var fard Meme
			fard.ID = ID
			json.Unmarshal(byteValue, &fard)
			Memes = append(Memes,fard)
			ID = ID + 1
		}

		return Memes
	} else { 
		Memes = Oldmemes
		for _, file := range files{

			jsonFile, err := os.Open(filepath.Join("data",file))
			check(err)

			defer jsonFile.Close()

			byteValue, _ := ioutil.ReadAll(jsonFile)
			var fard Meme
			fard.ID = ID
			json.Unmarshal(byteValue, &fard)
			Memes = append(Memes,fard)
			ID = ID + 1
		}

		return
	}

}

func filterMemes (oldfiles []string, files []string)(newfiles []string){


	difference := make([]string, 0) //create difference slice to store the difference of two slices
	// Iterate over slice1
	for _, val1 := range files { //nested for loop to check if two values are equal
		found := false
		// Iterate over slice2
		for _, val2 := range oldfiles {
			if val1 == val2 {
				found = true
				break
			}
		}
		if !found {
			difference = append(difference, val1)
		}
	}
	return difference

}

func savememes(){
	var fard Meme
	file, _ := json.MarshalIndent(fard, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)
	b, err := json.Marshal(fard)
	check(err)
	fmt.Println(fard)
	fmt.Println(string(b[:]))
	var fart2 Meme
	err = json.Unmarshal(b,&fart2)
	fmt.Println(fart2)
}


func getBuffers(memesNoBuffer []Meme, cycle int)(Memes []Meme){

	f1, err := os.Open(filepath.Join("data","snd","fard.mp3"))
	if err != nil {
		log.Fatal(err)
	}

	_, format, err := mp3.Decode(f1)
	if err != nil {
		log.Fatal(err)
	}
	fardrate := format.SampleRate

	if cycle == 0 {	speaker.Init(fardrate, fardrate.N(time.Second/10))}


	for _, meme := range memesNoBuffer {
		if meme.buffer != nil {fmt.Println("buffer found, skipping")
			Memes=append(Memes, meme)
		continue} else {


			f, err := os.Open(filepath.Join("data","snd",meme.SoundFile))
			check(err)
			streamer, format, err := mp3.Decode(f)

			check(err)
			resampled := beep.Resample(4,format.SampleRate,fardrate,streamer)

			meme.buffer = beep.NewBuffer(format)
			meme.buffer.Append(resampled)
			streamer.Close()
			Memes=append(Memes, meme)
	}}

	return Memes
}

