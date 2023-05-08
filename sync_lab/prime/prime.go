package prime

type Prime = int

const MinPrime Prime = 2

type Service struct {
	calcNum int
	primes  []Prime
}

func CreateService(calcNum int) *Service {
	return &Service{
		calcNum: calcNum,
		primes:  []Prime{MinPrime, 3},
	}
}

func (svc *Service) IsPrime(number int) bool {
	svc.appendPrimes(number)
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

func (svc *Service) appendPrimes(number int) {
	last := svc.primes[len(svc.primes)-1]
	if number <= last {
		return
	} else if number > last {
		for i := last + 2; i <= number; i += 2 {
			svc.appendIfPrime(i)
		}
	}
}

func (svc *Service) appendIfPrime(number int) bool {
	for _, prime := range svc.primes {
		if number < prime*prime {
			break
		}
		if number%prime == 0 {
			return false
		}
	}
	svc.primes = append(svc.primes, number)
	return true
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
