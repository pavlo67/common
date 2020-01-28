import { ep } from '../swagger_convertor';

let cfg = {};

function init(data) {
    if (!(data instanceof Object)) {
        return;
    }

    if ('common' in data) {
        cfg.common = data.common;
    }

    if ('backend' in data) {
        cfg.authorizeEp = ep(data.backend, "authorize");
        cfg.setCredsEp  = ep(data.backend, "set_creds");
        cfg.getCredsEp  = ep(data.backend, "get_creds");
    }

    if ('eventBus' in data) {
        cfg.eventBus = data.eventBus;
    }
}

export { cfg, init };
