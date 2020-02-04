export default {NewShowHide}

function NewShowHide(postfix) {
    return new ShowHideDetails(postfix);
}

class ShowHideDetails {
    constructor(postfix) {
        this.postfix = postfix || "";
        this.postfixRegExp = new RegExp(postfix + "$")
    }

    showContent(ev) {
        // if (this.lastDiv != null) {
        //     this.lastDiv.style.display = "none";
        // }
        this.showDiv(ev);
    }

    showDiv(ev) {
        this.divBlock = setTimeout(() => {
            let x = document.getElementById(ev.target.id + this.postfix);
            if (x) x.style.display = "block";
        }, 500);
    }

    hideContent(ev) {
        clearTimeout(this.divBlock);

        let x = document.getElementById(ev.target.id.replace(this.postfixRegExp, "") + this.postfix);
        if (x) x.style.display = "none";

        // if (x) {
        //     lastDiv = document.getElementById("news_content_id_" + id);
        //     if (event.relatedTarget == null || "news_content_id_" + id !== event.relatedTarget.id) {
        //         x.style.display = "none";
        //     }
        // }
    }


}


//
// function hideOwnContent() {
//     clearTimeout(divBlock);
//     let x = document.getElementById(event.target.id);
//     if (x != null) {
//         // if (event.relatedTarget == null || (event.relatedTarget.parentElement !== x && event.relatedTarget.src == null)) {
//         if (event.relatedTarget == null || (event.relatedTarget.parentElement !== x && event.relatedTarget.parentElement.parentElement !== x)) {
//             x.style.display = "none";
//         }
//     }
// }
//
//
