package lottotest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type TestResult struct {
	Identifier
	Duration time.Duration `json:"duration"`
	Sent     int           `json:"sent"`
	Received int           `json:"received"`
	Success  bool          `json:"success"`
}

func (r TestResult) SaveToDisk(filename string) error {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error marshaling json: %v", err)
	}
	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
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
	startTime         time.Time
	idCounter         int
}

func NewTestCollection(i Identifier) *TestCollection {
	tc := &TestCollection{
		ID:        i,
		Results:   []TestResult{},
		startTime: time.Now(),
	}
	return tc
}

func (tc *TestCollection) AddResultToTestCollection(result TestResult) error {
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
	fileName := path.Join("results", fmt.Sprint(result.Identifier.GetID()), fmt.Sprintf("%s.json", tc.nextID()))
	if err := result.SaveToDisk(fileName); err != nil {
		return fmt.Errorf("error saving to disk: %v", err)
	}
	return nil
}

func (tc *TestCollection) nextID() string {
	tc.idCounter++
	return fmt.Sprintf("%d_%d", tc.startTime.Unix(), tc.idCounter)
}
