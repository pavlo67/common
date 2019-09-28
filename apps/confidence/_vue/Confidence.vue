<template>
  <div id="confidence">
    {{ user.nickname }}
    <br>
    <span v-if="user.id">
      <button v-on:click="signOut">Вийти</button>
    </span>
    <span v-else>
      <input v-model="inputLogin"><input v-model="inputPassword">
      <button v-on:click="signIn">Авторизуватись</button>
    </span>

  </div>
</template>

<script>
    const unauthorizedUser = {nickname: "<unauthorized>"};
    const authorizedUser = {id: 1, nickname: "pavlo"};

    // getUserFromAuth -------------------------------------------------------
    function getUserFromAuth(login, password, cb) {
        fetch('http://localhost:3333/confidence/v1/auth/auth', { 
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
  }
</style>
