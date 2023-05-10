# newbing

new bing API

## new bing chat

```golang

    import "github.com/KendoCross/newbing"
    
    // 通过 cookie 初始化
    bingChat, err := newbing.NewChat(" new bing cookie ")
    if err != nil {
        t.Error(err)
        return
    }
    // 发起聊天
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
    
```

## new bing image generation

```golang

    import "github.com/KendoCross/newbing"
    
    bingImgGen := newbing.NewImgGen("new bing cookie")
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
    
```
