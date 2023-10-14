<!-- 作者详情 -->

<template>
    <div class="container" @dragover.prevent @drop="handleDrop">
        <div class="image-introduce">
            <img v-if="author.Image" :src="displayImage()" alt="Image" class="rounded-image" />
            <h2 class="author-name">{{ author.AuthorName }}</h2>
        </div>
        <textarea class="introduce" v-model="author.Introduce" @blur="uploadAuthorInfo"></textarea>

        <div v-if="showSuccess" class="success-message">
            Image uploaded successfully!
        </div>
    </div>
</template>

<script>
import api from '@/components/api.js';

export default ({
    name: 'Author',
    data() {
        return {
            author: {
                AuthorName: '',
                Introduce: '',
                Image: '',
            },
            showSuccess: false,
        };
    },
    mounted() {
        console.log('进入作者详情页面');
        this.author.AuthorName = this.$route.params.author;

        this.fetchAuthor();
    },
    methods: {
        displayImage() {
            return `data:image/jpeg;base64, ${this.author.Image}`;
        },
        uploadAuthorInfo() {
            const formData = new FormData();
            formData.append('authorName', this.author.AuthorName);
            formData.append('introduce', this.author.Introduce);

            // console.log("Send Image Str:" + this.author.Image)
            console.log("Send Image introduce:" + this.author.Introduce)

            formData.append('image', this.author.Image);

            api.post(`/authorUpload`, formData)
                .then(_ => {
                    this.showSuccess = true;
                    setTimeout(() => {
                        this.showSuccess = false;
                    }, 1500);
                    this.resetForm();
                })
        },
        // 处理图片拖入
        handleDrop(event) {
            event.preventDefault();

            const files = event.dataTransfer.files;

            if (files.length > 0) {
                const file = files[0];

                if (!file.type.startsWith('image/')) {
                    console.log("你他娘的放的什么几把");
                    return;
                }

                const reader = new FileReader();
                // 这几把是异步操作
                reader.onload = () => {
                    this.author.Image = reader.result.split(',')[1]; // 去除数据头部信息
                    // console.log(this.author.Image)
                    console.log("读取成功")
                    this.uploadAuthorInfo();
                };
                reader.readAsDataURL(file);
            }
        },
        fetchAuthor() {
            // 用作者名查询
            api.get(`/author?authorName=${this.author.AuthorName}`)
                .then(response => {
                    this.author = response.data;
                    console.log("查到：");
                    console.log(this.author);
                })
        }
    }
})
</script>

<style scoped>
.container {
    /* background-color: #f5f5f5; */
    padding: 50px;
    border-radius: 5px;
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.2);
    display: flex;
    flex-direction: column;
    /* justify-content: center; */
    align-items: center;
    /* height: 100vh; */
}

.author-name {
    font-size: 24px;
    font-weight: bold;
    text-align: center;
    margin: 20px;
}

.introduce {
    margin-bottom: 10px;
    resize: none;
    /* border: none; */
    height: 200px;
    width: 70%;
    border-radius: 20px;
    padding: 20px;
    font-size: 20px;
    font-weight: 900;
}

.image-introduce p {
    font-size: 16px;
    color: #333;
}

.rounded-image {
    width: 200px;
    height: 200px;
    border-radius: 50%;
}

.success-message {
    margin-top: 20px;
    padding: 10px;
    background-color: #00b3ff;
    color: white;
    border-radius: 10px;
    text-align: center;
}
</style>