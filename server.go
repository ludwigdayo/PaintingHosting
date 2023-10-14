package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
)

// 路径
var localDirPath = "./ImagesData"
var dbPath = "./db/imageshow-server-database.db"
var frontPath = "./dist"

// 返回作者信息
func authorinfoHandle(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头部
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许来自任意域的请求
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	authorName := r.URL.Query().Get("authorName")

	author, err := getAuthorInfo(authorName)
	if err != nil {
		log.Println("Invalid authorName parameter", err)
		http.Error(w, "Failed to query author", http.StatusInternalServerError)
		return
	}

	log.Println(author.AuthorName + author.Introduce + author.Image + "=> 发的几把玩意")

	// 将查询结果转换为 JSON
	jsonData, err := json.Marshal(author)
	if err != nil {
		log.Println("Failed to marshal JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 发送 JSON 响应
	w.Write(jsonData)

	log.Println("Success to send " + authorName + "'s info")
}

// 接收作者信息
func authorinfoUploadHandle(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头部
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许来自任意域的请求
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	err := r.ParseMultipartForm(32 << 20) // 设置最大内存限制为32MB
	if err != nil {
		log.Println("Failed to parse multipart form:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 获取表单字段的值
	authorName := r.FormValue("authorName")
	introduce := r.FormValue("introduce")
	imageStr := r.FormValue("image")

	// 作者头像直接存数据库
	author := Author{
		AuthorName: authorName,
		Introduce:  introduce,
		Image:      imageStr,
	}

	log.Println(author.Image + "=> upload author Image")

	err = saveAuthorInfo(author)
	if nil != err {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Failed to save author info:", err)
	}

	// 返回成功的响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
	log.Println("Author info uploaded successfully!")
}

// 返回交互信息
func interactInfoHandle(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头部
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许来自任意域的请求
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	imageIdStr := r.URL.Query().Get("imageId")

	// 将查询参数转换为整数
	imageId, err := strconv.Atoi(imageIdStr)
	if err != nil {
		log.Println("Invalid imageId parameter", err)
		http.Error(w, "Invalid imageId parameter", http.StatusBadRequest)
		return
	}

	interact, err := getInteractInfo(imageId)
	if err != nil {
		log.Println("Invalid imageId parameter", err)
		http.Error(w, "Failed to query interact", http.StatusInternalServerError)
		return
	}

	// 将查询结果转换为 JSON
	jsonData, err := json.Marshal(interact)
	if err != nil {
		log.Println("Failed to marshal JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 发送 JSON 响应
	w.Write(jsonData)

	log.Println("Success to send image " + imageIdStr + " info")
}

// 接收交互信息
func interactInfoUploadHandle(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头部
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许来自任意域的请求
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	err := r.ParseMultipartForm(32 << 20) // 设置最大内存限制为32MB
	if err != nil {
		log.Println("Failed to parse multipart form:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 获取表单字段的值
	imageIdStr := r.FormValue("ImageId")
	likeStr := r.FormValue("Like")
	favoriteStr := r.FormValue("Favorite")
	commentStr := r.FormValue("Comment")

	// 将查询参数转换为整数
	imageId, err := strconv.Atoi(imageIdStr)
	if err != nil {
		log.Println("Invalid imageId parameter", err)
		http.Error(w, "Invalid imageId parameter", http.StatusBadRequest)
		return
	}

	like, err := strconv.Atoi(likeStr)
	if err != nil {
		log.Println("Invalid like parameter", err)
		http.Error(w, "Invalid like parameter", http.StatusBadRequest)
		return
	}

	favorite, err := strconv.Atoi(favoriteStr)
	if err != nil {
		log.Println("Invalid favorite parameter", err)
		http.Error(w, "Invalid favorite parameter", http.StatusBadRequest)
		return
	}

	interact := Interact{
		ImageId:  imageId,
		Like:     like,
		Favorite: favorite,
		Comment:  commentStr,
	}

	err = saveInteractInfo(interact)
	if nil != err {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Failed to save interact info:", err)
	}

	// 返回成功的响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
	log.Println("Interact info uploaded successfully!")
}

// 生成新的文件名
func generateFileName(originalFileName string) string {
	// 使用 UUID 生成唯一的文件名
	newUUID := uuid.New()
	fileExtension := filepath.Ext(originalFileName)
	newFileName := newUUID.String() + fileExtension
	return newFileName
}

// 保存图片文件
func saveImageToFile(imageBytes []byte, fileName string) error {
	// 创建新的文件
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将图片文件内容写入新文件
	_, err = file.Write(imageBytes)
	if err != nil {
		return err
	}

	return nil
}

// 处理上传请求
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头部
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许来自任意域的请求
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	err := r.ParseMultipartForm(32 << 20) // 设置最大内存限制为32MB
	if err != nil {
		log.Println("Failed to parse multipart form:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 获取表单字段的值
	imageName := r.FormValue("imageName")
	author := r.FormValue("author")
	createTime := r.FormValue("createTime")
	story := r.FormValue("story")
	priceStr := r.FormValue("price")
	tags := r.FormValue("tags")

	// 获取上传的图片文件
	file, handler, err := r.FormFile("image")
	if err != nil {
		// 处理错误
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Failed to get image file: %v", err)
		return
	}
	defer file.Close()

	// 读取图片文件内容
	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		// 处理错误
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to read image file: %v", err)
		return
	}

	// 生成新的文件名，避免冲突
	newFileName := generateFileName(handler.Filename)
	filePath := localDirPath + newFileName

	// 保存图片文件
	err = saveImageToFile(imageBytes, filePath)
	if err != nil {
		// 处理错误
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to save image file: %v", err)
		return
	}

	priceNum, err := strconv.Atoi(priceStr)
	if err != nil {
		// 处理错误
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid price: %v", err)
		return
	}

	// 构建 Image 对象
	image := Image{
		ID:         0,
		ImageName:  imageName,
		Author:     author,
		CreateTime: createTime,
		Story:      story,
		Price:      priceNum,
		Path:       filePath,
		Tags:       tags,
	}

	// 存储图片元数据到数据库中
	err = SaveImageToDatabase(image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to save image to database: %v", err)
		return
	}

	// 返回成功的响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
	log.Println("Image uploaded successfully!")
}

// 处理拉取请求
func imagesHandler(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头部
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许来自任意域的请求
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 解析查询参数
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	isThumbnailStr := r.URL.Query().Get("isThumbnail")

	// 将查询参数转换为整数
	start, err := strconv.Atoi(startStr)
	if err != nil {
		log.Println("Invalid start parameter", err)
		http.Error(w, "Invalid start parameter", http.StatusBadRequest)
		return
	}

	end, err := strconv.Atoi(endStr)
	if err != nil {
		log.Println("Invalid end parameter", err)
		http.Error(w, "Invalid end parameter", http.StatusBadRequest)
		return
	}

	isThumbnail, err := strconv.ParseBool(isThumbnailStr)
	if err != nil {
		log.Println("Invalid isThumbnailStr parameter", err)
		http.Error(w, "Invalid isThumbnailStr parameter", http.StatusBadRequest)
		return
	}

	// 调用查询函数
	imageList, err := QueryImages(start, end, isThumbnail)
	if err != nil {
		log.Println("Failed to query images:", err)
		http.Error(w, "Failed to query images", http.StatusInternalServerError)
		return
	}

	// 将查询结果转换为 JSON
	jsonData, err := json.Marshal(imageList)
	if err != nil {
		log.Println("Failed to marshal JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 发送 JSON 响应
	w.Write(jsonData)

	log.Println("Success to send image list")
}

func main() {
	// 资源判断
	paths := []string{localDirPath, dbPath, frontPath}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("Error: Path %s does not exist\n", path)
			return
		} else {
			log.Printf("Path %s exists\n", path)
		}
	}

	// 设置自定义的用法信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -p port -h address\n", os.Args[0])
		flag.PrintDefaults()
	}

	// 请注意 这里的端口同时用于数据传输与提供前端页面 请务必与前端对应
	port := 80
	ipAddress := "0.0.0.0"
	flag.IntVar(&port, "p", 80, "Specify the port number. The default is 80. It must correspond to the front end.")
	flag.StringVar(&ipAddress, "h", "0.0.0.0", "Specify host address. Default is IP (0.0.0.0)")
	flag.Parse()
	address := ipAddress + ":" + strconv.Itoa(port)

	// 配置静态文件目录
	staticDir, _ := filepath.Abs(frontPath)

	// 作者详情
	http.HandleFunc("/author", authorinfoHandle)

	// 上传作者信息
	http.HandleFunc("/authorUpload", authorinfoUploadHandle)

	// 上传交互信息
	http.HandleFunc("/interactUpload", interactInfoUploadHandle)

	// 作品交互
	http.HandleFunc("/interact", interactInfoHandle)

	// 返回图像
	http.HandleFunc("/images", imagesHandler)

	// 处理上传
	http.HandleFunc("/upload", uploadHandler)

	// 设置静态文件处理器
	http.Handle("/", http.FileServer(http.Dir(staticDir)))

	log.Println("Server run on:", address)
	log.Println("If access fails please check permissions")

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
}
