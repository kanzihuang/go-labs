package primeservice

const MinPrime Prime = 2

type Prime = int
type PrimeState = int

const PrimeYes PrimeState = 1
const PrimePending PrimeState = 0
const PrimeNo PrimeState = -1

type Task struct {
	base    int
	numbers []PrimeState
}

func (task *Task) getNumber(index int) int {
	return task.base + index
}

type Service struct {
	calcNum      int
	primes       []Prime
	numbers      []PrimeState
	countPerTask int
}

func CreateService(calcNum int) *Service {
	return &Service{
		calcNum:      calcNum,
		primes:       []Prime{MinPrime},
		numbers:      []PrimeState{},
		countPerTask: 100,
	}
}

func (svc *Service) IsPrime(number int) bool {
	for {
		task := svc.getTask(number)
		if len(task.numbers) == 0 {
			break
		}
		svc.runTask(task)
		svc.migratePrimes()
	}
	return svc.findPrime(number) >= 0
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

func (svc *Service) getTask(number int) *Task {
	lastPrime := svc.primes[len(svc.primes)-1]
	base := lastPrime + len(svc.numbers) + 1
	startIndex := len(svc.numbers)
	for num := base; num <= number && num <= lastPrime*lastPrime; num++ {
		svc.numbers = append(svc.numbers, PrimePending)
	}
	return &Task{base: base, numbers: svc.numbers[startIndex:]}
}

func (svc *Service) runTask(task *Task) {
	for i := range task.numbers {
		task.numbers[i] = svc.checkPrime(task.getNumber(i))
	}
}

func (svc *Service) checkPrime(number int) PrimeState {
	if number > svc.primes[len(svc.primes)-1]*svc.primes[len(svc.primes)-1] {
		panic("已算出的质数数量不足，无法计算新的质数，应减小任务中的待计算的整数数量")
	}
	for _, prime := range svc.primes {
		if number < prime*prime {
			break
		}
		if number%prime == 0 {
			return PrimeNo
		}
	}
	return PrimeYes
}

func (svc *Service) migratePrimes() {
	lastPrime := svc.primes[len(svc.primes)-1]
	for i, state := range svc.numbers {
		if state == PrimeYes {
			svc.primes = append(svc.primes, lastPrime+i+1)
		} else if state == PrimePending {
			break
		}
	}
	svc.numbers = svc.numbers[svc.primes[len(svc.primes)-1]-lastPrime:]
}
