package test

import (
	"context"
	"testing"

	"github.com/KendoCross/newbing"
)

func TestSMS(t *testing.T) {
	bingChat, err := newbing.NewChat("new bing cookie ")
	if err != nil {
		t.Error(err)
		return
	}
	ans, err := bingChat.Chat(context.TODO(), "你是chatGPT吗？")
	if err != nil {
		t.Error(err)
		return
	}
	println(ans)

	ans, err = bingChat.Chat(context.TODO(), "你能做什么？")
	if err != nil {
		t.Error(err)
		return
	}
	println(ans)

}
