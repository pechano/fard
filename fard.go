package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
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
)


type Meme struct{
	SoundFile string `json:"file"`
Img string `json:"img"`
Title string `json:"title"`
	buffer *beep.Buffer
ID int `json:"entry"`
}

func homePage (w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Welcome to meme heaven")
	fmt.Println("Endpoint Hit: homePage")
}

type Todo struct {
    Title string
    Done  bool
}



type TodoPageData struct {
    PageTitle string
    Todos     []Todo
}


func main() {
//Load in template related to uploads
var templates = template.Must(template.ParseFiles("newmeme.html"))

// Display the named template
	display := func (w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}
	uploadFile :=func (w http.ResponseWriter, r *http.Request) {
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

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

	uploadHandler := func (w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}
	//set up handleFuncs for server and restart thereof
	//and initiate the loop that will allow for restarts of the server once the /shutdown endpoint is hit

 	cycle := 0

for {

   myRouter := mux.NewRouter().StrictSlash(true)

	myServer := http.Server{Addr: ":10000", Handler: myRouter}

myRouter.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))    // Write response body
        if err := myServer.Close(); err != nil {
            log.Fatal(err)
        }
    })


		if cycle > 0 {fmt.Println("Current reboots:",cycle)}
		
//read files to create the meme collection
		fmt.Println("Looking for new memes")
discoverMemes()
	fmt.Println("Gathering Memes")
	memelist := getlist()
	fmt.Println("Preparing meme database")
	memesNoBuffer := preparememes(memelist)
	fmt.Println("Buffering meme")
	Memes := getBuffers(memesNoBuffer,cycle)
		cycle++
	fmt.Println("Memes ready:")
	for _, meme := range Memes{
		fmt.Println(meme.Title)
		}

	returnAllMemes := func (w http.ResponseWriter, r *http.Request)(){
		fmt.Println("Endpoint Hit: returnAllMemes")
		json.NewEncoder(w).Encode(Memes)
	}
	fard := func (w http.ResponseWriter, r *http.Request)(){
		fmt.Println("Endpoint Hit: fard")
		vars := mux.Vars(r)
		key := vars["id"]
		ID, err := strconv.Atoi(key)
		if err != nil {fmt.Println("Error during conversion")
			return}
if ID < len(Memes) {
		fart := Memes[ID].buffer.Streamer(0,Memes[ID].buffer.Len() )
		speaker.Play(fart)
		} else {fmt.Println("Out of range request made")}
	}

	fileserver := http.FileServer(http.Dir("./data"))
myRouter.PathPrefix("/data").Handler(http.StripPrefix("/data",fileserver))
	myRouter.HandleFunc("/Memes", returnAllMemes)
	myRouter.HandleFunc("/fard/{id}", fard)
//Fill in the main page template
	tmpl := template.Must(template.ParseFiles("index.html"))
	myRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, Memes)
	})

	myRouter.HandleFunc("/upload", uploadHandler)

	    if err := myServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
    log.Printf("Finished")
	}
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
		dataPath := filepath.Join("data",f.Name)
		sndPath := filepath.Join("data","snd",f.Name)
		imgPath := filepath.Join("data","img",f.Name)
		//
			if strings.HasSuffix(f.Name,".json") { 
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
			}
	
			if strings.HasSuffix(f.Name,".mp3") { 
		outFile, err := os.OpenFile(sndPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()) 

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
			}
			if strings.HasSuffix(f.Name,".jpg") { 
		outFile, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()) 

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
			}
os.Remove(f.Name)
	}
os.Remove(archive)	
	fmt.Println("Files extracted")
	}
}


//En funktion som tager en sample-rate og et filnavn og returnerer en streamer der kan spiller igen og igen?
//Sample rate er unødvendigt, da FARD altid være være udgangspunktet.

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
	panic(e)
    }
}

func preparememes(files []string)(Memes []Meme){

	ID := 0
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
	f, err := os.Open(filepath.Join("data","snd",meme.SoundFile))
		check(err)
	streamer, format, err := mp3.Decode(f)

		check(err)
	resampled := beep.Resample(4,format.SampleRate,fardrate,streamer)

	meme.buffer = beep.NewBuffer(format)
meme.buffer.Append(resampled)
	streamer.Close()
		Memes=append(Memes, meme)
}

	return Memes
}
