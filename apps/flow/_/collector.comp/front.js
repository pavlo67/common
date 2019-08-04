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

function unavailableChoice(id){
    document.getElementById("from_import_id_" + id + " none").checked = true;
    document.getElementById("from_import_id_" + id + " accept").setAttribute("disabled","disabled");
    document.getElementById("from_import_id_" + id + " hide").setAttribute("disabled","disabled");
    document.getElementById("from_import_id_" + id + " delete").setAttribute("disabled","disabled");

    document.getElementById("one_object_import_id_" + id).innerText="";
}

function oneObjectImport() {
    let choice = event.target.id;
    choice = choice.replace(/^one_object_import_id_/, "");
    let url = $("#path_note_blank").val() + choice;
    unavailableChoice(choice);
    // alert("will be import: " + choice)
    let win = window.open(url, '_blank');
    win.focus();
}

function selectAllImport() {
    let choice = event.target.id;
    choice = choice.replace(/^select_all_as_/, "");
    if (choice !== '') {
        let selectors = document.querySelectorAll(`[id^="from_import_id_"]`);
        for (let el of selectors) {
            if (el.value === choice) {
                el.checked = true;
            }
        }
    }
}

function importTest() {

    let elId = event.target.id;
    let li = listenersIndex.importTest;
    let epElement = li.id;
    let regexp = '^' + epElement.slice(0, -1);
    let re = new RegExp(regexp, "i");
    let id = elId.replace(re, "");
    // alert("Click id:" + id);
    let ep = endpointsIndex.importTest;
    if (ep) {     // && ep.method === 'post'
        $.ajaxSetup({
            headers: {
                'Content-Type': 'application/json',
            },
        });
        $.post(
            ep["server_path"],
            JSON.stringify({id: id}),
            function(data) {
                console.log(data);
                if (data.hasOwnProperty("Info")) {
                    saveMessageToLocalStorage(data.Info, true);
                    if (data.hasOwnProperty("Redirect")) {
                        _fRedirect(data.Redirect);
                    } else {
                        window.location.reload();
                    }
                } else {
                    message(data.Error, false, false);
                }
            });

    } else {
        console.error('??? endpointsIndex.createFount', endpointsIndex);
    }

}

function createFount(event, params){
    let data = getData(params["form_id"]);

    //
    // let err = false;
    // let nil = '';
    // let title = $("#" + "title").val();
    // if (title.length < minTitleLength){
    //     nil += "Title is short! ";
    //     err = true;
    // }
    // let url = $("#" + "url").val();
    // if (url.length < 5){
    //     nil += "Url is short! ";
    //     err = true;
    // }
    // let direct = $("#" + "direct").val();
    // let tags = $("#" + "tags").val();
    // let type = $("#" + "type").val();
    //
    //
    //
    // if (err) {
    //     message(nil, false, false)
    // } else {
    // JSON.stringify({title: title, url: url, direct: direct, tags: tags, type: type}),

    let ep = endpointsIndex.createFount;
    if (ep) {     // && ep.method === 'post'
        $.ajaxSetup({
            headers: {
                'Content-Type': 'application/json',
            },
        });
        $.post(
            ep["server_path"],
            JSON.stringify(data),
            function(data) {
                console.log(data);
                if (data.hasOwnProperty("Info")) {
                    // alert(data.Info);
                    // localStorage.setItem('success_message', data.Info);
                    saveMessageToLocalStorage(data.Info, true);
                    if (data.hasOwnProperty("Redirect")) {
                        _fRedirect(data.Redirect);
                    } else {
                        window.location.reload();
                    }
                } else {
                    message(data.Error, false, false)
                }
            });

    } else {
        console.error('??? endpointsIndex.createFount', endpointsIndex);
    }
}

function testRegexp() {
    let id = $("#" + "id").val();
    let url = $("#" + "url").val();
    let start = $("#" + "import_start_regexp").val();
    let finish = $("#" + "import_finish_regexp").val();
    let split = $("#" + "import_split_regexp").val();
    let tagsList = $("#" + "import_tags_list").val();
    let title = $("#" + "title_regexp").val();
    let ep = endpointsIndex.fountTestSetting;
    if (ep) {     // && ep.method === 'post'
        $.ajaxSetup({
            headers: {
                'Content-Type': 'application/json',
            },
        });
        $.post(
            ep["server_path"],
            JSON.stringify({
                id: id,
                url: url,
                import_start_regexp: start,
                import_finish_regexp: finish,
                import_split_regexp: split,
                import_tags_list: tagsList,
                title_regexp: title,
            }),
            function(data) {
                console.log(data);
                if (data.hasOwnProperty("Info")) {
                    let w = window.open('about:blank');
                    w.document.open();
                    w.document.write(data.Info);
                    w.document.close();
                } else {
                    message(data.Error, false, false);
                }
            });

    } else {
        console.error('??? endpointsIndex.fountTestSetting', endpointsIndex);
    }

}


