import Vue from 'vue';
import Router from 'vue-router';

// import './ecosystem/registerServiceWorker';  // import store from './ecosystem/store';

import App       from './App.vue';
import parts     from './parts';
import swagger   from '../../notebook/notebook_actions/api-docs/swagger';
import {convert} from '../../../components_js/swagger_convertor';

Vue.use(Router);
Vue.config.productionTip = false;

let backend = convert(swagger);

let routes  = [];
let menu    = [];
for (let p of parts) {
    if (typeof p === "function") {
        p({backend, eventBus: App.eventBus});
    } else {
        routes.push(p);
        if (p.inMenu) menu.push(p);
    }
}

let router = new Router({ routes });
let vue    = new Vue({
  data: { menu },
  router,
  ...App,
}).$mount('#app');

for (let p of parts) {
  if (typeof p === "function") {
    p({vue});
  }
}


