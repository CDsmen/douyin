package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/CDsmen/douyin/dal"
	"github.com/CDsmen/douyin/myjwt"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	latestTime := c.Query("latest_time")
	strToken := c.Query("token")

	// 参数值检测（非空检测 + 防sql注入）
	if FilteredSQLInject(latestTime, 0) || FilteredSQLInject(strToken, 0) {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "参数值不符合要求"},
		})
		return
	}

	if latestTime == "0" || latestTime == "" {
		latestTime = strconv.FormatInt(time.Now().Unix(), 10)
	}
	if len([]rune(latestTime)) > 11 {
		latestTime = latestTime[0:10]
	}
	// token不存在
	err := myjwt.FindToken(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	// latestTime = strconv.FormatInt(time.Now().Unix(), 10)

	// 解析token
	claim, err := myjwt.VerifyAction(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	// 从数据库获取发布列表
	var videosList = []Video{}
	err = dal.DB.Raw("CALL feed(?, ?)", claim.UserID, latestTime).Scan(&videosList).Error
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Mysql feed error"},
		})
		return
	}

	// 补充user
	nextTime := time.Now().Unix()
	for id := range videosList {
		var user User
		err = dal.DB.Raw("CALL user_info(?)", videosList[id].Userid).Scan(&user).Error
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
			return
		}
		videosList[id].Author = user
		if videosList[id].CreateTime < nextTime {
			nextTime = videosList[id].CreateTime
		}
	}

	// videosList = DemoVideos
	if len(videosList) == 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videosList,
		})
	} else {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			NextTime:  nextTime,
			VideoList: videosList,
		})
	}
}
