package ohio_ring_buffer2

type Scanner struct {
	rb                  *RingBuffer
	whitespaceIsContent bool
}
