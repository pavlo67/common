let cfg = {};

function init(data) {
    cfg.router = data.router;
    // TODO: do it safely!!!
    cfg.listEp = window.location.protocol + "//" + window.location.hostname + data.backend.host + data.backend.endpoints.list.path;
}

export { cfg, init };
