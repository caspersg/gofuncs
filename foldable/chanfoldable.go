package foldable

import (
	"reflect"
)

// Channel as lazy list
// channels in go can be viewed as a lazy collection, therefore we can make a fold function for them
type Channel chan T

func (channel Channel) Foldl(init T, foldFunc func(result, next T) T) T {
	result := init
	if reflect.TypeOf(result).Name() == "Channel" {
		// if the result is a channel we need some special handling
		// for instance, we don't want to block on processing, so we process in a go func and return the result channel
		// this is a pretty major concession to channels, but they are treated specially in go
		// and there are no other special cases like this, so I think it's acceptable
		go func() {
			for item := range channel {
				// normally a fold would reassign the result here, but a channel is naturally mutable
				// so we're relying on the foldFunc to use Append to mutate the passed in result
				// this isn't very functional, but shouldn't be a real problem if the interface is used as intended
				foldFunc(result, item)
			}
			close(result.(Channel))
		}()
	} else {
		// if the result is some other type, normall handling will be fine
		for item := range channel {
			result = foldFunc(result, item)
		}
	}
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
