package foldable

import (
	"reflect"
	"testing"
)

func TestChannelAppend(t *testing.T) {
	expected := []int{5, 10, 15}
	channel := make(Channel)
	go func() {
		channel.Append(5).Append(10).Append(15)
		close(channel)
	}()
	got := []int{}
	for x := range channel {
		got = append(got, x.(int))
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestChannelMap(t *testing.T) {
	expected := []int{2, 4, 6}
	// create initial data
	channel := make(Channel)
	go func() {
		channel <- 1
		channel <- 2
		channel <- 3
		close(channel)
	}()
	// use gofuncs
	gotChannel := Map(
		channel,
		func(item T) T { return item.(int) * 2 }).(Channel)
	// get result
	got := []int{}
	for x := range gotChannel {
		got = append(got, x.(int))
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestChannelMapTwice(t *testing.T) {
	expected := []int{3, 5, 7}
	// create initial data
	channel := make(Channel)
	go func() {
		channel <- 1
		channel <- 2
		channel <- 3
		close(channel)
	}()
	// use gofuncs
	gotChannel := Map(
		Map(channel, func(item T) T { return item.(int) * 2 }).(Channel),
		func(item T) T { return item.(int) + 1 }).(Channel)
	// get result
	got := []int{}
	for x := range gotChannel {
		got = append(got, x.(int))
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestChannelParMap(t *testing.T) {
	expected := []int{2, 3, 4}
	// create initial data
	channel := make(Channel)
	go func() {
		channel <- 1
		channel <- 2
		channel <- 3
		close(channel)
	}()
	// use gofuncs
	gotChannel := ParMap(
		channel,
		func(item T) T {
			return item.(int) + 1
		}).(Channel)
	// get result
	got := []int{}
	for x := range gotChannel {
		got = append(got, x.(int))
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestChannelFilter(t *testing.T) {
	expected := []int{-1, -2, -53}
	// create initial data
	channel := make(Channel)
	go func() {
		channel <- -1
		channel <- 1
		channel <- -2
		channel <- 3
		channel <- -53
		close(channel)
	}()
	// use gofuncs
	gotChannel := Filter(
		channel,
		func(item T) bool { return item.(int) < 0 }).(Channel)
	// get result
	got := []int{}
	for x := range gotChannel {
		got = append(got, x.(int))
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestChannelLength(t *testing.T) {
	expected := 5
	// create initial data
	channel := make(Channel)
	go func() {
		channel <- -1
		channel <- 1
		channel <- -2
		channel <- 3
		channel <- -53
		close(channel)
	}()
	// use gofuncs
	got := Length(channel)
	// get result
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}

func TestListToChannel(t *testing.T) {
	expected := []int{2, 4, 6}
	// create initial data
	in := List{1, 2, 3}
	// use gofuncs
	gotChannel := Map(
		ToChannel(in),
		func(x T) T { return x.(int) * 2 },
	)
	// get result
	got := []int{}
	for x := range gotChannel.(Channel) {
		got = append(got, x.(int))
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("result == %v expected %v", got, expected)
	}
}
