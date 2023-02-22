package controller

import (
	"fmt"
	"github.com/CDsmen/douyin/dal"
	"github.com/CDsmen/douyin/myjwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	latesttime := c.Query("latest_time")
	strToken := c.Query("token")

	// token不存在
	err := myjwt.FindToken(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	// 解析token
	claim, err := myjwt.VerifyAction(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	// 从数据库获取发布列表
	var videosList = []Video{}
	err = dal.DB.Raw("CALL feed(?, ?)", claim.UserID).Scan(&latesttime).Error
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Mysql feed error"},
		})
		return
	}
	fmt.Println("videosList: ", videosList)

	// 补充user
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
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videosList,
		NextTime:  time.Now().Unix(),
	})
}
