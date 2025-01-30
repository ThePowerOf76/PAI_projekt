window.onload = function () {
    fetch("/l/getname", { method: "GET", redirect: "error" })
        .then(async function (res) {
            const json = await res.json();
            let c_prof = window.localStorage.getItem("curr_profile");
            if (c_prof == null || c_prof == 0) {
                document.getElementById("fight").disabled = false;
                document.getElementById("fight").value = "Wybierz profil";
                document.getElementById("fight").parentElement.action =
                    "/l/customise";
                document.getElementById("c_profile").innerHTML =
                    "<strong>Zalogowano jako " +
                    json.name +
                    ". Wybierz profil w ustawieniach</strong>";
            } else {
                document.getElementById("fight").disabled = false;
                document.getElementById("c_profile").innerHTML =
                    "<strong>Gracz: " + json.name + "</strong>";
            }
        })
        .catch((err) => {
            document.getElementById("fight").disabled = true;
            document.getElementById("c_profile").innerHTML =
                "<strong>By zagrać, zaloguj się</strong>";
        });
};
