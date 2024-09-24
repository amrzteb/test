package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"sync"
	"time"
)

var (
	concurrency  = flag.Int("c", 5, "Number of concurrent workers")
	requestCount int64
	countMutex   sync.Mutex
)

func generateRandom8DigitNumber() string {
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

func main() {
	// 解析命令行参数
	flag.Parse()
	// 定义两个 curl 命令
	commands := []string{
		`curl -x global.711proxy.net:50000 -U "ttt001-zone-custom-region-US:ttt001" https://ipinfo.io`,
		`curl -x global.711proxy.net:50000 -U "ttt002-zone-custom-region-US:ttt002" https://ipinfo.io`,
		`curl -x global.711proxy.net:50000 -U "5cb4da80-zone-custom-region-US:469b0254965a134db924ffecd23919a6" https://ipinfo.io`,
		`curl -x global.711proxy.net:50000 -U "zzzzzz-zone-custom-region-US:zzzzzz" https://ipinfo.io`,
		`curl -x global.711proxy.net:50000 -U "qwerqwer123-zone-custom-region-US:qwerqwer" https://ipinfo.io`,

		`curl -x global.711proxy.net:10000 -U "ttt001-zone-custom-region-US:ttt001" https://ipinfo.io`,
		`curl -x global.711proxy.net:10000 -U "ttt002-zone-custom-region-US:ttt002" https://ipinfo.io`,
		`curl -x global.711proxy.net:10000 -U "5cb4da80-zone-custom-region-US:469b0254965a134db924ffecd23919a6" https://ipinfo.io`,
		`curl -x global.711proxy.net:10000 -U "zzzzzz-zone-custom-region-US:zzzzzz" https://ipinfo.io`,
		`curl -x global.711proxy.net:10000 -U "qwerqwer123-zone-custom-region-US:qwerqwer" https://ipinfo.io`,

		`curl -x global.711proxy.net:12000 -U "ttt001-zone-custom-region-US:ttt001" https://ipinfo.io`,
		`curl -x global.711proxy.net:12000 -U "ttt002-zone-custom-region-US:ttt002" https://ipinfo.io`,
		`curl -x global.711proxy.net:12000 -U "5cb4da80-zone-custom-region-US:469b0254965a134db924ffecd23919a6" https://ipinfo.io`,
		`curl -x global.711proxy.net:12000 -U "zzzzzz-zone-custom-region-US:zzzzzz" https://ipinfo.io`,
		`curl -x global.711proxy.net:12000 -U "qwerqwer123-zone-custom-region-US:qwerqwer" https://ipinfo.io`,

		`curl -x global.711proxy.net:20000 -U "ttt001-zone-custom-region-US:ttt001" https://ipinfo.io`,
		`curl -x global.711proxy.net:20000 -U "ttt002-zone-custom-region-US:ttt002" https://ipinfo.io`,
		`curl -x global.711proxy.net:20000 -U "5cb4da80-zone-custom-region-US:469b0254965a134db924ffecd23919a6" https://ipinfo.io`,
		`curl -x global.711proxy.net:20000 -U "zzzzzz-zone-custom-region-US:zzzzzz" https://ipinfo.io`,
		`curl -x global.711proxy.net:20000 -U "qwerqwer123-zone-custom-region-US:qwerqwer" https://ipinfo.io`,

		`curl -x global.711proxy.net:20000 -U "5cb4da80-zone-custom-session-` + generateRandom8DigitNumber() + `-sessTime-5-sessAuto-1:469b0254965a134db924ffecd23919a6" https://ipinfo.io`,
		`curl -x global.711proxy.net:20000 -U "ttt001-zone-custom-session-` + generateRandom8DigitNumber() + `-sessTime-5-sessAuto-1:ttt001" https://ipinfo.io`,
		`curl -x global.711proxy.net:20000 -U "ttt002-zone-custom-session-` + generateRandom8DigitNumber() + `-sessTime-5-sessAuto-1:ttt002" https://ipinfo.io`,
		`curl -x global.711proxy.net:20000 -U "ttt003-zone-custom-session-` + generateRandom8DigitNumber() + `-sessTime-5-sessAuto-1:ttt003" https://ipinfo.io`,
		`curl -x global.711proxy.net:20000 -U "ttt004-zone-custom-session-` + generateRandom8DigitNumber() + `-sessTime-5-sessAuto-1:ttt004" https://ipinfo.io`,
		`curl -x global.711proxy.net:20000 -U "zzzzzz-zone-custom-session-` + generateRandom8DigitNumber() + `-sessTime-5-sessAuto-1:zzzzzz" https://ipinfo.io`,
		`curl -x global.711proxy.net:20000 -U "qwerqwer123-zone-custom-session-` + generateRandom8DigitNumber() + `-sessTime-5-sessAuto-1:qwerqwer123" https://ipinfo.io`,
	}

	// 使用 wait group 等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 启动指定数量的并发任务
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// 随机执行命令的主循环
			for {
				// 随机选择一个命令
				cmd := commands[rand.Intn(len(commands))]

				// 执行命令
				err := executeCommand(cmd)
				if err != nil {
					log.Printf("Worker %d: Error executing command: %v\n", workerID, err)
				} else {
					// 增加请求计数
					countMutex.Lock()
					requestCount++
					countMutex.Unlock()
				}

				// 延时一段时间以避免过快执行
				time.Sleep(2 * time.Second)
			}
		}(i)
	}

	// 启动一个 goroutine 定时输出请求数量
	go func() {
		for {
			time.Sleep(10 * time.Second)
			countMutex.Lock()
			fmt.Printf("Total requests in the last 10 seconds: %d\n", requestCount)
			countMutex.Unlock()
		}
	}()

	// 等待所有 goroutine 完成 (实际上因为没有 cancel，这里不会执行)
	wg.Wait()
}

// 执行 shell 命令并打印输出
func executeCommand(cmd string) error {
	fmt.Printf("Executing command: %s\n", cmd)

	// 创建命令实例
	command := exec.Command("sh", "-c", cmd)

	// 获取命令输出
	output, err := command.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Printf("Command output: %s\n", string(output))
	return nil
}
