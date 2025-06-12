package main

import (
	"bytes"
	"client/flight"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DroneStatusReq struct {
	OrderID        string `json:"orderID"`
	FlightCode     string `json:"flightCode"`
	Sn             string `json:"sn"`
	FlightStatus   string `json:"flightStatus"`
	ManufacturerID string `json:"manufacturerID"`
	UasID          string `json:"uasID"`
	TimeStamp      string `json:"timeStamp"`
	UasModel       string `json:"uasModel"`
	Coordinate     int    `json:"coordinate"`
	Longitude      int64  `json:"longitude"`
	Latitude       int64  `json:"latitude"`
	HeightType     int    `json:"heightType"`
	Height         int    `json:"height"`
	Altitude       int    `json:"altitude"`
	VS             int    `json:"VS"`
	GS             int    `json:"GS"`
	Course         int    `json:"course"`
	SOC            int    `json:"SOC"`
	RM             int    `json:"RM"`
	WindSpeed      int    `json:"windSpeed"`
	WindDirect     int    `json:"windDirect"`
	Temperture     int    `json:"temperture"`
	Humidity       int    `json:"humidity"`
}

type DroneStatusResp struct {
	Code int `json:"code"`
}

func main() {
	// 命令行参数定义
	latPtr := flag.Float64("lat", 22.8007210, "初始纬度（如 22.8007210）")
	lonPtr := flag.Float64("lon", 113.9530990, "初始经度（如 113.9530990）")
	bearingPtr := flag.Float64("bearing", 45.0, "初始飞行方向角度（0~360）")
	uavIdPtr := flag.String("id", "uav1", "无人机ID")

	flag.Parse()

	lat := *latPtr
	lon := *lonPtr
	bearing := *bearingPtr // 飞行方向（度），例如 45° 表示东北方向
	// 服务器 URL（请根据你的服务器地址修改）
	serverURL := "http://localhost:19999/api/drone/status"

	// 初始经纬度
	speed := 10.0 // 无人机速度（米/秒）
	interval := 1 // 每 1 秒计算一次位置
	orderID := *uavIdPtr
	flightCode := *uavIdPtr
	height := 150.0
	altitude := 160.0

	// 创建一个定时器，每秒执行一次
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// 持续发送请求
	id := 0
	for range ticker.C {
		// 构造要发送的数据
		droneStatusReq := DroneStatusReq{
			OrderID:        orderID,
			FlightCode:     flightCode,
			Sn:             flightCode,
			FlightStatus:   "Inflight",
			ManufacturerID: "112233",
			UasID:          "Uas-default",
			TimeStamp:      time.Now().Format("20060102150405"),
			UasModel:       "Uas-default",
			Coordinate:     1,
			Longitude:      int64(lon * 10000000), // 转换为整数
			Latitude:       int64(lat * 10000000), // 转换为整数
			HeightType:     1,
			Height:         int(height * 10),
			Altitude:       int(altitude * 10),
			VS:             int(speed * 10),
			GS:             0,
			Course:         int(bearing * 10),
			SOC:            100,
			RM:             10,
			WindSpeed:      10,
			WindDirect:     90,
			Temperture:     30,
			Humidity:       50,
		}

		jsonData, err := json.Marshal(droneStatusReq)
		if err != nil {
			break
		}

		// 发送 POST 请求
		resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("POST 请求失败:", err)
			continue
		}

		// 读取并打印服务器响应
		fmt.Printf("Step %d: 纬度: %.7f, 经度: %.7f\n", id, lat, lon)

		lat, lon = flight.GetNewLatLon(lat, lon, speed, bearing, interval)
		resp.Body.Close()

		id++
		if id%30 == 0 {
			bearing += 180.0
			if bearing >= 360.0 {
				bearing -= 360.0
			}
		}
	}
}
