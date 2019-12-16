<template>
  <div id="app">
    <div id="nav">

      <Confidence/>
      &nbsp;<br>

      <router-link v-for="item in menu" v-bind:key="item.path" :to="item.path">
        {{ item.title }}<br>
      </router-link>

    </div>

    <div id="view">
      <router-view/>
    </div>
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
    width: 300px;
    font-size: small;
    float: left;
    min-height: 100vh;
    text-align: left;
    background-color: #f2dede;
    a {
      color: #4414ff;
    }

  }

  .title {
    color: blue;
  }

  .small {
    font-size: small;
  }

  a {
    text-decoration: none;
    color: #820cff;
    /*
    &.router-link-exact-active {
     color: #42b983;
    }
    */
  }

  #view {
    margin-left: 320px;
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
