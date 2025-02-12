package controller

import (
	"math"
	"net/http"
	"strconv"
	"time"
	"x-tiktok/service"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
	NextTime  int64           `json:"next_time"`
}

// Feed 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
// 接口详情：
//
//	参数：latest_time		最新投稿时间戳，不填则表示当前时间
//		  token				用户登录后才带上
//	返回值：
//			status_code		0成功，其他失败
//			status_msg		返回状态描述
//			next_time		本次返回视频中发布时间最早的时间，作为下次请求的latest_time
//			video_list		返回视频列表
//				id			视频id
//				author		作者信息：和user一致
//				play_url	视频播放地址
//				cover_url	封面地址
//				favorite_count	视频的点赞总数
//				comment_count	视频的评论总数
//				is_favorite		true-已点赞，false-未点赞
//				title			视频标题
func Feed(c *gin.Context) {
	latestTime := c.Query("latest_time")
	//log.Println("返回视频的最新投稿时间戳:", latestTime)
	var convTime time.Time
	if latestTime != "0" {
		t, _ := strconv.ParseInt(latestTime, 10, 64)
		if t > math.MaxInt32 {
			convTime = time.Now()
		} else {
			convTime = time.Unix(t, 0)
		}
	} else {
		convTime = time.Now()
	}
	//log.Println("返回视频的最新投稿时间:", convTime)
	// 获取登录用户的 id（等待用户模块存入用户id到context）
	//userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	userId := c.GetInt64("userId")
	videoService := service.GetVideoServiceInstance()
	videos, nextTime, err := videoService.Feed(convTime, userId)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 1, StatusMsg: "刷新视频流失败"},
			VideoList: nil,
			NextTime:  nextTime.Unix(),
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "刷新视频流成功!"},
		VideoList: videos,
		NextTime:  nextTime.Unix(),
	})
}
