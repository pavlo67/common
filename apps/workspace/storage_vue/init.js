import { ep } from '../../components.js/swagger_convertor';

let cfg = {};

function init(data) {
    if (!(data instanceof Object)) {
        return;
    }

    cfg.router = data.router;
    // TODO: do it safely!!!

    cfg.listEp   = ep(data.backend, "list");
    cfg.readEp   = ep(data.backend, "read").replace("/{id}", "");
    cfg.saveEp   = ep(data.backend, "save");
    cfg.removeEp = ep(data.backend, "remove").replace("/{id}", "");

    cfg.tagsEp   = ep(data.backend, "tags");
    cfg.taggedEp = ep(data.backend, "tagged").replace("/{tag}", "");
}

export { cfg, init };
