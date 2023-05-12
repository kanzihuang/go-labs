package primeservice

import (
	"go-labs/utils"
	"runtime"
	"sync"
)

const MinPrime Prime = 2

type Prime = int

type Task struct {
	begin       int
	end         int
	priorPrimes []Prime
	taskPrimes  []Prime
	finished    bool
}

func (task *Task) IsEmpty() bool {
	return task.begin >= task.end
}

type QueryData struct {
	number     int
	resultChan chan bool
}

type Service struct {
	calcNum        int
	primes         []Prime
	countPerTask   int
	queryChan      chan QueryData
	resultChan     chan bool
	taskChan       chan *Task
	taskFinishChan chan bool
}

var svc *Service
var once sync.Once

func NewService() *Service {
	calcNum := runtime.GOMAXPROCS(8)
	once.Do(func() {
		svc = &Service{
			calcNum:        calcNum,
			primes:         []Prime{MinPrime},
			countPerTask:   100,
			queryChan:      make(chan QueryData, calcNum),
			resultChan:     make(chan bool, calcNum),
			taskChan:       make(chan *Task, calcNum),
			taskFinishChan: make(chan bool, calcNum),
		}
		for i := 0; i < calcNum; i++ {
			go svc.worker()
		}
		go svc.run()
	})
	return svc
}

func (svc *Service) worker() {
	for {
		task := <-svc.taskChan
		task.runTask()
		svc.taskFinishChan <- true
	}
}

func (svc *Service) run() {
OuterLoop:
	for {
		var lastTask, nextTask *Task
		var tasks []*Task
		// todo: 多个查询并发时，实现计算量小的查询先返回结果
		select {
		case queryData := <-svc.queryChan:
			lastIndex := len(svc.primes) - 1
			if queryData.number <= svc.primes[lastIndex] {
				queryData.resultChan <- svc.findPrime(queryData.number) >= 0
				break
			}
			for finished := false; finished == false; {
				if nextTask == nil {
					nextTask = svc.getNextTask(queryData.number, lastTask)
				}
				select {
				case svc.taskChan <- nextTask:
					if !nextTask.IsEmpty() {
						tasks = append(tasks, nextTask)
					}
					lastTask = nextTask
					nextTask = nil
				case <-svc.taskFinishChan:
					var i int
					var task *Task
					for i, task = range tasks {
						if task.finished {
							svc.primes = append(svc.primes, task.taskPrimes...)
							if task.end > queryData.number {
								finished = true
								queryData.resultChan <- svc.findPrime(queryData.number) >= 0
								continue OuterLoop
							}
						} else {
							break
						}
					}
					if i < len(tasks) {
						tasks = tasks[i:]
					} else {
						tasks = nil
					}
				}
			}
		}
	}
}

func (svc *Service) IsPrime(number int) bool {
	resultChan := make(chan bool)
	defer close(resultChan)
	svc.queryChan <- QueryData{number: number, resultChan: resultChan}
	return <-resultChan
}

func (svc *Service) findPrime(number int) int {
	if number == MinPrime {
		return 0
	} else if number%2 == 0 {
		return -1
	}
	firstIndex, lastIndex := 0, len(svc.primes)-1
	for {
		if number < svc.primes[firstIndex] || number > svc.primes[lastIndex] {
			return -1
		}
		middleIndex := (firstIndex + lastIndex) / 2
		middle := svc.primes[middleIndex]
		if number == middle {
			return middleIndex
		} else if number < middle {
			lastIndex = middleIndex
		} else if middleIndex > firstIndex {
			firstIndex = middleIndex
		} else {
			firstIndex = middleIndex + 1
		}
	}
}

func (svc *Service) GetPrimes(a, b int) []Prime {
	var primes []Prime

	if a == MinPrime {
		primes = append(primes, MinPrime)
	}
	if a%2 == 0 {
		a++
	}
	for i := a; i <= b; i += 2 {
		if svc.IsPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

func (svc *Service) getNextTask(number int, task *Task) *Task {
	var start, end int
	lastPrime := svc.primes[len(svc.primes)-1]
	if task == nil {
		start = lastPrime + 1
	} else {
		start = task.end
	}
	end = utils.Min(lastPrime*lastPrime+1, start+svc.countPerTask, number+1)
	return &Task{begin: start, end: end, priorPrimes: svc.primes}
}

func (task *Task) runTask() {
	for i := task.begin; i < task.end; i++ {
		if task.checkPrime(i) {
			task.taskPrimes = append(task.taskPrimes, i)
		}
	}
	task.finished = true
}

func (task *Task) checkPrime(number int) bool {
	lastPrime := task.priorPrimes[len(task.priorPrimes)-1]
	if number > lastPrime*lastPrime {
		panic("已算出的质数数量不足，无法计算新的质数，应减小任务中的待计算的整数数量")
	}
	for _, prime := range task.priorPrimes {
		if prime*prime > number {
			break
		}
		if number%prime == 0 {
			return false
		}
	}
	return true
}
