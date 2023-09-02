package vote

import "testing"

func TestClassService_CountVote(t *testing.T) {
	type fields struct {
		monitor          *Student
		assistantMonitor *Student
		students         []Student
		votes            []Vote
		candidates       []Candidate
	}
	tests := []struct {
		name                       string
		expectedMonitorID          int
		expectedAssistantMonitorID int
		fields                     fields
	}{
		{
			name:                       "支持票相当，反对票少的胜出",
			expectedMonitorID:          5,
			expectedAssistantMonitorID: 2,
			fields: fields{
				students: []Student{
					{id: 1, name: "s1"},
					{id: 2, name: "s2"},
					{id: 3, name: "s3"},
					{id: 4, name: "s4"},
					{id: 5, name: "s5"},
				},
				votes: []Vote{
					{supportID: 2, opposeID: 3},
					{supportID: 2, opposeID: 5},
					{supportID: 3, opposeID: 2},
					{supportID: 5, opposeID: 2},
					{supportID: 5, opposeID: 1},
				},
			},
		},
		{
			name:                       "支持票高的当选",
			expectedMonitorID:          2,
			expectedAssistantMonitorID: 3,
			fields: fields{
				students: []Student{
					{id: 1, name: "s1"},
					{id: 2, name: "s2"},
					{id: 3, name: "s3"},
				},
				votes: []Vote{
					{supportID: 2, opposeID: 2},
					{supportID: 2, opposeID: 2},
					{supportID: 3, opposeID: 2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &ClassService{
				monitor:          tt.fields.monitor,
				assistantMonitor: tt.fields.assistantMonitor,
				students:         tt.fields.students,
				votes:            tt.fields.votes,
				candidates:       tt.fields.candidates,
			}
			svc.CountVote()
			if svc.monitor.id != tt.expectedMonitorID {
				t.Error("预期选出班长：s5，实际选出班长:", svc.monitor.id)
			}
			if svc.assistantMonitor.id != tt.expectedAssistantMonitorID {
				t.Error("预期选出副班长：s2，实际选出副班长:", svc.assistantMonitor.id)
			}
		})
	}
}
