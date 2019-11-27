// preparation data and getting it back ------------------------------------------------------------

function prepareFields(fields) {
    let preparedFields = [];

    let top = [];
    let current = preparedFields;

    for (let f of array(fields)) {
        f.format = string(f.format);
        f.params = object(f.params);
        f.attributes = object(f.attributes);

        if (f.type === "]") {
            if (top.length > 0) {
                current = top.pop();
            }
            continue
        }

        current.push(f);

        if (f.type === "[") {
            top.push(current);
            current = (f.subFields = []);
        }
    }

    return preparedFields;
}

function getData(formId) {
    let data = {};

    let elements =  document.querySelectorAll(`[id^="${formId}"]`);
    for (let el of Array.from(elements)) {
        if ((el.type === "checkbox" || el.type === "radio") && !el.checked) continue;
        let id = el.id.substr(formId.length);

        // let ids = id.split("/");
        // let data_ = data;
        // for (let i = -1; ++i < ids.length - 1;) {
        //     let id = ids[i];
        //     data_ = (data_[id] = object(data_[id]));
        // }
        // data_[ids[ids.length - 1]] = el.value;

        data[id] = el.value;
    }

    log(data);

    return data;
}


// form/view elements ------------------------------------------------------------------------------

function generalEditParts(formId, elementKey, attributes){
    attributes = object(attributes);

    let attributesHtml = "";
    for (let k in attributes) {
        attributesHtml += " "  + escapeHtml(k) + '="' + escapeHtml(attributes[k]) + '"';
    }

    let elementId = formId + elementKey;
    return [elementId, `id="` + escapeHtml(elementId) + `"` + attributesHtml];
}

function htmlSelectEdit(elementId, general, params, options, selected) {
    params = object(params);

    registerOptions(params.model, elementId);

    options = params["add_blank"] ? [[], ...array(options)] : array(options);
    let body = '';
    let option;
    for (let i = 0; i < options.length; i++) {
        body += "<option";
        if (str(options[i][1])) {
            option = options[i][1];
            body += ` value="` + escapeHtml(options[i][1]) + `"`;
        } else {
            option = options[i][0];
        }
        if (option === selected) {
            body += " selected";
        }
        body += ">" + escapeHtml(options[i][0]) + "</option>\n";
    }
    return `<select ` + general + `>` + body + "</select>\n";
}

function htmlSelectView(options, selected) {
    options = array(options);
    for (let i = 0; i < options.length; i++) {
        optionI = array(options[i]);
        let option = str(optionI[1]) ? optionI[1] : optionI[0];
        if (option === selected) {
            return escapeHtml(optionI[0]);
        }
    }

    return "";
}


// view form ----------------------------------------------------------------------------------------

function fieldView(field, value, options, frontOps) {
    field = object(field);
    options = object(options);
    frontOps = object(frontOps);

    let types = ["password", "button", "hidden"];
    for (let v of types) {
        if (v === field.type) {
            return ["", ""]
        }
    }

    let attributesEscaped = "";
    for (let k in field.attributes) {
        attributesEscaped += " "  + escapeHtml(k) + '="' + escapeHtml(field.attributes[k]) + '"'
    }

    let resHTML = "";
    if (field.type === "select") {
        resHTML = htmlSelectView(options[field.key], value);
    } else if (field.format === "url") {
        let url = escapeHtml(value);
        resHTML = `<a href="` + url + `" target=_blank>` + url + `</a>`

    } else if (field.type === "text") {
        resHTML = field.format;           // !!! not escapeHtml()

    } else if (field.type === "checkbox") {
        if (value) {
            resHTML = "так"
        } else if (!field.params["not_empty"]) {
            resHTML = "ні"
        }
    } else if (field.params["not_empty"] && (value === "0" || str(value) === "")) {
        // shows nothing

    } else if (field.params["no_escape"]) {
        resHTML = value;

    } else {
        resHTML = escapeHtml(value);
    }
    return [escapeHtml(field.label), resHTML];
}

function setHtmlView(elTop, fields, data, options, frontOps, htmlButtons, listeners, paramsStr, paramsObj) {
    let elements = [];

    if (!htmlButtons) htmlButtons = "";
    let idEdit, idDelete;
    if (listeners.edit) {
        idEdit = buttonNameBase + ++buttonNum;
        htmlButtons += `[<a href="#" id="` + idEdit+ `">редаґувати</a>] &nbsp; `;
    }

    if (listeners.delete) {
        idDelete = buttonNameBase + ++buttonNum;
        htmlButtons += `[<a href="#" id="` + idDelete + `">вилучити</a>]`;
    }

    let htmlView = htmlViewFieldsList(true, prepareFields(fields), data, options, frontOps);
    elTop.innerHTML = htmlButtons + "\n<p>" + htmlView;

    if (listeners.edit) {
        document.getElementById(idEdit).addEventListener("click", ev => listeners.edit.action(ev, paramsStr, paramsObj));
    }
    if (listeners.delete) {

        document.getElementById(idDelete).addEventListener("click", ev => listeners.delete.action(ev, paramsStr, paramsObj));

        // log(4444444, idDelete, listeners.delete, document.getElementById(idDelete))

    }
}


