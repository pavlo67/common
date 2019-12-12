function init(cfg) {
  // console.log("data: ", cfg)
}

export default [
    {
        inMenu: true,
        path: '/workspace',
        name: 'workspace',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import(/* webpackChunkName: "workspace" */ './Data.vue'),
        title: 'мій каталог',
        init,
    },
    {
        path: '/workspace/item/:id',
        name: 'workspace item',
        component: () => import('./Item.vue'),
        title: 'ось!!!',
    },
];
