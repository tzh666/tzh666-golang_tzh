## GO 发起HTTPS请求调用接口

### 一、GET请求调用HTTPS接口

- tls.LoadX509KeyPair()方法读取证书路径，转换为证书对象；
- x509.NewCertPool()方法创建证书池；
- pool.AppendCertsFromPEM(caCrt)方法将根证书加入到证书池中。
- curl  --cacert /etc/kubernetes/pki/etcd/ca.crt  --cert /etc/kubernetes/pki/etcd/healthcheck-client.crt --key /etc/kubernetes/pki/etcd/healthcheck-client.key  https://192.168.1.70:2379/metrics -k

```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetHttps(url, caCertPath, certFile, keyFile string) ([]byte, error) {
	var (
		// 创建证书池及各类对象
		pool   *x509.CertPool // 我们要把一部分证书存到这个池中
		client *http.Client
		resp   *http.Response
		body   []byte
		err    error
		caCrt  []byte          // 根证书
		cliCrt tls.Certificate // 具体的证书加载对象
	)
	// 读取caCertPath
	caCrt, err = ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	// NewCertPool
	pool = x509.NewCertPool()

	// 解析一系列PEM编码的证书
	pool.AppendCertsFromPEM(caCrt)

	// 具体的证书加载对象
	cliCrt, err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// 把上面的准备内容传入 client
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
			},
		},
	}

	// Get 请求
	resp, err = client.Get(url)
	if err != nil {
		return nil, err
	}
	// 延时关闭
	defer resp.Body.Close()

	// 读取
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 延时关闭
	defer client.CloseIdleConnections()

	return body, nil
}

func main() {
	resp, err := GetHttps("https://192.168.1.70:2379/metrics", "./ca_file/ca.crt", "./ca_file/healthcheck-client.crt", "./ca_file/healthcheck-client.key")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp))
}
```



### 二、在 Header 中添加 token 的 HTTPS 请求

```go
func GetHttpsSkip(url, token string) ([]byte, error) {

	// 创建各类对象
	var (
		client  *http.Client
		request *http.Request
		resp    *http.Response
		body    []byte
		err     error
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
		return nil, nil
	}

	// 加入 token; Authorization、accessToken看你接口的请求头是什么了
	request.Header.Add("accessToken", token)
	resp, err = client.Do(request)
	if err != nil {
		log.Println("GetHttpSkip Response Error:", err)
		return nil, nil
	}

	// 延迟关闭
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ReadAll resp.Body Error:", err)
		return nil, nil
	}

	fmt.Println(resp.StatusCode)

	// 延迟关闭
	defer client.CloseIdleConnections()
	return body, nil
}

func main() {
	urlx := "https://xxxx"
	tokenx := "xxxx"

	resp, _ := GetHttpsSkip(urlx, tokenx)
	fmt.Println(string(resp))
}
```



#### 三、在 Header 中添加 token 的 HTTPS的POST 请求

```go
func main() {
	urlx := "https://xxx"
	tokenx := "xxxx"

	// 创建各类对象
	var (
		client  *http.Client
		request *http.Request
		resp    *http.Response
		body    []byte
		err     error
	)

	// 这里请注意，使用 InsecureSkipVerify: true 来跳过证书验证
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	
	// 参数
    s := []int{1}
	data := make(map[string]interface{})

	data["x"] = "x"
	data["x"] = "x"
	data["x"] = "x"
	data["x"] = "x"
	data["x"] = s
	data["x"] = 1
	data["x"] = "x"

	// Marshal
	bytesData, _ := json.Marshal(data)

	// 获取 request请求
	if request, err = http.NewRequest("POST", urlx, bytes.NewReader(bytesData)); err != nil {
		log.Println("GetHttpSkip Request Error:", err)
		return
	}

	// 加入 token; Authorization、accessToken看你接口的请求头是什么了
	request.Header.Add("accessToken", tokenx)
    // 请求格式json
	request.Header.Add("content-type", "application/json")
	if resp, err = client.Do(request); err != nil {
		log.Println("GetHttpSkip Response Error:", err)
		return
	}

	// 延迟关闭
	defer resp.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println("ReadAll resp.Body Error:", err)
		return
	}

	fmt.Println(string(body), resp.StatusCode)
	// 延迟关闭
	defer client.CloseIdleConnections()
}
```

#### 3.1、封装成函数

```go
type Params struct {
	UserName string `josn:"userName"`
	RealName string `josn:"realName"`
	Serial   string `josn:"serial"`
	Phone    string `josn:"phone"`
	RoleIds  []int  `josn:"roleIds"`
	OrgId    int    `josn:"orgId"`
	ImageUrl string `josn:"imageUrl"`
}

func PostHttpsSkip(urlx, tokenx string, param *Params) (body []byte, err error) {

	// 创建各类对象
	var (
		client  *http.Client
		request *http.Request
		resp    *http.Response
	)

	// 这里请注意，使用 InsecureSkipVerify: true 来跳过证书验证
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	// Marshal---->参数param
	bytesData, _ := json.Marshal(param)

	// 获取 request请求
	if request, err = http.NewRequest("POST", urlx, bytes.NewReader(bytesData)); err != nil {
		log.Println("GetHttpSkip Request Error:", err)
		return
	}

	// 加入 token; Authorization、accessToken看你接口的请求头是什么了
	request.Header.Add("accessToken", tokenx)
	request.Header.Add("content-type", "application/json")
	if resp, err = client.Do(request); err != nil {
		log.Println("GetHttpSkip Response Error:", err)
		return
	}

	// 延迟关闭
	defer resp.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println("ReadAll resp.Body Error:", err)
		return
	}
	// 延迟关闭
	defer client.CloseIdleConnections()
	return body, nil
}

func main() {
	urlx := "https://10.9.240.118:10220/uums/users"
	tokenx := "da1c840f473c4871a2cb29cb02f97bd4"

	data := &Params{
		UserName: "test031",
		RealName: "test031",
		Serial:   "111",
		Phone:    "111",
		RoleIds:  []int{1},
		OrgId:    1,
		ImageUrl: "",
	}

	x, _ := PostHttpsSkip(urlx, tokenx, data)
	fmt.Println(string(x))
}
```

