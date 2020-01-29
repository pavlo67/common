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
        cfg.listEp = ep(data.backend, "flow");
        cfg.readEp = ep(data.backend, "flow_read").replace("/{id}", "");
    }

    if ('eventBus' in data) {
        cfg.eventBus = data.eventBus;
        cfg.eventBus.$on('user', user => {
            // TODO ???
        });
    }

}

export { cfg, init };
