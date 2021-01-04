<template>
  <div id="auth">
    <div class="title"><b>Авторизація</b></div>

    <span v-if="user">
      {{ user.Creds.nickname }} &nbsp; <button v-on:click="signOut">Вийти</button>
      <!-- &nbsp; <button v-on:click="checkIn">?</button> -->
    </span>
    <span v-else>
      <input v-model="inputLogin"    style="width:80px;margin-right:5px;">
      <input v-model="inputPassword" style="width:80px;margin-right:5px;" type="password">
      <button v-on:click="signIn">»</button>
      <br>забули-сте пароль?
    </span>
  </div>
</template>

<script>
  const whoAmI = 'я???';

  // import b     from '../basis';
  import {cfg} from './init';

  // -------------------------------------------------------------------------------------------------------

  let menuTitle = whoAmI;

  function setUser(user) {
    if (user instanceof Object && user.Creds instanceof Object) {
      user.groups = []; // TODO!!!

      menuTitle = user.Creds.nickname;
      cfg.common.user = user;

    } else {
      menuTitle = whoAmI;
      cfg.common.user = null;

    }

    // console.log("CFG.eventBus TO SET USER CREDS: ", 'eventBus' in cfg);
    if (cfg.eventBus) {
      console.log("USER TO BE EMITTED: ", cfg.common.user);
      cfg.eventBus.$emit("user", cfg.common.user);
    }
    if (cfg.common instanceof Object && cfg.common.vue instanceof Object) {
      cfg.common.vue.$forceUpdate();
    }
  }

  function restoreUser() {
    let user;

    let userStored = localStorage.getItem('user');
    if (userStored && userStored.length > 0) {

      try {
        user = JSON.parse(userStored);
      } catch (err) {
        console.error("can't parse user data from local storage(%s): %s", userStored, err);
      }
    }

    setUser(user);
  }

  function saveUser(user) {
    setUser(user);

    localStorage.setItem('user', cfg.common.user ? JSON.stringify(cfg.common.user) : null);
  }

  // -------------------------------------------------------------------------------------------------------

  let first = true;
  let user  = (cfg.common && cfg.common.user) || {};

  export default   {
    preface: 'хто тут:',

    title() {
      // TODO!!! remove the kostyl (move this initiation in common init() or somewhere looks like to)
      if (first) {
        restoreUser();
        user = cfg.common.user;

        first = false;
      }
      return menuTitle;
    },

    data: () => {
      return {
        title: menuTitle,
        inputLogin: "",
        inputPassword: "",
        user,
      };
    },
    methods: {
      signIn: function () {
        fetch(cfg.authorizeEp, {
          method: 'POST', // *GET, POST, PUT, DELETE, etc.
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({login: this.inputLogin, password: this.inputPassword}),
          mode: 'cors', // no-cors, cors, *same-origin

          // cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
          // credentials: 'same-origin', // include, *same-origin, omit
          // redirect: 'follow', // manual, *follow, error
          // referrer: 'no-referrer', // no-referrer, *client
        }).then(response => {
          return response.json();
        }).then(data => {
          if (data instanceof Object) {
            saveUser(data.user);

          } else {
            console.log("what is the data from /authorize?", data)
            saveUser(null);
          }

          this.user = user = cfg.common.user;
        });
      },

      // checkIn: function() {
      //   fetch(cfg.getCredsEp, {
      //     method: 'POST',
      //     headers: {
      //       'content-type' : 'application/json',
      //       'authorization': cfg.jwt,
      //     },
      //     // mode: 'cors',
      //   }).then(response => {
      //     return response.json();
      //   }).then(data => {
      //     console.log(777777777, data);
      //   });
      // },

      signOut: function () {
        saveUser(null);
        this.user = user = cfg.common.user;

      },
    },
  };


</script>

<style lang="scss">
</style>

