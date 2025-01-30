class Keyframe {
    constructor(x, y, angle) {
        this.x = x;
        this.y = y;
        this.angle = angle;
    }
}

const gdoc_head = new Image(134, 196);
const gdoc_body = new Image(114, 88);
const gdoc_tail = new Image(86, 148);
const cursor = new Image(32, 32);
let p_colour, s_colour;
let head_offset = 2;
let hasPowerup = false;
let powerupSpawned = false;
let uid;
var context;
let head_x = 100;
let head_y = 100;
let speed = 5;
let h_off_x, h_off_y;
let angle = Math.PI / 2;
let score_text = "Wynik: ";
let score = 0;
let key = 0;
let keyframes = [];
let lastMouseState = [0, 0];
let powerupState = [-1, -1];
let damageIndicator = [-1, -1, -1, -1];
let start = false;
let bgmusplaying = false;
let sfxplayed = false;
var bgmus = new Audio("/static/media/bgmusic.mp3");
var hit = new Audio("/static/media/hit.mp3");
var death = new Audio("/static/media/death.mp3");
const BASE_KEYFRAME_DIFFERENCE = 8;
const TRIGGER_RADIUS = 30;
let SEGMENTS;
let alive = true;
let curr_keyframe_difference = BASE_KEYFRAME_DIFFERENCE;
let prf;
window.onload = function () {
    fetch("/l/getname", { method: "GET", redirect: "error" }).then(
        async function (res) {
            if (res.redirected) {
                window.location.href = "/";
            } else {
                const json = await res.json();
                uid = json.uid;
            }
        },
    );
    let c_prof = window.localStorage.getItem("curr_profile");
    if (c_prof == null) {
        window.location.href = "/";
    }
    fetch("/l/profile/" + c_prof, { method: "GET", redirect: "error" }).then(
        async function (response) {
            if (response.status == 204) {
                window.localStorage.setItem("curr_profile", 0);
                window.location.href = "/l/customise";
            }
            pro = await response.json();
            prf = pro.profile;

            cursor.src = "/static/media/" + prf.cursor + ".png";
            p_colour = prf.pcolour;
            s_colour = prf.scolour;
            SEGMENTS = prf.segments;

            bgmus.loop = true;
            death.loop = false;

            var c = document.getElementById("central_canvas");
            c.width = window.innerWidth;
            c.height = window.innerHeight;
            context = c.getContext("2d");
            if (prf.bgtype == "colour") {
                document.body.style.background = prf.bgcontent;
            } else if (prf.bgtype == "static") {
                document.body.style.backgroundImage =
                    "url('/static/media/sbg/" + prf.bgcontent + ".png')";
            } else {
                document.body.style.backgroundImage =
                    "url('/static/media/dbg/" + prf.bgcontent + ".gif')";
            }
            onmousemove = (event) => {
                start = true;
                lastMouseState[0] = event.clientX;
                lastMouseState[1] = event.clientY;
                h_off_y = lastMouseState[1] - head_y;
                h_off_x = lastMouseState[0] - head_x;
                if (h_off_y == 0) {
                    if (h_off_x < 0) {
                        angle = -0.5 * Math.PI;
                    } else {
                        angle = 0.5 * Math.PI;
                    }
                } else if (h_off_x == 0) {
                    if (h_off_y > 0) {
                        angle = Math.PI;
                    } else {
                        angle = 0;
                    }
                } else {
                    angle = Math.atan(Math.abs(h_off_y) / Math.abs(h_off_x));
                    if (h_off_y < 0 && h_off_x < 0) {
                        angle -= Math.PI / 2;
                    } else if (h_off_y > 0 && h_off_x < 0) {
                        angle = -Math.PI / 2 - angle;
                    } else if (h_off_y > 0 && h_off_x > 0) {
                        angle += Math.PI / 2;
                    } else {
                        angle = Math.PI / 2 - angle;
                    }
                }
            };
            let text = "Rusz myszką by zacząć";
            context.font = "60px bold Arial";
            context.textAlign = "center";
            if (prf.bgtype == "colour") {
                context.fillStyle = invertColor(prf.bg_content);
            } else if (prf.bgtype == "static") {
                if (prf.bg_content == "windows") {
                    context.fillStyle = "black";
                } else {
                    context.fillStyle = "white";
                }
            } else {
                if (prf.bg_content == "static") {
                    context.fillStyle = "yellow";
                } else {
                    context.fillStyle = "white";
                }
            }

            context.fillText(
                text,
                window.innerWidth / 2,
                window.innerHeight / 2,
            );

            setInterval(function () {
                if (start) {
                    if (bgmusplaying == false && prf.music == true) {
                        bgmus.play();
                        bgmusplaying = true;
                    }
                    if (alive) {
                        h_off_y = lastMouseState[1] - head_y;
                        h_off_x = lastMouseState[0] - head_x;
                        key++;
                        if (key % 10 == 0) {
                            score += 10;
                        }
                        if (key == 200) {
                            if (!powerupSpawned && !hasPowerup) {
                                spawnPowerup();
                            }
                            key = 0;
                        }
                        drawSerpent();
                        drawScore();
                        if (powerupSpawned) {
                            drawPowerup();
                        }
                        if (damageIndicator[0] != -1) {
                            drawExplosion();
                        }
                        drawMouse();
                    } else {
                        if (prf.music == true) {
                            bgmus.pause();
                        }
                        if (prf.sfx == true && !sfxplayed) {
                            death.play();
                            sfxplayed = true;
                        }
                        context.reset();
                        drawScore();
                        drawDeathText();
                        key++;
                        if (key == 500) {
                            fetch("/l/scores", {
                                method: "POST",
                                body: JSON.stringify({
                                    pid: parseInt(c_prof),
                                    score: score,
                                }),
                                headers: {
                                    "Content-type":
                                        "application/json; charset=UTF-8",
                                },
                            }).then((res) => {
                                window.location.href = "/scoreboard";
                            });
                        }
                    }
                }
            }, 10);
        },
    );
};
function drawExplosion() {
    context.beginPath();
    context.fillStyle = "rgba(255,0,0," + damageIndicator[3] + ")";
    context.arc(
        damageIndicator[0],
        damageIndicator[1],
        damageIndicator[2],
        0,
        Math.PI * 2,
    );
    context.stroke();
    context.fill();
    damageIndicator[2]++;
    damageIndicator[3] = damageIndicator[3] - 0.01;
    if (damageIndicator[3] < 0) {
        damageIndicator[0] = -1;
    }
}
function distance(x1, x2, y1, y2) {
    let x = x2 - x1;
    let y = y2 - y1;
    return Math.sqrt(x * x + y * y);
}
function drawPowerup() {
    if (
        distance(
            lastMouseState[0],
            powerupState[0],
            lastMouseState[1],
            powerupState[1],
        ) < TRIGGER_RADIUS
    ) {
        hasPowerup = true;
        powerupSpawned = false;
    }
    context.beginPath();
    context.fillStyle = s_colour;
    context.rect(powerupState[0], powerupState[1], 30, 30);
    context.stroke();
    context.fill();
}
function spawnPowerup() {
    powerupState[0] = Math.random() * (window.innerWidth - 50) + 25;
    powerupState[1] = Math.random() * (window.innerHeight - 50) + 25;
    powerupSpawned = true;
}
function drawDeathText() {
    let death_text = "Zostałeś zjedzony";
    context.save();
    context.font = "60px bold Arial";
    context.textAlign = "center";
    if (prf.bgtype == "colour") {
        context.fillStyle = invertColor(prf.bg_content);
    } else if (prf.bgtype == "static") {
        if (prf.bg_content == "windows") {
            context.fillStyle = "black";
        } else {
            context.fillStyle = "white";
        }
    } else {
        if (prf.bg_content == "static") {
            context.fillStyle = "yellow";
        } else {
            context.fillStyle = "white";
        }
    }

    context.fillText(
        death_text,
        window.innerWidth / 2,
        window.innerHeight / 2 + 20 * Math.sin(key * 0.01),
    );
    context.restore();
}
function invertColor(hex) {
    var clr = hex;
    clr = clr.substring(1);
    clr = parseInt(clr, 16);
    clr = 0xffffff ^ clr;
    clr = clr.toString(16);
    clr = ("000000" + clr).slice(-6);
    return "#" + clr;
}
function drawScore() {
    score_text = "Wynik: " + score;
    context.save();
    context.font = "40px bold Arial";
    if (prf.bgtype == "colour") {
        context.fillStyle = invertColor(prf.bgcontent);
    } else if (prf.bgtype == "static") {
        if (prf.bgcontent == "windows") {
            context.fillStyle = "black";
        } else {
            context.fillStyle = "white";
        }
    } else {
        if (prf.bgcontent == "static") {
            context.fillStyle = "yellow";
        } else {
            context.fillStyle = "white";
        }
    }
    context.textAlign = "center";
    context.fillText(score_text, window.innerWidth / 2, 40);
    context.restore();
}
function drawMouse() {
    if (Math.sqrt(h_off_x * h_off_x + h_off_y * h_off_y) < TRIGGER_RADIUS) {
        alive = false;
        console.log("died");
    }

    context.beginPath();
    if (hasPowerup) {
        context.fillStyle = s_colour;
    } else {
        context.fillStyle = p_colour;
    }
    context.arc(lastMouseState[0], lastMouseState[1], 10, 0, Math.PI * 2);
    context.stroke();
    context.fill();
    context.save();
    context.translate(lastMouseState[0], lastMouseState[1]);
    context.drawImage(
        cursor,
        -cursor.width / 3,
        -cursor.height / 3,
        cursor.width / 1.5,
        cursor.height / 1.5,
    );
    context.restore();
}
function drawSerpent() {
    context.reset();
    drawRotatedImage(gdoc_head, head_x, head_y, angle, 0, 0);
    head_y -= Math.cos(angle) * speed;
    head_x += Math.sin(angle) * speed;
    keyframes.push(new Keyframe(head_x, head_y, angle));
    for (
        var i = keyframes.length - head_offset * curr_keyframe_difference - 1;
        i >= curr_keyframe_difference;
        i -= curr_keyframe_difference
    ) {
        if (i == curr_keyframe_difference) {
            drawRotatedImage(
                gdoc_tail,
                keyframes[i].x,
                keyframes[i].y,
                keyframes[i].angle,
                0,
                0,
            );
            keyframes.splice(
                0,
                keyframes.length - curr_keyframe_difference * SEGMENTS,
            );
        } else {
            drawRotatedImage(
                gdoc_body,
                keyframes[i].x,
                keyframes[i].y,
                keyframes[i].angle,
                0,
                0,
            );
            if (
                hasPowerup &&
                distance(
                    keyframes[i].x,
                    lastMouseState[0],
                    keyframes[i].y,
                    lastMouseState[1],
                ) < TRIGGER_RADIUS
            ) {
                hasPowerup = false;
                damageIndicator[0] = keyframes[i].x;
                damageIndicator[1] = keyframes[i].y;
                damageIndicator[2] = 1;
                damageIndicator[3] = 1;
                score += 1000;
                if (prf.sfx == true) {
                    hit.play();
                }
            }
        }
    }
    speed += 0.001;
    if (speed > 7) {
        head_offset = 1;
    }
    if (speed > 11 && curr_keyframe_difference > 3) {
        curr_keyframe_difference =
            BASE_KEYFRAME_DIFFERENCE - Math.floor((speed - 9) / 2);
    } else if (speed > 30 && curr_keyframe_difference == 3) {
        curr_keyframe_difference = 2;
    } else if (speed > 45) {
        curr_keyframe_difference = 1;
    }
}
gdoc_head.src = "/static/media/head.png";
gdoc_body.src = "/static/media/body.png";
gdoc_tail.src = "/static/media/tail.png";

function drawRotatedImage(image, x, y, angle, ox, oy) {
    context.save();
    context.translate(x + ox, y + oy);
    context.rotate(angle);
    context.drawImage(
        image,
        -image.width / 2,
        -image.height / 2,
        image.width,
        image.height,
    );
    context.restore();
}
