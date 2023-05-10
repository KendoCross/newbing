package newbing

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

type BingImgGen struct {
	cookies string
}

func NewImgGen(cookies string) (imgGen *BingImgGen) {
	imgGen = &BingImgGen{
		cookies: cookies,
	}
	return
}

var HEADERS = map[string]string{
	"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"accept-language": "en-US,en;q=0.9",
	"cache-control":   "max-age=0",
	"content-type":    "application/x-www-form-urlencoded",
	"referrer":        "https://www.bing.com/images/create/",
	"origin":          "https://www.bing.com",
	"user-agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.63",
}

func (gen *BingImgGen) GenImgAync(ctx context.Context, prompt string) (ch <-chan []string, err error) {
	resp, err := gen.postCreate(ctx, prompt, 4)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusFound {
		// 尝试一次rt3
		resp, err = gen.postCreate(ctx, prompt, 3)
		if err != nil {
			return
		}
		if resp.StatusCode != http.StatusFound {
			err = fmt.Errorf("CURL %s http.StatusCode = %d", "createPath", resp.StatusCode)
			return
		}
	}
	redirectPath := resp.Header.Get("Location")
	locations := strings.Split(redirectPath, "&id=")
	requestId := locations[len(locations)-1]

	tmpCh := make(chan []string)
	ch = tmpCh
	redirectPath = fmt.Sprintf("https://www.bing.com/images/create/async/results/%s?q=%s", requestId, url.QueryEscape(prompt))
	go func() {

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Waiting for results done")
				close(tmpCh)
				return
			default:
				rsts, errIn := gen.getresults(redirectPath)
				if errIn != nil {
					err = errIn
					return
				}
				if len(rsts) > 0 {
					tmpCh <- rsts
					break
				}
				time.Sleep(time.Second)
			}
		}
	}()

	return
}

func (gen *BingImgGen) postCreate(ctx context.Context, prompt string, rt int) (resp *http.Response, err error) {
	createPath := fmt.Sprintf("https://www.bing.com/images/create?q=%s&rt=%d&FORM=GENCRE", url.QueryEscape(prompt), rt)
	payload := url.Values{
		"q":  {prompt},
		"qs": {"ds"},
	}
	bodyReader := strings.NewReader(payload.Encode())

	request, err := http.NewRequestWithContext(ctx, "POST", createPath, bodyReader)
	if err != nil {
		return
	}
	for k, v := range HEADERS {
		request.Header.Add(k, v)
	}
	request.Header.Add("x-forwarded-for", fmt.Sprintf("13.%d.%d.%d", 103+rand.Intn(4), rand.Intn(255)+1, rand.Intn(255)+1))
	request.Header.Add("cookie", gen.cookies)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err = client.Do(request)
	if err != nil {
		return
	}
	return
}

func (gen *BingImgGen) getresults(redirectPath string) (rsts []string, err error) {
	request, err := http.NewRequest("GET", redirectPath, nil)
	if err != nil {
		return
	}
	for k, v := range HEADERS {
		request.Header.Add(k, v)
	}
	request.Header.Add("cookie", gen.cookies)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("CURL %s http.StatusCode = %d", redirectPath, resp.StatusCode)
		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		return
	}

	re, err := regexp2.Compile(`src="([^"]+)"`, 0)
	if err != nil {
		return
	}
	rsts = gen.findAllString(re, string(body))

	return

}

func (gen *BingImgGen) findAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		rstUrl := strings.Split(m.String(), "?w=")[0]
		rstUrl = strings.ReplaceAll(rstUrl, "src=\"", "")
		matches = append(matches, rstUrl)
		m, _ = re.FindNextMatch(m)
	}
	return matches
}
