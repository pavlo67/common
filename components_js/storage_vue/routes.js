import StorageIndex  from './StorageIndex.vue';
import StorageTagged from './StorageTagged.vue';
import { init }      from './init';

export default [
    init,
    {
        inMenu:    true,
        path:      '/storage',
        name:      'StorageIndex',
        title:     StorageIndex.title,
        component: StorageIndex,
        // component: () => import(/* webpackChunkName: "data" */ './StorageIndex.vue'),
    },
    {
        path:      '/storage/:tag',
        name:      'StorageTagged',
        component: StorageTagged,
    },
    {
        path:     '/storage_item/:id',
        name:     'StorageItem',
        component: () => import('./StorageItem.vue'),
    },
    {
        path:      '/storage_item_import',
        name:      'StorageItemImport',
        component: () => import('./StorageImport.vue'),
    },
];
