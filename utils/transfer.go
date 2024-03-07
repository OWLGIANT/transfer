package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getCsv(gz bool) []string {
	currentDir, err := os.Getwd()
	if err != nil {
		logrus.Warnf("Error getting current directory: %v", err)
		return []string{}
	}
	var csvSlice []string
	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Warnf("Error walking path: %v", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileName := info.Name()
		if strings.Contains(fileName, "csv") {
			if gz {
				if strings.Contains(fileName, "gz") {
					csvSlice = append(csvSlice, fileName)
				}
			} else {
				if !strings.Contains(fileName, "gz") {
					csvSlice = append(csvSlice, fileName)
				}
			}
		}
		return nil
	})
	if err != nil {
		logrus.Warnf("Error walking directory: %v", err)
	}
	if len(csvSlice) == 0 {
		logrus.Warnf("%s文件夹内没有CSV文件--", currentDir)
	}
	return csvSlice
}
func getCsvDayTime(name string) (data time.Time) {
	var err error
	a := strings.Split(name, ".")[0]
	for _, v := range strings.Split(a, "_") {
		data, err = time.Parse(time.DateOnly, v)
		if err == nil {
			return data
		}
	}
	return
}
func ifCanCT(fileName string) bool {
	fileDayTime := getCsvDayTime(fileName)
	if time.Now().After(fileDayTime.AddDate(0, 0, 1)) {
		return true
	}
	return false
}
func CompressAndTransferLoop() {
	for _, v := range getCsv(false) { //csv压缩
		if ifCanCT(v) {
			//1.压缩
			if compress(v) {
				logrus.Infof("%s 文件压缩成功", v)
				gzipFilePath := fmt.Sprintf("%v.gz", v)
				//2.传输
				if ScpFile(gzipFilePath) {
					logrus.Infof("%v 传输成功", v)

					//3.移除.gz文件
					if err := os.Remove(gzipFilePath); err != nil {
						logrus.Errorf("%v 删除gz文件失败", v)
					} else {
						logrus.Infof("%v 删除gz文件成功", v)
						//4.移除csv文件
						if err := os.Remove(v); err != nil {
							logrus.Errorf("%v 删除csv文件失败", v)
						} else {
							logrus.Infof("%v 删除csv文件成功", v)
						}
					}
				} else {
					logrus.Errorf("%v 传输失败", v)
				}
			} else {
				logrus.Errorf("%s 文件压缩失败", v)
			}
		}
	}
}
