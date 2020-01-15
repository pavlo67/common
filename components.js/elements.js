// exports.int        = int;
//
// exports.nn         = nn;
// exports.str        = str;
// exports.cut        = cut;
// exports.escapeXml  = escapeXml;  // npm i html-entities
//
// exports.array      = array;
// exports.first      = first;
// exports.second     = second;
// exports.order      = order;
//
// exports.object     = object;
// exports.copy       = copy;

export default {
    href,
}

function href(url, title) {
    if (!title) title = url;

    return "<a href=\"" + url + "\" target=_blank>" + title+ "</a>";
}

