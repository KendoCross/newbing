package newbing

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type BingChat struct {
	cookies      string
	chatReq      ChatRequest
	invocationId int
}

func NewChat(cookies string) (chatRoom *BingChat, err error) {
	chatRoom = &BingChat{
		cookies: cookies,
	}
	err = chatRoom.newConversation(context.Background())
	if err != nil {
		return
	}
	return
}

func (c *BingChat) Chat(ctx context.Context, askMsg string) (ans string, err error) {
	defer func() {
		c.invocationId++

	}()

	chatHub := newChatHub(c.chatReq)
	resp, err := chatHub.SendChatMsg(ctx, c.invocationId, askMsg)
	if err != nil {
		return
	}
	defer chatHub.Close()

	if len(resp.Item.Messages) > 1 {
		ans = resp.Item.Messages[1].Text
	}

	return
}

func (c *BingChat) newConversation(ctx context.Context) (err error) {
	const URL = "https://www.bing.com/turing/conversation/create"

	request, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return
	}
	request.Header.Set("accept", "application/json")
	request.Header.Set("content-type", "application/json")
	request.Header.Set("accept-language", "en-US,en;q=0.9")
	request.Header.Set("sec-ch-ua", `"Not_A Brand";v="99", "Microsoft Edge";v="109", "Chromium";v="109"`)
	request.Header.Set("sec-ch-ua-arch", `"x86"`)
	request.Header.Set("sec-ch-ua-bitness", `"64"`)
	request.Header.Set("sec-ch-ua-full-version", `"112.0.1722.68"`)
	request.Header.Set("sec-ch-ua-full-version-list", `"Chromium";v="112.0.5615.138", "Microsoft Edge";v="112.0.1722.68", "Not:A-Brand";v="99.0.0.0"`)
	request.Header.Set("sec-ch-ua-mobile", `?0`)
	request.Header.Set("sec-ch-ua-model", ``)
	request.Header.Set("sec-ch-ua-platform", `"Windows"`)
	request.Header.Set("sec-ch-ua-platform-version", `"15.0.0"`)
	request.Header.Set("sec-fetch-dest", `empty`)
	request.Header.Set("sec-fetch-mode", `cors`)
	request.Header.Set("sec-fetch-site", `same-origin`)
	request.Header.Set("x-ms-BingChat-request-id", uuid.New().String())
	request.Header.Set("x-ms-useragent", "azsdk-js-api-BingChat-factory/1.0.0-beta.1 core-rest-pipeline/1.10.0 OS/Win32")
	request.Header.Set("cookie", c.cookies)
	request.Header.Set("Referer", "https://www.bing.com/search?form=MY02AA&OCID=MY02AA&pl=launch&q=Bing+AI&showconv=1")
	request.Header.Set("Referrer-Policy", "origin-when-cross-origin")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("CURL %s http.StatusCode = %d", URL, resp.StatusCode)
		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	response := struct {
		ConversationId        string `json:"conversationId"`
		ClientId              string `json:"clientId"`
		ConversationSignature string `json:"conversationSignature"`
		Result                struct {
			Value   string `json:"value"`
			Message any    `json:"message"`
		} `json:"result"`
	}{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	if response.Result.Value != "Success" {
		err = fmt.Errorf("CURL %s : %s", URL, string(body))
		return
	}

	c.chatReq = ChatRequest{
		ConversationId:        response.ConversationId,
		ClientId:              response.ClientId,
		ConversationSignature: response.ConversationSignature,
	}
	return
}
