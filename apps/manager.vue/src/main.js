import Vue from 'vue';
import App from './App.vue';
import router from './router';

// import store from './ecosystem/store';
import './ecosystem/registerServiceWorker';

import routes from './parts';

Vue.config.productionTip = false;


var appManager = new Vue({
  data: { routes, aaa: "asdgfttry" },
  router,
  // store,
  // render: h => h(App),
  ...App,
}).$mount('#app');



