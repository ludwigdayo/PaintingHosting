<!-- 单个画作的详情 -->

<template>
    <div class="container" v-if="imageData.image">
        <!-- 图片信息 -->
        <div class="container-inner-info">
            <div class="image-wrapper">
                <img :src="getImageSrc(imageData.image)" alt="Image" class="rounded-image" />
                <div id="image-story">{{ imageData.metadata.story }}</div>
            </div>
            <div class="image-info">
                <table>
                    <tr class="clickable">
                        <td>
                            名称
                        </td>
                        <td>
                            {{ imageData.metadata.imageName }}
                        </td>
                    </tr>
                    <tr class="clickable" @click="goToAuthor(imageData.metadata.author)">
                        <td>
                            作者
                        </td>
                        <td class="clickableTag">
                            {{ imageData.metadata.author }}
                        </td>
                    </tr>
                    <tr class="clickable">
                        <td>
                            时间
                        </td>
                        <td>
                            {{ imageData.metadata.createTime }}
                        </td>
                    </tr>
                    <tr class="clickable">
                        <td>
                            价格
                        </td>
                        <td>
                            {{ imageData.metadata.price }}
                        </td>
                    </tr>
                    <tr class="clickable">
                        <td>
                            标签
                        </td>
                        <td>
                            <span v-for="tag in imageData.metadata.tags.split(',')" class="clickableTag" :key="tag"
                                @click="handleTagClick(tag)">
                                {{ tag }}
                            </span>
                        </td>
                    </tr>
                </table>
            </div>
        </div>

        <!-- 互动信息 -->
        <div class="container-inner-interact">
            <div class="interact-wrapper">
                <div class="interact-item" @click="addLike">
                    <i class="far fa-heart" style="cursor: pointer;"></i>
                    <span class="interact-item-text">{{ imageInteract.Like }}</span>
                </div>

                <div class="interact-item" @click="addFavorite">
                    <i class="far fa-thumbs-up" style="cursor: pointer;"></i>
                    <span class="interact-item-text">{{ imageInteract.Favorite }}</span>
                </div>
            </div>

            <div class="interact-item comment-container">
                <p class="comment-content" v-for="(comment, index) in splitComments" :key="index">{{ comment }}</p>

                <input type="text" class="comment-container-input" v-model="commentInput" placeholder="发表你的感想"
                    @keyup.enter="dealInteractInfo">
            </div>

            <!-- 不好看 -->
            <!-- <div v-if="showSuccess" class="success-message">
                Interact uploaded successfully!
            </div> -->
            <!-- 只显示错误信息就够了 -->
            <div v-if="showError" class="error-message">
                Interact upload failed!
            </div>
        </div>
    </div>
</template>

<script>
import api from '@/components/api.js';

