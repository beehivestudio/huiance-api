package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/astaxie/beego/logs"
)

/* 获取文件名 */
func GetFileNameFromUrl(url string) string {
	s := strings.Split(url, "/")
	if l := len(s); l > 0 {
		se := strings.Split(s[l-1], "?")
		return se[0]
	}
	return ""
}

/* 查询文件是否存在 */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/* 下载文件 */
func DownloadFile(url, path string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logs.Info("Download file failed! url msg :%s, resp :%v", url, resp.StatusCode)
		return errors.New("download file failed")
	}

	out, err := os.Create(path)
	if err != nil {
		logs.Error("Create file path failed errmsg :%s", err.Error())
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logs.Error("Download file failed errmsg :%s", err.Error())
		return err
	}

	return nil
}
