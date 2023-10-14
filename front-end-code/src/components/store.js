import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

const store = new Vuex.Store({
    state: {
        // 已有的全部标签
        tags: [],

        // 在“大厅”或者“详情”页选中的标签
        selectedTag: '',

        // 其实这里没写完，应该加个翻页的功能，而不是把查询的图片范围固定成0-100
        start: 0,
        end: 100,
    },
    mutations: {
        setTags(state, tags) {
            state.tags = tags;
        },
        setSelectedTag(state, tag) {
            state.selectedTag = tag;
        },
    },
});

export default store;