function exportRegexp() {

    let id = $("#" + "id").val();
    let url = $("#" + "url").val();
    let type = $("#" + "type").val();
    let start = $("#" + "import_start_regexp").val();
    let finish = $("#" + "import_finish_regexp").val();
    let split = $("#" + "import_split_regexp").val();
    let tagsList = $("#" + "import_tags_list").val();
    let title = $("#" + "title_regexp").val();
    let ep = endpointsIndex.fountSettings;
    if (ep) {     // && ep.method === 'post'
        $.ajaxSetup({
            headers: {
                'Content-Type': 'application/json',
            },
        });
        $.post(
            ep["server_path"],
            JSON.stringify({
                id: id,
                url: url,
                type: type,
                import_start_regexp: start,
                import_finish_regexp: finish,
                import_split_regexp: split,
                import_tags_list: tagsList,
                title_regexp: title,
            }),
            function(data) {
                console.log(data);
                if (data.hasOwnProperty("Info")) {
                    saveMessageToLocalStorage(data.Info, true);
                    if (data.hasOwnProperty("Redirect")) {
                        _fRedirect(data.Redirect);
                    } else {
                        window.location.reload();
                    }
                } else {
                    message(data.Error, false, false);
                }
            });

    } else {
        console.error('??? endpointsIndex.fountSettings', endpointsIndex);
    }

}


function updateFount(event, params){
    let data = getData(params["form_id"]);
    //
    // let err = false;
    // let nil = '';
    // let id = $("#" + "id").val();
    // let title = $("#" + "title").val();
    // if (title.length < minTitleLength){
    //     nil += "Title is short! ";
    //     err = true;
    // }
    // let url = $("#" + "url").val();
    // if (url.length < 5){
    //     nil += "Url is short! ";
    //     err = true;
    // }
    // let direct = $("#" + "direct").val();
    // let tags = $("#" + "tags").val();
    // let type = '';
    // let elType = $("#" + "type");
    // if (elType != null){
    //     type = elType.val();
    // }
    // let IDT = $("#" + "import_details_type").val();
    // let IDPs = $("#" + "import_details_params").val();
    // JSON.stringify({
    //     id: id,
    //     title: title,
    //     url: url,
    //     direct: direct,
    //     tags: tags,
    //     type: type,
    //     import_details_type: IDT,
    //     import_details_params: IDPs,
    // }),



    let ep = endpointsIndex.updateFount;
    if (ep) {     // && ep.method === 'post'
        $.ajaxSetup({
            headers: {
                'Content-Type': 'application/json',
            },
        });
        $.post(
            ep["server_path"],
            JSON.stringify(data),
            function(data) {
                console.log(data);
                if (data.hasOwnProperty("Info")) {
                    saveMessageToLocalStorage(data.Info, true);
                    if (data.hasOwnProperty("Redirect")) {
                        _fRedirect(data.Redirect);
                    } else {
                        window.location.reload();
                    }
                } else {
                    message(data.Error, false, false);
                }
            });

    } else {
        console.error('??? endpointsIndex.createFount', endpointsIndex);
    }

}

function deleteFount(){

    let elId = event.target.id;
    let li = listenersIndex.deleteFount;
    let epElement = li.id;
    let regexp = '^' + epElement.slice(0, -1);
    let re = new RegExp(regexp, "i");
    let id = elId.replace(re, "");

    if (confirm('Підтвердіть видалення запису: ' + id)){
        let ep = endpointsIndex.deleteFount;
        $.post(
            // ep["server_path"].slice(0,-3) + id,
            routerPath(ep["server_path"], id),
            function(data) {
                console.log(data);
                if (data.hasOwnProperty("Info")) {
                    saveMessageToLocalStorage(data.Info, true);
                    if (data.hasOwnProperty("Redirect")) {
                        _fRedirect(data.Redirect);
                    } else {
                        window.location.reload();
                    }
                } else {
                    message(data.Error, false, false);
                }
            });
    }

}

function authorsFromImport(){
    let author = $("#" + "select_author").val();
    let ep = endpointsIndex.importFlows;
    _fRedirect(ep["server_path"] + "?author" + '=' + author);
    // alert('click: ' + author)
}

