package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type testStatus int

const (
	StatusUnknown testStatus = iota
	StatusSkipped
	StatusFail
	StatusPass
	StatusWarning
)

var (
	statusToString = map[testStatus]string{
		StatusUnknown: "Unknown",
		StatusSkipped: "Skipped",
		StatusFail:    "Fail",
		StatusPass:    "Pass",
		StatusWarning: "Warning",
	}

	stringToStatus = map[string]testStatus{
		"unknown": StatusUnknown,
		"skipped": StatusSkipped,
		"skip":    StatusSkipped,
		"fail":    StatusFail,
		"failed":  StatusFail,
		"pass":    StatusPass,
		"passed":  StatusPass,
		"warning": StatusWarning,
		"warn":    StatusWarning,
	}
)

func (s testStatus) String() string {
	if v, ok := statusToString[s]; ok {
		return v
	}
	return "Invalid"
}

func parseStatus(s string) (testStatus, bool) {
	v, ok := stringToStatus[strings.ToLower(strings.TrimSpace(s))]
	return v, ok
}

func (s testStatus) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "Invalid" {
		return nil, fmt.Errorf("cannot marshal invalid testStatus value: %d", int(s))
	}
	return json.Marshal(str)
}

// UnmarshalJSON accepts either a string ("Pass") or a number (3).
func (s *testStatus) UnmarshalJSON(data []byte) error {
	// Try string first.
	var asStr string
	if err := json.Unmarshal(data, &asStr); err == nil {
		if v, ok := parseStatus(asStr); ok {
			*s = v
			return nil
		}
		return fmt.Errorf("invalid status string %q", asStr)
	}

	var asNum int
	if err := json.Unmarshal(data, &asNum); err == nil {
		v := testStatus(asNum)
		if _, ok := statusToString[v]; ok {
			*s = v
			return nil
		}
		return fmt.Errorf("invalid status number %d", asNum)
	}
	// Neither string nor number.
	return fmt.Errorf("status must be string or number, got: %s", string(bytes.TrimSpace(data)))
}

func (s testStatus) MarshalText() ([]byte, error) {
	str := s.String()
	if str == "Invalid" {
		return nil, fmt.Errorf("cannot marshal invalid testStatus value: %d", int(s))
	}
	return []byte(str), nil
}

func (s *testStatus) UnmarshalText(text []byte) error {
	if v, ok := parseStatus(string(text)); ok {
		*s = v
		return nil
	}
	return fmt.Errorf("invalid status text %q", string(text))
}
