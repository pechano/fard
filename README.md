# Fardserver

## Installing and running
The code is generally a mess of hacky spaghetti, but the server can be run by cloning the repo, navigating to the folder and running
```
go build
./fard
```
During initial startup the ipv4 address is shown in the log. If you hosting on the same machine, navigate to ```localhost:10000``` to begin.
The core functionality of sound playback and uploading should work, but most of the special functions require additional ressources that cannot be distributed through github.

## Endpoint description

The fardserver exposes several endpoints that can be used for scripting purposes. They can be reached using ```curl``` or equivalent. 

**EXAMPLE:** ```curl fardserver.lan:10000/ding```

>/refreshmemes

Returns a JSON struct of all memes currently loaded into the fardserver. Each meme has their own unique ID, which is the integer that can be used in
>/fard/{int}

This will make the server play the sound with the ID from ```/refreshmemes```
>/getloops

Works the same as ```/refreshmemes```but with returns a JSON struct with all the loops currently loaded
>/loop/{int}

Same logic as the ```/fard```endpoint. Using the ID provided by ```/getloops``` tell the server to start playing a loop.
>/stoploop

Tells the server to stop playing sounds. This is very useful when bloood starts flowing from your eardrums.

>/ding

Makes the server play a "ding" sound. It is very convenient to bind a keyboard key to this function.
