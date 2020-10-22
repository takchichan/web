package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()
}

type UsrInfo struct {
	NickName string `json:"nickName"`
	Location string `json:"location"`
}

type ChatRequest struct {
	Peer      string `json:"peer"`
	Initiator bool   `json:"initiator"`
	Signal    string `json:"signal"`
	Action    int    `json:"action"`
}

var key []byte = []byte("passphrasewhichneedstobe32bytes!")

type redisHook struct {
}

const (
	ChatActionSendOffer = 1 + iota
	ChatActionSendAnswer
	ChatActionReceiveOffer
	ChatActionReceiveAnswer
	ChatActionConnected
	ChatActionTerminated
)

func actionName(action int) string {
	name := "unknow action"

	switch action {
	case ChatActionSendOffer:
		name = "sent out offer signal"
	case ChatActionSendAnswer:
		name = "sent out answer signal"
	case ChatActionReceiveOffer:
		name = "recived offer signal"
	case ChatActionReceiveAnswer:
		name = "recived answer signal"
	case ChatActionConnected:
		name = "chat connected"
	case ChatActionTerminated:
		name = "chat terminate"
	}

	return name
}

const (
	ChatStateInit = 1 + iota
	ChatStateWaitForReceivingAnswer
	ChatStateReceivedAnswer
	ChatStateWaitForSendingAnswer
	ChatStateSentAnswer
	ChatStateConnected
)

func stateName(state int) string {
	name := "unknow state"

	switch state {
	case ChatStateInit:
		name = "init state"
	case ChatStateWaitForReceivingAnswer:
		name = "wait for receiving answer signal state"
	case ChatStateWaitForSendingAnswer:
		name = "wait for sending answer signal state"
	case ChatStateReceivedAnswer:
		name = "received answer signal state"
	case ChatStateSentAnswer:
		name = "sent answer signal state"
	case ChatStateConnected:
		name = "connected state"
	}

	return name
}

type ChatInfo struct {
	Peer      string `json:"Peer"`
	State     int    `json:"State"`
	Offer     string `json:"Offer"`
	Answer    string `json:"Answer"`
	Initiator bool   `json:"Initiator"`
}

var rdb *redis.Client

func getChatInfo(ctx context.Context, usr string) (error, *ChatInfo) {
	strCmd := rdb.HGet(ctx, "chatInfo", usr)
	val, err := strCmd.Result()

	glog.Errorf("val: %s, err: %+v\n", val, err)

	if err != nil && err != redis.Nil {
		return err, nil
	}

	info := &ChatInfo{}

	if err == redis.Nil {
		info.State = ChatStateInit
		return nil, info
	}

	glog.Errorf("info: %+v\n", info)

	err = json.Unmarshal([]byte(val), info)

	return err, info
}

func setChatInfo(ctx context.Context, usr string, info *ChatInfo) error {
	val, _ := json.Marshal(info)
	cmd := rdb.HSet(ctx, "chatInfo", usr, string(val))

	ret, err := cmd.Result()

	glog.Errorf("ret: %d, err: %+v\n", ret, err)

	return err
}

