package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited ")
}

func InitRedis() {

	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

}

func InitMysql()  {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢启动阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		},
	)

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{
		Logger: newLogger,
	})

	fmt.Println("mysql inited")

}

const (
	PublishKey = "websocket"
)

// 发布消息
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("pubish ....", msg)
	err = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)	
	}
	return err
}

// 订阅消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Sub", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)	
		return "", err
	}
	fmt.Println("sub ...", msg.Payload)
	return msg.Payload, err
}