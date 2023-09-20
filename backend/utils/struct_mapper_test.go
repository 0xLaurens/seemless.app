package utils

import (
	"encoding/json"
	"laurensdrop/data"
	"testing"
)

func TestValidJsonMapsToObject(t *testing.T) {
	msg, _ := json.Marshal(data.Request{
		Type: "PeerJoin",
		From: "John",
	})

	req := &data.Request{}

	err := MapJsonToStruct(msg, req)
	if err != nil {
		return
	}

	if req.Type == "" || req.From == "" {
		t.Errorf("Struct Did not map properly %v", req)
	}
}

func TestValidJsonMissingFieldsWorks(t *testing.T) {
	msg, _ := json.Marshal(data.Request{
		Type: "PeerJoin",
	})

	req := &data.Request{}

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

	req := &data.Request{}

	err := MapJsonToStruct(msg, req)
	if err == nil {
		t.Error("Should throw invalid json error")
	}
}
