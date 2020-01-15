import Vue from 'vue';

import DataItemView from '../../workspace/data_vue/DataItemView';
import DataList     from '../../workspace/data_vue/DataList';
import TagsIndex    from '../../workspace/data_vue/TagsIndex';

Vue.component(DataItemView.name, DataItemView);
Vue.component(DataList.name,     DataList);
Vue.component(TagsIndex.name,    TagsIndex);

import home    from './home/routes';
import storage from '../../workspace/storage_vue/routes';
import flow    from '../../workspace/flow_vue/routes';

let routes = [
  ...home,
  ...storage,
  ...flow,
];

export default routes;
