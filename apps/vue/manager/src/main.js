import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import './registerServiceWorker';

import routes from './routes';


Vue.config.productionTip = false;


new Vue({
  data: { routes },
  router,
  store,
  // render: h => h(App),
  ...App,
}).$mount('#app');



