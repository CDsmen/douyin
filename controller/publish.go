package controller

import (
	"bytes"
	"fmt"
	"github.com/CDsmen/douyin/dal"
	"github.com/CDsmen/douyin/myjwt"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// GenerateVideoCover 获取封面
func GenerateVideoCover(inFileName string, frameNum int, coverName string) string {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}

	filePath := "public/video_cover/" + coverName + ".jpg"
	outFile, err := os.Create(filePath)
	fmt.Println(filePath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, buf)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	coverCul := SeverIp + ":8080" + "/static/video/" + coverName + ".jpg"
	return coverCul
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

	// 存储到本地
	filename := fmt.Sprintf("%v", claim.UserID) + fmt.Sprintf("%v", rand.Int63())
	filePath := "public/video/" + filename + ".mp4"
	playUrl := SeverIp + ":8080" + "/static/video/" + filename + ".mp4"
	fmt.Println("playurl: ", playUrl)
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

	// 提取第一帧作为封面 拿到coverurl
	coverurl := GenerateVideoCover(filePath, 1, filename)

	// 保存进数据库
	err = dal.DB.Raw("CALL add_vedio(?, ?, ?, ?)", claim.UserID, title, playUrl, coverurl).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Mysql add_video error",
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
	err = dal.DB.Raw("CALL list_vedio(?)", userid).Scan(&videosList).Error
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
	for id, _ := range videosList {
		videosList[id].Author = user
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videosList,
	})
}
