package lottotest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
)

type ResultStore struct {
	Tests          map[uint32]*TestCollection
	ResultIDs      []Identifier
	MVersions      []string
	IosVersions    []string
	Environments   []string
	TestNames      []string
	PreviousFilter Identifier
}

func NewResultStore() *ResultStore {
	rs := &ResultStore{}
	rs.Tests = map[uint32]*TestCollection{}
	rs.ResultIDs = []Identifier{}
	rs.MVersions = []string{}
	rs.IosVersions = []string{}
	rs.Environments = []string{}
	rs.TestNames = []string{}
	rs.PreviousFilter = Identifier{}
	return rs
}

var mVersions = []string{"v1", "v2", "v3", "v4", "v5"}
var iosVersions = []string{"v1", "v2", "v3", "v4", "v5"}
var environments = []string{"vcloud", "fusion", "openstack"}
var testNames = []string{"Markshimmer", "Followerscarlet", "Mindrampant", "Walkerstump", "Herondawn"}

func (rs *ResultStore) AddRandom(num int) {
	r := rand.New(rand.NewSource(1234))
	for i := 0; i < num; i++ {
		x := TestResult{}
		x.TestName = testNames[r.Intn(len(testNames))]
		x.MothershipVersion = mVersions[r.Intn(len(mVersions))]
		x.IncludeOSVersion = iosVersions[r.Intn(len(iosVersions))]
		x.Environment = environments[r.Intn(len(environments))]
		x.Duration = time.Duration(r.Int63n(100000000000))
		x.Sent = r.Int()
		x.Received = r.Int()
		x.Success = randomdata.Boolean()
		rs.AddResult(x, true)
	}
}

func (rs *ResultStore) FilterIDs(i Identifier) []Identifier {
	output := []Identifier{}
	for _, existingID := range rs.ResultIDs {
		if existingID.matches(i) {
			output = append(output, existingID)
		}
	}
	rs.PreviousFilter = i
	return output
}

func (rs *ResultStore) AddIdentifier(i Identifier) {
	if !contains2(rs.ResultIDs, i) {
		rs.ResultIDs = append(rs.ResultIDs, i)
		// Add all filter fields to ResultStore
		if !contains2(rs.MVersions, i.MothershipVersion) {
			rs.MVersions = append(rs.MVersions, i.MothershipVersion)
		}
		if !contains2(rs.IosVersions, i.IncludeOSVersion) {
			rs.IosVersions = append(rs.IosVersions, i.IncludeOSVersion)
		}
		if !contains2(rs.Environments, i.Environment) {
			rs.Environments = append(rs.Environments, i.Environment)
		}
		if !contains2(rs.TestNames, i.TestName) {
			fmt.Println("Adding test: ", i.TestName)
			rs.TestNames = append(rs.TestNames, i.TestName)
		}
	}
}

func (rs *ResultStore) AddResult(newResult TestResult, save bool) {
	rs.AddIdentifier(newResult.Identifier)
	id := newResult.Identifier.GetID()
	if _, ok := rs.Tests[id]; !ok {
		rs.Tests[id] = NewTestCollection(newResult.Identifier)
	}
	if save {
		if err := rs.Tests[id].AddResultToTestCollectionAndSaveToFile(newResult); err != nil {
			fmt.Printf("Error adding test: %v", err)
		}
	} else {
		rs.Tests[id].AddResultToTestCollection(newResult)
	}

}

func (rs *ResultStore) LoadFiles(files string) error {
	if err := filepath.Walk(files, func(path string, info os.FileInfo, err error) error {
		if os.IsNotExist(err) {
			fmt.Printf("path %s does not exist\n", path)
			return nil
		}
		if !info.IsDir() {
			fmt.Printf("reading file: %s\n", path)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("error reading json file %s: %v", path, err)
				return nil
			}
			result := TestResult{}
			if err = json.Unmarshal(data, &result); err != nil {
				fmt.Printf("error unmarshaling file %s: %v", path, err)
			}
			rs.AddResult(result, false)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("error reading dir: %v", err)
	}
	return nil
}

func contains2(slice interface{}, e interface{}) bool {
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)
		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Interface() == e {
				return true
			}
		}
	}
	return false
}
