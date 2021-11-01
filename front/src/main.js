import Vue from 'vue'
import App from './App.vue'
import router from './router'

Vue.config.productionTip = false
import VueRouter from 'vue-router'
Vue.use(VueRouter)

new Vue({
  render: h => h(App),
  router
}).$mount('#app')
