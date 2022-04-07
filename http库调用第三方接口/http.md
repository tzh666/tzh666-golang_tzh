## GO 发起HTTP请求调用接口

### 一、Go发起GET请求

#### 1.1、不带参数的GET请求

```go
func HttpGet(url string) error {
	// 请求xx网站首页
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// 延迟关闭
	defer resp.Body.Close()

	//
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// 请求结果
	fmt.Println(string(body))

	// 请求头
	fmt.Println(resp.Header)

	// 请求相应码
	fmt.Println(resp.Status)

	return nil
}
```

#### 1.2、带参数的GET请求

- 静态的参数，即放到url中的参数手动拼接url，然后发送GET请求即可，也是调上面的HttpGet函数即可

```go
HttpGet("http://www.baidu.com?name=Paul_Chan&age=26")
```

- 动态参数的GET请求，代码如下

```sh
// URL+参数拼接
func GetUrlPath(dUrl string, param map[string]string) (urlPath string, err error) {
	// 参数设置对象
	params := url.Values{}

	// 请求URL
	Url, err := url.Parse(dUrl)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 参数设置, K V形式, 可以设置多个
	for k, v := range param {
		params.Set(k, v)
	}
	// 有中文参数,这个方法会进行URLEncode, 没有中文也加上这两句
	Url.RawQuery = params.Encode()
	//携带参数以后的URL
	urlPath = Url.String()

	return urlPath, nil
}

func main() {
	param := make(map[string]string, 0)

	param["name"] = "Paul_Chan"
	param["age"] = "0"

	URL, err := GetUrlPath("http://www.baidu.com", param)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(URL)
}
```

#### 1.3、解析JSON类型的返回结果

```go
// 定义一个结构体,接收返回的json数据
type result struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

func main() {
	// 请求JSON返回的数据：{"age":18,"name":"tom"}, 自定义个接口模拟下数据即可
	resp, err := http.Get("http://127.0.0.1:8000/test")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 延时关闭
	defer resp.Body.Close()

	// 返回的结果resp.Body
	body, _ := ioutil.ReadAll(resp.Body)

	var res result
	// 把请求到的数据Unmarshal到res中
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println(err)
		return
	}
    
    // 拿到数据以后,就可以自行处理数据了,比如存到数据库里
	fmt.Printf("%#v", res)
}
```

#### 1.4、GET请求添加请求头

- 有时需要在请求的时候设置头参数、cookie之类的数据，就可以使用http.Do方法

```go
client := &http.Client{}

	url := "https://www.kuaidaili.com/free/inha"

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Print(err)
		return
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36")
	request.Header.Add("Host", "www.kuaidaili.com")

	resp, err := client.Do(request)
	if err != nil {
		fmt.Print(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(body))
```



### 二、Go发起POST请求

#### 2.1、方式一

```go
func main() {
	urlP := "http://127.0.0.1:8080/loginJSON"

	// urlValues ---> map[string][]string
	urlValues := url.Values{}

	// 参数
	urlValues.Add("user", "q1mi")
	urlValues.Add("password", "123456")

    // 通过http.PostForm请求第三方接口
	resp, err := http.PostForm请求第三方接口(urlP, urlValues)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	// 返回结果
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))

	// 状态码
	fmt.Println(resp.Status)
}
```

#### 2.2、方式二

```go
func main() {
	urlP := "http://127.0.0.1:8080/loginJSON"

	// 参数
	urlValues := url.Values{
		"user":     {"q1mi"},
		"password": {"123456"},
	}

	// 编码
	reqBody := urlValues.Encode()

	// 通过http.Post请求第三方接口
	resp, err := http.Post(urlP, "text/html", strings.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
    
    fmt.Println(resp.StatusCode)
}
```

#### 2.3、发送JSON数据的post请求

- 使用client

```go
func main() {
	urlP := "http://127.0.0.1:8080/loginJSON"

	client := &http.Client{}

	// 参数
	data := make(map[string]interface{})
	data["user"] = "q1mi"
	data["password"] = "123456"

	// Marshal
	bytesData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", urlP, bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body), resp.StatusCode)
}
```

#### 2.3、发送JSON数据的post请求

- 不用client的post请求

```go
func main() {
	urlP := "http://127.0.0.1:8080/loginJSON"

	// 参数
	data := make(map[string]interface{})
	data["user"] = "q1mi"
	data["password"] = "123456"

	// Marshal
	bytesData, _ := json.Marshal(data)

	// 使用http.Post请求
	resp, _ := http.Post(urlP, "application/json", bytes.NewReader(bytesData))
	body, _ := ioutil.ReadAll(resp.Body)

	// 打印返回结果
	fmt.Println(string(body), resp.StatusCode)

}
```