function oneObjectAccept() {
    // let elId = event.target.id;
    // let linker_server_http = listenersIndex.oneObjectAccept;
    // let epElement = linker_server_http.id;
    // let regexp = '^' + epElement.slice(0, -1);
    // let re = new RegExp(regexp, "i");
    // let id = elId.replace(re, "");
    //
    // let noteId = "from_import_id" + id;
    // let items.comp = document.querySelectorAll(`[id^="${noteId}"]`);
    // // alert("id: " + id + items.comp.item(0).id);
    // let importNotes = new Array();
    // if (items.comp.length > 0) {
    //     for (let i = 0; i < items.comp.length; i++) {
    //         if (document.getElementById(items.comp.item(i).id).checked) {
    //             let regexp = '^' + "from_import_id";
    //             let re = new RegExp(regexp, "i");
    //             let id = items.comp.item(i).id.replace(re, "");
    //             re = new RegExp(' (accept|delete|hide|none)$', "i");
    //             let myAct = id.match(re);
    //             if (myAct.length > 1) {
    //                 let idRow = id.replace(re, "");
    //                 importNotes.push({id: idRow, action: myAct[1]});
    //             }
    //         }
    //     }
    //     // alert(importNotes[0]);
    //     let ep = endpointsIndex.acceptImport.id;
    //     if (ep) {     // && ep.method === 'post'
    //         $.ajaxSetup({
    //             headers: {
    //                 'Content-Type': 'application/json',
    //             },
    //         });
    //         $.post(
    //             ep,
    //             JSON.stringify(importNotes),
    //             function (data) {
    //                 console.log(data);
    //                 if (data.hasOwnProperty("Info")) {
    //                     alert(data.Info);
    //                     unavailableChoice(id);
    //                     // localStorage.setItem('success_message', data.Info);
    //                     // window.location.reload();
    //                     // $(window).scrollTop(0);
    //                 } else {
    //                     alert(data.Error);
    //                     // document.getElementById('error_message').className = 'alert';
    //                     // document.getElementById('error_message').innerText = data.Error;
    //                 }
    //             });
    //
    //     } else {
    //         console.error('??? endpointsIndex.acceptImport', endpointsIndex);
    //     }
    //
    // }
    //
}

function importFlow(){
    let elId = event.target.id;
    let li = listenersIndex.importFlow;
    let epElement = li.id;
    let regexp = '^' + epElement.slice(0, -1);
    let re = new RegExp(regexp, "i");
    let id = elId.replace(re, "");
    let ep = routerPath(endpointsIndex.importFlow["server_path"], id);
      // alert("will be import call=" + ep);
    if (ep) {     // && ep.method === 'post'
        $.get(
            ep,
            function (data) {
                console.log(data);
                if (data.hasOwnProperty("Info")) {
                    saveMessageToLocalStorage(data.Info, true);
                    // localStorage.setItem('success_message', data.Info);
                    if (data.hasOwnProperty("Redirect")) {
                        _fRedirect(data.Redirect);
                    } else {
                        window.location.reload();
                    }
                } else {
                    message(data.Error, false, false);
                }
            });

    } else {
        console.error('??? endpointsIndex.importFlow', endpointsIndex);
    }

}

function acceptImport() {
    // let importNotes = new Array();
    // let noteId = "from_import_id" + '*';
    // let items.comp = (noteId[noteId.length - 1] === '*')
    //     ? document.querySelectorAll(`[id^="${noteId.substr(0, noteId.length - 1)}"]`)
    //     : document.querySelectorAll(`[id="${noteId}"]`);
    //
    // if (items.comp.length > 0) {
    //     for (let i = 0; i < items.comp.length; i++) {
    //         if (document.getElementById(items.comp.item(i).id).checked) {
    //             let regexp = '^' + "from_import_id";
    //             let re = new RegExp(regexp, "i");
    //             let id = items.comp.item(i).id.replace(re, "");
    //             re = new RegExp(' (accept|delete|hide|none)$', "i");
    //             let myAct = id.match(re);
    //             if (myAct.length > 1) {
    //                 let idRow = id.replace(re, "");
    //                 importNotes.push({id: idRow, action: myAct[1]});
    //             }
    //         }
    //     }
    //     // alert(importNotes[0]);
    //     let ep = fieldsIndex.acceptImport.id;
    //     if (ep) {     // && ep.method === 'post'
    //         $.ajaxSetup({
    //             headers: {
    //                 'Content-Type': 'application/json',
    //             },
    //         });
    //         $.post(
    //             ep,
    //             JSON.stringify(importNotes),
    //             function (data) {
    //                 console.log(data);
    //                 if (data.hasOwnProperty("Info")) {
    //                     // alert(data.Info);
    //                     localStorage.setItem('success_message', data.Info);
    //                     window.location.reload();
    //                     $(window).scrollTop(0);
    //                 } else {
    //                     document.getElementById('error_message').className = 'alert';
    //                     document.getElementById('error_message').innerText = data.Error;
    //                 }
    //             });
    //
    //     } else {
    //         console.error('??? endpointsIndex.acceptImport', endpointsIndex);
    //     }
    //
    // }
    // // alert(importNotes);

}



function _fRedirect(url) {
    window.location = url;
    window.setTimeout(function(){window.location.reload()},2000)
}
