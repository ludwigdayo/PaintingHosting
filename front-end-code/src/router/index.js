import Vue from 'vue'
import Router from 'vue-router'

// 画作展示页面 
import Home from '@/components/Home'

// 画作上传页面
import Upload from '@/components/Upload'

// 单个画作详情页面
import Details from '@/components/Details'

// 作者页面
import Author from '@/components/Author'

// 告诉 Vue.js 在应用中使用 Vue Router 插件，并进行必要的初始化和配置
Vue.use(Router)

// export default 导出对象 一个vue文件只能有一句
export default new Router({
  routes: [
    {
      path: '/',
      redirect: '/home' // 默认跳转到Home
    },
    {
      path: '/home',
      name: 'Home',
      component: Home
    },
    {
      path: '/upload',
      name: 'Upload',
      component: Upload
    },
    {
      path: '/details/:id',
      name: 'Details',
      component: Details
    },
    {
      path: '/author/:author',
      name: 'Author',
      component: Author
    }
  ]
})
