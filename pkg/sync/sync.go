package sync

import (
	"context"
	"fmt"
	"github.com/xxxapi/slave/pkg/config"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Sync() error {
	// 初始化配置
	conf := config.GetConfig()

	// API endpoint URL
	apiEndpoint := conf.Sync.Endpoint

	// Get remote version from API endpoint
	resp, err := http.Get(apiEndpoint)
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	remoteVersion := strings.TrimSpace(string(body))

	// 尝试从 version.txt 文件中获取本地版本号，如果失败则将其设置为 "1.0"
	localVersionFile := "./data/version.txt"
	localVersionBytes, err := os.ReadFile(localVersionFile)
	var localVersion string
	if err != nil {
		fmt.Println("未检测到版本文件，即将开始全量同步，此过程耗费时间较长，请耐心等待")
		localVersion = "1.0"
	} else {
		localVersion = strings.TrimSpace(string(localVersionBytes))
	}
	// Compare remote and local versions
	if remoteVersion > localVersion {
		// Initialize Minio client
		endpoint := conf.Sync.MinioEndpoint
		accessKeyID := conf.Sync.AccessKey
		secretAccessKey := conf.Sync.SecretKey
		useSSL := conf.Sync.UseSSL
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			fmt.Println(err)
		}

		// Download all objects in the remote folder
		bucketName := conf.Sync.BucketName
		ctx := context.TODO()
		objectCh := minioClient.ListObjects(ctx, bucketName,
			minio.ListObjectsOptions{
				Recursive: true,
			},
		)
		for object := range objectCh {
			if object.Err != nil {
				fmt.Println(object.Err)
			}
			// Check if local file already exists and skip if it does
			localFilePath := "./data/" + object.Key
			if _, err := os.Stat(localFilePath); err == nil {
				continue
			}
			// 下载文件并显示下载进度
			fmt.Printf("Downloading %s...\n", object.Key)
			err = minioClient.FGetObject(context.Background(), bucketName, object.Key, localFilePath, minio.GetObjectOptions{})
			if err != nil {
				fmt.Println(err)
			}
		}

		// 下载完成后打印消息
		fmt.Println("All objects downloaded successfully.")

		// Update local version file with remote version
		err = os.WriteFile(localVersionFile, []byte(remoteVersion), 0644)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("文件已拉去至远程版本", remoteVersion)
	} else {
		fmt.Println("本地版本", localVersion, "已经是最新的了")
	}
	return nil
}

func UpdatePeriodically() {
	for {
		time.Sleep(60 * time.Minute)
		err := Sync()
		if err != nil {
			fmt.Println("更新失败 :", err)
		}
	}
}
