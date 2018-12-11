package lottotest

import (
	"time"
)

type TestResult struct {
	Identifier
	Duration time.Duration `json:"duration"`
	Sent     int           `json:"sent"`
	Received int           `json:"received"`
	Success  bool          `json:"success"`
}

type TestCollection struct {
	ID                Identifier
	Results           []TestResult
	SuccessPercentage float32
	NumRuns           int
	Passed            int
	Failed            int
	TotalSent         int
	TotaltReceived    int
}

func NewTestCollection(firstResult TestResult) *TestCollection {
	tc := &TestCollection{
		ID:      firstResult.Identifier,
		Results: []TestResult{},
	}
	tc.AddResultToTestCollection(firstResult)
	return tc
}

func (tc *TestCollection) AddResultToTestCollection(result TestResult) {
	tc.NumRuns++
	switch result.Success {
	case true:
		tc.Passed++
	default:
		tc.Failed++
	}
	tc.SuccessPercentage = float32(tc.Passed) / float32(tc.NumRuns) * 100
	tc.TotalSent += result.Sent
	tc.TotaltReceived += result.Received
	tc.Results = append(tc.Results, result)
}
