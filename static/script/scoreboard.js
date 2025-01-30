window.onload = async function () {
    const response = await fetch("/scores", {
        method: "GET",
        headers: {
            "Content-type": "application/json; charset=UTF-8",
        },
    });
    const rjson = await response.json();
    let scores = rjson.scores;
    for (let i = 0; i < scores.length; i++) {
        let prf = document.createElement("div");
        prf.innerHTML = genScore(scores[i], i + 1);
        prf.className = "score";
        document.getElementById("scoreboard").appendChild(prf);
    }
};
function genScore(score, index) {
    let sc = document.createElement("p");
    sc.innerHTML = index + ". ";
    let crsr = document.createElement("img");
    if (score.cursor == "classic") {
        crsr.src = "/static/media/classic.png";
    } else if (score.cursor == "sharp") {
        crsr.src = "/static/media/sharp.png";
    } else {
        crsr.src = "/static/media/time.png";
    }
    crsr.style.backgroundColor = score.colour;
    return (
        sc.outerHTML +
        crsr.outerHTML +
        " <p>" +
        score.name +
        " Wynik: " +
        score.score +
        "</p>"
    );
}
