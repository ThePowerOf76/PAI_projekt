let list;

window.onload = function () {
    document.forms["profile"].style.display = "none";
    fetch("/l/profilelist", { method: "GET" }).then(async function (response) {
        const aaa = await response.json();
        list = aaa.profiles;
        if (list == null || list.length == 0) {
            document.getElementById("pickP").disabled = true;
            document.getElementById("editP").disabled = true;
            document.getElementById("delP").disabled = true;
        }
        document.forms["profile"].addEventListener("submit", (e) => {
            e.preventDefault();
            let backg;
            let form = document.forms["profile"];
            if (form["bg"].value == "colour") {
                backg = form["cbg"].value;
            } else if (form["bg"].value == "static") {
                backg = form["sbg"].value;
            } else {
                backg = form["dbg"].value;
            }
            let obj = {
                pid: parseInt(form["pid"].value),
                oid: 0,
                name: form["name"].value,
                cursor: form["cursor"].value,
                bgtype: form["bg"].value,
                bgcontent: backg,
                pcolour: form["player"].value,
                scolour: form["powerup"].value,
                segments: parseInt(form["segments"].value),
                music: form["music"].checked,
                sfx: form["sfx"].checked,
            };
            if (form["pid"].value == -1) {
                fetch("/l/profile", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(obj),
                }).then(async function (response) {
                    const obj = await response.json();
                    window.localStorage.setItem("curr_profile", obj.pid);
                    alert("Nowy profil dodany");
                    window.location.href = "/";
                });
            } else {
                fetch("/l/profile/" + form["pid"].value, {
                    method: "PATCH",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(obj),
                }).then(function (response) {
                    alert("Profil zaktualizowany");
                    location.reload();
                });
            }
        });
    });
};
function bgChange() {
    let bg_style = document.forms["profile"]["bg"].value;
    if (bg_style == "colour") {
        document.forms["profile"]["cbg"].disabled = false;
        document.forms["profile"]["dbg"].disabled = true;
        document.forms["profile"]["sbg"].disabled = true;
        document.forms["profile"]["cbg"].hidden = false;
        document.forms["profile"]["dbg"].hidden = true;
        document.forms["profile"]["sbg"].hidden = true;
    } else if (bg_style == "static") {
        document.forms["profile"]["cbg"].disabled = true;
        document.forms["profile"]["dbg"].disabled = true;
        document.forms["profile"]["sbg"].disabled = false;
        document.forms["profile"]["cbg"].hidden = true;
        document.forms["profile"]["dbg"].hidden = true;
        document.forms["profile"]["sbg"].hidden = false;
    } else {
        document.forms["profile"]["cbg"].disabled = true;
        document.forms["profile"]["dbg"].disabled = false;
        document.forms["profile"]["sbg"].disabled = true;
        document.forms["profile"]["cbg"].hidden = true;
        document.forms["profile"]["dbg"].hidden = false;
        document.forms["profile"]["sbg"].hidden = true;
    }
}
function newProfile() {
    document.forms["profile"].hidden = false;
    document.forms["profile"].disabled = false;
    document.forms["profile"].reset();
    document.forms["profile"].style.display = "flex";
    document.forms["profile"]["classic"].checked = true;
    document.forms["profile"]["sbg"].hidden = true;
    document.forms["profile"]["sbg"].disabled = true;
    document.forms["profile"]["dbg"].hidden = true;
    document.forms["profile"]["dbg"].disabled = true;
    document.forms["profile"]["cbg"].hidden = true;
    document.forms["profile"]["cbg"].disabled = true;
    document.forms["profile"]["pid"].value = -1; // -1 is new player
    document.forms["all_profiles"].hidden = true;
    document.forms["all_profiles"].disabled = true;
    document.forms["all_profiles"].reset();
}
function loadProfiles() {
    document.forms["profile"].reset();
    document.forms["profile"].hidden = true;
    document.forms["profile"].disabled = true;
    document.forms["profile"].style.display = "none";
    document.forms["all_profiles"].hidden = false;
    document.forms["all_profiles"].disabled = false;
    document.forms["all_profiles"].reset();
    document.forms["all_profiles"]["profiles"].innerHTML = "";
    for (var i = 0; i < list.length; i++) {
        let prf = document.createElement("option");
        prf.value = list[i].pid;
        prf.innerHTML = list[i].name;
        document.forms["all_profiles"]["profiles"].appendChild(prf);
    }
}
function pick(e) {
    e.preventDefault();
    window.localStorage.setItem(
        "curr_profile",
        document.forms["all_profiles"]["profiles"].value,
    );
    window.location.href = "/";
}
function pickProfileForm() {
    loadProfiles();
    document.forms["all_profiles"]["subm"].value = "Zatwierdź";
    document.forms["all_profiles"].addEventListener("submit", pick);
}
function del(e) {
    e.preventDefault();
    console.log(list);
    console.log(document.forms["all_profiles"]["profiles"].value);
    let text = "Czy na pewno chcesz usunąć profil?";
    if (confirm(text) == true) {
        fetch(
            "/l/profile/" + document.forms["all_profiles"]["profiles"].value,
            {
                method: "DELETE",
            },
        ).then(function (response) {
            if (
                window.localStorage.getItem("curr_profile") ==
                document.forms["all_profiles"]["profiles"].value
            ) {
                window.localStorage.setItem("curr_profile", 0);
            }
            alert("Profil usunięty");
            location.reload();
        });
    }
}
function delProfileForm() {
    loadProfiles();
    document.forms["all_profiles"]["subm"].value = "Usuń profil";
    document.forms["all_profiles"].addEventListener("submit", del);
}

