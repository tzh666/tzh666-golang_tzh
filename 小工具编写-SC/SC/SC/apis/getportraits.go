package apis

import (
	"UserInsert/models"
	"UserInsert/setting"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// 调用接口,搞到结构体中【文件看getportraits.json】
func GetPortraits(param *models.Portraits) (body []byte) {
	scUrl := fmt.Sprintf("https://%s:10220/whale-openapi/portraits/query", setting.Conf.SenseCity.SCIP)
	// 调用接口
	// 创建各类对象
	var (
		client    *http.Client
		request   *http.Request
		resp      *http.Response
		err       error
		bytesData []byte
		// token   string
	)

	// 这里请注意，使用 InsecureSkipVerify: true 来跳过证书验证
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	// Marshal---->参数param, 序列化数据给接口当参数传递 bytesData是json数据
	if bytesData, err = json.Marshal(param); err != nil {
		logrus.Error("Marshal param Error:", err)
		return
	}

	// 获取 request请求
	if request, err = http.NewRequest("POST", scUrl, bytes.NewReader(bytesData)); err != nil {
		log.Println("GetHttpSkip Request Error:", err)
		return
	}

	request.Header.Add("accessToken", GetToken())
	request.Header.Add("content-type", "application/json")
	if resp, err = client.Do(request); err != nil {
		log.Println("GetHttpSkip Response Error:", err)
		return
	}

	// 延迟关闭
	defer request.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println("ReadAll resp.Body Error:", err)
		return
	}

	// 延迟关闭
	defer client.CloseIdleConnections()

	// 返回调用接口得到的数据
	return body
}

func DownLoadImage(url, imgPath, fileName string) (written int64, err error) {

	// 创建各类对象
	var (
		client  *http.Client
		request *http.Request
		resp    *http.Response
	)

	// 这里请注意，使用 InsecureSkipVerify: true 来跳过证书验证
	client = &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}

	// 获取 request请求
	request, err = http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("GetHttpSkip Request Error:", err)
		return
	}

	// 设置请求头
	request.Header.Add("Content-Type", "application/json;charset=UTF-8")
	resp, err = client.Do(request)
	if err != nil {
		log.Println("GetHttpSkip Response Error:", err)
		return
	}
	// 延迟关闭
	defer resp.Body.Close()
	defer client.CloseIdleConnections()

	// 获得get请求响应的reader对象
	readerx := bufio.NewReaderSize(resp.Body, 32*1024)

	// 创建文件
	file, err := os.Create(imgPath + fileName)
	if err != nil {
		panic(err)
	}

	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	return io.Copy(writer, readerx)
}
