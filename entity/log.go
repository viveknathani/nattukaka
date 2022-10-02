package entity

import "encoding/json"

type Log struct {
	Level     string          `json:"string"`
	Ts        float64         `json:"ts"`
	Path      string          `json:"path"`
	Message   string          `json:"message"`
	RequestID string          `json:"requestID"`
	Method    string          `json:"method"`
	IP        string          `json:"ip"`
	Info      json.RawMessage `json:"info"`
}
