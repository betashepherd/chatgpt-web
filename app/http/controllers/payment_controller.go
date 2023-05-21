package controllers

import (
	"chatgpt-web/library/lfs"
	"chatgpt-web/library/util"
	"chatgpt-web/pkg/logger"
	"chatgpt-web/pkg/model/user"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"

	"github.com/betashepherd/xunhupay"
	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type PaymentController struct {
	BaseController
}

func NewPaymentController() *PaymentController {
	return &PaymentController{}
}

type PayRequest struct {
	Username string `binding:"required" json:"username"` // 商户订单号
	Plan     string `binding:"required" json:"plan"`     //订单支付金额
}

func (c *PaymentController) Pay(ctx *gin.Context) {
	var req PayRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	plans := map[string]string{"plan1": "5.00", "plan30": "30.00", "plan90": "90.00"}
	plansExpire := map[string]int64{"plan1": 24 * 3600, "plan30": 30 * 24 * 3600, "plan90": 90 * 24 * 3600}

	if _, ok := plans[req.Plan]; !ok {
		c.ResponseJson(ctx, http.StatusOK, "请选择套餐", nil)
		return
	}
	appId := "201906157182"                               //Appid
	appSecret := "35e08fa719b288dc8754af05f1700f78"       //密钥
	var host = "https://api.xunhupay.com/payment/do.html" //跳转支付页接口URL
	client := xunhupay.NewHuPi(&appId, &appSecret)        //初始化调用

	//支付参数，appid、time、nonce_str和hash这四个参数不用传，调用的时候执行方法内部已经处理
	params := map[string]string{
		"version":        "1.1",
		"trade_order_id": util.GetCurrentTime().Format("20060102150405000"),
		"total_fee":      plans[req.Plan],
		"title":          "VIP会员_" + req.Plan,
		"notify_url":     "https://ai.bgton.cn/payment/notify",
		"return_url":     "https://ai.bgton.cn",
		"attach":         util.Base64Encode([]byte(req.Username)),
	}

	execute, err := client.Execute(host, params) //执行支付操作
	if err != nil {
		question := fmt.Sprintf("payerr_%s.json", util.GetCurrentTime().Format("20060102150405000"))
		lfs.DataFs.SaveDataFile(question, []byte(err.Error()), "pay")
		c.ResponseJson(ctx, http.StatusOK, err.Error(), nil)
		return
	}

	type PayResponse struct {
		OpenId    int    `json:"openid"`
		UrlQrCode string `json:"url_qrcode"`
		Url       string `json:"url"`
		ErrCode   int    `json:"err_code"`
		ErrMsg    string `json:"err_msg"`
		Hash      string `json:"hash"`
	}

	var pr PayResponse
	err = json.Unmarshal([]byte(execute), &pr)

	if err != nil {
		question := fmt.Sprintf("payerr_%s.json", util.GetCurrentTime().Format("20060102150405000"))
		lfs.DataFs.SaveDataFile(question, []byte(err.Error()), "pay")
		c.ResponseJson(ctx, http.StatusOK, err.Error(), nil)
		return
	}

	if pr.ErrCode != 0 {
		c.ResponseJson(ctx, http.StatusOK, "pay api error, "+pr.ErrMsg, nil)
		return
	}

	question := fmt.Sprintf("paysucc_%s.json", util.GetCurrentTime().Format("20060102150405000"))
	lfs.DataFs.SaveDataFile(question, []byte(execute), "pay")

	data := gin.H{
		"username":   req.Username,
		"url_qrcode": pr.UrlQrCode,
		"url":        pr.Url,
	}

	nowTime := util.GetCurrentTime().Unix()
	expire := 3600 + plansExpire[req.Plan]
	ou, err := user.GetByName(req.Username)
	if err != nil && err == gorm.ErrRecordNotFound {
		//新增
		pwd, _ := util.NewPwd(10)
		data["password"] = pwd
		if _, err := user.InitUser(req.Username, pwd, req.Username, nowTime+expire); err != nil {
			logger.Info("create user error:", err)
			c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	} else {
		//续费
		if ou.ExpireTimestamp > nowTime {
			ou.ExpireTimestamp += expire
		} else {
			ou.ExpireTimestamp = nowTime + expire
		}
		ou.Save()
	}
	c.ResponseJson(ctx, http.StatusOK, "", data)
}

type PayNotifyForm struct {
	TradeOrderId  string `form:"trade_order_id" binding:"required"` //商户订单号
	TotalFee      string `form:"total_fee" binding:"required"`      //订单支付金额
	TransactionId string `form:"transaction_id"`                    //支付平台内部订单号
	OpenOrderId   string `form:"open_order_id"`                     //虎皮椒内部订单号
	OrderTitle    string `form:"order_title"`                       //订单标题
	Status        string `form:"status"`                            //订单状态 目前固定值为：OD
	Attach        string `from:"attach"`                            //附加信息
	AppId         string `form:"appid"`                             //支付渠道ID
	TimeStamp     string `form:"time"`                              //时间戳
	NonceStr      string `form:"nonce_str"`                         //随即字符串
	Hash          string `form:"hash"`                              //签名
}

func (c *PaymentController) Notify(ctx *gin.Context) {
	var req PayNotifyForm
	if err := ctx.ShouldBind(&req); err != nil {
		c.ResponseJson(ctx, http.StatusOK, err.Error(), nil)
		return
	}

	params := map[string]string{}
	params["trade_order_id"] = req.TradeOrderId
	params["total_fee"] = req.TotalFee
	params["transaction_id"] = req.TransactionId
	params["open_order_id"] = req.OpenOrderId
	params["order_title"] = req.OrderTitle
	params["status"] = req.Status
	attach, _ := util.Base64Decode(&req.Attach)
	params["attach"] = string(attach)
	params["appid"] = req.AppId
	params["time"] = req.TimeStamp
	params["nonce_str"] = req.NonceStr
	appId := "201906157182"                         //Appid
	appSecret := "35e08fa719b288dc8754af05f1700f78" //密钥
	client := xunhupay.NewHuPi(&appId, &appSecret)  //初始化调用
	reSign := client.Sign(params)

	if req.Hash != reSign {
		ctx.Writer.Write([]byte("sign error"))
		return
	}

	pjs, _ := json.Marshal(req)
	question := fmt.Sprintf("paynotify_%s.json", util.GetCurrentTime().Format("20060102150405000"))
	lfs.DataFs.SaveDataFile(question, pjs, "pay")

	ou, err := user.GetByName(params["attach"])
	if err != nil && err == gorm.ErrRecordNotFound {
		ctx.Writer.Write([]byte("fail"))
	} else {
		ou.Stat = 0 // 激活
		ou.Save()
	}

	ctx.Writer.Write([]byte("success"))
}
