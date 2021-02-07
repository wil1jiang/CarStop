package model

// Kratos hello kratos.
type Kratos struct {
	Hello string
}

type Article struct {
	ID int64
	Content string
	Author string
}

type Data struct {
	Caruser string `json:"caruser"`
	Carcode string `json:"carCode"`
	Stopaddr string `json:"stopaddr"`
	PassType string `json:"inOrOut"`
	PassTime string `json:"passTime"`
	ImagePath string `json:"imagePath"`
}

type MQCarInfo struct {
	ProjectUuid string `json:"projectUuid"`
	Topic string `json:"topic"`
	Data  `json:"data"`
}

type Config struct{
	ProjectID string
	Stopaddr string
}

type CarInterface struct{
	Ip  string
	Port int
}

type Rabbitmq struct {
	Exchange string
	Topic string
	User string
	Password string
	Ip  string
	Port int
	Vhost string
}