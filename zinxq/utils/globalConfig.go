package utils

import (
	"encoding/json"
	"fmt"
	"github.com/winterqin/zinxq/ziface"
	"io/ioutil"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数也可以通过 用户根据 zinxq.json来配置
*/
type GlobalConfig struct {
	TcpServer ziface.IServer //当前Zinx的全局Server对象
	Host      string         //当前服务器主机IP
	TcpPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称
	IPVersion string         //当前Zinx版本号

	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConnectionNum int    //当前服务器主机允许的最大链接个数
}

/*
定义一个全局的对象
*/
var Config *GlobalConfig

// Reload 读取用户的配置文件
func (g *GlobalConfig) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &Config)
	if err != nil {
		panic(err)
	}
}

func init() {
	Config = &GlobalConfig{
		TcpServer:        nil,
		Host:             "127.0.0.1",
		TcpPort:          8080,
		Name:             "winterServer",
		IPVersion:        "tcp4",
		MaxPacketSize:    512,
		MaxConnectionNum: 128,
	}
	//Config.Reload()
}
