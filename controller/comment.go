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

	// 发布评论
	videoid := c.Query("video_id")
	if actionType == "1" {
		text := c.Query("comment_text")
		err = dal.DB.Raw("CALL add_comment(?, ?, ?)", videoid, claim.UserID, text).Scan(&comment).Error
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
		commentid := c.Query("comment_id")
		err = dal.DB.Raw("CALL del_comment(?, ?)", videoid, commentid).Scan(&comment).Error
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

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	video_id := c.Query("video_id")
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
	fmt.Println(claim)

	// 从数据库获取发布列表
	var commentsList = []Comment{}
	err = dal.DB.Raw("CALL list_comment(?)", video_id).Scan(&commentsList).Error
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Mysql list_comment error"},
		})
		return
	}
	fmt.Println("commentsList: ", commentsList)

	// 补充user
	for id, _ := range commentsList {
		var user User
		err = dal.DB.Raw("CALL user_info(?)", commentsList[id].Userid).Scan(&user).Error
		if err != nil {
			c.JSON(http.StatusOK, VideoListResponse{
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
}