func stateTransfer(action int, peer string, signal string, info *ChatInfo) (error, *ChatInfo) {
	glog.Errorf("%s, %s, info: %+v", actionName(action), peer, info)

	switch info.State {
	case ChatStateInit:
		{
			info.Peer = peer

			if action == ChatActionReceiveOffer {
				info.Offer = signal
				info.Initiator = false
				info.State = ChatStateWaitForSendingAnswer
			} else if action == ChatActionSendOffer {
				info.Offer = signal
				info.Initiator = true
				info.State = ChatStateWaitForReceivingAnswer
			} else {
				return fmt.Errorf("%s, while at %s", actionName(action), stateName(info.State)), nil
			}
		}
	case ChatStateWaitForReceivingAnswer:
		{
			if peer != info.Peer {
				return fmt.Errorf("peer not match, %s vs %s", peer, info.Peer), nil
			}

			if action == ChatActionReceiveAnswer {
				info.Answer = signal
				info.State = ChatStateReceivedAnswer
			} else {
				return fmt.Errorf("%s, while at %s", actionName(action), stateName(info.State)), nil
			}
		}
	case ChatStateWaitForSendingAnswer:
		{
			if peer != info.Peer {
				return fmt.Errorf("peer not match, %s vs %s", peer, info.Peer), nil
			}

			if action == ChatActionSendAnswer {
				info.Answer = signal
				info.State = ChatStateSentAnswer
			} else {
				return fmt.Errorf("%s, while at %s", actionName(action), stateName(info.State)), nil
			}
		}
	case ChatStateReceivedAnswer:
		{
			if peer != info.Peer {
				return fmt.Errorf("peer not match, %s vs %s", peer, info.Peer), nil
			}

			if action == ChatActionConnected {
				info.State = ChatStateConnected
			} else {
				return fmt.Errorf("%s, while at %s", actionName(action), stateName(info.State)), nil
			}
		}
	case ChatStateSentAnswer:
		{
			if peer != info.Peer {
				return fmt.Errorf("peer not match, %s vs %s", peer, info.Peer), nil
			}

			if action == ChatActionConnected {
				info.State = ChatStateConnected
			} else {
				return fmt.Errorf("%s, while at %s", actionName(action), stateName(info.State)), nil
			}
		}
	case ChatStateConnected:
		{
			if peer != info.Peer {
				return fmt.Errorf("peer not match, %s vs %s", peer, info.Peer), nil
			}

			if action == ChatActionTerminated {
				info.Peer = ""
				info.Offer = ""
				info.Answer = ""
				info.State = ChatStateInit
			} else {
				return fmt.Errorf("%s, while at %s", actionName(action), stateName(info.State)), nil
			}
		}
	}

	return nil, info
}

func (p redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	glog.Errorf("BeforeProcess: %+v", cmd)
	return ctx, nil
}

func (p redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	glog.Errorf("AfterProcess: %+v", cmd)
	return nil
}

func (p redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	glog.Errorf("BeforeProcessPipeline: %+v", cmds)
	return ctx, nil
}

func (p redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	glog.Errorf("AfterProcessPipeline: %+v", cmds)
	return nil
}

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: ":6379",
	})

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)

	// rdb.AddHook(redisHook{})

	rdb.Ping(ctx)
}

func GinUsrList(ctx *gin.Context) {
	usrList := make([]UsrInfo, 0)

	usrList = append(usrList, UsrInfo{
		NickName: "李师傅",
		Location: "水星",
	})

	usrList = append(usrList, UsrInfo{
		NickName: "陈师傅",
		Location: "水星",
	})

	usrList = append(usrList, UsrInfo{
		NickName: "唐师傅",
		Location: "水星",
	})

	ctx.JSON(200, usrList)
}

func encrypt(text string) (error, string) {
	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)

	// if there are any errors, handle them
	if err != nil {
		return err, ""
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)

	// if any error generating new GCM
	// handle them
	if err != nil {
		return err, ""
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())

	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err, ""
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.

	encryptedText := gcm.Seal(nonce, nonce, []byte(text), nil)

	return nil, base64.StdEncoding.EncodeToString(encryptedText)
}

func decrypt(encryptedText string) (error, string) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)

	if err != nil {
		return err, ""
	}

	c, err := aes.NewCipher(key)

	if err != nil {
		return err, ""
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return err, ""
	}

	nonceSize := gcm.NonceSize()

	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext length %d, while nonce size is %d", len(ciphertext), nonceSize), ""
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return err, ""
	}

	return nil, string(plaintext)
}