function htmlViewFieldsList(wrapInTable, fields, data, options, frontOps, altMode, altKey) {
    if (altMode) {
        data = object(object(data)[altKey]);
        fields = fields.filter(f => (f.params.alt === altKey));
    }

    let htmlView = "";
    let htmlTitle, htmlRes;
    for (let f of fields){
        for (let value of array(data[f.key])) {
            if (f.subFields instanceof Array) {
                htmlTitle = escapeHtml(f.label);
                if (f.format === "alt") {
                    value = object(value);
                    htmlRes = htmlViewFieldsList(false, f.subFields, value, options, frontOps, true, data[f.key + ".alt"]);
                } else {
                    htmlRes = htmlViewFieldsList(false, f.subFields, value, options, frontOps);
                }

            } else {
                [htmlTitle, htmlRes] = fieldView(f, value, options, frontOps);
            }
            if ((htmlRes === "") && ((f.params["not_empty"]) || !htmlTitle)) {
                continue
            }

            if (htmlTitle) {
                htmlTitle = "<small>" + htmlTitle + ":</small>\n";
            }

            htmlView += wrapInTable
                      ? "<tr><td>\n" + htmlTitle + "</td><td>&nbsp;</td><td>" + htmlRes + "\n</td></tr>\n"
                      : (htmlView ? "\n<p>" : "") + (htmlTitle ? htmlTitle + "<br>" : "")  + htmlRes;
        }
    }

    return wrapInTable ? inTable(htmlView, 5, "#ffffcc") : htmlView;
}


// edit form ---------------------------------------------------------------------------------------

let formNum = 0;
const formNameBase = "edit_";
const formNameRegexp = new RegExp('^' + formNameBase + '\\d+_');

function getFormId(id) {
    let res = id.match(formNameRegexp);
    return res ? res[0] : undefined;
}

let buttonNum = 0;
const buttonNameBase = "button_";

function setHtmlEdit(elTop, fields, data, options, frontOps, htmlButtons, listeners, paramsStr, paramsObj) {
    let formId = formNameBase + ++formNum + "_";
    let elements = [];

    if (!htmlButtons) htmlButtons = "";
    let idView, idDelete;
    if (listeners.view) {
        idView = buttonNameBase + ++buttonNum;
        htmlButtons += `[<a href="#" id="` + idView + `">переглянути (без збереження змін)</a>] &nbsp; `;
    }

    if (listeners.delete) {
        idDelete = buttonNameBase + ++buttonNum;
        htmlButtons += `[<a href="#" id="` + idDelete + `">вилучити</a>]`;
    }

    let htmlEdit = htmlEditFieldsList(formId, prepareFields(fields), data, options, frontOps, elements, listeners);
    elTop.innerHTML = htmlButtons + "\n<p>" + inTable("<tr><td>" +  htmlEdit + "</td></tr>", 7, "#ffffcc");

    for (let e of elements) {
        for (let l in listeners) {
            if (listeners[l].id == e[0]) {
                let el = document.getElementById(e[1]);
                let action = listeners[l].action;
                for (eventType of array(listeners[l].events)) {
                    if (eventType === "init") {
                        action({target: el}, paramsStr, paramsObj);
                    } else {
                        el.addEventListener(eventType, ev => action(ev, paramsStr, paramsObj));
                    }
                }
            }
        }
    }

    if (listeners.view) {
        document.getElementById(idView).addEventListener("click", ev => listeners.view.action(ev, paramsStr, paramsObj));
    }
    if (listeners.delete) {
        document.getElementById(idDelete).addEventListener("click", ev => listeners.delete.action(ev, paramsStr, paramsObj));
    }

    addToOptions("person_id", "aaaaa", "333");
}

