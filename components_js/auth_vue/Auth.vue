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
  const unauthorized = undefined;
  const whoAmI       = 'я???';

  import b     from '../basis';
  import {cfg} from './init';

  // -------------------------------------------------------------------------------------------------------

  let menuTitle = whoAmI;
  let user;

  function setUser(user) {
    if (user instanceof Object && user.Creds instanceof Object) {
      console.log("CFG.eventBus TO SET USER CREDS: ", 'eventBus' in cfg);

      if (cfg.eventBus) {
        console.log("USER CREDS TO EMIT JWT: ", user.Creds);

        cfg.eventBus.$emit("jwt", user.Creds.jwt);
      }
      menuTitle = user.Creds.nickname;

      if (cfg.vue instanceof Object && typeof cfg.vue.$forceUpdate === "function") cfg.vue.$forceUpdate();
      return user;
    } else {
      menuTitle = whoAmI;
      if (cfg.vue instanceof Object && typeof cfg.vue.$forceUpdate === "function") cfg.vue.$forceUpdate();
      return unauthorized;
    }
  }

  function restoreUser() {
    let user;

    let userStored = localStorage.getItem('user');
    if (userStored) {
      try {
        user = JSON.parse(userStored);
      } catch (err) {
        console.error("can't parse user data from local storage(%s): %s", userStored, err);
      }
    }

    return setUser(user);
  }

  function saveUser(user) {
    localStorage.setItem('user', JSON.stringify(user));

    return setUser(user);
  }

  // -------------------------------------------------------------------------------------------------------

  let first = true;

  export default   {
    preface: 'хто тут:',

    created() {
    },

    title() {
      if (first) {
        user = restoreUser();
        console.log("CREATED!", user);
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
      getUserFromAuth: function(login, password, cb) {
      },

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
            this.user = saveUser(data.user);
          } else {
            console.log("what is the data from /authorize?", data)
          }
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
        this.user = saveUser(unauthorized);
      },
    },
  };


</script>

<style lang="scss">
</style>



<!--<style lang="scss">-->
<!--  #auth {-->
<!--    background-color: #ffffff;-->
<!--    padding: 5px;-->
<!--  }-->
<!--</style>-->
