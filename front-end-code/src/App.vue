 <!-- 根组件 -->
 
<template>
  <div id="app" style="overflow-x:hidden">
    <vue-header></vue-header>
    <!-- 路由的内容 路由差不多是动态面板 -->
    <div class='content'>
      <router-view />
    </div>
    <!-- <vue-footer></vue-footer> -->
  </div>
</template>

<script>
import Header from "@/components/Header";
import Footer from "@/components/Footer";

export default {
  // 给组件命名，方便在开发工具和调试中进行识别和查找 可以不用
  name: 'App',

  // 注册和引入组件
  components: {
    // 他娘的用header,footer Vue不认
    'vue-header': Header,
    'vue-footer': Footer
  }
}
</script>

<style scoped>
html,
body {
  overflow-x: hidden;
}

.container {
  display: flex;
  flex-direction: column;
}

.header {
  margin-bottom: 20px;
}

.content {
  flex: 0;
  margin-top: 30px;
  margin-bottom: 30px;
  width: 100%;
  padding: 0 20px;
  box-sizing: border-box;

  display: flex;
  flex-direction: column;
  /* align-items: center; */
}

.footer {
  margin-top: 20px;
}
</style>
