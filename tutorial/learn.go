package tutorial

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Notification struct {
	UserID  int64
	Message string
}

type SendResult struct {
	UserID  int64
	Success bool
}

// FCM 호출이라고 가정: 0.1초 걸리고, 10명 중 1명은 실패
func sendPush(n Notification) bool {
	time.Sleep(100 * time.Millisecond)
	return n.UserID%10 != 0
}

// 일꾼: jobs에서 꺼내서 보내고, 결과를 results에 넣음
func worker(jobs <-chan Notification, results chan<- SendResult, wg *sync.WaitGroup) {
	// 함수가 끝날 때 실행해줘, Java의 Finally
	defer wg.Done()
	for n := range jobs { // jobs가 close될 때까지 계속 꺼냄
		ok := sendPush(n)
		// <- 실행문에서 넣기/꺼내기
		// 요리사: 선반에 올려둠
		results <- SendResult{UserID: n.UserID, Success: ok}
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World! from GoLang<h1>")
}

func learnMain() {
	// val logging = mutableMapOf<String, MutableList<String>>()
	logging := make(map[string][]string)

	// logging.getOrPut("user1") { mutableListOf() }.addAll(listOf("10:00", "12:00"))
	// append는 "기존 리스트 + 새 값들 = 새 리스트"를 돌려주는 함수

	// 1단계, logging["user1"] 읽기 / map에서 “user1"을 찾기 -> 없는 key를 조회하면 에러 대신 빈 값을 주는 게 Go 규칙
	// 2단계, append([], "10:00", "12:00") 계산 / 빈 리스트 뒤에 두 값을 붙여서 ["10:00", "12:00"]
	// 3단계, logging["user1"] = ["10:00", "12:00"] / 그 결과를 map의 “user1” 자리에 저장해요.
	logging["user1"] = append(logging["user1"], "10:00", "12:00")
	logging["user2"] = append(logging["user2"], "09:30")

	// slice 를 활용하여 정렬할 것
	// for ((key, value) in logging) { ... }
	for id, times := range logging {
		fmt.Printf("사용자 %s의 접속기록 %v\n", id, times)
	}

	names := []string{"a", "b", "c"}

	// Go 에서는 안쓰는 변수가 있으면 안됨 "여기 값이 오긴 하는데 안 쓸 거야"
	for _, name := range names {
		fmt.Println(name)
	}

	// [GOROUTINE]
	// 상황: 배송 시작 푸시 알림 50건 보내기
	// 주문 50건이 배송을 시작해서 사용자에게 푸시 알림(FCM)을 보내야 해요. 알림 한 건 보내는 데 0.1초가 걸려요.

	// 왜 굳이 channel로 일을 나누지? -> "공유하지 말고 전달하라"
	// 5명이 같은 next 변수를 동시에 고치면 같은 주문을 두 번 처리하거나 건너뛰는 일이 생김,막으려면 자물쇠(Mutex)를 걸어야 하고 코드가 복잡
	const workerCount = 5

	// 작업 대기열 (버퍼 100) -> 더쉽게:“Notification을 담을 수 있는 칸이 100개인 꽂이(대기열)를 만들어라”
	// chan Notification: 이 통로로는 Notification만 오갈 수 있다
	// 100: 꺼내 가는 사람이 없어도 최대 100장까지 꽂아둘 수 있다
	jobs := make(chan Notification, 100)
	// 결과 대기열 -> 더쉽게 -> 완성음식 선반
	results := make(chan SendResult, 100)

	// 일꾼 5명 출근시키기
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// 생산자, 작업을 넣고 다 넣으면 close
	go func() {
		for uid := int64(1); uid <= 50; uid++ {
			jobs <- Notification{UserID: uid, Message: "주문하신 상품의 배송이 시작되었습니다"}
		}
		close(jobs)
	}()

	// 일꾼들이 퇴근하면 results를 닫기
	go func() {
		wg.Wait()
		close(results)
	}()

	success, fail := 0, 0

	// 카운터: 하나씩 가져가서 셈
	for r := range results {
		if r.Success {
			success++
		} else {
			fail++
		}
	}
	fmt.Printf("성공 %d건, 실패 %d건\n", success, fail)

	http.HandleFunc("/", helloHandler)
	fmt.Println("9090포트에서 Go 서버가 실행되었습니다.")

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println(err)
	}
}
