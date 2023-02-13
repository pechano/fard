package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)


type meme struct{
file string
img string
title string
	buffer *beep.Buffer
}



func main() {
	var fard meme
	fard.file, fard.img, fard.title = "fard.mp3","poop.jpg","reverb fard"

	var fart meme
	fart.file, fart.img, fart.title = "fart.mp3","poop.jpg","wet fart"

var memes []meme
	memes = append(memes,fard)
	memes = append(memes,fart)

	f, err := os.Open(fard.file)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	fardrate := format.SampleRate
	f2, err := os.Open(fart.file)
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
	var choice int
for {
		fmt.Printf("Pick an option and press [ENTER] to <<<FARD>>>! \n")
		for i:=0;i<len(memes);i++{
		fmt.Printf("%d: %s \n",i+1,memes[i].title)
	}
		fmt.Scan(&choice)
		choice = choice -1
		if choice >= len(memes){fmt.Println("Please pick a valid option") } else if
		choice < 0 {fmt.Println("Please pick a valid option") } else if
				reflect.TypeOf(choice) !=reflect.TypeOf(1) {fmt.Println("Please pick a valid option")} else 
		{
				fart := memes[choice].buffer.Streamer(0,memes[choice].buffer.Len() )
				speaker.Play(fart)
		}}



}



//En funktion som tager en sample-rate og et filnavn og returnerer en streamer der kan spiller igen og igen?
//Sample rate er unødvendigt, da FARD altid være være udgangspunktet.
