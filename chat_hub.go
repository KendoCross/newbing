package newbing

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ChatHub struct {
	conn    *websocket.Conn
	chatReq ChatRequest
}

func newChatHub(req ChatRequest) *ChatHub {
	hub := &ChatHub{
		chatReq: req,
	}
	return hub
}

func (c *ChatHub) SendChatMsg(ctx context.Context, invocationId int, message string) (resp ChatResponse, err error) {
	msgCh, err := c.sendChatMsg(ctx, invocationId, message)
	if err != nil {
		return
	}
	resp = <-msgCh
	return
}

func (c *ChatHub) sendChatMsg(ctx context.Context, invocationId int, message string) (msgCh <-chan ChatResponse, err error) {
	ch := make(chan ChatResponse)
	err = c.openWS(ctx, ch)
	if err != nil {
		return
	}
	msgCh = ch

	data := map[string]any{
		"arguments": []any{
			map[string]any{
				"source": "cib",
				"optionsSets": []any{
					"nlu_direct_response_filter",
					"deepleo",
					"disable_emoji_spoken_text",
					"responsible_ai_policy_235",
					"enablemm",
					"galileo",
					"visualcreative",
					"enbfpr",
					"dv3sugg",
				},
				"allowedMessageTypes": []any{
					"Chat",
					"InternalSearchQuery",
					"InternalSearchResult",
					"Disengaged",
					"InternalLoaderMessage",
					"RenderCardRequest",
					"AdsQuery",
					"SemanticSerp",
					"GenerateContentQuery",
					"SearchQuery",
				},
				"sliceIds": []string{
					"winmuid1tf",
					"0427btfirstc",
					"ssoverlap0",
					"sswebtop1",
					"forallv2pc",
					"sbsvgoptcf",
					"winstmsg2tf",
					"contansperf",
					"ctrlconvcss",
					"0427visual_b",
					"420langdsats0",
					"420bics0",
					"0329resps0",
					"425bicpctrl",
					"425bfpr",
					"424dagslnv1s0",
				},
				"isStartOfSession": invocationId == 0,
				"message": map[string]any{
					"author":      "user",
					"inputMethod": "Keyboard",
					"text":        message,
					"messageType": "Chat",
				},
				"conversationSignature": c.chatReq.ConversationSignature,
				"participant": map[string]any{
					"id": c.chatReq.ClientId,
				},
				"conversationId": c.chatReq.ConversationId,
			},
		},
		"invocationId": strconv.Itoa(invocationId),
		"target":       "chat",
		"type":         4,
	}
	marshal, _ := json.Marshal(data)
	err = c.sendWSMessage(marshal)
	if err != nil {
		return
	}
	return
}

