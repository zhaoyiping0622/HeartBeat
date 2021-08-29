package main

import "encoding/json"

type Contact interface {
	json.Marshaler
	json.Unmarshaler
	Send(message string) error
}
