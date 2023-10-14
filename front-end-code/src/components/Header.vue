<template>
    <div class="container">
        <header class="header-wrap">
            <nav class="nav">
                <!-- 点击大厅具有重置页面功能 -->
                <div @click="reset" class="nav-link">大厅</div>
                <router-link to="/upload" class="nav-link">上传</router-link>
            </nav>

            <div class="header-bar">
                <div v-if="$route.path === '/home'">
                    <!-- 不行，没法通知Home去拉取图片 -->
                    <!-- <button @click="setLastPage">上一页</button> -->
                    <!-- <button @click="setNextPage">下一页</button> -->

                    <div class="custom-select">
                        <select v-model="selectedTagStr" class="select-box">
                            <option value="">全部标签</option>
                            <option v-for="tag in tags" :key="tag" :value="tag" :selected="tag === selectedTag">
                                {{ tag }}
                            </option>
                        </select>
                        <span class="arrow"></span>
                    </div>
                </div>
            </div>
        </header>
    </div>
</template>

<script>
import { mapMutations, mapState } from 'vuex';

export default {
    name: 'Header',
    computed: {
        // 把store里的标签展示出来
        // tags是全部标签，selectedTag是已选择的标签
        ...mapState(['tags', 'selectedTag']),
        selectedTagStr: {
            get() {
                return this.selectedTag;
            },
            set(newTag) {
                this.$store.commit('setSelectedTag', newTag);
            }
        }
    },
    methods: {
        ...mapMutations(['setNextPage', 'setLastPage']),
        backtrack() {
            // 原本想写个返回
            // 太麻烦了不写了
        },
        reset() {
            this.$store.commit('setSelectedTag', '');
            if (this.$route.path !== '/home')
                this.$router.push('/home');
            console.log("重置页面");
        }
    }
}
</script>

<style scoped>
.container {
    background-color: #fff;
    padding: 10px;
}

.header-wrap {
    background-color: #fff;
    padding: 6px;
    transition: box-shadow 0.3s ease;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 999;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 10px 10px rgba(0, 0, 0, 0.1);
}

.header-wrap:hover {
    opacity: 1;
}

.nav {
    display: flex;
    justify-content: center;
    align-items: center;
}

.nav-link {
    text-decoration: none;
    color: #333;
    font-size: 18px;
    font-weight: bold;
    padding: 5px;
    margin: 0 10px;
    transition: color 0.3s ease;
    cursor: pointer;
}

.header-bar {
    display: flex;
}

.nav-link:hover {
    color: #00b3ff;
}

.custom-select {
    position: relative;
    display: inline-block;
}

.select-box {
    appearance: none;
    background-color: #fff;
    border: 1px solid #ccc;
    border-radius: 5px;
    padding: 6px 10px;
    font-size: 15px;
    color: #333;
    width: 300px;
}

.arrow {
    position: absolute;
    top: 50%;
    right: 10px;
    transform: translateY(-50%);
    width: 0;
    height: 0;
    border-left: 6px solid transparent;
    border-right: 6px solid transparent;
    border-top: 6px solid #666;
}
</style>
