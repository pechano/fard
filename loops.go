package main

import (
	"encoding/json"
	"net/http"
)


type loop struct{
	Name string `json:"Name"`
Playing bool `json:"Playing"`
Filename string `json:"Filename"`
}
type playerinfo struct{
		loops []loop
	}


func (looplist *playerinfo)fakeloops()(){

		loop1  := loop{Name: "metal pipe", Playing: true, Filename: "pipe.mp3"}
looplist.loops= append(looplist.loops, loop1)
		loop2  := loop{Name: "reverb fard", Playing: false, Filename: "fard.mp3"}
looplist.loops= append(looplist.loops, loop2)
}

func (looplist playerinfo) sendloops(w http.ResponseWriter, r *http.Request)(){

			json.NewEncoder(w).Encode(looplist.loops)

}

