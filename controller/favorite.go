package controller

import (
	"github.com/CDsmen/douyin/dal"
	"github.com/CDsmen/douyin/myjwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	strToken := c.Query("token")
	actionType := c.Query("action_type")

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

	// userID是否存在
	var comment Comment
	err = dal.DB.Raw("CALL user_info(?)", claim.UserID).Scan(&comment.User).Error
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}

	// 点赞
	videoid := c.Query("video_id")
	if actionType == "1" {
		err = dal.DB.Raw("CALL add_favorite(?, ?)", claim.UserID, videoid).Scan(&comment).Error
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Favorite succeeded"},
			Comment:  comment,
		})
	} else { // 取消点赞
		err = dal.DB.Raw("CALL del_favorite(?, ?)", claim.UserID, videoid).Scan(&comment).Error
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Delete favorite succeeded"},
		})
	}

	//if _, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, Response{StatusCode: 0})
	//} else {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userid := c.Query("user_id")
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

	// 鉴权不通过
	if strconv.FormatInt(claim.UserID, 10) != userid {
		c.String(http.StatusOK, "Userid != token")
		return
	}

	// 从数据库获取发布列表
	var videosList = []Video{}
	err = dal.DB.Raw("CALL list_favorite(?)", userid).Scan(&videosList).Error
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Mysql list_favorite error"},
		})
	}

	// ??
	for id, _ := range videosList {
		var user User
		err = dal.DB.Raw("CALL user_info(?)", userid).Scan(&user).Error
		if err != nil {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
		}
		videosList[id].Author = user
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videosList,
	})

	//c.JSON(http.StatusOK, VideoListResponse{
	//	Response: Response{
	//		StatusCode: 0,
	//	},
	//	VideoList: DemoVideos,
	//})
}
