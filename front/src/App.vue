<template>
  <div id="app">
    <Xterm />
  </div>
</template>

<script>
import Xterm from './components/Xterm.vue'
export default {
  name: 'App',
  components: {
    Xterm
  }
}
</script>

<style scoped>
.active {
  color: #000;
  background: rgb(211, 24, 24);
}
</style>