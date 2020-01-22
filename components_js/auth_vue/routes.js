import Confidence  from './Auth.vue';
import { init }    from './init';

export default [
    init,
    {
        inMenu:    true,
        path:      '/auth',
        name:      'Confidence',
        title:     Confidence.title,
        component: Confidence,
    },
];
