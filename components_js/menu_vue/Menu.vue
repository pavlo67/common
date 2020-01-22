<template>
  <div>
    <router-link v-for="item in menuItems" v-bind:key="item.path" :to="item.path">
      {{ item.title() }}<br>
    </router-link>
  </div>
</template>

<script>
  export default {
    name: 'Menu',
    props: ['menuItems'],
    // methods: {
    //   forceUpdate() {
    //     this.$forceUpdate();  // Notice we have to use a $ here
    //   }
    // },
  };

</script>

<style lang="scss">
</style>
