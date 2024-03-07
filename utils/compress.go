package utils

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func compress(csvFilePath string) bool {
	// 压缩后的 Gzip 文件路径
	gzipFilePath := fmt.Sprintf("%v.gz", csvFilePath)

	// 打开 CSV 文件
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		logrus.Warnf("[compress] Error opening CSV file: %v", err)
		return false
	}
	csvFileInfo, err := csvFile.Stat()
	if err != nil {
		logrus.Warnf("[compress] Error csvFile.Stat() get csvFileInfo: %v", err)
		return false
	}
	if csvFileInfo.Size() <= 0 {
		logrus.Warnf("[compress] Error  压缩原文件大小为0: %v", err)
		return false
	}
	defer csvFile.Close()
	// 创建 Gzip 文件
	gzipFile, err := os.Create(gzipFilePath)
	if err != nil {
		logrus.Warnf("[compress] Error creating Gzip file: %v", err)
		return false
	}
	defer gzipFile.Close()

	// 创建 Gzip Writer
	gzipWriter := gzip.NewWriter(gzipFile)
	defer gzipWriter.Close()

	// 创建 CSV Reader
	csvReader := csv.NewReader(csvFile)

	// 创建 CSV Writer
	csvWriter := csv.NewWriter(gzipWriter)
	defer csvWriter.Flush()

	// 逐行读取 CSV 文件并写入 Gzip 文件
	for {
		record, ferr := csvReader.Read()
		if ferr != nil {
			break // 读取完成或出错时退出循环
		}
		if werr := csvWriter.Write(record); err != nil {
			logrus.Warnf("Error writing to Gzip file: %v", werr)
			return false
		}
	}
	logrus.Warnf("CSV file compressed to Gzip successfully.")
	return true
}
