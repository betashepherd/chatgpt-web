package controllers

import (
	"chatgpt-web/library/lfs"
	"chatgpt-web/library/util"
	"chatgpt-web/pkg/model/user"
	"chatgpt-web/pkg/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"

	"chatgpt-web/config"
	"chatgpt-web/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/sashabaranov/go-openai"
)

var chatModels = []string{openai.GPT432K0314, openai.GPT4, openai.GPT40314, openai.GPT432K, openai.GPT3Dot5Turbo, openai.GPT3Dot5Turbo0301}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ChatController 首页控制器
type ChatController struct {
	BaseController
}

// NewChatController 创建控制器
func NewChatController() *ChatController {
	return &ChatController{}
}

// Index 首页
func (c *ChatController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func (c *ChatController) CompletionWsStream(ctx *gin.Context) {
	wsClient, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	defer wsClient.Close()

	wsResp := map[string]interface{}{
		"code":     http.StatusOK,
		"errorMsg": "",
		"data":     nil,
	}

	cnf := config.LoadConfig()
	gptConfig := openai.DefaultConfig(cnf.ApiKey)

	if cnf.Proxy != "" {
		transport := &http.Transport{}

		if strings.HasPrefix(cnf.Proxy, "socks5h://") {
			// 创建一个 DialContext 对象，并设置代理服务器
			dialContext, err := newDialContext(cnf.Proxy[10:])
			if err != nil {
				panic(err)
			}
			transport.DialContext = dialContext
		} else {
			// 创建一个 HTTP Transport 对象，并设置代理服务器
			proxyUrl, err := url.Parse(cnf.Proxy)
			if err != nil {
				panic(err)
			}
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
		// 创建一个 HTTP 客户端，并将 Transport 对象设置为其 Transport 字段
		gptConfig.HTTPClient = &http.Client{
			Transport: transport,
		}

	}

	// 自定义gptConfig.BaseURL
	if cnf.ApiURL != "" {
		gptConfig.BaseURL = cnf.ApiURL
	}

	client := openai.NewClientWithConfig(gptConfig)

	var request openai.ChatCompletionRequest
	if err := wsClient.ReadJSON(&request); err != nil {
		wsResp["code"] = http.StatusInternalServerError
		wsResp["errorMsg"] = err.Error()
		wsClient.WriteJSON(wsResp)
		return
	}
	logger.Info(request)
	if len(request.Messages) == 0 {
		wsResp["code"] = http.StatusBadRequest
		wsResp["errorMsg"] = "request messages required"
		wsClient.WriteJSON(wsResp)
		return
	}
	if request.Messages[0].Role != "system" {
		newMessage := append([]openai.ChatCompletionMessage{
			{Role: "system", Content: cnf.BotDesc},
		}, request.Messages...)
		request.Messages = newMessage
	}

	request.Model = cnf.Model
	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		wsResp["code"] = http.StatusInternalServerError
		wsResp["errorMsg"] = err.Error()
		wsClient.WriteJSON(wsResp)
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			//fmt.Println("\nStream finished")
			fmt.Println("\n")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(response.Choices[0].Delta.Content)
		/**
		"reply":    resp.Choices[0].Message.Content,
		"messages": append(request.Messages, resp.Choices[0].Message),
		*/
		wsClient.WriteJSON(wsResp)
	}

	return
}

// Completion 回复
func (c *ChatController) Completion(ctx *gin.Context) {
	var request openai.ChatCompletionRequest
	if err := ctx.BindJSON(&request); err != nil {
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var authUser *user.User
	if iter, ok := ctx.Get("authUser"); ok {
		authUser = iter.(*user.User)
	}

	rjs, _ := json.Marshal(request)
	logger.Info(authUser.Name, "__", string(rjs))

	if len(request.Messages) == 0 {
		c.ResponseJson(ctx, http.StatusBadRequest, "request messages required", nil)
		return
	}

	cnf := config.LoadConfig()
	gptConfig := openai.DefaultConfig(cnf.ApiKey)

	if cnf.Proxy != "" {
		transport := &http.Transport{}

		if strings.HasPrefix(cnf.Proxy, "socks5h://") {
			// 创建一个 DialContext 对象，并设置代理服务器
			dialContext, err := newDialContext(cnf.Proxy[10:])
			if err != nil {
				panic(err)
			}
			transport.DialContext = dialContext
		} else {
			// 创建一个 HTTP Transport 对象，并设置代理服务器
			proxyUrl, err := url.Parse(cnf.Proxy)
			if err != nil {
				panic(err)
			}
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
		// 创建一个 HTTP 客户端，并将 Transport 对象设置为其 Transport 字段
		gptConfig.HTTPClient = &http.Client{
			Transport: transport,
		}

	}

	// 自定义gptConfig.BaseURL
	if cnf.ApiURL != "" {
		gptConfig.BaseURL = cnf.ApiURL
	}

	client := openai.NewClientWithConfig(gptConfig)
	if request.Messages[0].Role != "system" {
		newMessage := append([]openai.ChatCompletionMessage{
			{Role: "system", Content: cnf.BotDesc},
		}, request.Messages...)
		request.Messages = newMessage
	}

	if types.Contains(chatModels, cnf.Model) {
		request.Model = cnf.Model
		resp, err := client.CreateChatCompletion(ctx, request)
		if err != nil {
			c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		answer := gin.H{
			"reply":    resp.Choices[0].Message.Content,
			"messages": append(request.Messages, resp.Choices[0].Message),
		}

		ajs, _ := json.Marshal(answer)
		question := fmt.Sprintf("question_%s_%s.json", request.Messages[len(request.Messages)-1].Content, util.GetCurrentTime().Format("20060102150405000"))
		subDir := fmt.Sprintf("chat/%s", authUser.Name)
		lfs.DataFs.SaveDataFile(question, ajs, subDir)
		c.ResponseJson(ctx, http.StatusOK, "", answer)

	} else {
		prompt := ""
		for _, item := range request.Messages {
			prompt += item.Content + "/n"
		}
		prompt = strings.Trim(prompt, "/n")

		logger.Info("request prompt is %s", prompt)
		req := openai.CompletionRequest{
			Model:            cnf.Model,
			MaxTokens:        cnf.MaxTokens,
			TopP:             cnf.TopP,
			FrequencyPenalty: cnf.FrequencyPenalty,
			PresencePenalty:  cnf.PresencePenalty,
			Prompt:           prompt,
		}

		resp, err := client.CreateCompletion(ctx, req)
		if err != nil {
			c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		c.ResponseJson(ctx, http.StatusOK, "", gin.H{
			"reply": resp.Choices[0].Text,
			"messages": append(request.Messages, openai.ChatCompletionMessage{
				Role:    "assistant",
				Content: resp.Choices[0].Text,
			}),
		})
	}

}

type dialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)

func newDialContext(socks5 string) (dialContextFunc, error) {
	baseDialer := &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}

	if socks5 != "" {
		// split socks5 proxy string [username:password@]host:port
		var auth *proxy.Auth = nil

		if strings.Contains(socks5, "@") {
			proxyInfo := strings.SplitN(socks5, "@", 2)
			proxyUser := strings.Split(proxyInfo[0], ":")
			if len(proxyUser) == 2 {
				auth = &proxy.Auth{
					User:     proxyUser[0],
					Password: proxyUser[1],
				}
			}
			socks5 = proxyInfo[1]
		}

		dialSocksProxy, err := proxy.SOCKS5("tcp", socks5, auth, baseDialer)
		if err != nil {
			return nil, err
		}

		contextDialer, ok := dialSocksProxy.(proxy.ContextDialer)
		if !ok {
			return nil, err
		}

		return contextDialer.DialContext, nil
	} else {
		return baseDialer.DialContext, nil
	}
}

var questions = cmap.New[openai.ChatCompletionRequest]()

func (c *ChatController) Question(ctx *gin.Context) {
	var request openai.ChatCompletionRequest
	if err := ctx.BindJSON(&request); err != nil {
		logrus.Error(err.Error())
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if len(request.Messages) == 0 {
		c.ResponseJson(ctx, http.StatusBadRequest, "request messages required", nil)
		return
	}

	var authUser *user.User
	if iter, ok := ctx.Get("authUser"); ok {
		authUser = iter.(*user.User)
	}
	cnf := config.LoadConfig()
	if request.Messages[0].Role != "system" {
		newMessage := append([]openai.ChatCompletionMessage{
			{Role: "system", Content: cnf.BotDesc},
		}, request.Messages...)
		request.Messages = newMessage
	}

	rjs, _ := json.Marshal(request)
	logger.Info(authUser.Name, "__", string(rjs))

	questions.Set(authUser.Name, request)
	answer := gin.H{
		"id":       uuid.New().String(),
		"reply":    "",
		"messages": request.Messages,
	}
	c.ResponseJson(ctx, http.StatusOK, "", answer)
}

func (c *ChatController) Reply(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")
	ctx.Header("Access-Control-Allow-Origin", "*")

	var authUser *user.User
	if iter, ok := ctx.Get("authUser"); ok {
		authUser = iter.(*user.User)
	}

	request, ok := questions.Get(authUser.Name)
	if !ok {
		c.ResponseJson(ctx, http.StatusBadRequest, "no request exists", nil)
		return
	}
	questions.Remove(authUser.Name)

	if len(request.Messages) == 0 {
		c.ResponseJson(ctx, http.StatusBadRequest, "request messages required", nil)
		return
	}
	cnf := config.LoadConfig()
	gptConfig := openai.DefaultConfig(cnf.ApiKey)

	if cnf.Proxy != "" {
		transport := &http.Transport{}

		if strings.HasPrefix(cnf.Proxy, "socks5h://") {
			// 创建一个 DialContext 对象，并设置代理服务器
			dialContext, err := newDialContext(cnf.Proxy[10:])
			if err != nil {
				panic(err)
			}
			transport.DialContext = dialContext
		} else {
			// 创建一个 HTTP Transport 对象，并设置代理服务器
			proxyUrl, err := url.Parse(cnf.Proxy)
			if err != nil {
				panic(err)
			}
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
		// 创建一个 HTTP 客户端，并将 Transport 对象设置为其 Transport 字段
		gptConfig.HTTPClient = &http.Client{
			Transport: transport,
		}

	}

	// 自定义gptConfig.BaseURL
	if cnf.ApiURL != "" {
		gptConfig.BaseURL = cnf.ApiURL
	}

	client := openai.NewClientWithConfig(gptConfig)
	if request.Messages[0].Role != "system" {
		newMessage := append([]openai.ChatCompletionMessage{
			{Role: "system", Content: cnf.BotDesc},
		}, request.Messages...)
		request.Messages = newMessage
	}

	request.Model = cnf.Model
	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return
	}
	defer stream.Close()

	reply := ""
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			js, _ := json.Marshal(gin.H{"data": "--_--xfsdkjfkjsdfjdksjfkdsjfksdjkfjsdkdjf"})
			fmt.Fprintf(ctx.Writer, "data: %s\n\n", string(js))
			ctx.Writer.Flush()
			break
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break
		}

		reply += response.Choices[0].Delta.Content
		js, _ := json.Marshal(gin.H{"data": response.Choices[0].Delta.Content})
		fmt.Fprintf(ctx.Writer, "data:%s\n\n", string(js))
		ctx.Writer.Flush()
	}

	chatCompletionMessage := openai.ChatCompletionMessage{Role: "assistant", Content: reply}
	answer := gin.H{
		"reply":    reply,
		"messages": append(request.Messages, chatCompletionMessage),
	}

	ajs, _ := json.Marshal(answer)
	question := fmt.Sprintf("question_%s_%s.json", request.Messages[len(request.Messages)-1].Content, util.GetCurrentTime().Format("20060102150405000"))
	subDir := fmt.Sprintf("chat/%s", authUser.Name)

	lfs.DataFs.SaveDataFile(question, ajs, subDir)

	return
}
