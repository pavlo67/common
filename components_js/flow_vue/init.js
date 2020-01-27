import { ep } from '../swagger_convertor';

let cfg = {};

function init(data) {
    if (!(data instanceof Object)) {
        return;
    }

    if ('router' in data) {
        cfg.router = data.router;
    }

    if ('backend' in data) {
        cfg.listEp = ep(data.backend, "flow");
        cfg.readEp = ep(data.backend, "flow_read").replace("/{id}", "");
    }

    if ('eventBus' in data) {
        cfg.eventBus = data.eventBus;
        cfg.eventBus.$on('user', user => {
            cfg.eventBus.$on('user', user => { cfg.user = user; });
        });
    }

    if ('vue' in data) {
        cfg.vue = data.vue;
    }

}

export { cfg, init };
