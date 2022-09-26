package middleware

import (
	"fmt"
	"os"
	"time"

	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	//日志存放的位置
	filePath := "log/log"

	scr, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755) //0755表示用户具有读/写/执行权限
	if err != nil {
		fmt.Println("err:", err)
	}
	//初始化日志器
	logger := logrus.New()
	//将日志输出到scr中
	logger.Out = scr
	//设置日志等级
	logger.SetLevel(logrus.DebugLevel) //DebugLevel为最低等级

	logWriter, _ := retalog.New(
		filePath+"%Y%m%d.log",                  //年月日
		retalog.WithMaxAge(7*24*time.Hour),     //存放的最大时限
		retalog.WithRotationTime(24*time.Hour), //日志分割的时间点
	)

	//写入内容,都写入到logWriter
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", //时间输出格式
	})

	logger.AddHook(Hook)

	return func(c *gin.Context) {
		//设置开始时间
		startTime := time.Now()
		c.Next()
		stopTime := time.Since(startTime).Milliseconds()
		spendTime := fmt.Sprintf("%d ms", stopTime)
		hostName, err := os.Hostname() //主机名
		if err != nil {
			hostName = "unknown"
		}
		statusCode := c.Writer.Status()    //状态码
		clientIP := c.ClientIP()           //客户端IP
		userAgent := c.Request.UserAgent() //代理
		dataSize := c.Writer.Size()        //返回的文件大小
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method   //请求方法
		path := c.Request.RequestURI //请求链接
		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"status":    statusCode,
			"SpendTime": spendTime,
			"Ip":        clientIP,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})
		//系统内部错误
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode > 500 {
			entry.Error()
		} else if statusCode > 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
