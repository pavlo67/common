<template>
  <div id="confidence">
    {{ user }}
    <br>
    <button v-on:click="signIn">Авторизуватись</button>
  </div>
</template>

<script>
    let n;
    let user;

    let userOld = localStorage.getItem('user');
    if (userOld) {
      try {
        user = JSON.parse(userOld);
      } catch (err) {
        console.error(err);
      }
    }

    if (user instanceof Object) {
        n = user.id
    } else {
        n = 1;
        user = {id:n, nickname:"pavlo"};
    }

    localStorage.setItem('user', JSON.stringify(user));

    export default {
        name: 'Confidence',
        data: () => {
            return {
                user: localStorage.getItem('user'),
            };
        },

        methods: {
            signIn: function () {
                this.user = JSON.stringify({id:++n, nickname:"pavlo"});
                localStorage.setItem('user', this.user);
            }
        },
    };
</script>

<style lang="scss">
  #confidence {
    background-color: #ffffff;
  }
</style>
