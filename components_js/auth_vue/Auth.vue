<template>
  <div id="auth">
    <span v-if="user">
      {{ user.nickname }}
      <br><button v-on:click="signOut">Вийти</button>
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
  import b       from '../basis';
  import { cfg } from './init';

  const unauthorized = undefined;
  const whoAmI       = 'хто я?';

  let title = whoAmI;
  let user  = restoreUser();

  function setUser(user) {
    if (user instanceof Object) {
      title = user.nickname;
      if (cfg.vue instanceof Object && typeof cfg.vue.$forceUpdate === "function") cfg.vue.$forceUpdate();

      console.log(2222222222, title);

      return user;
    } else {
      title = whoAmI;
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

    console.log(111111111, user);

    return setUser(user);
  }

  function saveUser(user) {
    localStorage.setItem('user', JSON.stringify(user));

    return setUser(user);
  }

  let exported = {
    title() {
      return title;
    },

    data: () => {
      return {
        title,
        inputLogin: "",
        inputPassword: "",
        user,
      };
    },
    methods: {
      getUserFromAuth: function(login, password, cb) {
        cb({user: {nickname: "pavlo"}});

        // fetch('http://localhost:3333/confidence/auth/auth', {
        //   method: 'POST', // *GET, POST, PUT, DELETE, etc.
        //   headers: {
        //     'Content-Type': 'application/json',
        //   },
        //   body: JSON.stringify({values: {login, password}}),
        //   mode: 'cors', // no-cors, cors, *same-origin
        //
        //   // cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        //   // credentials: 'same-origin', // include, *same-origin, omit
        //   // redirect: 'follow', // manual, *follow, error
        //   // referrer: 'no-referrer', // no-referrer, *client
        // }).then(response => {
        //   return response.json();
        // }).then(data =>
        //   cb(data)
        // );
      },

      signIn: function () {
        this.getUserFromAuth(this.inputLogin, this.inputPassword, data => {
          this.user = saveUser(data.user);
        });
      },

      signOut: function () {
        this.user = saveUser(unauthorized);
      },
    },
  };

  export default exported;

</script>

<style lang="scss">
</style>



<!--<style lang="scss">-->
<!--  #auth {-->
<!--    background-color: #ffffff;-->
<!--    padding: 5px;-->
<!--  }-->
<!--</style>-->
