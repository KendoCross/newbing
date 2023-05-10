package test

import (
	"context"
	"testing"
	"time"

	"github.com/KendoCross/newbing"
)

func TestImgGen(t *testing.T) {
	bingImgGen := newbing.NewImgGen("new bing cookie ")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	imgURLsCh, err := bingImgGen.GenImgAync(ctx, "功夫熊猫")
	if err != nil {
		t.Error(err)
		return
	}

	imgURLs := <-imgURLsCh
	for _, url := range imgURLs {
		println(url)
	}
}
