package primeservice

import (
	"math"
	"sync"
)

const MinPrime Prime = 2

type Prime = int

type Task struct {
	begin  int
	end    int
	primes []Prime
}

func (task *Task) IsEmpty() bool {
	return task.begin >= task.end
}

func (task *Task) getNumber(index int) int {
	return task.begin + index
}

type Service struct {
	calcNum          int
	primes           []Prime
	countPerTask     int
	runTaskWaitGroup sync.WaitGroup
}

func CreateService(calcNum int) *Service {
	return &Service{
		calcNum:      calcNum,
		primes:       []Prime{MinPrime},
		countPerTask: 100,
	}
}

func (svc *Service) IsPrime(number int) bool {
	for {
		tasks := svc.getTasks()
		svc.runTasks(tasks)
		svc.migratePrimes(tasks)
		if len(tasks) == 0 || number < tasks[len(tasks)-1].end {
			break
		}
	}
	return svc.findPrime(number) >= 0
}

func (svc *Service) getTasks() []*Task {
	var lastTask *Task = nil
	tasks := make([]*Task, 0, svc.calcNum)
	for i := 0; i < svc.calcNum; i++ {
		task := svc.getNextTask(lastTask)
		if task.IsEmpty() {
			break
		}
		tasks = append(tasks, task)
		lastTask = task
	}
	return tasks
}

func (svc *Service) runTasks(tasks []*Task) {
	for _, task := range tasks {
		svc.runTaskWaitGroup.Add(1)
		go svc.runTask(task)
	}
	svc.runTaskWaitGroup.Wait()
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

func Max(numbers ...int) int {
	result := math.MinInt
	for _, num := range numbers {
		if num > result {
			result = num
		}
	}
	return result
}
func Min(numbers ...int) int {
	result := math.MaxInt
	for _, num := range numbers {
		if num < result {
			result = num
		}
	}
	return result
}

func (svc *Service) getNextTask(task *Task) *Task {
	var start, end int
	lastPrime := svc.primes[len(svc.primes)-1]
	if task == nil {
		start = lastPrime + 1
	} else {
		start = task.end
	}
	end = Min(lastPrime*lastPrime, start+svc.countPerTask)
	return &Task{begin: start, end: end}
}

func (svc *Service) runTask(task *Task) {
	defer svc.runTaskWaitGroup.Done()
	for i := task.begin; i < task.end; i++ {
		if svc.checkPrime(i) {
			task.primes = append(task.primes, i)
		}
	}
}

func (svc *Service) checkPrime(number int) bool {
	if number > svc.primes[len(svc.primes)-1]*svc.primes[len(svc.primes)-1] {
		panic("已算出的质数数量不足，无法计算新的质数，应减小任务中的待计算的整数数量")
	}
	for _, prime := range svc.primes {
		if prime*prime > number {
			break
		}
		if number%prime == 0 {
			return false
		}
	}
	return true
}

func (svc *Service) migratePrimes(tasks []*Task) {
	for _, task := range tasks {
		svc.primes = append(svc.primes, task.primes...)
	}
}
