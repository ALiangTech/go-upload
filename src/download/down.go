package download

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// 解析url 获取文件名 完整的路径名称
func GetDownloadFilename(cdn string) (string, string) {
	cdnurl, err := url.Parse(cdn)
	if err != nil {
		fmt.Println("获取文件名失败")
		return "获取文件名失败", "获取文件名失败"
	}
	splitCdnUrl := strings.Split(cdnurl.Path, "/")
	endIndex := len(splitCdnUrl) - 1
	u := splitCdnUrl[0:]
	return u[endIndex], cdnurl.Path
}
func Download(url string) {
	filename, _ := GetDownloadFilename(url)
	fmt.Println("下载开始--", filename)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("下载出错--", filename)
	}
	defer resp.Body.Close()
	f, err := os.Create("src/file/" + filename)
	io.Copy(f, resp.Body)
}
