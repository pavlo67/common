import b    from '../basis';
import {ep} from '../swagger_convertor';

let cfg = {};

function init(data) {
    if (!(data instanceof Object)) {
        return;
    }

    if ('common' in data) {
        cfg.common = data.common;
    }

    if ('backend' in data) {
        // TODO: do it safely!!!

        cfg.readEp   = ep(data.backend, "read").replace("/{id}", "");
        cfg.saveEp   = ep(data.backend, "save");
        cfg.removeEp = ep(data.backend, "remove").replace("/{id}", "");

        cfg.recentEp = ep(data.backend, "recent");
        cfg.tagsEp   = ep(data.backend, "tags");
        cfg.taggedEp = ep(data.backend, "tagged");

        cfg.exportEp = ep(data.backend, "export");

    }

    if ('eventBus' in data) {
        cfg.eventBus = data.eventBus;
        cfg.eventBus.$on('user', user => {
            // TODO ???
        });
    }

}

export { cfg, init };
