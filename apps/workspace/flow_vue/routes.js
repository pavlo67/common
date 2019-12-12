import Flow from './Flow.vue';

// import { init } from './Flow.vue';
// console.log(Flow)

export default [
    {
        inMenu: true,
        path: '/flow',
        name: 'flow',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        // component: () => import(/* webpackChunkName: "flow" */ './Flow.vue'),
        component: Flow,
        title: 'новини',
        init: Flow.methods.init,
    },
    {
        path: '/flow/select/:id',
        name: 'flow_selector',
        component: () => import('./Select.vue'),
        title: 'новини',
    },
];
