"use strict";
async function getLoops() {
    const url = "../getloops";
    const response = await fetch(url);
    const jsonData = await response.json();
    var maymay = document.getElementById("loopsm8");
    var temp = document.createElement("loopsm8");
    temp.setAttribute("id", "loopsm8");
    for (var object of jsonData) {
        console.log(object);
        if (object.ID == undefined) {
            object.ID = "0";
        }
        ; // the first entry requires a bit of extra work for some reason
        const title = object.title;
        const pp = document.createElement("div");
        pp.textContent = title;
        pp.setAttribute("class", "container2");

        pp.innerHTML = 
        "<div class='title'>  "+object.Name+" </div>"+
                "<div>"+
            "<img src=/data/loops/" + object.Img + " alt='poop' style='width:100%'"+"onclick='loopFard(" + object.ID + ")'>"+
                "</div>"
        temp?.appendChild(pp);
    }
    temp?.setAttribute("class", "wrapper");
    maymay?.replaceWith(temp);
}
function loopFard(id) {
    var url = "../loop/";
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url + id, true);
    xhr.send();
}
function stopLoop() {
    var url = "../stoploop";
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.send();
}
async function getLongplay() {
    const url = "../getlong";
    const response = await fetch(url);
    const jsonData = await response.json();
    var maymay = document.getElementById("loopsm8");
    var temp = document.createElement("loopsm8");
    temp.setAttribute("id", "loopsm8");
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
        pp.innerHTML = "<img src=/data/loops/" + object.Img + " alt='poop' style='width:100%'>" +
            "<button class='btn' onclick='loopFard(" + object.ID + ")'>" + object.Name + "</button>";
        temp?.appendChild(pp);
    }
    temp?.setAttribute("class", "wrapper");
    maymay?.replaceWith(temp);
}


