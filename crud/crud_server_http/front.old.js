
function nova_probam (id, vis) {
    for (var i in forma.elements) {

        if (forma.elements[i].type === "radio"
            && forma.elements[i].value == vis
            && forma.elements[i].name == ("proba_" + id)
        ) {
            forma.elements[i].click();
            return;
        }
    }
}

function quaere_in (nomen) {
    for (var i in forma.elements) {
        if (forma.elements[i].tagName == "INPUT" && forma.elements[i].name == nomen) return forma.elements[i].value;
    }
    return undefined;
}

function cresce_nomen (el, id, sinistre) {		// !!! ������ ������� ���� - ��� parentNode

    var parens = el.parentNode.parentNode;
    var vis    = parens.innerHTML;
    var vis_   = vis.split("<br><nobr>");

    var vis1 = quaere_in ("nomen_" + id);
    if (!vis1) vis1 = vis_[0];

    var specimen = new RegExp(vis1.replace(/[^\w�-��-�]+/g, "\\\W*").replace("amp", "(amp)?").replace("quot", "(quot)?"), "i");

    var nodes    = parens.parentNode.childNodes;
    var num      = nodes.length;
    var i_sing   = 0;
    for (var i = -1; ++i < nodes.length;) {
        if (nodes[i] == parens) {
            i_sing = i + 8;
            break;
        }
    }

    var singula = nodes[i_sing].innerHTML;

    var in_sing = singula.search(specimen);

    // alert(singula);
    // alert(in_sing);

    if (in_sing >= 0) {
        var in_sing_m   = singula.match(specimen);

        var in_sing_p, in_sing_p_v;

        if (sinistre) {
            in_sing_p   = singula.substr(0, in_sing);
            in_sing_p_v = in_sing_p.match(/[\w�-��-�]+[^\w�-��-�]*$/);

        } else {
            in_sing_p   = singula.substr(in_sing + in_sing_m[0].length);
            in_sing_p_v = in_sing_p.match(/^[^\w�-��-�]*[\w�-��-�]+/);

        }

        if (in_sing_p_v) {
            var vis_p    = vis_[1].replace(/<input[^>]*>/, "").replace(/<br>\\s*$/, "");
            var vis_nova = (sinistre ? in_sing_p_v[0] + vis1 : vis1 + in_sing_p_v[0]).replace(/\s+$/, "").replace(/"/, "&quot;");
            var nomen    = "nomen_" + id;
            parens.innerHTML = vis_nova + "<br><nobr>" + vis_p + "<br><input name=" + nomen + " id=" + nomen + " style=\"width:150px;\" value=\"" + vis_nova + "\">";
            nova_probam (id, "b");
            var novum    = document.getElementById(nomen);
            if (novum) novum.focus();
        }
    }
}

function curta_nomen (el, id, sinistre) {

    var parens = el.parentNode.parentNode;   // !!! ������ ������� ���� - ��� parentNode

    var vis    = parens.innerHTML;
    var vis_   = vis.split("<br><nobr>");

    var vis1 = quaere_in ("nomen_" + id);
    if (!vis1) vis1 = vis_[0];

    var vis_p    = vis_[1].replace(/<input[^>]*>/, "").replace(/<br>\\s*$/,"");;
    var vis_nova = sinistre
        ? vis1.replace(/^([^\w�-��-�]+|\d+\s*|[A-Za-z�-��-�]+\s*)/, "").replace(/"/, "&quot;")
        : vis1.replace(/([^\w�-��-�]+|\s*\d+|\s*[A-Za-z�-��-�]+)$/, "").replace(/"/, "&quot;");
    var nomen    = "nomen_" + id;
    parens.innerHTML = vis_[0] + "<br><nobr>" + vis_p + "<br><input name=" + nomen + " id=" + nomen + " style=\"width:150px;\" value=\"" + vis_nova + "\">";
    nova_probam (id, "b");
    var novum    = document.getElementById(nomen);
    if (novum) 	novum.focus();
}



function adde_ser (el, id, fabr) {
    var vis1 = quaere_in ("nomen_" + id);
    if (!vis1) {
        var parens = el.parentNode;
        var vis    = parens.innerHTML;			// !!! ������ ������� ���� - ���� parentNode
        var vis_   = vis.split("<br><nobr>");
        vis1       = vis_[0].replace(/\s+/g, "+");
    }

    var fenestra = window.open("/c/k.pl?c_=thes&c_=propr&c_=.vis_nova_ser&c_=" + vis1 + "&fabr=" + fabr, "", "location,width=800,height=550,top=0,left=400,scrollbars=1"); fenestra.focus();
    nova_probam (id, "n");
}

var fam = "";

function adde_propr (el, id) {
    var vis1 = quaere_in ("nomen_" + id);
    if (!vis1) {
        var parens = el.parentNode;				// !!! ������ ������� ���� - ���� parentNode
        var vis    = parens.innerHTML;
        var vis_   = vis.split("<br><nobr>");
        vis1       = vis_[0];
    }

    var lim    = vis1.search("::");
    if (lim >= 0) {
        fam  = vis1.substr(lim + 2).replace(/\s+/g, "+");
        vis1 = vis1.substr(0, lim).replace(/\s+/g, "+");		// !!! ������� ��� � ���������������

    } else {
        vis1 = vis1.replace(/\s+/g, "+");

    }

    var fenestra = window.open("/c/k.pl?c_=thes&c_=propr&c_=.vis_nova_&c_=" + vis1 + "&fam=" + fam, "", "location,width=800,height=550,top=0,left=400,scrollbars=1"); fenestra.focus();
    nova_probam (id, "n");
}

function adde_gen (el, id) {
    var vis1 = quaere_in ("nomen_" + id);
    if (!vis1) {
        var parens = el.parentNode;				// !!! ������ ������� ���� - ���� parentNode
        var vis    = parens.innerHTML;
        var vis_   = vis.split("<br><nobr>");
        vis1       = vis_[0];
    }

    var sigla_par;
    var lim    = vis1.search("::");
    if (lim >= 0) {
        sigla_par = vis1.substr(lim + 2).replace(/\s+/g, "+");
        vis1      = vis1.substr(0, lim).replace(/\s+/g, "+");		// !!! ������� ��� � ���������������

    } else {
        vis1      = vis1.replace(/\s+/g, "+");
        sigla_par = "";

    }

    var fenestra = window.open("/c/k.pl?c_=thes&c_=gen&c_=.vis_nova_&c_=" + vis1 + "&sigla_par=" + sigla_par, "", "location,width=800,height=550,top=0,left=400,scrollbars=1"); fenestra.focus();
    nova_probam (id, "n");
}
