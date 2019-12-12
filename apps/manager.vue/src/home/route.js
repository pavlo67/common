import Home from './Home.vue';

function init(cfg) {
  // console.log("home!: ", cfg)
}

export default {
  inMenu: true,
  path: '/',
  name: 'home',
  component: Home,
  title: 'ой, мамо, де я?',
  init,
};
