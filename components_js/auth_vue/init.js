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
        cfg.eventBus.$on('jwt', jwt => {
            cfg.jwt = jwt;

            console.log(666666666, cfg.jwt);
        });
    }

    if ('vue' in data) {
        cfg.vue = data.vue;
    }

}

export { cfg, init };
