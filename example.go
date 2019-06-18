package main

import (
	"fmt"
	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"log"
	"mime"
	"os"
	"strings"
)

func main() {
	endpoint := "192.168.35.193:9000"
	accessKeyID := "LFN2RZK0YLOCLXFUZ29H"
	secretAccessKey := "V7xgxS8fF+zRWNiXJXBcfhC7RwTgrSaemtfO75Eo"
	useSSL := false
	/*----------------------------------------------Initialize minio client object.-------------------------------------------------------------*/
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}
	/*----------------------------------------------Create Bucket.-------------------------------------------------------------*/
	bucketName := "test1"
	exist, err := minioClient.BucketExists(bucketName)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if !exist {
		err = minioClient.MakeBucket(bucketName, "us-east-1")
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = minioClient.SetBucketPolicy(bucketName, fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucketName))
		if err != nil {
			err = errors.WithStack(err)
			return
		}
	}

	/*----------------------------------------------Upload File.-------------------------------------------------------------*/
	objectName := "liufan.jpg"
	filePath := "./test.jpg"
	//提取文件后缀类型
	var ext string
	if pos := strings.LastIndexByte(objectName, '.'); pos != -1 {
		ext = objectName[pos:]
		if ext == "." {
			ext = ""
		}
	}
	//返回文件扩展类型
	contentType := mime.TypeByExtension(ext)
	/*----------------------------------------------Put Object.-------------------------------------------------------------*/
	/*n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType:contentType})
	if err != nil {
		log.Fatalln(err)
	}*/
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	n, err := minioClient.PutObject(bucketName, objectName, file, fileInfo.Size(), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	path := "http://" + endpoint + "/" + bucketName + "/" + objectName
	fmt.Println(path)
	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

	/*----------------------------------------------Get Object.-------------------------------------------------------------*/
	/*object, err := minioClient.GetObject("test", "liufan.jpg", minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	localFile, err := os.Create("./get.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err = io.Copy(localFile, object); err != nil {
		fmt.Println(err)
		return
	}*/
}
