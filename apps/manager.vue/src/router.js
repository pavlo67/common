import Vue from 'vue';
import Router from 'vue-router';

import routes from './parts';

Vue.use(Router);

export default new Router({ routes });  // .map(_ => _.route)