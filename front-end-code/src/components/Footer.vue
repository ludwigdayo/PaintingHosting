<template>
    <div class="container">
        <footer class="footer">
            <p>一个图片托管展示服务器</p>
        </footer>
    </div>
</template>

<script>
export default {
    name: 'Footer',
};
</script>

<style scoped>
.container {
    display: flex;
    flex-direction: column;
    min-height: 100%;
}

.content {
    flex: 1;
}

.footer {
    background-color: #ffffff;
    padding: 6px;
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 -10px 10px rgba(0, 0, 0, 0.1);
}

/* Sticky Footer styles */
.content {
    margin-bottom: 25px;
    /* Footer 的高度 */
}

.footer {
    font-size: 12px;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 18px;
}
</style>