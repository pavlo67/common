import Vue from 'vue';

import DataItemView from '../../notebook/data_vue/DataItemView';
import DataList     from '../../notebook/data_vue/DataList';
import TagsIndex    from '../../notebook/data_vue/TagsIndex';

Vue.component(DataItemView.name, DataItemView);
Vue.component(DataList.name,     DataList);
Vue.component(TagsIndex.name,    TagsIndex);

import home    from './home/routes';
import storage from '../../notebook/storage_vue/routes';
import flow    from '../../notebook/flow_vue/routes';

let routes = [
  ...home,
  ...storage,
  ...flow,
];

export default routes;
