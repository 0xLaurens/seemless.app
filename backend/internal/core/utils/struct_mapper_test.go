package utils

import (
	"encoding/json"
	"testing"
)

type MessageMock struct {
	Type string
	From string
}

func TestValidJsonMapsToObject(t *testing.T) {
	msg, _ := json.Marshal(MessageMock{
		Type: "PeerJoin",
		From: "John",
	})

	req := &MessageMock{}

	err := MapJsonToStruct(msg, req)
	if err != nil {
		return
	}

	if req.Type == "" || req.From == "" {
		t.Errorf("Struct Did not map properly %v", req)
	}
}

func TestValidJsonMissingFieldsWorks(t *testing.T) {
	msg, _ := json.Marshal(MessageMock{
		Type: "PeerJoin",
	})

	req := &MessageMock{}

	err := MapJsonToStruct(msg, req)
	if err != nil {
		return
	}
	if req.Type == "" {
		t.Errorf("Struct Did not map properly %v", req)
	}

	if req.From != "" {
		t.Errorf("Struct Did not map properly %v", req)
	}
}

func TestInvalidJsonThrowsErr(t *testing.T) {
	msg, _ := json.Marshal("type: 'PeerLeave'")

	req := &MessageMock{}

	err := MapJsonToStruct(msg, req)
	if err == nil {
		t.Error("Should throw invalid json error")
	}
}
