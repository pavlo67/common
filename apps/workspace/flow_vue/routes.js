import Flow     from './Flow.vue';
import { init } from './init';

export default [
    init,
    {
        inMenu: true,
        path: '/flow',
        name: 'flow',
        component: Flow,  // component: () => import(/* webpackChunkName: "flow" */ './Flow.vue'),
        title: Flow.title,
    },
    {
        path: '/flow/select/:id',
        name: 'flow_selector',
        component: () => import('./Select.vue'),
    },
];
