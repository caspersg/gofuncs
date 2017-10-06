package foldable

// Channel as lazy list
// channels in go can be viewed as a lazy collection, therefore we can make a fold function for them
type Channel chan T

func (channel Channel) Foldl(init T, foldFunc func(result, next T) T) T {
	result := init
	go func() {
		for item := range channel {
			// normally a fold would reassign the result here, but a channel is naturally mutable
			// so we're relying on the foldFunc to use Append to mutate the passed in result
			// this isn't very functional, but shouldn't be a real problem if the interface is used as intended
			foldFunc(result, item)
		}
	}()
	return result
}

func (channel Channel) Init() Foldable {
	var ch Channel = make(chan T)
	return ch
}

func (channel Channel) Append(item T) Foldable {
	// as above, this mutates the passed in channel, which is the only way to use channels in go
	channel <- item
	return channel
}
