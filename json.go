package main

import "time"

// DeviceContent JSON Struct for marshalling
type DeviceContent struct {
	Manufacturer  string        `json:"manufacturer"`
	DeviceName    string        `json:"deviceName"`
	Model         string        `json:"model"`
	UniqueID      string        `json:"uniqueID"`
	DeviceLocale  string        `json:"deviceLocale"`
	DeviceCountry string        `json:"deviceCountry"`
	UserAgent     string        `json:"userAgent"`
	Measurements  []Measurement `json:"measurements"`
}

// Measurement JSON struct for marshalling
type Measurement struct {
	UUID     string    `json:"uuid"`
	TrialID  string    `json:"trialId"`
	Trial    int       `json:"trial"`
	Subject  string    `json:"subject"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration int       `json:"duration"`
	Target   int       `json:"target"`
	Response int       `json:"response"`
	Correct  bool      `json:"correct"`
}
