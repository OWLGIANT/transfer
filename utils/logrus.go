package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

const Index = "bpy_35.75.62.181"

type CHook struct{}

func NewCHook() *CHook {
	return &CHook{}
}

func (hook CHook) Fire(entry *logrus.Entry) error {
	entry.Data["index"] = Index
	log, err := entry.Bytes()
	if err != nil {
		fmt.Println("Fire entry Error :", err)
	}
	go LogInsert(http.MethodPost, "http://57.180.139.206/logrus/insert", log)
	return nil
}

func (hook CHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func LogInit() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{}) //json格式数据
	logrus.AddHook(NewCHook())
}

func LogInsert(method, url string, preload []byte) (response []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(preload))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return
	}
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	response, err = io.ReadAll(resp.Body)
	return
}
