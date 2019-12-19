<template>
  <div id="app">
    <div id="nav">

      <Confidence/>
      &nbsp;<br>

      <router-link v-for="item in menu" v-bind:key="item.path" :to="item.path">
        {{ item.title }}<br>
      </router-link>

    </div>

    <!--<div id="message">-->
    <!--</div>-->

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

  .control {
    color: blue;
    font-size: xx-small;
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

  #message {
    margin-left: 320px;
    padding: 10px;
    text-align: center;
    font-size: small;
    background-color: #BC6060;
    position: absolute;
    visibility: hidden;
  }

  #view {
    margin-left: 320px;
    padding: 0px 10px 10px 10px;
    text-align: left;
    font-size: small;
  }

</style>


<script>
  import Vue from 'vue';
  import Confidence from '../../confidence/_vue/Confidence.vue';

  let eventBus = new Vue();

  function show(id, html) {
    let el = document.getElementById(id);
    if (el) {
      el.innerHTML = html;
      el.style.position = "relative";
      el.style.visibility = "visible";
    }
  }

  function hide(id) {
    let el = document.getElementById(id);
    if (el) {
      el.style.visibility = "hidden";
      el.style.position = "absolute";
    }
  }

  export default {
    eventBus,
    mounted() {
      eventBus.$on('message', message => {
        show("message", "received: " + message);
        setTimeout(() => { hide("message"); }, 3000);
      });

      eventBus.$emit('message', "!!!");
    },
    data: () => {
      return {
        message: '',
      };
    },

    components: {
      Confidence,
    },
  };
</script>
