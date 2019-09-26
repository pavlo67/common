<template>
  <div id="app">
    <div id="nav">

      <Confidence/>
      &nbsp;<br>

      <router-link v-for="item in routes" v-bind:key="item.path" :to="item.path">
        {{ item.title }}<br>
      </router-link>
      <br>&nbsp;<br> {{ aaa }}

    </div>
    <router-view/>
  </div>
</template>

<style lang="scss">
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
#nav {
  padding: 10px;
  float: left;
  width: 300px;
  min-height: 100vh;
  text-align: left;
  background-color: #f2dede;
  a {
    text-decoration: none;
    color: #2c3e50;
    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>


<script>
  import Confidence from '../../confidence/_vue/Confidence.vue';

  export default {
    components: {
      Confidence,
    },
  };



</script>
