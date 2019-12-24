<template>
  <div id="app">
    <div id="nav">

      <Confidence/>
      &nbsp;<br>

      <router-link v-for="item in menu" v-bind:key="item.path" :to="item.path">
        {{ item.title }}<br>
      </router-link>

    </div>

    <div id="message">
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

  .time {
    color: brown;
  }

  .control {
    color: blue;
    font-size: xx-small;
  }

  .small {
    font-size: small;
  }

  .smaller {
    font-size: x-small;
  }

  .smallest {
    font-size: xx-small;
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
    margin: 0px 0px 10px 330px;
    padding: 10px;
    text-align: center;
    font-size: small;
    border-color: green;
    border-width: 1px;
    border-style: solid;
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

  function show(id, html, color) {
    let el = document.getElementById(id);
    if (el) {
      el.innerHTML = html;
      el.style.position = "relative";
      el.style.visibility = "visible";
      el.style["border-color"] = color || "green";
    }
  }

  function hide(id) {
    let el = document.getElementById(id);
    if (el) {
      el.style.visibility = "hidden";
      el.style.position = "absolute";
    }
  }

  export { eventBus };

  export default {
    mounted() {
      eventBus.$on('message', message => {
        show("message", message);
        setTimeout(() => { hide("message"); }, 3000);
      });
      eventBus.$on('error', message => {
        show("message", message, "red");
        setTimeout(() => { hide("message"); }, 3000);
      });
    },

    components: {
      Confidence,
    },
  };
</script>
