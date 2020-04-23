package main

import (
	"experiment-judge-server/config"
	"experiment-judge-server/message"
	"experiment-judge-server/model"
	"experiment-judge-server/util/redis"
	"fmt"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net"
)

var (
	cfg = pflag.StringP("config", "c", "", "server config file path.")
)

func main() {
	pflag.Parse()

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()

	// init redis
	redis.Client.Init()
	defer redis.Client.Close()

	// init consumer
	client := message.GetKafkaClient()
	client.Consumer()
	defer client.Close()

	// 创建 listener
	listener, err := net.Listen(viper.GetString("type"), viper.GetString("url"))
	if err != nil {
		log.Errorf(err, "Error listening")
		return //终止程序
	}
	log.Infof("tcp server listen in %s", viper.GetString("port"))
	// 监听并接受来自客户端的连接
	for {
		_, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		// TODO 由 api server 探测使用，判断该 judge 是否还存活
	}
}
