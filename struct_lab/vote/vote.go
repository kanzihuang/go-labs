package vote

import (
	"math/rand"
	"sort"
	"time"
)

type Student struct {
	id   int
	name string
}

type Vote struct {
	supportID int
	opposeID  int
}

type Candidate struct {
	*Student
	supportCount int
	opposeCount  int
}

type ClassService struct {
	monitor          *Student
	assistantMonitor *Student
	students         []Student
	votes            []Vote
	candidates       []Candidate
}

func (svc *ClassService) Vote() {
	svc.votes = make([]Vote, len(svc.students))
	for i, voter := range svc.students {
		svc.votes[i] = svc.GetVote(voter, svc.students)
	}
}

func (svc *ClassService) CountVote() {
	svc.candidates = make([]Candidate, len(svc.students))
	for i := range svc.students {
		svc.candidates[i].Student = &svc.students[i]
	}

	for _, vote := range svc.votes {
		for _, student := range svc.students {
			if vote.supportID == student.id {
				svc.CountSupportVote(student)
			}
			if vote.opposeID == student.id {
				svc.CountOpposeVote(student)
			}
		}
	}

	sort.SliceStable(svc.candidates, func(i, j int) bool {
		left := svc.candidates[i]
		right := svc.candidates[j]
		switch {
		case left.supportCount > right.supportCount:
			return true
		case left.supportCount == right.supportCount && left.opposeCount < right.opposeCount:
			return true
		default:
			return false
		}
	})

	svc.monitor = svc.candidates[0].Student
	svc.assistantMonitor = svc.candidates[1].Student
}

func (svc *ClassService) GetVote(voter Student, candidates []Student) (vote Vote) {
	rand.Seed(int64(time.Now().Nanosecond()))
	supportINX := rand.Intn(len(candidates))
	opposeINX := rand.Intn(len(candidates))
	vote.supportID = candidates[supportINX].id
	vote.opposeID = candidates[opposeINX].id
	return
}

func (svc *ClassService) CountSupportVote(student Student) {
	candidate := svc.getCandidate(student)
	candidate.supportCount++
}
func (svc *ClassService) CountOpposeVote(student Student) {
	candidate := svc.getCandidate(student)
	candidate.opposeCount++
}

func (svc *ClassService) getCandidate(student Student) *Candidate {
	for i := range svc.candidates {
		if svc.candidates[i].id == student.id {
			return &svc.candidates[i]
		}
	}
	return nil
}