function fieldEdit(formId, field, value, options, frontOps)  {
    field = object(field);
    options = object(options);
    frontOps = object(frontOps);

    if ((field.type === "view") || (field.type === "text")) {
        return ["", ...fieldView(field, value, options, frontOps)];
    }

    let [elementId, general] = generalEditParts(formId, field.key, field.attributes);

    let htmlTitle = escapeHtml(field.label);
    let htmlRes = "";

    if (field.type === "password") {
        htmlRes = `<input style="width:100%" type="password" ` + general + ` />`;
    } else if (field.type === "select") {
        htmlRes = htmlSelectEdit(elementId, general, field.params, options[field.key], value);
    } else if (field.type === "checkbox") {
        htmlRes = `<input type="checkbox" ` + general + (value ? " checked" : "") + `/>`;
    } else {
        value = escapeHtml(value);
        if (field.type === "button") {
            htmlRes = `<input type="button" ` + general + ` data-form_id="` + escapeHtml(formId) + `" data-value="` + value + `" value="` + htmlTitle + `" />`;
            htmlTitle = ""
        } else if (field.type === "hidden") {
            htmlRes = `<input type="hidden" ` + general + ` value="` + value + `" /> `;
            htmlTitle = ""
        } else if (field.type === "textarea") {
            htmlRes = `<textarea style="width:100%" ` + general + ` rows=` + field.format + `>` + value + `</textarea>`;
        } else if (field.format === "number") {
            let parameters = ' step="' + (field.params.step ? field.params.step : "1") + '"';
            if ("min" in field.params) { parameters += ' min="' + field.params.min + '"'; }
            if ("max" in field.params) { parameters += ' max="' + field.params.min + '"'; }
            htmlRes = `<input type="number"` + parameters + general + ` value="` + value + `" />`;
        } else if ((field.format === "date") || (field.format === "time") || (field.format === "datetime") || (field.format === "email") || (field.format === "url") || (field.format === "color")) {
            htmlRes = `<input type="` + field.format + `"` + general + ` value="` + value + `" />`;
        } else {
            let t = field.type ? ' type="' + field.type + '"' : "";
            htmlRes = `<input ` + t + ` style="width:100%" ` + general + ` value="` + value + `" />`;
        }
    }

    return [elementId, htmlTitle, htmlRes];
}

function htmlEditFieldsList(formId, fields, data, options, frontOps, elements, listeners, altField, altValue) {
    data = object(data);

    let htmlEdit = "";
    let divAttributes = ' class="ut"';
    let elementId, htmlTitle, htmlRes, alt, altList;
    if (altField) {
        alt = {};
        altList = [];
    }

    for (let f of fields){

        let divId = "div_" + formId + f.key;
        let htmlEditField = "";
        let htmlActions = [];
        let listNumber;

        if (f.params["multiply"]) {
            htmlActions.push(htmlMultiply(formId, divId, f, data[f.key], options, frontOps, elements, listeners));
            listNumber = 0;
        }
        if (f.params["create_new"]) {
            htmlActions.push(htmlCreateNew(formId, divId, f, data[f.key], options, frontOps, elements, listeners));
            listNumber = 0;
        }


        let values = ("alt" in f.params) ? object(data[f.params.alt])[f.key]  : data[f.key];
        if (!(values instanceof Array)) {
            values = [values];
        }

        for (let value of values) {
            let listPrefix = (listNumber === undefined) ? "" : listNumber + "/";

            if (f.subFields instanceof Array) {
                htmlTitle = escapeHtml(f.label);
                elementId = divId;
                if (f.format === "alt") {
                    value = object(value);
                    htmlRes = inTable("<tr><td>" + htmlEditFieldsList(formId + f.key + "/" + listPrefix, f.subFields, value, options, frontOps, elements, listeners, f.key + ".alt", data[f.key + ".alt"]) + "</td></tr>", 5);
                } else {
                    htmlRes = inTable("<tr><td>" + htmlEditFieldsList(formId + f.key + "/" + listPrefix, f.subFields, value, options, frontOps, elements, listeners) + "</td></tr>", 5);
                }
            } else {

                [elementId, htmlTitle, htmlRes] = fieldEdit(formId + listPrefix, f, value, options, frontOps);
                if ((htmlRes === "") && (f.params["not_empty"])) continue;

                elements.push([f.key, elementId]);
            }

            if (htmlTitle !== "") {
                htmlTitle = "<small>" + (htmlTitle ? htmlTitle + ":" : "") + "</small> \n";
            }

            htmlEditField += `<div id="` + escapeHtml(divId) + `"` + divAttributes + ">\n" +
                htmlTitle + htmlRes +
                "</div>\n";

            listNumber++;
        }

        if (htmlActions.length > 0) {
            htmlEditField += "\n<p>" + htmlActions.join("\n<p>");
        }

        if (altField) {
            let isNew = true;
            for (let al of altList) {
                if (al[0] === f.params.alt) {
                    isNew = false;
                    break;
                }
            }
            if (isNew) altList.push([f.params.alt, f.label]);
            alt[f.params.alt] = (alt[f.params.alt] ? alt[f.params.alt] : "") + htmlEditField;
        } else {
            htmlEdit += htmlEditField;
        }
    }

    return altField ? htmlAltEdit(formId, altField, altValue, altList, alt, elements, listeners) : htmlEdit;
}

