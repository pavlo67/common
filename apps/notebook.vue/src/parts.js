import Vue from 'vue';

import DataItemView from '../../../components_js/data_vue/DataView';
import DataItemEdit from '../../../components_js/data_vue/DataEdit';
import DataList     from '../../../components_js/data_vue/DataList';
import DataListTags from '../../../components_js/data_vue/DataListTags';
import VueTagsInput from '@johmun/vue-tags-input';

Vue.component(DataItemView.name,   DataItemView);
Vue.component(DataItemEdit.name,   DataItemEdit);
Vue.component(DataList.name,       DataList);
Vue.component(DataListTags.name,   DataListTags);
Vue.component("vue-tags-input", VueTagsInput);

import auth    from '../../../components_js/auth_vue/routes';
import home    from './home/routes';
import notebook from '../../../components_js/notebook_vue/routes';

let routes = [
  ...auth,
  ...home,
  ...notebook,
];

export default routes;
