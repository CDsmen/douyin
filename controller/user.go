package controller

import (
	"github.com/CDsmen/douyin/dal"
	"github.com/CDsmen/douyin/myjwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var userId int64
	err := dal.DB.Raw("CALL register(?, ?)", username, password).Scan(&userId).Error
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Mysql Register Failed"},
		})
		return
	} else {
		if userId == 0 {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
			})
			return
		} else {
			// 获得token
			claims := &myjwt.JWTClaims{
				UserID:   userId,
				Username: username,
				Password: password,
			}
			claims.IssuedAt = time.Now().Unix()
			claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(myjwt.ExpireTime)).Unix()
			signedToken, err := myjwt.GetToken(claims)
			if err != nil {
				c.String(http.StatusNotFound, err.Error())
				return
			}

			myjwt.TakenGetMap(signedToken)

			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   userId,
				Token:    signedToken,
			})
			return
		}
	}

	//if _, exist := usersLoginInfo[token]; exist { // 检测token
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	//	})
	//} else { // 新增用户
	//	atomic.AddInt64(&userIdSequence, 1)
	//	newUser := User{
	//		Id:   userIdSequence,
	//		Name: username,
	//	}
	//	usersLoginInfo[token] = newUser
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   userIdSequence,
	//		Token:    username + password,
	//	})
	//}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var userId int64
	err := dal.DB.Raw("CALL login(?, ?)", username, password).Scan(&userId).Error

	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Mysql Login Failed"},
		})
		return
	} else {
		if userId == 0 {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User not exist"},
			})
			return
		} else {
			// 获得token
			claims := &myjwt.JWTClaims{
				UserID:   userId,
				Username: username,
				Password: password,
			}
			claims.IssuedAt = time.Now().Unix()
			claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(myjwt.ExpireTime)).Unix()
			signedToken, err := myjwt.GetToken(claims)
			if err != nil {
				c.String(http.StatusNotFound, err.Error())
				return
			}

			myjwt.TakenGetMap(signedToken)

			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   userId,
				Token:    signedToken,
			})
		}
	}

	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   user.Id,
	//		Token:    token,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}

// 调用mysql的存储过程"user_info" 参数为：user_id
func UserInfo(c *gin.Context) {
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

	var user User
	err = dal.DB.Raw("CALL user_info(?)", userid).Scan(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
		return
	}

	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 0},
	//		User:     user,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}
