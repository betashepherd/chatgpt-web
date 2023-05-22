import React, { useState } from "react";
import {Flex, FlexItem, Input, Button, toast, Modal, Image} from "@chatui/core";
import { useNavigate } from "react-router-dom";
import css from "../../App.module.css";
import "./index.less";
import "@chatui/core/dist/index.css";
import "@chatui/core/es/styles/index.less";
import "md-editor-rt/lib/style.css";
import {pay} from "../../services/port";
interface PaymentFormState {
  username: string,
  plan: string,
}
const Payment = () => {

  const navigate = useNavigate();

  const [paymentForm, setPaymentForm] = useState<PaymentFormState>({
    username: "",
    plan: "plan30",
  });

  const [paytip, setPaytip] = useState("");
  const [payimage, setPayimage] = useState("");

  const submitPayment = async () => {
    if (paymentForm.username == "" && paymentForm.plan == "") {
      return toast.show("请检查账号和套餐", undefined);
    }

    const res = await pay(paymentForm);
    if (res.data.code === 200) {
      console.log(res.data.data);
      var tip = "用户名：(" + res.data.data.username + ')，';
      if (res.data.data.password != undefined && res.data.data.password != "") {
        tip += "密码：(" + res.data.data.password + ") ，请牢记。";
      }
      tip += '支付宝扫下面二维码，开启AI之旅(等待1分钟左右)。:)'+ "\n";
      setPaytip(tip);
      setPayimage(res.data.data.url_qrcode);
    }
  };

  const handleInputChange = (event: any) => {
    const { name, value } = event.target;
    var sn ;
    if (name == 'plan1' || name == 'plan30' || name == 'plan90') {
      sn = 'plan';
    } else {
      sn = name;
    }
    setPaymentForm({ ...paymentForm, [sn]: value });
  };

  const goLogin = async () => {
    navigate("/login");
  };

  return (
    <div className={css.app}>
      <Flex center direction={"column"} style={{ background: "var(--gray-7)" }}>
        <FlexItem
          flex={"1"}
          style={{ marginLeft: "1em" }}
          className="form-Item"
        >
          <div className="login-header">
            <h3>开通VIP会员</h3>
          </div>
        </FlexItem>
        <FlexItem
          flex={"6"}
          style={{ marginLeft: "1em" }}
          className="form-Item"
        >
          <div className={css.m_top}>手机号：
            <input
                className="input-item"
                type="text"
                name="username"
                id="username"
                value={paymentForm.username}
                onChange={(e:any) => handleInputChange( e)}
                placeholder="请输入账号"
            />
          </div>
          <div className={css.m_top}>套餐：
              <input type="radio" value="plan30" name="plan30" id="plan30" checked={paymentForm.plan === 'plan30'} onChange={(e:any) => handleInputChange( e)} /> 30天（推荐 RMB 30.00）
              <input type="radio" value="plan90" name="plan90" id="plan90" checked={paymentForm.plan === 'plan90'} onChange={(e:any) => handleInputChange( e)} /> 90天（RMB 90.00）
              <input type="radio" value="plan1" name="plan1" id="plan1" checked={paymentForm.plan === 'plan1'} onChange={(e:any) => handleInputChange( e)} /> 1天（体验 RMB 5.00）
          </div>
          <div className={css.m_top}>
            <Button color="primary" size={"md"} onClick={() => submitPayment()}>
              开通
            </Button>
            <Button color="primary" size={"sm"} onClick={() => goLogin()}>
              返回
            </Button>
          </div>

          <div className={css.m_top}>
            {paytip}<br/>
            <Image src={payimage} id="payimage" lazy={true} />
          </div>
        </FlexItem>
      </Flex>



    </div>
  );
};

export default Payment;
