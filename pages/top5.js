
async function gettop5() {
    const url = "./top5refresh";
    const response = await fetch(url);
    const htmlData = await response.text();


    const current = document.getElementById("sausage");
    const p = document.createElement("div");
    p.id = "sausage";
    p.innerHTML = htmlData;
    current?.replaceWith(p);


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


        pp.innerHTML = "<div class='title'>  " + object.title + " </div>" + "<div>" + "<img src=/data/img/" + object.img + " alt='poop' onclick='fardFunc(" + object.ID + ")'style='width:100%'>" +

            "<button class='btn2' onclick=browserFard('" + object.file + "')><i class='fa-solid fa-headphones'></i></button>" + "</div>";
        temp?.appendChild(pp);
    }
    temp?.setAttribute("class", "wrapper");
    maymay?.replaceWith(temp);
}
