<!-- 大厅：展示所有画作 -->

<template>
    <div class="container">
        <div class="image-row" v-for="(row, index) in imageRows" :key="index">
            <div class="image-container" v-for="image in row" :key="image.metadata.id">
                <div class="image-container-inner">
                    <div class="image-wrapper" @mouseover="showImageMessage($event, image)"
                        @mouseleave="hideImageMessage($event, image)">
                        <img :src="getImageSrc(image.image)" alt="Image" class="rounded-image"
                            @click="goToDetails(image.metadata.id)" />
                    </div>
                    <div class="image-info">
                        <p class="image-name">{{ image.metadata.imageName }}</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- 每个图片的信息 -->
        <div class="image-message" v-if="imageNow.isHovered" :style="{ top: imageNow.topValue, left: imageNow.leftValue }">
            作者：{{ imageNow.image.metadata.author }}，时间： {{ imageNow.image.metadata.createTime }}
        </div>
    </div>
</template>

<script>
import api from '@/components/api.js';
import { mapState } from 'vuex';

export default {
    name: 'Home',
    data() {
        return {
            // 这是从后端拉取的全部图片对象
            imageDataList: [],

            // 记录有没有图片被盯上
            imageNow: {
                isHovered: false,
                image: null,
            },
        };
    },
    mounted() {
        // 组件加载时获取图片数据
        this.fetchImages(this.start, this.end);
    },
    computed: {
        // 已选择的标签，显示图片的范围
        ...mapState(['selectedTag', 'start', 'end']),

        // 切分图片数组
        imageRows() {
            // 将图片数据列表按行分组，每行最多显示4个图片
            // "rows"数组中存储了"filteredImage"数组按行切分后的所有行数据
            const rows = [];
            const rowSize = 4;
            for (let i = 0; i < this.filteredImage.length; i += rowSize) {
                rows.push(this.filteredImage.slice(i, i + rowSize));
            }
            return rows;
        },

        // 进行按标签的过滤
        filteredImage() {
            if (this.selectedTag !== '') {
                console.log("按标签过滤");
                console.log("此时的过滤：" + this.selectedTag);
                return this.imageDataList.filter(imageData => {
                    const tags = imageData.metadata.tags.split(','); // 将字符串转换为数组
                    return tags.includes(this.selectedTag); // 检查数组中是否包含选定的标签
                });
            }
            return this.imageDataList;
        },
    },
    methods: {
        // 把图片的标签遍历出来
        getTagsFromImageDataList() {
            const uniqueTags = [...new Set(
                this.imageDataList.flatMap(imageData => imageData.metadata.tags.split(','))
            )];

            // 塞到store里给Header.vue读取
            this.$store.commit('setTags', uniqueTags);
            console.log(uniqueTags);
        },

        // 给base64的图片加个头，为了给vue识别
        getImageSrc(base64String) {
            // console.log(base64String);
            return `data:image/png;base64, ${base64String}`;
        },

        // 从后端获取图片对象列表
        fetchImages(start, end) {
            api.get(`/images?start=${start}&end=${end}&isThumbnail=true`)
                .then(response => {
                    this.imageDataList = response.data;
                    console.log(this.imageDataList);

                    this.getTagsFromImageDataList();
                })
        },

        // 悬停时显示部分图片信息
        showImageMessage(event, image) {
            // console.log(image.metadata.id + '显示图片信息');
            this.imageNow.image = image;
            this.imageNow.isHovered = true;

            // 把鼠标位置记录下来 为了让那个几把元素显示到鼠标附近
            const posX = event.clientX + 'px';
            const posY = event.clientY + 'px';
            this.imageNow.topValue = posY;
            this.imageNow.leftValue = posX;
        },
        hideImageMessage(event, image) {
            // console.log(image.metadata.id + '隐藏图片信息');
            this.imageNow.image = null;
            this.imageNow.isHovered = false;
        },

        // 跳转到详情页面
        goToDetails(id) {
            this.$router.push(`/details/${id}`);
        }
    }
};
</script>

<style scoped>
.image-row {
    display: flex;
}

.image-container {
    flex: 0 0 25%;
    padding: 10px;
    position: relative;
    margin-top: 15px;
    margin-bottom: 15px;
}

.image-container-inner {
    width: 300px;
    height: 200px;
}

.image-wrapper {
    position: relative;
    border-radius: 8px;
    overflow: hidden;
    width: 300px;
    height: 180px;
    box-shadow: -5px -5px 5px rgba(0, 0, 0, 0.2);
}

.rounded-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.image-info {
    position: relative;
    padding: 0px;
    color: #000000;
    text-align: center;
    font-size: 20px;
}

.image-message {
    position: absolute;
    background-color: #272727;
    color: #fff;
    padding: 5px;
    font-size: 18px;
    border-radius: 5px;
    z-index: 999;
}
</style>