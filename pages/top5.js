


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

setInterval(gettop5, 5000);
