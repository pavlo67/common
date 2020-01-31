import Auth   from './Auth.vue';
import {init} from './init';

export default [
    init,
    {
        inMenu:    true,
        path:      '/auth',
        name:      'Confidence',
        preface:   Auth.preface,
        title:     Auth.title,
        component: Auth,
    },
];
