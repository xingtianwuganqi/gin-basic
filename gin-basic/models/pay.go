package models

// 用户订阅
type Subscription struct {
	BaseModel
	UserId                uint   `json:"userId" form:"userId"`
	ProductId             string `json:"ProductId" form:"ProductId"`
	Platform              string `json:"platform" form:"platform" gorm:"size:20"`                                        // 平台'apple/google'
	OriginalTransactionId string `json:"originalTransactionId" form:"originalTransactionId" gorm:"size:200;uniqueIndex"` // 平台唯一订阅标识 -- Apple
	PurchaseToken         string `json:"purchaseToken" form:"purchaseToken" gorm:"size:200"`                             // -- Google
	ExpireTime            int64  `json:"expireTime" form:"expireTime"`                                                   // 到期时间
	Status                int    `json:"status" form:"status"`                                                           // 订阅状态 'active/expired/canceled'
	AutoRenew             int    `json:"autoRenew" form:"autoRenew" gorm:"default:1"`
}

type Transaction struct {
	BaseModel
	UserID                uint   `json:"userId" form:"userId"`
	ProductId             string `json:"productId" form:"productId"`
	Platform              string `json:"platform" form:"platform" gorm:"size:20;uniqueIndex:uk_tx"`
	TransactionID         string `json:"TransactionID" form:"TransactionID" gorm:"size:100;uniqueIndex:uk_tx"`
	OriginalTransactionID string `json:"originalTransactionId" form:"originalTransactionId" gorm:"size:100;index"`
	PurchaseToken         string `json:"purchaseToken" form:"purchaseToken" gorm:"size:200"` // -- Google
	PurchaseTime          int64  `json:"purchaseTime" form:"purchaseTime"`                   // 订阅开始时间
	ExpireTime            int64  `json:"expireTime" form:"expireTime"`                       // 到期时间

}

/*
CREATE TABLE subscription_usages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    period_start BIGINT NOT NULL,
    period_end BIGINT NOT NULL,
    used_count INT DEFAULT 0,
    limit_count INT DEFAULT 50,
    updated_at TIMESTAMP,
    UNIQUE KEY uk_user_period (user_id, period_start)
);
*/

type SubscriptionUsage struct {
	BaseModel
	UserId       uint  `json:"userId" form:"userId" gorm:"index"`
	PeriodStart  int64 `json:"periodStart" form:"periodStart" gorm:"index"` //当前周期开始
	PeriodEnd    int64 `json:"periodEnd" form:"periodEnd" gorm:"index"`     //当前周期结束
	UsedCount    int   `json:"usedCount" form:"usedCount"`                  //已用次数
	PendingCount int   `json:"pendingCount" form:"pendingCount"`            //预占次数
	LimitCount   int   `json:"limitCount" form:"limitCount"`                //总额度
	UsageType    uint  `json:"usageType" form:"usageType"`                  // 👈 核心字段 "free" 0 / "subscription" 1 / "extra" 2
}

const (
	UsageLedgerStatusReserved  = "reserved"
	UsageLedgerStatusCommitted = "committed"
	UsageLedgerStatusReleased  = "released"
)

type UsageLedger struct {
	BaseModel
	UserId        uint   `json:"userId" form:"userId" gorm:"index"`
	UsageId       uint   `json:"usageId" form:"usageId" gorm:"index"`
	AnalysisRunId uint   `json:"analysisRunId" form:"analysisRunId" gorm:"uniqueIndex"`
	UsageCost     int    `json:"usageCost" form:"usageCost"`
	UsageType     uint   `json:"usageType" form:"usageType"`
	Status        string `json:"status" form:"status" gorm:"size:20;index"` // reserved / committed / released
	ReservedAt    int64  `json:"reservedAt" form:"reservedAt"`
	SettledAt     int64  `json:"settledAt" form:"settledAt"`
}

type NotificationLog struct {
	ID   uint   `gorm:"primaryKey"`
	UUID string `gorm:"size:100;uniqueIndex"`
}

type UserCredit struct {
	BaseModel
	UserId       uint `json:"userId" form:"userId" gorm:"uniqueIndex"`
	TotalCredits int  `json:"totalCredits" form:"totalCredits"`
	UsedCredits  int  `json:"usedCredits" form:"usedCredits"`
}

type CreditLog struct {
	BaseModel
	UserId        uint   `json:"userId" form:"userId" gorm:"index"`
	TransactionId uint   `json:"transactionId" form:"transactionId" gorm:"index"`
	AnalysisRunId uint   `json:"analysisRunId" form:"analysisRunId" gorm:"index"`
	ChangeAmount  int    `json:"changeAmount" form:"changeAmount"`    // +10 / +20 / -1 / -2
	Reason        string `json:"reason" form:"reason" gorm:"size:50"` // purchase/analyze/refund
	BalanceAfter  int    `json:"balanceAfter" form:"balanceAfter"`
}
