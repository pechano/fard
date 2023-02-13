package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)


type meme struct{
	File string `json:"file"`
Img string `json:"img"`
Title string `json:"title"`
	buffer *beep.Buffer
}




func main() {

	//prepare memes from json:


	var fard meme
	fard.File, fard.Img, fard.Title = "fard.mp3","poop.jpg","reverb fard"

	var fart meme
	fart.File, fart.Img, fart.Title = "fart.mp3","poop.jpg","wet fart"
	b, err := json.Marshal(fard)
	fmt.Println(fard)
	fmt.Println(string(b[:]))
	var fart2 meme
	err = json.Unmarshal(b,&fart2)
fmt.Println(fart2)
var memes []meme
	memes = append(memes,fard)
	memes = append(memes,fart)

	f, err := os.Open(fard.File)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	fardrate := format.SampleRate
	f2, err := os.Open(fart.File)
	if err != nil {
		log.Fatal(err)
	}

	streamer2, format2, err := mp3.Decode(f2)
	if err != nil {
		log.Fatal(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	resampled := beep.Resample(4,format2.SampleRate,fardrate,streamer2)

	memes[0].buffer = beep.NewBuffer(format)
	memes[1].buffer = beep.NewBuffer(format2)
memes[0].buffer.Append(streamer)
memes[1].buffer.Append(resampled)
	streamer.Close()
	streamer2.Close()


	//This part reads the meme collection and presents them in menu form as a numbered list.
	//The user inputs a number with ENTER and the corresponding sound will play.
	var choice int
for {
		fmt.Printf("Pick an option and press [ENTER] to <<<FARD>>>! \n")
		for i:=0;i<len(memes);i++{
		fmt.Printf("%d: %s \n",i+1,memes[i].Title)
	}
		fmt.Scan(&choice)
		choice = choice -1
		if choice >= len(memes){fmt.Println("Please pick a valid option") } else if
		choice < 0 {fmt.Println("Please pick a valid option") } else 
		{
				fart := memes[choice].buffer.Streamer(0,memes[choice].buffer.Len() )
				speaker.Play(fart)
		}}



}



//En funktion som tager en sample-rate og et filnavn og returnerer en streamer der kan spiller igen og igen?
//Sample rate er unødvendigt, da FARD altid være være udgangspunktet.
