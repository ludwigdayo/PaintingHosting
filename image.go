package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nfnt/resize"
)

// 图片元数据
type Image struct {
	ID         int    `json:"id"`
	ImageName  string `json:"imageName"`
	Author     string `json:"author"`
	CreateTime string `json:"createTime"`
	Story      string `json:"story"`
	Price      int    `json:"price"`
	Path       string `json:"image"`
	Tags       string `json:"tags"`
}

// 包含图片本体
type ImageData struct {
	Metadata Image  `json:"metadata"`
	ImageStr string `json:"image"`
}

// CREATE TABLE IMAGES(
// ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
// IMAGENAME TEXT NOT NULL,
// AUTHOR TEXT NOT NULL,
// CREATETIME TEXT NOT NULL,
// STORY TEXT NOT NULL,
// PRICE INT NOT NULL,
// PATH TEXT NOT NULL,
// TAGS TEXT NOT NULL
// );

// 查询图像
func QueryImages(start int, end int, isThumbnail bool) ([]ImageData, error) {
	query := fmt.Sprintf("SELECT ID, IMAGENAME, AUTHOR, CREATETIME, STORY, PRICE, PATH, TAGS FROM IMAGES LIMIT %d OFFSET %d", end-start, start)

	// 只查一条则使用ID查询，这个是不完善的解决方法
	if end-start == 0 {
		query = fmt.Sprintf("SELECT ID, IMAGENAME, AUTHOR, CREATETIME, STORY, PRICE, PATH, TAGS FROM IMAGES WHERE ID = %d", start)
	}

	db, err := openDB()
	if err != nil {
		log.Println("Open database failed", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Println("Failed to query")
		return nil, err
	}
	defer rows.Close()

	var imageList []ImageData
	for rows.Next() {
		var id int
		var imageName string
		var author string
		var createTime string
		var story string
		var price int
		var path string
		var tags string

		err := rows.Scan(&id, &imageName, &author, &createTime, &story, &price, &path, &tags)
		if err != nil {
			log.Println("Failed to scan row: ", err)
			continue
		}

		// 打开图片文件
		file, err := os.Open(path)
		if err != nil {
			log.Println("Failed to open image file: ", err)
			continue
		}
		defer file.Close()

		var imgBytes []byte

		if isThumbnail {
			// 解码原始图像文件
			img, _, err := image.Decode(file)
			if err != nil {
				log.Println("Failed to decode image: ", err)
				continue
			}

			// 调整图像大小为缩略图尺寸
			thumbnail := resize.Thumbnail(1080, 720, img, resize.Lanczos3)

			// 将缩略图转换为字节数组
			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, thumbnail, nil)
			if err != nil {
				log.Println("Failed to convert to jpeg: ", err)
				continue
			}

			imgBytes = buf.Bytes()
		} else {
			bs, err := ioutil.ReadAll(file)
			if err != nil {
				log.Println("Failed to read image file: ", err)
				continue
			}

			imgBytes = bs
		}

		// 将图像文件内容编码为Base64字符串
		base64Str := base64.StdEncoding.EncodeToString(imgBytes)

		// 创建 Image 结构体对象
		imageData := ImageData{
			Metadata: Image{
				ID:         id,
				ImageName:  imageName,
				Author:     author,
				CreateTime: createTime,
				Story:      story,
				Price:      price,
				Path:       "",
				Tags:       tags,
			},
			ImageStr: base64Str,
		}

		imageList = append(imageList, imageData)
	}

	return imageList, nil
}

// 元数据存到数据库
func SaveImageToDatabase(imageMetaData Image) error {
	// 执行插入语句将图片数据插入数据库
	sql := `INSERT INTO IMAGES (IMAGENAME, AUTHOR, CREATETIME, STORY, PRICE, PATH, TAGS) VALUES (?, ?, ?, ?, ?, ?, ?)`

	db, err := openDB()
	if err != nil {
		log.Println("Open database failed", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println("Failed to prepare SQL statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(imageMetaData.ImageName, imageMetaData.Author, imageMetaData.CreateTime, imageMetaData.Story, imageMetaData.Price, imageMetaData.Path, imageMetaData.Tags)
	if err != nil {
		log.Println("Failed to insert row into IMAGES table:", err)
		return err
	}

	return nil
}

type Author struct {
	AuthorName string
	Introduce  string
	Image      string
}

// CREATE TABLE AUTHORS(
// AUTHORNAME TEXT PRIMERY KEY UNIQUE NOT NULL,
// INTRODUCE TEXT NOT NULL,
// IMAGE TEXT NOT NULL
// );

func getAuthorInfo(authorName string) (Author, error) {
	db, err := openDB()
	if err != nil {
		log.Println("Open database failed", err)
		return Author{}, err
	}
	defer db.Close()

	query := "SELECT * FROM AUTHORS WHERE AUTHORNAME = ?"
	row := db.QueryRow(query, authorName)

	author := Author{}
	err = row.Scan(&author.AuthorName, &author.Introduce, &author.Image)
	if err != nil {
		return author, err
	}

	// log.Println("getAuthorImage:" + author.Image)

	return author, nil
}

func saveAuthorInfo(author Author) error {
	db, err := openDB()
	if err != nil {
		log.Println("Open database failed", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT OR REPLACE INTO AUTHORS(AUTHORNAME, INTRODUCE, IMAGE) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	log.Println(author.Image + "=> inserted author Image")

	// 执行插入语句
	_, err = stmt.Exec(author.AuthorName, author.Introduce, author.Image)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Author info inserted successfully!")

	return nil
}

type Interact struct {
	ImageId  int
	Like     int
	Favorite int
	Comment  string
}

// CREATE TABLE INTERACT(
// IMAGEID INT PRIMERY KEY UNIQUE NOT NULL,
// LIKE INT NOT NULL,
// FAVORITE INT NOT NULL,
// COMMENT TEXT NOT NULL
// );

// 保存交互信息到数据库
func saveInteractInfo(interact Interact) error {
	db, err := openDB()
	if err != nil {
		log.Println("Open database failed", err)
		return err
	}
	defer db.Close()

	// 就算没有这一行也能插入
	sql := "INSERT OR REPLACE INTO INTERACT(IMAGEID, LIKE, FAVORITE, COMMENT) VALUES(?,?,?,?)"

	// 执行更新操作
	_, err = db.Exec(sql, interact.ImageId, interact.Like, interact.Favorite, interact.Comment)
	if err != nil {
		log.Println("Failed to update row:", err)
		return err
	}

	// 更新成功
	log.Println("Row updated successfully!")

	return nil
}

// 获取图片的交互信息
func getInteractInfo(imageId int) (Interact, error) {
	db, err := openDB()
	if err != nil {
		log.Println("Open database failed", err)
		return Interact{}, err
	}
	defer db.Close()

	query := "SELECT * FROM INTERACT WHERE IMAGEID = ?"
	row := db.QueryRow(query, imageId)

	interact := Interact{}
	err = row.Scan(&interact.ImageId, &interact.Like, &interact.Favorite, &interact.Comment)
	if err != nil {
		return interact, err
	}

	return interact, nil
}