const altDivSuffix = "_div";

function htmlAltEdit(formId, altField, altValue, altList, alt, elements, listeners) {
    if (altList.length <= 0) {
        return "";
    }

    let altSelected = false;
    for (let e of altList) {
        if (e[0] === altValue) {
            altSelected = true;
            break;
        }
    }
    if (!altSelected) altValue = altList[0][0];

    let altKey = "alt" + ++buttonNum;
    let elementID = formId + altField;

    let htmlAltEdit = "", htmlAltTitleList = [];
    for (let e of altList) {
        let altID = elementID + "_" + e[0];
        let divID = altID + altDivSuffix;

        let [checked, divStyle] = (e[0] === altValue)
                                ? [' checked = "checked"', "visibility:visible;position:static;"]
                                : ["", "visibility:hidden;position:absolute;"];

        htmlAltTitleList.push(`<input id="${altID}" type="radio" name="${elementID}" value="${e[0]}"${checked}> ${e[1] || e[0]}`);
        htmlAltEdit +=  `<div id="${divID}" style="${divStyle}">` + alt[e[0]] + "</div>\n";

        elements.push([altKey, altID])
    }

    let htmlAltTitle = htmlAltTitleList.join(" &nbsp; ");
    listeners[altKey] = {"id": altKey, "events": "change", "action": actionAlt};

    return htmlAltTitle + htmlAltEdit;
}

function actionAlt(event) {
    let radios = document.querySelectorAll(`input[type=radio][name="${event.target.name}"]`);
    for (let el of Array.from(radios)) {
        let elDiv = document.getElementById(el.id + altDivSuffix)
        if (elDiv) {
            [elDiv.style.visibility, elDiv.style.position] = (el.id === event.target.id)
            ? ["visible", "static"] : ["hidden", "absolute"];
        }
    }
}

function htmlCreateNew(formId, divId, f, value, options, frontOps, elements, listeners) {
    let idCreateNew = buttonNameBase + formId + ++buttonNum;
    let keyCreateNew = "create_new" + buttonNum;
    elements.push([keyCreateNew, idCreateNew]);

    // let idCreateNewContent = "create_new_content_" + ++buttonNum;

    let htmlCreateNew = `[<a id="` + idCreateNew + `" style="cursor:pointer;">` +
        (object(f.params)["create_new"] || 'створити новий запис у довіднику') + `</a>]\n`;
        // + `<div id="${idCreateNewContent}" style="visibility:hidden; position:absolute;"></div>`;
    listeners[keyCreateNew] = {
        "id":     keyCreateNew,
        "events": "click",
        "action": ev => actionCreateNew(f.params.genus, divId, f.params["create_new_title"]),
    };

    return htmlCreateNew;
}

function actionCreateNew(genus, divId, title) {
    let elDiv = document.getElementById(divId);
    if (elDiv) {
        let newNode = document.createElement("p");
        elDiv.appendChild(newNode);

        callListener("items.comp", "loadableItemFulfill", {target: newNode}, {genus, title});
    }
}


function htmlMultiply(formId, divId, f, value, options, frontOps, elements, listeners) {
    let idMultiply = buttonNameBase + formId + ++buttonNum;
    let keyMultiply = "multiply" + buttonNum;
    elements.push([keyMultiply, idMultiply]);

    let htmlMultiply = `[<a id="` + idMultiply + `" style="cursor:pointer;">` +
        (object(f.params)["multiply"] || 'додати ще одне поле у формі') + `</a>]\n`;
    listeners[keyMultiply] = {
        "id":     keyMultiply,
        "events": "click",
        "action": actionMultiply(formId, divId, f, value, options, frontOps, elements, listeners),
    };

    return htmlMultiply;
}

function actionMultiply(formId, divId, f, data, options, frontOps, elements, listeners) {
    let indexMultiply = 0;
    return ev => {
        ++indexMultiply;

        let htmlRes;
        let listPrefix = indexMultiply + "/";
        if (f.subFields instanceof Array) {
            htmlRes = inTable("<tr><td>" + htmlEditFieldsList(formId + f.key + "/" + listPrefix, f.subFields, data, options, frontOps, elements, listeners) + "</td></td>", 5);
        } else {
            htmlRes = fieldEdit(formId + listPrefix, f, data, options, frontOps)[2];
        }

        let elDiv = document.getElementById(divId);

        if (elDiv) {
            let newNode = document.createElement("p");
            newNode.innerHTML = htmlRes;
            elDiv.appendChild(newNode);
        }

        return false;
    };
}