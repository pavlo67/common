import b from "../basis";

function createdAt(item) {
    if (!(item instanceof Object && item.History instanceof Array)) return "";

    for (let h of item.History) {
        if (["saved", "created"].includes(h.Key)) {
            return b.dateStr(h.DoneAt);
        }
    }

    return "";
}

export {createdAt};