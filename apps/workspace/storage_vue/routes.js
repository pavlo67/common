import Storage  from './Storage.vue';
import { init } from './init';

export default [
    init,
    {
        inMenu: true,
        path: '/storage',
        name: 'storage',
        component: Storage,  // component: () => import(/* webpackChunkName: "data" */ './Storage.vue'),
        title: Storage.title,
    },
    {
        path: '/storage/item_import',
        name: 'storage_item_import',
        component: () => import('./Import.vue'),
    },
    {
        path: '/storage/item/:id',
        name: 'storage_item',
        component: () => import('./Item.vue'),
    },
];
