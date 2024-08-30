package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-co-op/gocron/v2"
)

// HelloHandler は "Hello, World!" を返すエンドポイント
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Hello, World!"}
	fmt.Println("Request received")
	json.NewEncoder(w).Encode(response)
}

func main() {
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func(a string, b time.Time) {
				// 現在時刻を出力
				fmt.Println("Hello, World!", a, b)
			},
			"hello",
			time.Now(),
		),
	)
	if err != nil {
		// handle error
	}

	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()

	// HTTP サーバーの設定
	http.HandleFunc("/hello", HelloHandler)

	// ゴルーチンで HTTP サーバーを実行
	go func() {
		fmt.Println("HTTP server is running on port 8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Error starting HTTP server:", err)
		}
	}()

	// メインゴルーチンをブロック
	select {}

	// // block until you are ready to shut down
	// time.Sleep(time.Minute)

	// // when you're done, shut it down
	// err = s.Shutdown()
	// if err != nil {
	// 	// handle error
	// }
}
