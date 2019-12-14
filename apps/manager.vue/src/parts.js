import Vue from 'vue';

import DataItemView from '../../workspace/data_vue/DataItemView';

console.log(999, DataItemView)

Vue.component(DataItemView.name, DataItemView);

import home    from './home/routes';
import storage from '../../workspace/storage_vue/routes';
import flow    from '../../workspace/flow_vue/routes';


let routes = [
  ...home,
  ...storage,
  ...flow,
];

export default routes;
