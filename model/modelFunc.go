package model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func GetArea() ([]Area,error) {
	var Areas []Area

	//先从redis缓存中获取数据
	conn:=GlobalRedis.Get()
	redisbyts,err:=redis.Bytes(conn.Do("get","Areas"))
	if err != nil {
		fmt.Println("redis.Bytes(conn.Do err",err)
		return nil,err
	}
	json.Unmarshal(redisbyts,&Areas)
	//fmt.Println(Areas)
	defer conn.Close()
	if len(Areas)!=0 {
		fmt.Println("从redis数据库获取数据")
	}else {
		err:=GlobalDB.Find(&Areas).Error
		if err != nil {
			fmt.Println("获取数据错误",err)
			return nil,err
		}
		Areasbytes,err:=json.Marshal(Areas)
		if err != nil {
			fmt.Println("json序列化",err)
			return nil,err
		}
		conn.Do("set","Areas",Areasbytes)
		fmt.Println("从mysql数据库获取数据")
	}

	return Areas,nil
}