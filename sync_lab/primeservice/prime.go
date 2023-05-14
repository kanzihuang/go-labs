package primeservice

import (
	"fmt"
	"go-labs/utils"
	"runtime"
	"sync"
)

const MinPrime Prime = 2

type Prime = int
type PrimeSection struct {
	begin  int
	end    int
	primes []Prime
}
type Task struct {
	PrimeSection
	priorPrimes []Prime
	finished    bool
}

func (s *PrimeSection) Append(right *PrimeSection) {
	if s.end != right.begin {
		panic(fmt.Sprintf("添加的质数区间不连续：(%d - %d), (%d - %d)\n",
			s.begin, s.end, right.begin, right.end))
	}
	s.end = right.end
	s.primes = append(s.primes, right.primes...)
}

func (s *PrimeSection) InSection(number int) bool {
	return number >= s.begin && number < s.end
}

func (task *Task) IsEmpty() bool {
	return task.begin >= task.end
}

type Query struct {
	number     int
	resultChan chan<- bool
	finished   bool
}

type Service struct {
	primes         PrimeSection
	queries        []*Query
	tasks          []*Task
	queryChan      chan *Query
	taskChan       chan *Task
	taskFinishChan chan *struct{}
}

const countPerTask int = 100

var calcNum int = runtime.GOMAXPROCS(0)

// var calcNum int = 2
var svc *Service
var once sync.Once

func NewService() *Service {
	once.Do(func() {
		svc = &Service{
			primes:         PrimeSection{begin: MinPrime, end: MinPrime + 1, primes: []Prime{MinPrime}},
			queryChan:      make(chan *Query, calcNum),
			taskChan:       make(chan *Task, calcNum),
			taskFinishChan: make(chan *struct{}, calcNum),
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
		svc.taskFinishChan <- &struct{}{}
	}
}
func recvFromChannel[T any](channel <-chan *T, waiting bool) *T {
	if waiting {
		select {
		case data := <-channel:
			return data
		}
	} else {
		select {
		case data := <-channel:
			return data
		default:
			return nil
		}
	}
}

func (svc *Service) run() {
	var lastTask, nextTask *Task
	for {
		// 接收新查询（如查询队列为空，等待新的查询）
		if q := recvFromChannel(svc.queryChan, svc.queries == nil); q != nil {
			svc.queries = append(svc.queries, q)
		}

		// 获取新的任务（如条件不满足，返回空任务）
		if nextTask == nil || nextTask.IsEmpty() {
			nextTask = svc.primes.getNextTask(lastTask)
		}

		// 接收任务完成情况（如后续任务为空，等待存量任务完成）
		if nil != recvFromChannel(svc.taskFinishChan, nextTask.IsEmpty()) {
			svc.checkTasksFinished()
			svc.checkQueriesFinished()
		}

		// 发送任务（该任务非空，如 channel 已满，直接返回 false）
		if svc.sendTask(nextTask) {
			lastTask, nextTask = nextTask, nil
		}
	}
}

func (svc *Service) sendTask(nextTask *Task) bool {
	select {
	case svc.taskChan <- nextTask:
		svc.tasks = append(svc.tasks, nextTask)
		return true
	default:
		return false
	}
}

func (svc *Service) checkTasksFinished() {
	var task *Task
	finishedCount := 0
	for _, task = range svc.tasks {
		if task.finished {
			svc.primes.Append(&task.PrimeSection)
			finishedCount++
		} else {
			break
		}
	}
	if finishedCount < len(svc.tasks) {
		svc.tasks = svc.tasks[finishedCount:]
	} else {
		svc.tasks = nil
	}
}

func (svc *Service) checkQueriesFinished() {
	finishedCount := 0
	for _, query := range svc.queries {
		if !query.finished && svc.primes.InSection(query.number) {
			query.finished = true
			query.resultChan <- svc.primes.findPrime(query.number) >= 0
		}
		if query.finished {
			finishedCount++
		}
	}
	if finishedCount == len(svc.queries) {
		svc.queries = nil
	}
}

func (svc *Service) IsPrime(number int) bool {
	if number < MinPrime {
		return false
	}
	resultChan := make(chan bool)
	defer close(resultChan)
	svc.queryChan <- &Query{number: number, resultChan: resultChan}
	return <-resultChan
}

func (s *PrimeSection) findPrime(number int) int {
	if !s.InSection(number) {
		return -1
	}
	firstIndex, lastIndex := 0, len(s.primes)-1
	for {
		if number < s.primes[firstIndex] || number > s.primes[lastIndex] {
			return -1
		}
		middleIndex := (firstIndex + lastIndex) / 2
		middle := s.primes[middleIndex]
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

func (s *PrimeSection) getNextTask(task *Task) *Task {
	var start, end int
	if task == nil {
		start = s.end
	} else {
		start = task.end
	}
	lastPrime := s.primes[len(s.primes)-1]
	end = utils.Min(lastPrime*lastPrime+1, start+countPerTask)
	return &Task{PrimeSection: PrimeSection{begin: start, end: end}, priorPrimes: s.primes}
}

func (task *Task) runTask() {
	for i := task.begin; i < task.end; i++ {
		if task.checkPrime(i) {
			task.primes = append(task.primes, i)
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

func FilterPrime(n int) bool {
	if n < 2 {
		return false
	}
	isnp := make([]bool, n+1) // is not prime: 不是素数
	isnp[0], isnp[1] = true, true
	for i := 2; i <= n; i++ {
		for j := 2; i*j <= n; j++ {
			if isnp[j] == false {
				isnp[i*j] = true
				if i%j == 0 {
					break
				}
			}
		}
	}
	return !isnp[n]
}

func FilterPrime2(num int) bool {
	if num < 2 {
		return false
	}
	st := make([]bool, num+1)
	primes := make([]int, num/2)
	//1不是质数也不是合数
	for i := 2; i <= num; i++ {
		if !st[i] {
			primes = append(primes, i) //没有被筛去,说明是质数
		}
		for j := 0; j < len(primes) && i*primes[j] <= num; j++ {
			st[i*primes[j]] = true //筛去合数
			if i%primes[j] == 0 {
				break //核心操作,保证了O(n)的复杂度
			}
		}

	}
	return !st[num]
}
