package readfile

import (
	"bufio"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// 用户信息文件读取,返回一个字符串切片
func ReadFile(filePath string) (fileContent []string, err error) {
	var (
		file *os.File
		line string
		r    *bufio.Reader
	)

	// 读取一个文件的内容
	if file, err = os.Open(filePath); err != nil {
		logrus.Error("open file err:", err)
		return
	}

	// 处理结束后关闭文件
	defer file.Close()

	// 使用bufio读取
	r = bufio.NewReader(file)

	for {
		// 以分隔符形式读取,比如此处设置的分割符是\n,则遇到\n就返回,且包括\n本身 直接返回字符串
		line, err = r.ReadString('\n')
		fileContent = append(fileContent, line)

		// 读取到末尾退出
		if err == io.EOF {
			break
		}
		// 有错误直接返回
		if err != nil {
			logrus.Error("Error read user profile:", err)
			break
		}
	}
	return fileContent, nil
}
