// 入口js文件
// 其他文件导出对象,到这里统一导入使用,整个页面本质上都是vuejs构建的

import Vue from 'vue'
import App from './App'
import router from './router'
import store from '@/components/store';
import '@fortawesome/fontawesome-free/css/all.css'

// 禁用生产环境下的提示信息
Vue.config.productionTip = false

new Vue({
  el: '#app',

  // 页面的路由
  router: router,

  store: store,

  // App包含一切
  components: { App },

  // 整个页面都是<App/>中的内容
  template: '<App/>'
})
