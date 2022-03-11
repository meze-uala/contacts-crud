package service

import (
	"contacts-crud/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPushContactToSNSService_PublishContactIDToSNS(t *testing.T) {

	snsIClient := &internal.MockSNS{}
	service := NewPushContactToSNSService(snsIClient.SnsClient)
	result, err := service.PublishContactIDToSNS("1")

	assert.NotNil(t, result)
	assert.NoError(t, err)
}
