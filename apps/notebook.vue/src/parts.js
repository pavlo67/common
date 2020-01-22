import Vue from 'vue';

import DataItemView from '../../../components_js/data_vue/DataItemView';
import DataList     from '../../../components_js/data_vue/DataList';
import TagsIndex    from '../../../components_js/data_vue/TagsIndex';

Vue.component(DataItemView.name, DataItemView);
Vue.component(DataList.name,     DataList);
Vue.component(TagsIndex.name,    TagsIndex);

import auth    from '../../../components_js/auth_vue/routes';
import home    from './home/routes';
import storage from '../../../components_js/storage_vue/routes';

let routes = [
  ...auth,
  // ...home,
  // ...storage,
];

export default routes;
