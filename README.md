# newbing

new bing API

## example

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
