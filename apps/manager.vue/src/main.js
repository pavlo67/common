import Vue from 'vue';
import App from './App.vue';
import Router from 'vue-router';
import './ecosystem/registerServiceWorker';  // import store from './ecosystem/store';

import parts from './parts';

import swagger     from '../../workspace/ws_routes/api-docs/swagger';
import { convert } from '../../components.js/swagger_convertor';

Vue.use(Router);
Vue.config.productionTip = false;

let inits  = [];
let routes = [];
let menu   = [];
for (let p of parts) {
    if (typeof p === "function") {
        inits.push(p);
    } else {
        routes.push(p);
        if (p.inMenu) menu.push(p);
    }
}

let backend = convert(swagger);
let router  = new Router({ routes });
for (let i of inits) i({router, backend});

new Vue({
  data: { menu },
  router,
  ...App,
}).$mount('#app');



