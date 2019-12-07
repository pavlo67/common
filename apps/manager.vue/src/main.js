import Vue from 'vue';
import App from './App.vue';
import Router from 'vue-router';
import './ecosystem/registerServiceWorker';  // import store from './ecosystem/store';

import routes from './parts';

import swagger from '../../workspace/routes/api-docs/swagger';
import swaggerConvertor from './lib/swagger_convertor';

let endpoints = swaggerConvertor(swagger);

for (let r of routes) {
  if (typeof r.init === "function") r.init(endpoints);
}

Vue.use(Router);

let router = new Router({ routes });  // .map(_ => _.route)

Vue.config.productionTip = false;

// let appManager =
new Vue({
  data: { routes },
  router,
  // store,
  // render: h => h(App),
  ...App,
}).$mount('#app');



