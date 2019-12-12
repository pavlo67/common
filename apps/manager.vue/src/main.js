import Vue from 'vue';
import App from './App.vue';
import Router from 'vue-router';
import './ecosystem/registerServiceWorker';  // import store from './ecosystem/store';

import routes from './parts';

import swagger from '../../workspace/ws_routes/api-docs/swagger';
import swaggerConvertor from '../../components.js/swagger_convertor';

Vue.use(Router);
Vue.config.productionTip = false;

let router = new Router({ routes });  // .map(_ => _.route)
let menuItems = routes.filter(_ => _.inMenu);

let backend = swaggerConvertor(swagger);
for (let r of routes) {
  if (typeof r.init === "function") r.init({router, backend});
}


// let appManager =
new Vue({
  data: { menuItems },
  router,
  // store,
  // render: h => h(App),
  ...App,
}).$mount('#app');



