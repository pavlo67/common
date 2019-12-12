let cfg = {};

function init(data) {
    cfg.router = data.router;
    // TODO: do it safely!!!
    cfg.listEp = window.location.protocol + "//" + window.location.hostname + data.backend.host + data.backend.endpoints.flow.path;
    cfg.readEp = window.location.protocol + "//" + window.location.hostname + data.backend.host + data.backend.endpoints.flow_read.path.replace("/{id}", "");
}

export { cfg, init };
