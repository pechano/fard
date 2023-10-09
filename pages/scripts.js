"use strict";
var torpedo;
function fardFunc(id) {
    var url = "./fard/";
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url + id, true);
    xhr.send();
}
function browserFard(soundFile) {
    let cleanText = soundFile.replace(/\"/g, "");
    var basePath = "/data/snd/";
    var path = basePath.concat(cleanText);
    var meme = new Audio(path);
    meme.play();
}
function reFard() {
    var url = "./shutdown";
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.send();
    setTimeout(function () {
        window.location.reload();
    }, 2000);
}
async function getStatus() {
    const url = "./status";
    const response = await fetch(url);
    const jsonData = await response.json();
    const currentMemeStatus = jsonData.current_memes;
    const current = document.getElementById("current");
    const p = document.createElement("current_status");
    p.textContent = currentMemeStatus;
    current?.replaceWith(p);
    const nextMemeStatus = jsonData.next_memes;
    const next = document.getElementById("next");
    const pp = document.createElement("next_status");
    pp.textContent = nextMemeStatus;
    next?.replaceWith(pp);
}
async function getMemes() {
    const url = "./refreshmemes";
    const response = await fetch(url);
    const jsonData = await response.json();
    var maymay = document.getElementById("maymay");
    var temp = document.createElement("maymay");
    temp.setAttribute("id", "maymay");
    for (var object of jsonData) {
        console.log(object);
        if (object.ID == undefined) {
            object.ID = "0";
        }
        ; // the first entry requires a bit of extra work for some reason
        const title = object.title;
        const pp = document.createElement("div");
        pp.textContent = title;
        pp.setAttribute("class", "container");
        pp.innerHTML = "<img src=/data/img/" + object.img + " alt='poop' style='width:100%'>" +
            "<button class='btn' onclick='fardFunc(" + object.ID + ")'>" + object.title + "</button>" +
            "<button class='btn2' onclick=browserFard('" + object.file + "')><i class='fa-solid fa-headphones'></i></button>";
        temp?.appendChild(pp);
    }
    temp?.setAttribute("class", "wrapper");
    maymay?.replaceWith(temp);
}
async function getFilteredMemes(value) {
    if (value == "") {
        getMemes();
        return;
    }
    ;
    const url = "./filtermemes/";
    var term = value;
    const response = await fetch(url + term);
    const jsonData = await response.json();
    var temp = document.createElement("div");
    temp.setAttribute("id", "maymay");
    for (var object of jsonData) {
        if (object.ID == undefined) {
            object.ID = "0";
        }
        ; // the first entry requires a bit of extra work for some reaso
        globalThis.torpedo = jsonData[0].ID;
        const title = object.title;
        const pp = document.createElement("div");
        pp.textContent = title;
        pp.innerHTML = "<div class='container'>" +
            "<img src=/data/img/" + object.img + " alt='poop' style='width:100%'>" +
            "<button class='btn' onclick='fardFunc(" + object.ID + ")'>" + object.title + "</button>" +
            "<button class='btn2' onclick=browserFard('" + object.file + "')><i class='fa-solid fa-headphones'></i></button>" +
            "</div>";
        temp?.appendChild(pp);
    }
    temp?.setAttribute("class", "wrapper");
    var maymay = document.getElementById("maymay");
    maymay?.replaceWith(temp);
}