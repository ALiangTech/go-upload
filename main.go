package main

import (
	f "chrome-cdn-upload/src/download"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/upyun/go-sdk/upyun"
)

var up = upyun.NewUpYun(&upyun.UpYunConfig{
	Bucket:   "jsdeliver",
	Operator: "a17602545401234",
	Password: "YlQ9foPJUHkL2BycLEBojTTvIg4HdpKw",
})
var cdn = "http://jsdeliver.test.upcdn.net"

func send(url string, wg *sync.WaitGroup, c chan string) {
	defer wg.Done()
	// testurl := "https://cdn.jsdelivr.net/npm/vue@2.6.14/dist/vue.min.js"
	f.Download(url)
	fmt.Println("上传开始")
	filename, savePath := f.GetDownloadFilename(url) // 存储路径
	err := up.Put(&upyun.PutObjectConfig{
		Path:      savePath,
		LocalPath: "src/file/" + filename,
	})
	if err == nil {
		var build strings.Builder
		build.WriteString(cdn)
		build.WriteString(savePath)
		remoteCdnUrl := build.String()
		fmt.Println("上传完成", remoteCdnUrl)
		c <- remoteCdnUrl
		fmt.Println("上传完成--", remoteCdnUrl)
	}
}
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./src/html/*")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"title": "CDN",
		})
	})
	type UploadJsonBody struct {
		Urls []string `form:"urls" json:"urls" uri:"urls" xml:"urls" binding:"required"`
	}

	r.POST("/upload", func(c *gin.Context) {
		var wg sync.WaitGroup
		cs := make(chan string)
		var jsonUrls UploadJsonBody
		var callbackUrls []string
		c.ShouldBindJSON(&jsonUrls)
		urls := jsonUrls.Urls
		// 并发下载 和 上传
		for _, url := range urls {
			wg.Add(1)
			go send(url, &wg, cs)
		}
		go func() {
			wg.Wait()
			close(cs)
		}()
		for i := range cs {
			fmt.Println("23", i)
			callbackUrls = append(callbackUrls, i)
		}
		c.JSON(200, gin.H{
			"message": callbackUrls,
		})
	})
	r.Run(":8089")
}
