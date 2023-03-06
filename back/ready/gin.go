package ready

import (
	"fmt"

	"ViewLog/back/global"
	"ViewLog/back/middleware"
	"ViewLog/back/router"

	"github.com/gin-gonic/gin"
)

func Gin() {
	r := gin.Default()

	r.Use(middleware.Trace())
	r.Use(middleware.Recovery())

	r.Static("static", "front/static")
	r.LoadHTMLGlob("front/view/*")

	router.Router(r)

	// writeFile()

	if err := r.Run(fmt.Sprintf("%s:%d", global.Conf.Host, global.Conf.Port)); err != nil {
		panic(err)
	}
}

// func writeFile() {
// 	fp := "D:\\1_liuxiaobo\\testlog\\log.txt"

// 	// 写入10w行数据
// 	f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		fmt.Println("open file err=", err)
// 		return
// 	}
// 	defer f.Close()

// 	// 循环写入
// 	for i := 1; i <= 10000; i++ {
// 		s := fmt.Sprintf("%010d ## \r\n", i)
// 		f.WriteString(s)
// 	}
// }