func GinChatStatus(ctx *gin.Context) {
	cookie, err := ctx.Cookie("interstellar")

	var info UsrInfo

	if err != nil {
		glog.Errorf("err: %+v", err)

		nickName := ctx.Query("nickName")
		location := ctx.Query("location")

		glog.Errorf("nickName: %s, location: %s", nickName, location)

		if nickName == "" || location == "" {
			ctx.JSON(500, gin.H{})
			return
		}

		info = UsrInfo{
			NickName: nickName,
			Location: location,
		}

		data, _ := json.Marshal(info)

		err, cookie = encrypt(string(data))

		if err != nil {
			ctx.JSON(500, gin.H{})
			return
		}

		ctx.SetCookie("interstellar", cookie, 3000, "/", "", false, true)
	} else {
		err, data := decrypt(cookie)

		if err != nil {
			ctx.JSON(500, gin.H{})
			return
		}

		err = json.Unmarshal([]byte(data), &info)

		if err != nil {
			ctx.JSON(500, gin.H{})
			return
		}

		glog.Errorf("nickName: %s, location: %s", info.NickName, info.Location)
	}

	err, chatInfo := getChatInfo(ctx, info.NickName)
	glog.Errorf("err: %+v, info: %+v\n", err, chatInfo)

	ctx.JSON(200, chatInfo)
}

func GinChatRequest(ctx *gin.Context) {
	cookie, err := ctx.Cookie("interstellar")

	if err != nil {
		ctx.JSON(500, gin.H{})
		return
	}

	err, data := decrypt(cookie)

	if err != nil {
		ctx.JSON(500, gin.H{})
		return
	}

	info := UsrInfo{}

	err = json.Unmarshal([]byte(data), &info)

	if err != nil {
		ctx.JSON(500, gin.H{})
		return
	}

	glog.Errorf("nickName: %s, location: %s", info.NickName, info.Location)

	req := ChatRequest{}

	err = ctx.BindJSON(&req)

	if err != nil {
		glog.Errorf("parse args failed, for: %s\n", err.Error())

		ctx.JSON(500, gin.H{
			"msg": err.Error(),
		})

		return
	}

	glog.Errorf("args: %+v\n", req)

	// TODO 同时锁住彼此
	// lock
	err, hostChatInfo := getChatInfo(ctx, info.NickName)
	err, hostChatInfo = stateTransfer(req.Action, req.Peer, req.Signal, hostChatInfo)
	glog.Errorf("err: %+v, info: %+v\n", err, hostChatInfo)

	err = setChatInfo(ctx, info.NickName, hostChatInfo)
	glog.Errorf("err: %+v, info: %+v\n", err, hostChatInfo)

	// 发出offer时，对方接受offer
	if req.Action == ChatActionSendOffer {
		glog.Errorf("set peer: %s", req.Peer)

		err, peerChatInfo := getChatInfo(ctx, req.Peer)
		err, peerChatInfo = stateTransfer(ChatActionReceiveOffer, info.NickName, req.Signal, peerChatInfo)
		glog.Errorf("err: %+v, info: %+v\n", err, peerChatInfo)

		err = setChatInfo(ctx, req.Peer, peerChatInfo)
		glog.Errorf("err: %+v, info: %+v\n", err, peerChatInfo)
	}

	// 发出answer时，对方接受answer
	if req.Action == ChatActionSendAnswer {
		glog.Errorf("set peer: %s", req.Peer)

		err, peerChatInfo := getChatInfo(ctx, req.Peer)
		err, peerChatInfo = stateTransfer(ChatActionReceiveAnswer, info.NickName, req.Signal, peerChatInfo)
		glog.Errorf("err: %+v, info: %+v\n", err, peerChatInfo)

		err = setChatInfo(ctx, req.Peer, peerChatInfo)
		glog.Errorf("err: %+v, info: %+v\n", err, peerChatInfo)
	}

	// unlock

	ctx.JSON(200, gin.H{})
}

func main() {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	corsConf := cors.DefaultConfig()
	// corsConf.AllowAllOrigins = true
	corsConf.AllowOrigins = []string{"http://localhost:8089", "http://localhost:8080"}
	corsConf.AllowCredentials = true
	router.Use(cors.New(corsConf))

	router.GET("/api/v1/usr", GinUsrList)
	router.GET("/api/v1/chat", GinChatStatus)
	router.POST("/api/v1/chat", GinChatRequest)

	glog.Info("Starting up")

	router.Run(":8081")

	glog.Info("Shutting down")
}
