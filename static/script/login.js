fetch("/l/getname", { method: "GET", redirect: "error" })
    .then(function (res) {
        if (!res.redirected) {
            document.getElementById("log").href = "/logout";
            document.getElementById("log").innerHTML = "Wyloguj się";
        }
    })
    .catch((err) => {});
