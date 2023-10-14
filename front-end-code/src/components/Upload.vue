<template>
    <div>
        <label for="imageName">图像名称:</label>
        <input type="text" id="imageName" v-model="imageName">

        <label for="author">作者:</label>
        <input type="text" id="author" v-model="author">

        <label for="createTime">创作时间:</label>
        <input type="text" id="createTime" v-model="createTime">

        <label for="story">故事:</label>
        <textarea id="story" v-model="story"></textarea>

        <label for="price">价格:</label>
        <input type="number" id="price" v-model="price">

        <label for="image">Image:</label>
        <input type="file" id="image" ref="imageInput" @change="handleImageInputChange">

        <label for="tags">标签(逗号分隔):</label>
        <input type="text" id="tags" v-model="tags">

        <button @click="submitImage">Upload</button>

        <div v-if="showSuccess" class="success-message">
            Image uploaded successfully!
        </div>
    </div>
</template>

<script>
import api from '@/components/api.js';

export default {
    name: 'Upload',
    data() {
        return {
            imageName: '',
            author: '',
            createTime: '',
            story: '',
            price: 0,
            image: null,
            tags: '',
            showSuccess: false,
        };
    },
    methods: {
        handleImageInputChange(event) {
            this.image = event.target.files[0];
        },
        submitImage() {
            const formData = new FormData();
            formData.append('imageName', this.imageName);
            formData.append('author', this.author);
            formData.append('createTime', this.createTime);
            formData.append('story', this.story);
            formData.append('price', this.price);
            formData.append('image', this.image);
            formData.append('tags', this.tags);

            // 调试
            for (let pair of formData.entries()) {
                console.log(pair[0] + ', ' + pair[1]);
            }

            // 发送数据到后端
            api.post('/upload', formData)
                .then(response => {
                    this.showSuccess = true;
                    setTimeout(() => {
                        this.showSuccess = false;
                    }, 2000);
                    this.resetForm();
                })
        },
        resetForm() {
            this.imageName = '';
            this.author = '';
            this.createTime = '';
            this.story = '';
            this.price = '';
            this.image = null;
            this.tags = '';
        }
    },
};
</script>

<style scoped>
label {
    display: block;
    font-weight: bold;
    margin-bottom: 20px;
    text-align: left;
    margin: 5px;
    margin-top: 10px;
    margin-left: 0;
    font-size: 20px;
}

input[type="text"],
textarea {
    width: 100%;
    padding: 10px;
    margin: 2px;
    border: 1px solid #ccc;
    border-radius: 10px;
    border-color: rgba(0, 0, 0, 0.2);
    margin: 1px;
}

button {
    padding: 7px 30px;
    background-color: #00b3ff;
    color: white;
    border: none;
    border-radius: 7px;
    margin: 5px;
    margin-left: 0;
    cursor: pointer;
}

.success-message {
    margin-top: 20px;
    padding: 10px;
    background-color: #00b3ff;
    color: white;
    border-radius: 4px;
    text-align: center;
}
</style>