package lottotest

import (
	"fmt"
	"hash/fnv"
)

type Identifier struct {
	MothershipVersion string `json:"mothershipVersion"`
	IncludeOSVersion  string `json:"includeosVersion"`
	Environment       string `json:"environment"`
	TestName          string `json:"name"`
}

func (i Identifier) GetID() uint32 {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%s%s%s%s", i.MothershipVersion, i.IncludeOSVersion, i.Environment, i.TestName)))
	return h.Sum32()
}

func (i Identifier) matches(j Identifier) bool {
	if j.MothershipVersion != "" {
		if i.MothershipVersion != j.MothershipVersion {
			return false
		}
	}
	if j.IncludeOSVersion != "" {
		if i.IncludeOSVersion != j.IncludeOSVersion {
			return false
		}
	}
	if j.Environment != "" {
		if i.Environment != j.Environment {
			return false
		}
	}
	if j.TestName != "" {
		if i.TestName != j.TestName {
			return false
		}
	}
	return true
}
