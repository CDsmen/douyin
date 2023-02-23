package controller

import (
	"fmt"
	"github.com/CDsmen/douyin/dal"
	"github.com/CDsmen/douyin/myjwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	strToken := c.Query("token")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	text := c.Query("comment_text")
	commentId := c.Query("comment_id")

	// 参数值检测（非空检测 + 防sql注入）
	if FilteredSQLInject(strToken, 1) || FilteredSQLInject(videoId, 1) || FilteredSQLInject(actionType, 1) ||
		FilteredSQLInject(text, 0) || FilteredSQLInject(commentId, 0) {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "参数值不符合要求"},
		})
		return
	}

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
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}

	// 发布评论
	if actionType == "1" {
		err = dal.DB.Raw("CALL add_comment(?, ?, ?)", videoId, claim.UserID, text).Scan(&comment).Error
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Mysql Comment Pubish Failed"},
			})
			return
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 0, StatusMsg: "Publishing succeeded"},
				Comment:  comment,
			})
			return
		}
	} else { // 删除评论
		err = dal.DB.Raw("CALL del_comment(?)", commentId).Scan(&comment).Error
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Mysql Comment Delete Failed"},
			})
			return
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 0, StatusMsg: "Delete succeeded"},
			})
			return
		}
	}

	//if user, exist := usersLoginInfo[token]; exist {
	//	if actionType == "1" {
	//		text := c.Query("comment_text")
	//		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
	//			Comment: Comment{
	//				Id:         1,
	//				User:       user,
	//				Content:    text,
	//				CreateDate: "05-01",
	//			}})
	//		return
	//	}
	//	c.JSON(http.StatusOK, Response{StatusCode: 0})
	//} else {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//}
}

func CommentList(c *gin.Context) {
	strToken := c.Query("token")
	videoId := c.Query("video_id")

	// 参数值检测（非空检测 + 防sql注入）
	if FilteredSQLInject(strToken, 1) || FilteredSQLInject(videoId, 1) {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "参数值不符合要求"},
		})
		return
	}

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
	fmt.Println(claim)

	// 从数据库获取发布列表
	var commentsList = []Comment{}
	err = dal.DB.Raw("CALL list_comment(?)", videoId).Scan(&commentsList).Error
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Mysql list_comment error"},
		})
		return
	}

	// 补充user
	for id := range commentsList {
		var user User
		err = dal.DB.Raw("CALL user_info(?)", commentsList[id].Userid).Scan(&user).Error
		if err != nil {
			c.JSON(http.StatusOK, CommentListResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
			return
		}
		commentsList[id].User = user

	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentsList,
	})
	return
}
