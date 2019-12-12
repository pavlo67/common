import Data     from './Data.vue';
import { init } from './init';

export default [
    init,
    {
        inMenu: true,
        path: '/data',
        name: 'data',
        component: Data,  // component: () => import(/* webpackChunkName: "data" */ './Data.vue'),
        title: Data.title,
    },
    {
        path: '/data/item_import',
        name: 'data_item_import',
        component: () => import('./Import.vue'),
    },
    {
        path: '/data/item/:id',
        name: 'data_item',
        component: () => import('./Item.vue'),
    },
];
