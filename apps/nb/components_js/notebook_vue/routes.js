import Home from './Home.vue';
import ListTags   from './ListTags.vue';
import ListRecent from './ListRecent.vue';
import ListTagged from './ListTagged.vue';
import NoteNew    from './NoteNew.vue';
import NoteView   from './NoteView.vue';
import {init}     from './init';

export default [
    init,
    {
        inMenu: true,
        path: '/',
        name: 'home',
        component: Home,
        title: Home.title,
    },
    {
        inMenu:    true,
        path:      '/notebook/',
        name:      'ListRecent',
        title:     ListRecent.title,
        component: ListRecent,
    },
    // {
    //     inMenu:    true,
    //     path:      '/notebook/tags',
    //     name:      'ListTags',
    //     title:     ListTags.title,
    //     component: ListTags,
    //     // component: () => import(/* webpackChunkName: "data" */ './StorageIndex.vue'),
    // },
    {
        inMenu:    true,
        path:     '/notebook/note_new',
        name:     'NoteNew',
        title:     NoteNew.title,
        component: NoteNew,
    },
    {
        path:      '/notebook/tag/:tag',
        name:      'ListTagged',
        component: ListTagged,
    },
    {
        path:     '/notebook/note/:id',
        name:     'NoteView',
        component: NoteView,
    },
    {
        path:     '/notebook/note_edit/:id',
        name:     'NoteEdit',
        component: () => import('./NoteEdit.vue'),
    },
    {
        path:      '/notebook/note_import',
        name:      'NoteImport',
        component: () => import('./NoteImport.vue'),
    },
];
