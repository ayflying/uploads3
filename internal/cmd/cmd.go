package cmd

import (
	"context"
	"time"
	"uploads3/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// helpDescription 定义P2P命令的详细帮助信息
const helpDescription = `
S3上传工具使用帮助:
-p,path  本地文件夹路径
-u,upload_path  S3上传根路径
-w,worker  并发数，默认10，最大50
`

var (
	path        string
	uploadPath  string
	maxCount    int
	uploadCount int = 0
	// 管道长度
	workerCount = 10 // 固定并发数100

	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		// Description 提供命令的详细描述和使用帮助
		Description: helpDescription,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			parser, err = gcmd.Parse(g.MapStrBool{
				"p,path":        true,
				"u,upload_path": true,
				"w,worker":      true,
			})

			//s := g.Server()
			//s.Group("/", func(group *ghttp.RouterGroup) {
			//	group.Middleware(ghttp.MiddlewareHandlerResponse)
			//	group.Bind(
			//		hello.NewV1(),
			//	)
			//})
			//s.Run()
			workerCount = parser.GetOpt("worker", workerCount).Int()
			if workerCount > 50 {
				workerCount = 50
			}
			path = parser.GetOpt("path").String()
			uploadPath = parser.GetOpt("upload_path").String()
			if path == "" || uploadPath == "" {
				g.Log().Errorf(ctx, "path 或 upload_path 为空")
				return
			}

			S3(path)

			return nil
		},
	}
)

var UploadTask = make(chan string, workerCount*5)

func S3(path string) {
	list, _ := gfile.ScanDirFile(path, "*", true)
	maxCount = len(list)
	g.Log().Debugf(gctx.New(), "当前需要处理的文件数量：%v", len(list))
	go func() {
		for _, v := range list {
			UploadTask <- v
		}
	}()
	time.Sleep(1 * time.Second)
	startWorkers()
}

// 启动100个worker，持续处理任务
func startWorkers() {

	// 启动100个worker
	for i := 0; i < workerCount; i++ {
		ctx := gctx.New()
		go func() {
			// 持续从管道取任务，直到管道关闭且所有任务处理完毕
			for {
				select {
				case filename := <-UploadTask:
					//执行上传任务
					uploadToS3(ctx, filename)
				case <-ctx.Done():
					// 上下文取消时，退出循环
					return
				}
			}
		}()
	}

	// 等待所有任务处理完毕
	for {
		if len(UploadTask) == 0 {
			return
		}
	}

}

func uploadToS3(ctx context.Context, filename string) {
	//todo 实现上传到S3的逻辑
	uploadCount++

	filepath := gstr.Replace(filename, path, uploadPath)
	filepath = gstr.Replace(filepath, "\\", "/")
	g.Log().Debugf(ctx, "(%d,%d)上传到s3：%v", uploadCount, maxCount, filepath)
	//time.Sleep(grand.D(10*time.Millisecond, time.Second))
	f, _ := gfile.Open(filename)
	service.S3().PutObject(ctx, f, filepath)
	return
}
