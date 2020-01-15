import { ep } from '../../../components.js/swagger_convertor';

let cfg = {};

function init(data) {
    if (!(data instanceof Object)) {
        return;
    }

    cfg.router = data.router;
    cfg.listEp = ep(data.backend, "flow");
    cfg.readEp = ep(data.backend, "flow_read").replace("/{id}", "");
}

export { cfg, init };
