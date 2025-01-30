window.onload = function () {
    let RBBanner = document.getElementById("redirectbanner");
    let RBContent = RBBanner.children[0].innerHTML;
    if (RBContent == "") {
        RBBanner.hidden = true;
    } else {
        RBBanner.hidden = false;
    }
};
