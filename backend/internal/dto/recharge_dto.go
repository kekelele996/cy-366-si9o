package dto

// RechargeReq 会员充值请求。
type RechargeReq struct {
	UserID        uint    `json:"user_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=balance cash wechat alipay"`
	Remark        string  `json:"remark" binding:"omitempty,max=255"`
}

// BuyPackageReq 购买时长包请求。
type BuyPackageReq struct {
	PackageID     uint   `json:"package_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=balance cash wechat alipay"`
}
