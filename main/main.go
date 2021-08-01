package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"os"
	"os/signal"
)

func StartHttpServer(srv *http.Server) error {
	http.HandleFunc("/geek", HelloServer2)
	fmt.Println("http server start")
	err := srv.ListenAndServe()
	return err
}

func HelloServer2(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "home work!\n")
}

func main() {
	ctx := context.Background()
	// 定义 withCancel -> cancel() 方法 去取消下游的 Context
	ctx, cancel := context.WithCancel(ctx)
	// 使用 errgroup 进行 goroutine 取消
	group, errCtx := errgroup.WithContext(ctx)
	//http server
	srv := &http.Server{Addr: ":9090"}

	group.Go(func() error {
		return StartHttpServer(srv)
	})

	chanel := make(chan os.Signal, 1)
	signal.Notify(chanel)

	group.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				fmt.Println("http server stop")
				return srv.Shutdown(errCtx) // 关闭 http server
			case <-chanel: // 因为 kill -9 或其他而终止
				cancel()
			}
		}
	})

	if err := group.Wait(); err != nil {
		fmt.Println("group error: ", err)
	}
	fmt.Println("all group done!")

}
