package main

import (
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
    myRouter := mux.NewRouter().StrictSlash(true)

var templates = template.Must(template.ParseFiles("upload.html"))

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
	fmt.Println("Gathering Memes")
	memelist := getlist()
	fmt.Println("Preparing meme database")
	memesNoBuffer := preparememes(memelist)
	fmt.Println("Buffering meme")
	Memes := getBuffers(memesNoBuffer)
	fmt.Println("Memes ready:")
	for _, meme := range Memes{
		fmt.Println(meme.Title)}

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

	tmpl := template.Must(template.ParseFiles("index.html"))
	myRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, Memes)
	})

	myRouter.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":10000",myRouter))


	//This part reads the meme collection and presents them in menu form as a numbered list.
	//The user inputs a number with ENTER and the corresponding sound will play.
}
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
func getBuffers(memesNoBuffer []Meme)(Memes []Meme){

	f1, err := os.Open(filepath.Join("data","snd","fard.mp3"))
	if err != nil {
		log.Fatal(err)
	}

_, format, err := mp3.Decode(f1)
	if err != nil {
		log.Fatal(err)
	}
	fardrate := format.SampleRate

	speaker.Init(fardrate, fardrate.N(time.Second/10))


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