export default ({
    name: 'Details',
    data() {
        return {
            // 初始值给v-if判断，让图片还没拉下来时不要渲染
            imageData: { image: null },
            // 图片互动信息
            imageInteract: {
                Like: 0,
                Favorite: 0,
                Comment: '',
            },
            index: 0,
            showSuccess: false,
            showError: false,
            commentInput: '',
            interval: null,
        };
    },
    beforeRouteLeave(_, __, next) {
        // 在路由跳转前把定时器关了，因为定时器是全局的，vue-router跳转不会影响它执行
        clearInterval(this.interval);

        // 允许跳转
        next();
    },
    mounted() {
        this.index = this.$route.params.id;

        this.fetchImages();
        this.getInteract();

        // 定时拉取交互信息
        this.interval = setInterval(this.getInteract, 2000);

        // 参数不对他娘的连页面都不会进，所以不需要
        // console.log('参数不对');
        // this.$router.push('/home');
    },
    computed: {
        // this.imageInteract.Comment是通过','分割多个评论组成的，也即评论不能加','不然成两段
        splitComments() {
            return this.imageInteract.Comment.split(",");
        }
    },
    methods: {
        addLike() {
            console.log("Like+1");
            this.sendInteractInfo(false, this.imageInteract.Like + 1, this.imageInteract.Favorite, this.imageInteract.Comment);
        },
        addFavorite() {
            console.log("Favorite+1");
            this.sendInteractInfo(false, this.imageInteract.Like, this.imageInteract.Favorite + 1, this.imageInteract.Comment)
        },
        // 点击某标签的处理
        handleTagClick(tag) {
            // 跳到大厅按标签展示
            this.$store.commit('setSelectedTag', tag);
            this.$router.push('/home');
        },
        getImageSrc(base64String) {
            return `data:image/png;base64, ${base64String}`;
        },
        // 拉取图片
        fetchImages() {
            console.log('id:' + this.index);

            // 只查一条，并且不使用缩略图
            api.get(`/images?start=${this.index}&end=${this.index}&isThumbnail=false`)
                .then(response => {
                    this.imageData = response.data[0];
                    console.log('查到:' + this.imageData);
                })
        },
        // 得到交互信息
        getInteract() {
            api.get(`/interact?imageId=${this.index}`)
                .then(response => {
                    this.imageInteract = response.data;
                })
        },
        // 点击作者的处理
        goToAuthor(author) {
            this.$router.push(`/author/${author}`);
        },
        // 发送交互数据到后端
        sendInteractInfo(showSuccess, Like, Favorite, Comment) {
            const formData = new FormData();
            formData.append('ImageId', this.index);
            formData.append('Like', Like);
            formData.append('Favorite', Favorite);
            formData.append('Comment', Comment);

            api.post(`/interactUpload`, formData)
                .then(_ => {
                    // 本地数据通过远程更新下来，防止本地数据与远程不同步
                    this.getInteract();

                    if (showSuccess) {
                        this.showSuccess = true;
                        setTimeout(() => {
                            this.showSuccess = false;
                        }, 2000);
                    }

                    this.commentInput = '';
                })
                .catch(error => {
                    this.showError = true;
                    setTimeout(() => {
                        this.showError = false;
                    }, 2000);

                    console.error(error);
                });
        },
        // 处理评论的发送
        dealInteractInfo() {
            if (this.commentInput !== '') {
                // 把评论拼接起来
                let comment = this.imageInteract.Comment + ',' + this.commentInput;
                this.sendInteractInfo(true, this.imageInteract.Like, this.imageInteract.Favorite, comment);
            }
        }
    }
})
</script>

<style scoped>
.container-inner-info {
    display: flex;
    /* 横着摆放 */
    flex-direction: row;
    margin-top: 20px;
    margin-right: 10px;
    margin-left: 0;
    margin-bottom: 20px;
    font-family: "LiSu";
}

.container-inner-interact {
    position: relative;
    width: 100%;
    height: 100%;
    margin-bottom: 200px;
}

.image-wrapper {
    position: relative;
    /* border-radius: 2px; */
    overflow: hidden;
    max-width: 50%;
    height: 100%;
    box-shadow: -5px -5px 5px rgba(0, 0, 0, 0.2);
}

.rounded-image {
    width: 100%;
    /* height: 90%; */
    /* object-fit: cover; */
}

#image-story {
    width: 100%;
    text-align: left;
    font-size: 30px;
}

.image-info {
    width: 100%;
}

.clickable {
    cursor: pointer;
}

.clickable:hover {
    background-color: #9de2ff;
    border-radius: 50px;
}

table {
    font-size: 30px;
    text-align: center;
    width: 100%;
    /* border: 1px solid rgba(0, 0, 0, 0.5); */
    border-collapse: collapse;
    margin: 10px;
    font-weight: 700;
}

table tr {
    width: 100%;
    /* border: 1px solid rgba(0, 0, 0, 0.5); */
}

table td {
    padding: 5px;
    width: 50%;
    /* border: 1px solid rgba(0, 0, 0, 0.5); */
}

.clickableTag {
    display: inline-block;
    margin-right: 20px;
}

.clickableTag:hover {
    color: rgb(255, 255, 255, 0.7);
}

.interact-wrapper {
    display: flex;
    margin: 5px;
    margin-top: 0;
}

.interact-item-text {
    font-size: 20px;
    margin: 5px;
    margin-right: 20px;
}

.comment-container {
    position: relative;
    /* margin-bottom: 50px; */
    /* margin: 10px; */
}

.comment-content {
    border-color: rgba(0, 0, 0, 0);
    width: 100%;
    font-size: 20px;
}

.comment-container-input {
    width: 100%;
    height: 40px;
    border-radius: 50px;
    border-color: rgba(0, 0, 0, 0.2);
    padding: 2px;
    padding-left: 20px;
    font-size: 20px;
}

.success-message {
    margin-top: 20px;
    padding: 10px;
    background-color: #00b3ff;
    color: white;
    border-radius: 50px;
    text-align: center;
    width: 100%;
}

.error-message {
    margin-top: 20px;
    padding: 10px;
    background-color: #4e4e4e;
    color: white;
    border-radius: 50px;
    text-align: center;
    width: 100%;
}
</style>