const minTitleLength = 1;

let divBlock;
let lastDiv;

function showContent() {
    if (lastDiv != null) {
        lastDiv.style.display = "none";
    }
    showDiv(event);
}

function hideOwnContent() {
    clearTimeout(divBlock);
    let x = document.getElementById(event.target.id);
    if (x != null) {
        // if (event.relatedTarget == null || (event.relatedTarget.parentElement !== x && event.relatedTarget.src == null)) {
        if (event.relatedTarget == null || (event.relatedTarget.parentElement !== x && event.relatedTarget.parentElement.parentElement !== x)) {
            x.style.display = "none";
        }
    }
}

function hideContent() {
    clearTimeout(divBlock);
    let h5ID = event.target.id;
    let id = h5ID.replace(/^news_title_id_/, "");
    let x = document.getElementById("news_content_id_" + id);
    if (id !== '') {
        lastDiv = document.getElementById("news_content_id_" + id);
        if (event.relatedTarget == null || "news_content_id_" + id !== event.relatedTarget.id) {
            x.style.display = "none";
        }
    }
}

function showDiv(e) {
    divBlock = setTimeout(function () {
        let h5ID = e.target.id;
        let id = h5ID.replace(/^news_title_id_/, "");
        if (id !== '') {
            let x = document.getElementById("news_content_id_" + id);
            x.style.display = "block";
        }
    }, 500);
}
