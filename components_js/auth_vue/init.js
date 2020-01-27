import { ep } from '../swagger_convertor';

let cfg = {};

function init(data) {
    if (!(data instanceof Object)) {
        return;
    }

    if ('backend' in data) {
        cfg.authorizeEp = ep(data.backend, "authorize");
        cfg.setCredsEp  = ep(data.backend, "set_creds");
        cfg.getCredsEp  = ep(data.backend, "get_creds");
    }

    if ('eventBus' in data) {
        cfg.eventBus = data.eventBus;
        cfg.eventBus.$on('user', user => { cfg.user = user; });
    }

    if ('vue' in data) {
        cfg.vue = data.vue;
    }

}

export { cfg, init };
