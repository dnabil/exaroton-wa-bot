package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mau.fi/whatsmeow"
)

func TestWAClient_GetSetQRSub(t *testing.T) {
	wa := &waClient{}

	// Initially nil
	assert.Nil(t, wa.qrSub)
	assert.Nil(t, wa.getQRSub())

	// Set and verify channel
	sub := make(chan whatsmeow.QRChannelItem)
	wa.setQRSub(&sub)

	// Verify the channel is stored and retrieved correctly
	assert.NotNil(t, wa.qrSub)
	assert.NotNil(t, wa.getQRSub())

	// Send a test message through the channel
	testMsg := whatsmeow.QRChannelItem{
		Event: whatsmeow.QRChannelEventCode,
		Code:  "test-code",
	}

	go func() {
		wa.qrSub <- testMsg
	}()

	// Verify we can receive the message
	select {
	case msg := <-sub:
		assert.Equal(t, testMsg, msg)
	case <-time.After(time.Second):
		t.Error("Timeout waiting for message")
	}

	// Reset to nil
	wa.setQRSub(nil)
	assert.Nil(t, wa.qrSub)
	assert.Nil(t, wa.getQRSub())
}