function editProfileForm() {
    loadProfiles();
    document.forms["all_profiles"]["subm"].value = "Edytuj profil";
    document.forms["all_profiles"].removeEventListener("submit", del);
    document.forms["all_profiles"].removeEventListener("submit", pick);
    document.forms["all_profiles"].addEventListener("submit", (e) => {
        e.preventDefault();
        document.forms["profile"].hidden = false;
        document.forms["profile"].disabled = false;
        document.forms["profile"].reset();
        document.forms["profile"].style.display = "flex";
        let id = document.forms["all_profiles"]["profiles"].value;
        fetch("/l/profile/" + id, { method: "GET" }).then(
            async function (response) {
                const pobj = await response.json();
                document.forms["profile"]["pid"].value = id;
                document.forms["profile"]["name"].value = pobj.profile.name;
                document.forms["profile"]["cursor"].value = pobj.profile.cursor;
                document.forms["profile"]["bg"].value = pobj.profile.bgtype;
                if (pobj.profile.bgtype == "colour") {
                    document.forms["profile"]["cbg"].disabled = false;
                    document.forms["profile"]["dbg"].disabled = true;
                    document.forms["profile"]["sbg"].disabled = true;
                    document.forms["profile"]["cbg"].hidden = false;
                    document.forms["profile"]["dbg"].hidden = true;
                    document.forms["profile"]["sbg"].hidden = true;
                    document.forms["profile"]["cbg"].value =
                        pobj.profile.bgcontent;
                } else if (pobj.profile.bgtype == "static") {
                    document.forms["profile"]["cbg"].disabled = true;
                    document.forms["profile"]["dbg"].disabled = true;
                    document.forms["profile"]["sbg"].disabled = false;
                    document.forms["profile"]["cbg"].hidden = true;
                    document.forms["profile"]["dbg"].hidden = true;
                    document.forms["profile"]["sbg"].hidden = false;
                    document.forms["profile"]["sbg"].value =
                        pobj.profile.bgcontent;
                } else {
                    document.forms["profile"]["cbg"].disabled = true;
                    document.forms["profile"]["dbg"].disabled = false;
                    document.forms["profile"]["sbg"].disabled = true;
                    document.forms["profile"]["cbg"].hidden = true;
                    document.forms["profile"]["dbg"].hidden = false;
                    document.forms["profile"]["sbg"].hidden = true;
                    document.forms["profile"]["dbg"].value =
                        pobj.profile.bgcontent;
                }
                document.forms["profile"]["player"].value =
                    pobj.profile.pcolour;
                document.forms["profile"]["powerup"].value =
                    pobj.profile.scolour;
                document.forms["profile"]["segments"].value =
                    pobj.profile.segments;
                document.forms["profile"]["music"].checked = pobj.profile.music;
                document.forms["profile"]["sfx"].checked = pobj.profile.sfx;
                document.forms["all_profiles"].hidden = true;
                document.forms["all_profiles"].disabled = true;
            },
        );
    });
}
