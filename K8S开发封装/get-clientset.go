package utils

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// 获取kubeconfig文件的父路径
func homeDir() string {
	// 获取环境linux变量, HOME=/root
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	// 如果是Windows
	return os.Getenv("USERPROFILE")
}

func ReadKubeConf() (clientset *kubernetes.Clientset, err error) {
	var (
		config     *rest.Config // 可以在集群内部访问,也可以在集群外部访问。集群内部是在Pod中访问
		kubeconfig *string      // 集群外部是通过KubeConfig访问(Kubelet)
	)

	// 获取kubeconfig配置文件------这种方式是通过kubelet类型形式获取clientset
	if home := homeDir(); home != "" {
		// 如果home不等于空,在Linux上应该是等于/root,那么就拼接一个全路径
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(可选)kubeconfig文件的绝对路径")
	} else {
		// 如果config不在默认的路径,那么手动给一个,读取命令行传递进来的
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig文件的绝对路径")
	}

	// 解析文件
	flag.Parse()

	// 使用集群内部模式(ServiceAccount)获取clientset
	if config, err = rest.InClusterConfig(); err != nil {
		// 如果使用内部模式出错则使用外部模式 kubeconfig文件来创建集群配置
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}

	// 创建 Clientset 对象
	if clientset, err = kubernetes.NewForConfig(config); err != nil {
		panic(err.Error())
	}

	return clientset, nil
}
