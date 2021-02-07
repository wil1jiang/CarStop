package mq

import (
	"encoding/json"
	"github.com/prometheus/common/log"
	"github.com/streadway/amqp"
	"lifang/internal/model"
	"strings"
)

var conn *amqp.Connection
var channel *amqp.Channel
var hasMQ bool =false
var exchanges string

func SetupMQ(mqAddr string)(err error){
	if channel == nil{
		conn,err = amqp.Dial(mqAddr)
		if err!=nil{
			log.Info("rabbitmq dial failed!")
			return err
		}

		channel,err = conn.Channel()
		if err != nil{
			log.Info("rabbitmq connect init failed!")
		}
		hasMQ = true
	}
	return nil
}

func Publish(exchange string,topic string,carinfo model.MQCarInfo)(err error){
	if topic=="" || !strings.Contains(exchanges,exchange){
		err = channel.ExchangeDeclare(exchange,"topic",true,false,false,true,nil)
		if err!=nil{
			log.Info("channel ExchangeDeclare failed!")
			return err
		}
	}
	msg,_:= json.Marshal(carinfo)
	err = channel.Publish(exchange,topic,false,false,amqp.Publishing{
			Headers: nil,
			ContentType: "text/plain",
			Body: msg,
	})
	if err != nil{
		log.Info("rabbitmq publish failed!")
	}
	return nil
}

func Close(){
	channel.Close()
	conn.Close()
	hasMQ =false
}
