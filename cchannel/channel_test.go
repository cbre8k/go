package cchannel_test

import (
	"cchannel"
	"sync"
	"testing"
	"time"
)

// Basic test for send and receive on an unbuffered channel
func TestChannelBasicSendReceive(t *testing.T) {
	ch := cchannel.NewChannel(0, cchannel.Bidirectional) // Unbuffered channel
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		msg := "Hello"
		t.Log("Sending:", msg)
		ch.Send(msg)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 100) // Simulate delay in receiving
		msg, ok := ch.Receive()
		t.Logf("Received: %v, ok: %v", msg, ok)

		if msg != "Hello" || !ok {
			t.Errorf("Expected to receive 'Hello', got %v", msg)
		}
	}()

	wg.Wait()
}

// Test multiple sends and receives in a buffered channel
func TestChannelMultipleSendReceive(t *testing.T) {
	ch := cchannel.NewChannel(1, cchannel.Bidirectional) // Buffered channel with buffer size 1
	var wg sync.WaitGroup

	wg.Add(5)

	go func() {
		defer wg.Done()
		ch.Send("Message 1")
	}()

	go func() {
		defer wg.Done()
		ch.Send("Message 2")
	}()

	go func() {
		defer wg.Done()
		ch.Send("Message 3")
	}()

	go func() {
		defer wg.Done()
		ch.Send("Message 4")
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 1)
		for i := 0; i < 4; i++ {
			msg, ok := ch.Receive()
			t.Logf("Received: %v, ok: %v", msg, ok)
			if msg == nil || !ok {
				t.Errorf("Expected to receive a valid message, got %v", msg)
			}
		}
	}()

	wg.Wait()
}

// Test closing the channel and behavior after closure
func TestChannelClose(t *testing.T) {
	ch := cchannel.NewChannel(1, cchannel.Bidirectional) // Buffered channel with buffer size 1
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		ch.Send("Message before close")
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 100)
		msg, ok := ch.Receive()
		t.Logf("Received: %v, ok: %v", msg, ok)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 200)
		ch.Close()
		t.Log("Channel closed")

		// After the channel is closed, try sending and expect panic
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic after sending to a closed channel")
			}
		}()
		ch.Send("Message after close")
	}()

	wg.Wait()
}

// Test receiving from a closed channel
func TestReceiveFromClosedChannel(t *testing.T) {
	ch := cchannel.NewChannel(1, cchannel.Bidirectional) // Buffered channel
	ch.Close()

	// Receiving from a closed channel should return nil, false
	msg, ok := ch.Receive()
	if msg != nil || ok {
		t.Errorf("Expected to receive nil, false from a closed channel, got %v, %v", msg, ok)
	}
}

// Test sending on an unbuffered channel where receiver is delayed
func TestUnbufferedSendReceive(t *testing.T) {
	ch := cchannel.NewChannel(0, cchannel.Bidirectional) // Unbuffered channel
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 100) // Simulate delay in receiver
		msg, ok := ch.Receive()
		t.Logf("Received: %v, ok: %v", msg, ok)

		if msg != "Unbuffered Message" || !ok {
			t.Errorf("Expected to receive 'Unbuffered Message', got %v", msg)
		}
	}()

	go func() {
		defer wg.Done()
		t.Log("Sending unbuffered message")
		ch.Send("Unbuffered Message")
	}()

	wg.Wait()
}

// Test sending when buffer is full and receiver reads before send blocks
func TestBufferedChannelFullSendReceive(t *testing.T) {
	ch := cchannel.NewChannel(2, cchannel.Bidirectional) // Buffered channel with buffer size 2
	var wg sync.WaitGroup

	wg.Add(3)

	// Fill the buffer
	go func() {
		defer wg.Done()
		ch.Send("Buffered Message 1")
		ch.Send("Buffered Message 2")
		t.Log("Buffer full now")
	}()

	// Receiver should empty the buffer, allowing further sends
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 100)
		msg1, ok1 := ch.Receive()
		t.Logf("Received: %v, ok: %v", msg1, ok1)
		msg2, ok2 := ch.Receive()
		t.Logf("Received: %v, ok: %v", msg2, ok2)
		if msg1 != "Buffered Message 1" || !ok1 || msg2 != "Buffered Message 2" || !ok2 {
			t.Errorf("Expected to receive 'Buffered Message 1' and 'Buffered Message 2', got %v and %v", msg1, msg2)
		}
	}()

	// Another send that will block until the receiver reads from the buffer
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 150) // Ensure buffer is emptied before this send
		ch.Send("Buffered Message 3")
		t.Log("Sent 'Buffered Message 3'")
	}()

	wg.Wait()
}

// Test panic when sending to a closed unbuffered channel
func TestPanicOnSendToClosedUnbufferedChannel(t *testing.T) {
	ch := cchannel.NewChannel(0, cchannel.Bidirectional) // Unbuffered channel
	ch.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when sending to a closed unbuffered channel")
		}
	}()

	ch.Send("This should panic")
}

// Test receiving from an unbuffered channel after close
func TestReceiveAfterCloseUnbuffered(t *testing.T) {
	ch := cchannel.NewChannel(0, cchannel.Bidirectional) // Unbuffered channel
	ch.Close()

	// Receiving from a closed channel should return nil, false
	msg, ok := ch.Receive()
	if msg != nil || ok {
		t.Errorf("Expected to receive nil, false from closed unbuffered channel, got %v, %v", msg, ok)
	}
}

// Test panic when sending to a closed buffered channel
func TestPanicOnSendToClosedBufferedChannel(t *testing.T) {
	ch := cchannel.NewChannel(1, cchannel.Bidirectional) // Buffered channel
	ch.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when sending to a closed buffered channel")
		}
	}()

	ch.Send("This should panic")
}

// Test SendOnly type for a unidirectional send-only channel
func TestSendOnlyChannel(t *testing.T) {
	ch := cchannel.NewChannel(1, cchannel.SendOnly) // Send-only channel

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when receiving from a send-only channel")
		}
	}()

	// Attempting to receive from a send-only channel should panic
	ch.Receive()
}

// Test ReceiveOnly type for a unidirectional receive-only channel
func TestReceiveOnlyChannel(t *testing.T) {
	ch := cchannel.NewChannel(1, cchannel.ReceiveOnly) // Receive-only channel

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when sending to a receive-only channel")
		}
	}()

	// Attempting to send to a receive-only channel should panic
	ch.Send("This should panic")
}
