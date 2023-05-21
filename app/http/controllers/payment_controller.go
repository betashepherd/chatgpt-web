package controllers

import (
	"chatgpt-web/library/lfs"
	"chatgpt-web/library/util"
	"encoding/json"
	"fmt"
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

func (c *PaymentController) Pay(ctx *gin.Context) {
	appId := "201906157182"                               //Appid
	appSecret := "35e08fa719b288dc8754af05f1700f78"       //密钥
	var host = "https://api.xunhupay.com/payment/do.html" //跳转支付页接口URL
	client := xunhupay.NewHuPi(&appId, &appSecret)        //初始化调用

	//支付参数，appid、time、nonce_str和hash这四个参数不用传，调用的时候执行方法内部已经处理
	params := map[string]string{
		"version":        "1.1",
		"trade_order_id": "15652936798_" + util.GetCurrentTime().Format("20060102150405"),
		"total_fee":      "0.1",
		"title":          "VIP会员 - 30天 - 测试",
		"notify_url":     "https://ai.bgton.cn/payment/notify",
		"return_url":     "https://ai.bgton.cn",
		"attach":         "15652936798",
	}

	execute, err := client.Execute(host, params) //执行支付操作
	if err != nil {
		panic(err)
	}
	fmt.Println(execute) //打印支付结果
	c.ResponseJson(ctx, http.StatusOK, "", execute)
}

type PayNotifyForm struct {
	TradeOrderId  string  `form:"trade_order_id" binding:"required"` // 商户订单号
	TotalFee      float64 `form:"total_fee" binding:"required"`      //订单支付金额
	TransactionId string  `form:"transaction_id"`                    //支付平台内部订单号
	OpenOrderId   string  `form:"open_order_id"`                     //虎皮椒内部订单号
	OrderTitle    string  `form:"order_title"`                       //订单标题
	Status        string  `form:"status"`                            // 订单状态 目前固定值为：OD
	Attach        string  `from:"attach"`                            // 附加信息
	AppId         string  `form:"appid"`                             // 支付渠道ID
	TimeStamp     string  `form:"time"`                              // 时间戳
	NonceStr      string  `form:"nonce_str"`                         //随即字符串
	Hash          string  `form:"hash"`                              //签名
}

func (c *PaymentController) Notify(ctx *gin.Context) {
	pjs, _ := json.Marshal(ctx.Request.Body)
	question := fmt.Sprintf("paynotify_%s_%s.json", pjs, util.GetCurrentTime().Format("20060102150405000"))
	lfs.DataFs.SaveDataFile(question, pjs, "pay")
	var req PayNotifyForm
	if err := ctx.ShouldBind(&req); err != nil {
		c.ResponseJson(ctx, http.StatusOK, err.Error(), nil)
		return
	}
	c.ResponseJson(ctx, http.StatusOK, "", nil)
}
