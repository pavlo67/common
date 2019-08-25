import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import './registerServiceWorker';

Vue.config.productionTip = false;

let items = [
  { path: "/",      title: 'Home' },
  { path: "/about", title: 'About' }
];


new Vue({
  data: { items },
  router,
  store,
  // render: h => h(App),
  ...App,
}).$mount('#app');


// vm.items = items;

// console.log(vm.items);

