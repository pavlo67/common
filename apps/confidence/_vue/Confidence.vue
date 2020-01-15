<template>
  <div id="confidence">
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
    const unauthorizedUser = undefined;

    // getUserFromAuth -------------------------------------------------------
    function getUserFromAuth(login, password, cb) {
        fetch('http://localhost:3333/confidence/auth/auth', {
            method: 'POST', // *GET, POST, PUT, DELETE, etc.
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({values: {login, password}}),
            mode: 'cors', // no-cors, cors, *same-origin

            // cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
            // credentials: 'same-origin', // include, *same-origin, omit
            // redirect: 'follow', // manual, *follow, error
            // referrer: 'no-referrer', // no-referrer, *client
        })
        .then(response => {
            return response.json();
        }).then(data => cb(data));
    }

    // getUser ---------------------------------------------------------------
    function getUser() {
        let user;

        let userStored = localStorage.getItem('user');
        if (userStored) {
            try {
                user = JSON.parse(userStored);
            } catch (err) {
                console.error(err);
                return null;
            }
        }

        return user instanceof Object ? user : unauthorizedUser;
    }

    // setUser ---------------------------------------------------------------
    function setUser(user) {
        localStorage.setItem('user', JSON.stringify(user));
    }

    // -----------------------------------------------------------------------
    export default {
        name: 'Confidence',
        data: () => {
            return {
                user: getUser(),
            };
        },

        methods: {
            signIn: function () {
                getUserFromAuth(this.inputLogin, this.inputPassword, data => {
                    if (data instanceof Object && data.user instanceof Object) {
                        this.user = data.user;
                        setUser(this.user);
                    }
                });
            },
            signOut: function () {
                this.user = unauthorizedUser;
                setUser(this.user);
            },
        },
    };
</script>

<style lang="scss">
  #confidence {
    background-color: #ffffff;
    padding: 5px;
  }
</style>
