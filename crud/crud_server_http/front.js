function editInTable(e, data) {
    let parentNode = e.target.parentNode;
    let parts = parentNode.innerHTML.split("&nbsp;<a");
    let attributes = `data-key="${escapeHtml(data.key)}"`;
    for (let k in data) {
        if (k.match(/^id(_\w+)?$/)) {
            attributes += `data-${escapeHtml(k)}="${escapeHtml(data[k])}"`;
        }
    }
    // parentNode.innerHTML = `<input id="editedInTable" ${attributes} value="${escapeHtml(parts[0])}">`;
    parentNode.innerHTML = `<textarea id="editedInTable" ${attributes}>${escapeHtml(parts[0])}</textarea>`;

}

function saveEditedInTable(e, data) {

    // TODO: restrict by table id

    let elements =  document.querySelectorAll(`[id="editedInTable"]`);
    let dataREST = [];

    for (let el of Array.from(elements)) {
        let id = {};
        for (let k in el.dataset) {
            if (k.match(/^id(_\w+)?$/)) {
                id[k] = el.dataset[k];
            }
        }
        dataREST.push({
            value: el.value,
            key: el.dataset.key,
            id: id,
        });
    }

    // let formId = params["form_id"] ? params["form_id"] : getFormId(e.target.id);

    let ep = endpointsIndex.update_list_rest;
    if (!ep) {     // && ep.method === 'post'
        console.error('??? endpointsIndex.update_list_rest', endpointsIndex);
        return;
    }

    log(2222222, JSON.stringify(dataREST))

    $.post(routerPath(ep["server_path"], data.crud_type), JSON.stringify(dataREST), afterSave);

}


function afterSave(res) {
    if (res.hasOwnProperty("info")) {
        saveMessageToLocalStorage(res.info, true);
        if (res.hasOwnProperty("redirect")) {
            window.location = res.redirect;
        } else {
            window.location.reload();
        }
    } else {
        let errors = array(res.error);
        message(["не вдалось оновити записи:-(", ...errors], false, false);
    }
}
