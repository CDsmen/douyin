package controller

import (
	"fmt"
	"github.com/CDsmen/douyin/dal"
	"github.com/CDsmen/douyin/myjwt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	strToken := c.PostForm("token")
	title := c.PostForm("title")

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

	// 获取上传的文件
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "获取上传的文件 error",
		})
		return
	}

	// 读取文件数据
	fileBytes, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "读取文件数据 error",
		})
		return
	}
	defer fileBytes.Close()

	// 获取当前的工作目录
	wd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 存储到本地 ?? 更改filename
	filePath := wd + "/public/video/" + file.Filename
	playurl := SeverIp + ":8080" + "/static/video/" + file.Filename
	fmt.Println("playurl: ", playurl)
	outFile, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, fileBytes)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "io.Copy error",
		})
		return
	}

	//// 提取第一帧作为封面 拿到coverurl
	coverurl := "11"
	//inputFile, err := os.Open(filePath)
	//if err != nil {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}
	//defer inputFile.Close()
	//
	//// 创建转码器
	//tc := transcoder.New()
	//
	//// 设置输入文件
	//if err := tc.Input(inputFile); err != nil {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}
	//
	//// 提取第一帧作为封面图片
	//tc.OutputFormat("image2")
	//tc.OutputOptions("-vframe")

	// 保存进数据库
	err = dal.DB.Raw("CALL add_video(?, ?, ?, ?)", claim.Id, title, playurl, coverurl).Error
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Mysql add_video error"},
		})
		return
	}

	// 成功上传
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  file.Filename + " uploaded successfully",
	})

	//filename := filepath.Base(data.Filename)
	//user := usersLoginInfo[token]
	//finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	//saveFile := filepath.Join("./public/", finalName)
	//if err := c.SaveUploadedFile(data, saveFile); err != nil {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, Response{
	//	StatusCode: 0,
	//	StatusMsg:  finalName + " uploaded successfully",
	//})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
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
	err = dal.DB.Raw("CALL list_video(?)", userid).Scan(&videosList).Error
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Mysql list_video error"},
		})
		return
	}

	var user User
	err = dal.DB.Raw("CALL user_info(?)", userid).Scan(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
	for id := range videosList {
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
