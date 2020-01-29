import Vue from 'vue';
import Router from 'vue-router';
Vue.use(Router);
Vue.config.productionTip = false;


// backend configuration ---------------------------------------------------------------------------

// TODO: use config to find the appropriate swagger file

import swagger   from '../../api-docs/swagger';
import {convert} from '../../components_js/swagger_convertor';
let backend = convert(swagger);


// common components -------------------------------------------------------------------------------

import ActionSteps  from '../../components_js/helpers_vue/ActionSteps';
import DataItemView from '../../components_js/data_vue/DataView';
import DataItemEdit from '../../components_js/data_vue/DataEdit';
import DataList     from '../../components_js/data_vue/DataList';
import DataListTags from '../../components_js/data_vue/DataListTags';
import VueTagsInput from '@johmun/vue-tags-input';

Vue.component(ActionSteps.name,    ActionSteps);
Vue.component(DataItemView.name,   DataItemView);
Vue.component(DataItemEdit.name,   DataItemEdit);
Vue.component(DataList.name,       DataList);
Vue.component(DataListTags.name,   DataListTags);
Vue.component("vue-tags-input", VueTagsInput);


// application parts initiated with backend and event bus ------------------------------------------

import auth     from '../../components_js/auth_vue/routes';
import notebook from '../../components_js/notebook_vue/routes';

let parts = [
  ...auth,
  ...notebook,
];

let common = {};
let routes = [];
let menu   = [];
for (let p of parts) {
    if (typeof p === "function") {
        p({common, backend, eventBus: App.eventBus});
    } else {
        routes.push(p);
        if (p.inMenu) menu.push(p);
    }
}


// app start ---------------------------------------------------------------------------------------

import App from './App.vue';
let router = new Router({routes});
common.vue = new Vue({
  data: { menu },
  router,
  ...App,
}).$mount('#app');


