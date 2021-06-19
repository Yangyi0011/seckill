package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		// 构造商家和商品数据
		buildSellerAndGoodsData()
		wg.Done()
	}()
	go func() {
		// 构造参与秒杀的用户数据
		buildUserData()
		wg.Done()
	}()

	wg.Wait()
	log.Println("数据构造完成")
}

// 构造商家注册、登录，并添加秒杀商品所需的数据
func buildSellerAndGoodsData() {
	var (
		path string = "D:/Program Files/JMeter/测试文件/秒杀系统/商家注册和登录并添加商品.txt"
		file *os.File
		e    error
	)
	// 文件不存在则创建
	if file, e = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755); e != nil {
		log.Fatalln("os.OpenFile() err:", e)
	}
	defer file.Close()

	// 创建 *bufio.Writer
	writer := bufio.NewWriter(file)
	// 输出标题
	title := "username,password,kind,ip,name,img,originPrice,price,amount,stock,startTime,endTime\n"
	writer.WriteString(title)
	content := "jerry,123,1,127.0.0.1,商品-1,,100,9.9,100,100,2021-06-14 22:00:00,2021-07-01 00:00:00\n"
	// 输出内容
	writer.WriteString(content)
	// 清空缓冲区
	if e = writer.Flush(); e != nil {
		log.Fatalln("buildSellerAndGoodsData 清空缓冲区出错：", e)
	}
	log.Println("商家和商品数据构建成功")
}

// 构建参与秒杀的用户数据
func buildUserData() {
	var (
		registerUserFilePath            string = "D:/Program Files/JMeter/测试文件/秒杀系统/用户注册.txt"
		loginUserFilePath               string = "D:/Program Files/JMeter/测试文件/秒杀系统/用户登录并参与商品秒杀.txt"
		registerUserFile, loginUserFile *os.File
		userTotal                       int = 100000
		e                               error
	)
	// 文件不存在则创建
	if registerUserFile, e = os.OpenFile(registerUserFilePath, os.O_CREATE|os.O_WRONLY, 0755); e != nil {
		log.Fatalln("os.OpenFile(registerUserFile) err:", e)
	}
	if loginUserFile, e = os.OpenFile(loginUserFilePath, os.O_CREATE|os.O_WRONLY, 0755); e != nil {
		log.Fatalln("os.OpenFile(loginUserFile) err:", e)
	}
	defer loginUserFile.Close()
	defer registerUserFile.Close()

	// 创建 *bufio.Writer
	registerUserFileWriter := bufio.NewWriter(registerUserFile)
	loginUserFileWriter := bufio.NewWriter(loginUserFile)

	// 输出标题
	registerUserTitle := "username,password,kind,ip\n"
	loginUserTitle := "username,password,ip,goodsId\n"
	registerUserFileWriter.WriteString(registerUserTitle)
	loginUserFileWriter.WriteString(loginUserTitle)

	// 输出内容
	// 注册用户数据：customer_序号,密码,用户类别,ip地址
	registerUserContent := "customer_%d,123,0,%s\n"
	// 登录用户数据：customer_序号,密码,ip地址,秒杀商品id
	loginUserContent := "customer_%d,123,%s,1\n"

	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 存储已使用的用户序号
	numMap := make(map[int]bool, userTotal)
	// 循环输出内容
	for i := 1; i <= userTotal; i++ {
		// 用户序号
		var num int
		// 生成一个没有使用过的用户序号
		for {
			num = rand.Intn(userTotal*10) + 1
			if !numMap[num] {
				numMap[num] = true
				break
			}
		}

		// ip地址的生成为 [1~255].[1~255].[1~255].[1~255]
		ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255)+1, rand.Intn(255)+1, rand.Intn(255)+1, rand.Intn(255)+1)
		registerUserFileWriter.WriteString(fmt.Sprintf(registerUserContent, num, ip))
		loginUserFileWriter.WriteString(fmt.Sprintf(loginUserContent, num, ip))
	}
	// 清空缓冲区
	e = registerUserFileWriter.Flush()
	if e != nil {
		log.Fatalln("registerUserFileWriter 清空缓冲区出错：", e)
	}
	e = loginUserFileWriter.Flush()
	if e != nil {
		log.Fatalln("loginUserFileWriter 清空缓冲区出错：", e)
	}
	log.Println("用户数据构建成功")
}
