package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func GenTimer(userInfos []*UserInfo, timer *time.Timer) *time.Timer {
	if len(userInfos) == 0 {
		return nil
	}
	now := time.Now()
	countdown := userInfos[0].GetTimeout(now)
	for _, user := range userInfos {
		current := user.GetTimeout(now)
		if current.Microseconds() < countdown.Microseconds() {
			countdown = current
		}
	}
	if timer == nil {
		timer = time.AfterFunc(countdown, func() {
			now := time.Now()
			for _, user := range userInfos {
				if user.GetTimeout(now).Microseconds() <= 0 {
					for _, contact := range user.Contacts {
						contact.Send(user.EmergencyMessage)
					}
				}
			}
		})
	} else {
		timer.Reset(countdown)
	}
	return timer
}

func main() {
	userfile := "users.json"
	var userInfos []*UserInfo
	content, err := ioutil.ReadFile(userfile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &userInfos)
	if err != nil {
		panic(err)
	}
	lock := sync.Mutex{}
	timer := GenTimer(userInfos, nil)
	router := gin.Default()
	router.POST("/heartbeat", func(c *gin.Context) {
		lock.Lock()
		defer lock.Unlock()
		key := c.PostForm("key")
		now := time.Now()
		log.Printf("key: %v time: %v\n", key, now)
		for _, user := range userInfos {
			if user.Key == key {
				user.LastHeartBeartTime = time.Now()
			}
		}
		timer = GenTimer(userInfos, timer)
	})
	router.Run("0.0.0.0:12345")
}
