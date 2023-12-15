var torpedo;
function hitEndpoint(endpoint) {
    var url = "./"+endpoint;
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.send();
}



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
function changeFlakes(flakePath) {
    console.log(flakePath);

    var oldflake = document.getElementById("snowflakeContainer");
    var newflake = document.createElement("div");
  newflake.setAttribute("id","snowflakeContainer");
  newflake.setAttribute("style","display=block\;");
newflake.innerHTML = "<img src="+flakePath+" class='snowflake'></span>";
  oldflake?.replaceWith(newflake);
  generateSnowflakes();
}

async function getStatus() {
    const url = "./status";
    const response = await fetch(url);
    const jsonData = await response.json();
    const MemeStatus = jsonData.memes;
    const current = document.getElementById("memes");
    const p = document.createElement("memes");
    p.textContent = MemeStatus;
    current?.replaceWith(p);
    const LoopStatus = jsonData.loops;
    const next = document.getElementById("loops");
    const pp = document.createElement("loops");
    pp.textContent = LoopStatus;
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
        pp.setAttribute("class", "container2");


        pp.innerHTML = "<div class='title'>  "+object.title+" </div>" +"<img src=/data/img/" + object.img + " alt='poop' onclick='fardFunc(" + object.ID + ")'style='width:100%'>" +
            
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

function changeFlakes2(flakePath) {
    console.log(flakePath);
if (snowing == false){generateSnowflakes();
globalThis.snowing = true;
  }

    var oldflakes = document.getElementsByClassName("snowflake");
  for (var flake of oldflakes) {
    flake.attributes.item(0).nodeValue = '/data/flakes/'+flakePath;
  
  }


}


  function setAccessibilityState() {
    if (reduceMotionQuery.matches) {
      enableAnimations = false;
    } else {
      enableAnimations = true;
    }
  }

  function setup() {
    if (enableAnimations) {
      window.addEventListener("DOMContentLoaded", generateSnowflakes, false);
      window.addEventListener("resize", setResetFlag, false);
    }
  }

  function generateSnowflakes() {
    // get our snowflake element from the DOM and store it
    let originalSnowflake = document.querySelector(".snowflake");

    // access our snowflake element's parent container
    let snowflakeContainer = originalSnowflake.parentNode;
    snowflakeContainer.style.display = "block";

    // get our browser's size
    browserWidth = document.documentElement.clientWidth;
    browserHeight = document.documentElement.clientHeight;

    // create each individual snowflake
    for (let i = 0; i < numberOfSnowflakes; i++) {
      // clone our original snowflake and add it to snowflakeContainer
      let snowflakeClone = originalSnowflake.cloneNode(true);
      snowflakeContainer.appendChild(snowflakeClone);

      // set our snowflake's initial position and related properties
      let initialXPos = getPosition(50, browserWidth);
      let initialYPos = getPosition(50, browserHeight);
      let speed = (5 + Math.random() * 40) * delta;

      // create our Snowflake object
      let snowflakeObject = new Snowflake(
        snowflakeClone,
        speed,
        initialXPos,
        initialYPos
      );
      snowflakes.push(snowflakeObject);
    }

    // remove the original snowflake because we no longer need it visible
    snowflakeContainer.removeChild(originalSnowflake);

    requestAnimationFrame(moveSnowflakes);
  }
  function moveSnowflakes(currentTime) {
    delta = (currentTime - previousTime) / frame_interval;

    if (enableAnimations) {
      for (let i = 0; i < snowflakes.length; i++) {
        let snowflake = snowflakes[i];
        snowflake.update(delta);
      }
    }

    previousTime = currentTime;

    // Reset the position of all the snowflakes to a new value
    if (resetPosition) {
      browserWidth = document.documentElement.clientWidth;
      browserHeight = document.documentElement.clientHeight;

      for (let i = 0; i < snowflakes.length; i++) {
        let snowflake = snowflakes[i];

        snowflake.xPos = getPosition(50, browserWidth);
        snowflake.yPos = getPosition(50, browserHeight);
      }

      resetPosition = false;
    }

    requestAnimationFrame(moveSnowflakes);
  }

  function setTransform(xPos, yPos, scale, el) {
    el.style.transform = `translate3d(${xPos}px, ${yPos}px, 0) scale(${scale}, ${scale})`;
  }




  //
  // Constructor for our Snowflake object
  //

  //
  // A performant way to set your snowflake's position and size
  //

  //
  // The function responsible for creating the snowflake
  //

  //
  // Responsible for moving each snowflake by calling its update function
  //


  //
  // This function returns a number between (maximum - offset) and (maximum + offset)
  //
  function getPosition(offset, size) {
    return Math.round(-1 * offset + Math.random() * (size + 2 * offset));
  }

  //
  // Trigger a reset of all the snowflakes' positions
  //
  function setResetFlag(e) {
    resetPosition = true;
  }
  class Snowflake {
    constructor(element, speed, xPos, yPos) {
      // set initial snowflake properties
      this.element = element;
      this.speed = speed;
      this.xPos = xPos;
      this.yPos = yPos;
      this.scale = 1;

      // declare variables used for snowflake's motion
      this.counter = 0;
      this.sign = Math.random() < 0.5 ? 1 : -1;

      // setting an initial opacity and size for our snowflake
      this.element.style.opacity = (0.1 + Math.random()) / 3;
    }

    // The function responsible for actually moving our snowflake
    update(delta) {
      // using some trigonometry to determine our x and y position
      this.counter += (this.speed / 5000) * delta;
      this.xPos += (this.sign * delta * this.speed * Math.cos(this.counter)) / 40;
      this.yPos += Math.sin(this.counter) / 40 + (this.speed * delta) / 30;
      this.scale = 0.5 + Math.abs((10 * Math.cos(this.counter)) / 20);

      // setting our snowflake's position
      setTransform(
        Math.round(this.xPos),
        Math.round(this.yPos),
        this.scale,
        this.element
      );

      // if snowflake goes below the browser window, move it back to the top
      if (this.yPos > browserHeight) {
        this.yPos = -50;
      }
    }
  }// Array to store our Snowflake objects

  let snowflakes = [];

  // Global variables to store our browser's window size
  let browserWidth;
  let browserHeight;

  // Specify the number of snowflakes you want visible
  let numberOfSnowflakes = 50;

  // Flag to reset the position of the snowflakes
  let resetPosition = false;

  // Handle accessibility
  let enableAnimations = false;
  let reduceMotionQuery = matchMedia("(prefers-reduced-motion)");

  let frames_per_second = 60;
  let frame_interval = 1000 / frames_per_second;



  let previousTime = performance.now();
  let delta = 1;
  // Handle animation accessibility preferences
  setAccessibilityState();

  reduceMotionQuery.addListener(setAccessibilityState);

  //
  // It all starts here...
  //
  //setup();