func (c *ChatHub) openWS(ctx context.Context, ch chan ChatResponse) (err error) {
	const URL = "wss://sydney.bing.com/sydney/ChatHub"
	// 建立WebSocket连接
	conn, _, err := websocket.DefaultDialer.Dial(URL, nil)
	if err != nil {
		return
	}
	c.conn = conn
	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"protocol":"json","version":1}`+Split))
	if err != nil {
		return
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Listen WebSocketConnection done.")
				close(ch)
				return
			default:
				messageType, message, err := conn.ReadMessage()
				if err != nil {
					// 需要额外处理
					continue
				}
				switch messageType {
				case websocket.TextMessage:
					msg := string(message)
					splits := strings.Split(msg, Split)
					splits = lo.Filter(splits, func(item string, index int) bool {
						return strings.TrimSpace(item) != ""
					})
					for i := range splits {
						c.formatMsg(splits[i], ch)
					}
				}
			}
		}
	}()
	return
}

func (c *ChatHub) Close() (err error) {
	return
}
func (c *ChatHub) sendPongMsg() (err error) {
	err = c.sendWSMessage([]byte(`{"type":6}`))
	return
}

func (c *ChatHub) sendWSMessage(message []byte) (err error) {
	err = c.conn.WriteMessage(websocket.TextMessage, append(message, []byte(Split)...))
	return
}

func (c *ChatHub) formatMsg(message string, ch chan ChatResponse) (err error) {
	if !(strings.HasPrefix(message, "{\"type\":1") || strings.HasPrefix(message, "{\"type\":2")) {
		fmt.Printf("reviced message: %s\n", message)
	}
	var response = map[string]any{}
	err = json.Unmarshal([]byte(message), &response)
	if err != nil {
		return
	}

	t, ok := response["type"]
	if !ok {
		err = c.sendPongMsg()
		if err != nil {
			return
		}
		return
	}

	switch int64(t.(float64)) {
	case 1:
	case 2:
		res := ChatResponse{}
		err = json.Unmarshal([]byte(message), &res)
		if err != nil {
			return
		}
		ch <- res
	case 3:
	case 6:
		err = c.sendPongMsg()
	case 7:
	}
	return
}

type ChatRequest struct {
	ConversationId        string `json:"conversationId"`
	ClientId              string `json:"clientId"`
	ConversationSignature string `json:"conversationSignature"`
}

type ChatResponse struct {
	Type         int `json:"type"`
	InvocationId int `json:"invocationId,string"`
	Item         struct {
		Messages []struct {
			Text   string `json:"text"`
			Author string `json:"author"`
			From   struct {
				Id   string `json:"id"`
				Name any    `json:"name"`
			} `json:"from"`
			CreatedAt     time.Time `json:"createdAt"`
			Timestamp     time.Time `json:"timestamp"`
			Locale        string    `json:"locale"`
			Market        string    `json:"market"`
			Region        string    `json:"region"`
			Location      string    `json:"location"`
			LocationHints []struct {
				Country           string `json:"country"`
				CountryConfidence int    `json:"countryConfidence"`
				State             string `json:"state"`
				City              string `json:"city"`
				CityConfidence    int    `json:"cityConfidence"`
				ZipCode           string `json:"zipCode"`
				TimeZoneOffset    int    `json:"timeZoneOffset"`
				Dma               int    `json:"dma"`
				SourceType        int    `json:"sourceType"`
				Center            struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
					Height    any     `json:"height"`
				} `json:"center"`
				RegionType int `json:"regionType"`
			} `json:"locationHints"`
			MessageId uuid.UUID `json:"messageId"`
			RequestId uuid.UUID `json:"requestId"`
			Offense   string    `json:"offense"`
			Feedback  struct {
				Tag       any    `json:"tag"`
				UpdatedOn any    `json:"updatedOn"`
				Type      string `json:"type"`
			} `json:"feedback"`
			ContentOrigin string `json:"contentOrigin"`
			Privacy       any    `json:"privacy"`
			InputMethod   string `json:"inputMethod"`
			HiddenText    string `json:"hiddenText"`
			MessageType   string `json:"messageType"`
			AdaptiveCards []struct {
				Type    string `json:"type"`
				Version string `json:"version"`
				Body    []struct {
					Type    string `json:"type"`
					Inlines []struct {
						Type     string `json:"type"`
						IsSubtle bool   `json:"isSubtle"`
						Italic   bool   `json:"italic"`
						Text     string `json:"text"`
					} `json:"inlines,omitempty"`
					Text string `json:"text,omitempty"`
					Wrap bool   `json:"wrap,omitempty"`
					Size string `json:"size,omitempty"`
				} `json:"body"`
			} `json:"adaptiveCards"`
			SourceAttributions []struct {
				ProviderDisplayName string `json:"providerDisplayName"`
				SeeMoreUrl          string `json:"seeMoreUrl"`
				SearchQuery         string `json:"searchQuery"`
			} `json:"sourceAttributions"`
			SuggestedResponses []struct {
				Text        string    `json:"text"`
				Author      string    `json:"author"`
				CreatedAt   time.Time `json:"createdAt"`
				Timestamp   time.Time `json:"timestamp"`
				MessageId   string    `json:"messageId"`
				MessageType string    `json:"messageType"`
				Offense     string    `json:"offense"`
				Feedback    struct {
					Tag       any    `json:"tag"`
					UpdatedOn any    `json:"updatedOn"`
					Type      string `json:"type"`
				} `json:"feedback"`
				ContentOrigin string `json:"contentOrigin"`
				Privacy       any    `json:"privacy"`
			} `json:"suggestedResponses"`
			SpokenText string `json:"spokenText"`
		} `json:"messages"`
		FirstNewMessageIndex   int       `json:"firstNewMessageIndex"`
		SuggestedResponses     any       `json:"suggestedResponses"`
		ConversationId         string    `json:"conversationId"`
		RequestId              string    `json:"requestId"`
		ConversationExpiryTime time.Time `json:"conversationExpiryTime"`
		Telemetry              struct {
			Metrics   any       `json:"metrics"`
			StartTime time.Time `json:"startTime"`
		} `json:"telemetry"`
		ShouldInitiateConversation bool `json:"shouldInitiateConversation"`
		Result                     struct {
			Value          string `json:"value"`
			Message        string `json:"message"`
			ServiceVersion string `json:"serviceVersion"`
		} `json:"result"`
	} `json:"item,omitempty"`
}
