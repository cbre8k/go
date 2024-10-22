package cchannel

import (
	"testing"
	"time"
)

// Benchmark for the custom CChannel
func BenchmarkCChannel_SendReceive(b *testing.B) {
	ch := NewChannel(1, Bidirectional)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			ch.Send(i)
		}()

		_, ok := ch.Receive()
		if !ok {
			b.Error("Failed to receive value")
		}
	}
}

// Benchmark for Go's native channel
func BenchmarkNativeChannel_SendReceive(b *testing.B) {
	ch := make(chan int, 1)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			ch <- i
		}()

		_, ok := <-ch
		if !ok {
			b.Error("Failed to receive value")
		}
	}
}

// Benchmark for custom CChannel under heavy load (e.g., large buffer)
func BenchmarkCChannel_Buffered(b *testing.B) {
	ch := NewChannel(1024, Bidirectional)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch.Send(i)
		ch.Receive()
	}
}

// Benchmark for Go's native channel with large buffer
func BenchmarkNativeChannel_Buffered(b *testing.B) {
	ch := make(chan int, 1024)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- i
		<-ch
	}
}

// Benchmark closing the custom CChannel
func BenchmarkCChannel_Close(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := NewChannel(1, Bidirectional)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					// Ensure we panic with the expected message
					if r != "send on closed channel" {
						b.Errorf("unexpected panic: %v", r)
					}
				}
			}()

			for j := 0; j < 10; j++ {
				ch.Send(j) // This should panic after the channel is closed
			}
		}()

		time.Sleep(time.Microsecond) // Allow some time for the send to happen
		ch.Close()                   // Close the channel, should cause panic on subsequent Send calls
	}
}

// Benchmark closing Go's native channel
func BenchmarkNativeChannel_Close(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := make(chan int, 1)

		go func(val int) {
			defer func() {
				// Recover from panic without logging
				_ = recover() // Ignore the panic without logging
			}()

			ch <- val // Attempt to send the value
		}(i)

		time.Sleep(time.Microsecond) // Allow some time for the goroutine to send
		close(ch)                    // Close the channel, which may cause panic in the goroutine

		// Optionally, add a small delay to let goroutines finish
		time.Sleep(time.Microsecond)
	}
}
