package service

import (
	"errors"
	"time"

	"gin-basic/db"
	"gin-basic/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrInsufficientCredits = errors.New("insufficient credits")

var productCreditsMap = map[string]int{
	"20004": 10,
	"20003": 20,
}

func GetCreditsForProduct(productID string) int {
	return productCreditsMap[productID]
}

func IsConsumableProduct(productID string) bool {
	_, ok := productCreditsMap[productID]
	return ok
}

func getOrCreateUserCreditForUpdate(tx *gorm.DB, userID uint) (*models.UserCredit, error) {
	var credit models.UserCredit
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&credit).Error
	if err == nil {
		return &credit, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	credit = models.UserCredit{UserId: userID}
	if err := tx.Create(&credit).Error; err != nil {
		return nil, err
	}
	return &credit, nil
}

func GrantCredits(tx *gorm.DB, userID uint, transactionID uint, amount int) error {
	if amount <= 0 {
		return nil
	}

	credit, err := getOrCreateUserCreditForUpdate(tx, userID)
	if err != nil {
		return err
	}

	newTotal := credit.TotalCredits + amount
	if err := tx.Model(credit).Update("total_credits", newTotal).Error; err != nil {
		return err
	}

	return tx.Create(&models.CreditLog{
		UserId:        userID,
		TransactionId: transactionID,
		ChangeAmount:  amount,
		Reason:        "purchase",
		BalanceAfter:  newTotal - credit.UsedCredits,
	}).Error
}

func ConsumeCredits(userID uint, amount int) error {
	if amount <= 0 {
		return nil
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		return reserveCreditsForAnalysisTx(tx, userID, 0, amount)
	})
}

func RefundCredits(userID uint, amount int) error {
	if amount <= 0 {
		return nil
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		return refundCreditsForAnalysisTx(tx, userID, 0, amount)
	})
}

func reserveCreditsForAnalysisTx(tx *gorm.DB, userID uint, analysisRunID uint, amount int) error {
	credit, err := getOrCreateUserCreditForUpdate(tx, userID)
	if err != nil {
		return err
	}

	remaining := credit.TotalCredits - credit.UsedCredits
	if remaining < amount {
		return ErrInsufficientCredits
	}

	newUsed := credit.UsedCredits + amount
	if err := tx.Model(credit).Update("used_credits", newUsed).Error; err != nil {
		return err
	}

	return tx.Create(&models.CreditLog{
		UserId:        userID,
		AnalysisRunId: analysisRunID,
		ChangeAmount:  -amount,
		Reason:        "analyze",
		BalanceAfter:  credit.TotalCredits - newUsed,
	}).Error
}

func refundCreditsForAnalysisTx(tx *gorm.DB, userID uint, analysisRunID uint, amount int) error {
	credit, err := getOrCreateUserCreditForUpdate(tx, userID)
	if err != nil {
		return err
	}

	newUsed := credit.UsedCredits - amount
	if newUsed < 0 {
		newUsed = 0
	}
	if err := tx.Model(credit).Update("used_credits", newUsed).Error; err != nil {
		return err
	}

	return tx.Create(&models.CreditLog{
		UserId:        userID,
		AnalysisRunId: analysisRunID,
		ChangeAmount:  amount,
		Reason:        "refund",
		BalanceAfter:  credit.TotalCredits - newUsed,
	}).Error
}

func GetUserCreditSummary(userID uint) (*models.UserCredit, error) {
	var credit models.UserCredit
	err := db.DB.Where("user_id = ?", userID).First(&credit).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.UserCredit{UserId: userID}, nil
	}
	return &credit, err
}

func GetCreditLogs(userID uint, pageNum, pageSize int) ([]models.CreditLog, int64, error) {
	var logs []models.CreditLog
	var total int64

	query := db.DB.Model(&models.CreditLog{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&logs).Error
	return logs, total, err
}

type PurchaseLogItem struct {
	ID            uint   `json:"id"`
	CreatedAt     int64  `json:"created_at"`
	ChangeAmount  int    `json:"change_amount"`
	BalanceAfter  int    `json:"balance_after"`
	ProductId     string `json:"product_id"`
	TransactionId string `json:"transaction_id"`
}

func GetUserPurchaseLogs(userID uint, pageNum, pageSize int) ([]PurchaseLogItem, int64, error) {
	var total int64
	if err := db.DB.Model(&models.CreditLog{}).
		Where("user_id = ? AND reason = ?", userID, "purchase").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		ID            uint      `gorm:"column:id"`
		CreatedAt     time.Time `gorm:"column:created_at"`
		ChangeAmount  int       `gorm:"column:change_amount"`
		BalanceAfter  int       `gorm:"column:balance_after"`
		ProductId     string    `gorm:"column:product_id"`
		TransactionNo string    `gorm:"column:transaction_no"`
	}
	var rows []row
	err := db.DB.Table("credit_logs").
		Select("credit_logs.id, credit_logs.created_at, credit_logs.change_amount, credit_logs.balance_after, COALESCE(transactions.product_id,'') AS product_id, COALESCE(transactions.transaction_id,'') AS transaction_no").
		Joins("LEFT JOIN transactions ON transactions.id = credit_logs.transaction_id").
		Where("credit_logs.user_id = ? AND credit_logs.reason = ?", userID, "purchase").
		Order("credit_logs.created_at DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]PurchaseLogItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, PurchaseLogItem{
			ID:            row.ID,
			CreatedAt:     row.CreatedAt.UnixMilli(),
			ChangeAmount:  row.ChangeAmount,
			BalanceAfter:  row.BalanceAfter,
			ProductId:     row.ProductId,
			TransactionId: row.TransactionNo,
		})
	}
	return items, total, nil
}

func GetAllUserCredits(pageNum, pageSize int) ([]models.UserCredit, int64, error) {
	var credits []models.UserCredit
	var total int64

	query := db.DB.Model(&models.UserCredit{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&credits).Error
	return credits, total, err
}

func SendMembershipActivatedPush(userID uint, productID string) error {
	return nil
}
