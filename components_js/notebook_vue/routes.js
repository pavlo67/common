import ListTags   from './ListTags.vue';
import ListRecent from './ListRecent.vue';
import ListTagged from './ListTagged.vue';
import ItemNew    from './ItemNew.vue';
import ItemView   from './ItemView.vue';
import {init}     from './init';

export default [
    init,
    {
        inMenu:    true,
        path:      '/notebook/',
        name:      'ListRecent',
        title:     ListRecent.title,
        component: ListRecent,
    },
    {
        inMenu:    true,
        path:      '/notebook/tags',
        name:      'ListTags',
        title:     ListTags.title,
        component: ListTags,
        // component: () => import(/* webpackChunkName: "data" */ './StorageIndex.vue'),
    },
    {
        inMenu:    true,
        path:     '/notebook/item_new',
        name:     'ItemNew',
        title:     ItemNew.title,
        component: ItemNew,
    },
    {
        path:      '/notebook/tag/:tag',
        name:      'ListTagged',
        component: ListTagged,
    },
    {
        path:     '/notebook/item/:id',
        name:     'ItemView',
        component: ItemView,
    },
    {
        path:     '/notebook/item_edit/:id',
        name:     'ItemEdit',
        component: () => import('./ItemEdit.vue'),
    },
    {
        path:      '/notebook/item_import',
        name:      'ItemImport',
        component: () => import('./ItemImport.vue'),
    },
];
