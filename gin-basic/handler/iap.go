package handler

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"gin-basic/db"
	"gin-basic/logger"
	"gin-basic/models"
	"gin-basic/response"
	"gin-basic/service"
	"gin-basic/settings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ResponseStatusMsg = map[int64]string{
	21000: "App Store 的请求不是使用 HTTP POST 请求方法发出的。",
	21001: "App Store 不再发送此状态代码。",
	21002: "属性中的数据receipt-data格式不正确或服务遇到临时问题。",
	21003: "无法验证收据。",
	21004: "您提供的共享密钥与您帐户的文件共享密钥不匹配。",
	21005: "收据服务器暂时无法提供收据。",
	21006: "此收据有效，但订阅已过期。",
	21007: "这条回执是来自测试环境，但它是发送到生产环境进行验证的。",
	21008: "这条回执来自生产环境，但它被发送到测试环境进行验证。",
	21009: "内部数据访问错误。",
	21010: "用户帐户找不到或已被删除。",
}

type AppStoreApi struct {
	kid         string
	iss         string
	bid         string
	privatePath string
	token       string
	baseApiHost string
}

type AppleVerifyRequest struct {
	TransactionJWS string `json:"transaction_jws" form:"transaction_jws" binding:"required"`
	UserID         uint   `json:"user_id" form:"user_id"`
	ProductID      string `json:"product_id" form:"product_id"`
}

type InAppLookupResp struct {
	AppAccountToken       string `json:"appAccountToken"`
	BundleID              string `json:"bundleId"`
	Currency              string `json:"currency"`
	Environment           string `json:"environment"`
	ExpiresDate           int64  `json:"expiresDate"`
	InAppOwnershipType    string `json:"inAppOwnershipType"`
	IsUpgraded            bool   `json:"isUpgraded"`
	OfferDiscountType     string `json:"offerDiscountType"`
	OfferIdentifier       string `json:"offerIdentifier"`
	OfferPeriod           string `json:"offerPeriod"`
	OfferType             int    `json:"offerType"`
	OriginalPurchaseDate  int64  `json:"originalPurchaseDate"`
	OriginalTransactionID string `json:"originalTransactionId"`
	Price                 int64  `json:"price"`
	ProductID             string `json:"productId"`
	PurchaseDate          int64  `json:"purchaseDate"`
	Quantity              int    `json:"quantity"`
	RevocationDate        int64  `json:"revocationDate"`
	RevocationReason      int    `json:"revocationReason"`
	SignedDate            int64  `json:"signedDate"`
	Storefront            string `json:"storefront"`
	StorefrontID          string `json:"storefrontId"`
	SubscriptionGroupID   string `json:"subscriptionGroupIdentifier"`
	TransactionID         string `json:"transactionId"`
	TransactionReason     string `json:"transactionReason"`
	Type                  string `json:"type"`
	WebOrderLineItemID    string `json:"webOrderLineItemId"`
}

type JWSRenewalInfoResp struct {
	AutoRenewProductID        string `json:"autoRenewProductId"`
	AutoRenewStatus           int    `json:"autoRenewStatus"`
	Currency                  string `json:"currency"`
	Environment               string `json:"environment"`
	ExpirationIntent          int    `json:"expirationIntent"`
	GracePeriodExpiresDate    int64  `json:"gracePeriodExpiresDate"`
	IsInBillingRetryPeriod    bool   `json:"isInBillingRetryPeriod"`
	OfferIdentifier           string `json:"offerIdentifier"`
	OfferType                 int    `json:"offerType"`
	OriginalTransactionID     string `json:"originalTransactionId"`
	PriceIncreaseStatus       int    `json:"priceIncreaseStatus"`
	ProductID                 string `json:"productId"`
	RecentSubscriptionStartAt int64  `json:"recentSubscriptionStartDate"`
	RenewalDate               int64  `json:"renewalDate"`
	SignedDate                int64  `json:"signedDate"`
}

type InAppTransactionsResp struct {
	OriginalTransactionID string              `json:"originalTransactionId"`
	Status                int                 `json:"status"`
	SignedTransactionInfo string              `json:"signedTransactionInfo"`
	SignedRenewalInfo     string              `json:"signedRenewalInfo"`
	TransactionInfo       *InAppLookupResp    `json:"transactionInfo,omitempty"`
	RenewalInfo           *JWSRenewalInfoResp `json:"renewalInfo,omitempty"`
}

type SubscriptionGroupIdentifierItem struct {
	SubscriptionGroupIdentifier string                  `json:"subscriptionGroupIdentifier"`
	LastTransactions            []InAppTransactionsResp `json:"lastTransactions"`
}

type AppleVerifyResult struct {
	IsPro                       bool                              `json:"is_pro"`
	Environment                 string                            `json:"environment,omitempty"`
	ProductID                   string                            `json:"productId,omitempty"`
	TransactionID               string                            `json:"transactionId,omitempty"`
	OriginalTransactionID       string                            `json:"originalTransactionId,omitempty"`
	SubscriptionGroupIdentifier string                            `json:"subscriptionGroupIdentifier,omitempty"`
	ExpireTime                  int64                             `json:"expire_time,omitempty"`
	ClientTransaction           *InAppLookupResp                  `json:"client_transaction,omitempty"`
	TransactionInfo             *InAppLookupResp                  `json:"transaction_info,omitempty"`
	LatestTransaction           *InAppTransactionsResp            `json:"latest_transaction,omitempty"`
	SubscriptionGroups          []SubscriptionGroupIdentifierItem `json:"subscriptionGroups,omitempty"`
}

type AppleNotificationRequest struct {
	SignedPayload string `json:"signedPayload" form:"signedPayload" binding:"required"`
}

type AppleNotificationData struct {
	AppAppleID            int64  `json:"appAppleId"`
	BundleID              string `json:"bundleId"`
	BundleVersion         string `json:"bundleVersion"`
	Environment           string `json:"environment"`
	SignedRenewalInfo     string `json:"signedRenewalInfo"`
	SignedTransactionInfo string `json:"signedTransactionInfo"`
	Status                int    `json:"status"`
}

type AppleNotificationPayload struct {
	NotificationType string                 `json:"notificationType"`
	Subtype          string                 `json:"subtype"`
	NotificationUUID string                 `json:"notificationUUID"`
	Version          string                 `json:"version"`
	SignedDate       int64                  `json:"signedDate"`
	Data             *AppleNotificationData `json:"data,omitempty"`
}

type AppleNotificationResult struct {
	Verified         bool                      `json:"verified"`
	NotificationType string                    `json:"notificationType,omitempty"`
	Subtype          string                    `json:"subtype,omitempty"`
	NotificationUUID string                    `json:"notificationUUID,omitempty"`
	Environment      string                    `json:"environment,omitempty"`
	BundleID         string                    `json:"bundleId,omitempty"`
	SignedDate       int64                     `json:"signedDate,omitempty"`
	Payload          *AppleNotificationPayload `json:"payload,omitempty"`
	TransactionInfo  *InAppLookupResp          `json:"transactionInfo,omitempty"`
	RenewalInfo      *JWSRenewalInfoResp       `json:"renewalInfo,omitempty"`
}

type appleJWSHeader struct {
	Alg string   `json:"alg"`
	X5C []string `json:"x5c"`
}

func newAppStoreAPIFromEnv() (*AppStoreApi, error) {
	api := &AppStoreApi{
		kid:         settings.Conf.AppleKeys.AppleIapKeyId,
		iss:         settings.Conf.AppleKeys.AppleIapIssuerId,
		bid:         settings.Conf.AppleKeys.AppleBundleId,
		privatePath: settings.Conf.AppleKeys.AppleIapPrivateKeyPath,
		baseApiHost: settings.Conf.AppleKeys.AppleIapBaseUrl,
	}

	if api.baseApiHost == "" {
		api.baseApiHost = "https://api.storekit.itunes.apple.com"
	}

	if api.kid == "" || api.iss == "" || api.bid == "" || api.privatePath == "" {
		return nil, errors.New("apple iap config missing")
	}

	return api, nil
}

func (a *AppStoreApi) getToken() (string, error) {
	if a.token != "" {
		return a.token, nil
	}

	token := &jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"kid": a.kid,
			"alg": jwt.SigningMethodES256.Alg(),
		},
		Claims: jwt.MapClaims{
			"iss": a.iss,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour).Unix(),
			"aud": "appstoreconnect-v1",
			"bid": a.bid,
		},
		Method: jwt.SigningMethodES256,
	}

	privatePEM, err := os.ReadFile(a.privatePath)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(privatePEM)
	if block == nil {
		return "", errors.New("token: AuthKey must be a valid .p8 PEM file")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	ecdsaKey, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return "", errors.New("token: AuthKey must be of type ecdsa.PrivateKey")
	}

	a.token, err = token.SignedString(ecdsaKey)
	if err != nil {
		return "", err
	}

	return a.token, nil
}

func (a *AppStoreApi) GetTransactionInfo(transactionID string) (*InAppLookupResp, error) {
	apiURI := fmt.Sprintf("%s/inApps/v1/transactions/%s", a.baseApiHost, transactionID)
	req, err := http.NewRequest(http.MethodGet, apiURI, nil)
	if err != nil {
		return nil, err
	}

	var apiResp struct {
		SignedTransactionInfo string `json:"signedTransactionInfo"`
	}
	if err = a.doReq(req, &apiResp); err != nil {
		return nil, err
	}
	if apiResp.SignedTransactionInfo == "" {
		return nil, errors.New("apple transaction info not found")
	}

	return decodeJWSPayload[InAppLookupResp](apiResp.SignedTransactionInfo)
}

func (a *AppStoreApi) GetSubscriptionsStatus(originalTransactionID string) ([]SubscriptionGroupIdentifierItem, error) {
	apiURI := fmt.Sprintf("%s/inApps/v1/subscriptions/%s", a.baseApiHost, originalTransactionID)
	req, err := http.NewRequest(http.MethodGet, apiURI, nil)
	if err != nil {
		return nil, err
	}

	var apiResp struct {
		AppAppleID  int64                             `json:"appAppleId"`
		Environment string                            `json:"environment"`
		BundleID    string                            `json:"bundleId"`
		Data        []SubscriptionGroupIdentifierItem `json:"data"`
	}
	if err = a.doReq(req, &apiResp); err != nil {
		return nil, err
	}

	for i := range apiResp.Data {
		for j := range apiResp.Data[i].LastTransactions {
			tx := &apiResp.Data[i].LastTransactions[j]
			if tx.SignedTransactionInfo != "" {
				txInfo, decodeErr := decodeJWSPayload[InAppLookupResp](tx.SignedTransactionInfo)
				if decodeErr == nil {
					tx.TransactionInfo = txInfo
				}
			}
			if tx.SignedRenewalInfo != "" {
				renewalInfo, decodeErr := decodeJWSPayload[JWSRenewalInfoResp](tx.SignedRenewalInfo)
				if decodeErr == nil {
					tx.RenewalInfo = renewalInfo
				}
			}
		}
	}

	return apiResp.Data, nil
}

func (a *AppStoreApi) doReq(req *http.Request, out interface{}) error {
	token, err := a.getToken()
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		if i == 2 {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("app store api status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return json.Unmarshal(body, out)
}

func decodeJWSPayload[T any](signed string) (*T, error) {
	parts := strings.Split(signed, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid jws payload")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var out T
	if err = json.Unmarshal(payload, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func decodeJWSHeader(signed string) (*appleJWSHeader, error) {
	parts := strings.Split(signed, ".")
	if len(parts) < 3 {
		return nil, errors.New("invalid jws header")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	var header appleJWSHeader
	if err = json.Unmarshal(headerBytes, &header); err != nil {
		return nil, err
	}
	if header.Alg != jwt.SigningMethodES256.Alg() {
		return nil, errors.New("unexpected jws alg")
	}
	if len(header.X5C) == 0 {
		return nil, errors.New("missing x5c certificate chain")
	}

	return &header, nil
}

func getLeafCertificateFromHeader(header *appleJWSHeader) (*x509.Certificate, error) {
	leafDER, err := base64.StdEncoding.DecodeString(header.X5C[0])
	if err != nil {
		return nil, err
	}
	leafCert, err := x509.ParseCertificate(leafDER)
	if err != nil {
		return nil, err
	}
	return leafCert, nil
}

func verifyAndDecodeJWS[T any](signed string) (*T, error) {
	header, err := decodeJWSHeader(signed)
	if err != nil {
		return nil, err
	}

	leafCert, err := getLeafCertificateFromHeader(header)
	if err != nil {
		return nil, err
	}

	if _, err = jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodES256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return leafCert.PublicKey, nil
	}); err != nil {
		return nil, err
	}

	return decodeJWSPayload[T](signed)
}

func millisToTime(ms int64) time.Time {
	if ms <= 0 {
		return time.Time{}
	}
	return time.UnixMilli(ms)
}

func latestSubscriptionTransaction(groups []SubscriptionGroupIdentifierItem) *InAppTransactionsResp {
	var candidates []*InAppTransactionsResp
	for i := range groups {
		for j := range groups[i].LastTransactions {
			tx := &groups[i].LastTransactions[j]
			candidates = append(candidates, tx)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		left := int64(0)
		right := int64(0)
		if candidates[i].TransactionInfo != nil {
			left = candidates[i].TransactionInfo.ExpiresDate
			if left == 0 {
				left = candidates[i].TransactionInfo.SignedDate
			}
		}
		if candidates[j].TransactionInfo != nil {
			right = candidates[j].TransactionInfo.ExpiresDate
			if right == 0 {
				right = candidates[j].TransactionInfo.SignedDate
			}
		}
		return left > right
	})

	return candidates[0]
}

func verifyApplePurchase(req AppleVerifyRequest) (*AppleVerifyResult, error) {
	api, err := newAppStoreAPIFromEnv()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	clientTransaction, err := decodeJWSPayload[InAppLookupResp](req.TransactionJWS)
	if err != nil {
		return nil, err
	}
	logger.Logger.Info("verify apple purchase environment", zap.Any("client_transaction environment", clientTransaction.Environment))
	// 设置请求环境
	if clientTransaction.Environment == "Sandbox" {
		api.baseApiHost = settings.Conf.AppleKeys.AppleIapSandboxBaseUrl
	} else {
		api.baseApiHost = settings.Conf.AppleKeys.AppleIapBaseUrl
	}

	if clientTransaction.TransactionID == "" {
		return nil, errors.New("transactionId not found in transaction_jws")
	}

	transactionInfo, err := api.GetTransactionInfo(clientTransaction.TransactionID)
	if err != nil {
		return nil, err
	}
	if transactionInfo.OriginalTransactionID == "" {
		return nil, errors.New("originalTransactionId not found in transaction info")
	}

	subGroups, err := api.GetSubscriptionsStatus(transactionInfo.OriginalTransactionID)
	if err != nil {
		return nil, err
	}

	latestTransaction := latestSubscriptionTransaction(subGroups)
	if latestTransaction == nil || latestTransaction.TransactionInfo == nil {
		return nil, errors.New("latest subscription transaction not found")
	}

	expireTime := latestTransaction.TransactionInfo.ExpiresDate
	return &AppleVerifyResult{
		IsPro:                       expireTime > now.UnixMilli(),
		Environment:                 transactionInfo.Environment,
		ProductID:                   transactionInfo.ProductID,
		TransactionID:               transactionInfo.TransactionID,
		OriginalTransactionID:       transactionInfo.OriginalTransactionID,
		SubscriptionGroupIdentifier: latestTransaction.TransactionInfo.SubscriptionGroupID,
		ExpireTime:                  expireTime,
		ClientTransaction:           clientTransaction,
		TransactionInfo:             transactionInfo,
		LatestTransaction:           latestTransaction,
		SubscriptionGroups:          subGroups,
	}, nil
}

// verifyCreditPackage 处理次数包购买验证（Consumable 类型产品）
func verifyCreditPackage(c *gin.Context, req AppleVerifyRequest) {
	api, err := newAppStoreAPIFromEnv()
	if err != nil {
		logger.Logger.Error("apple iap config missing", zap.Error(err))
		response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
		return
	}

	clientTransaction, err := decodeJWSPayload[InAppLookupResp](req.TransactionJWS)
	if err != nil {
		logger.Logger.Error("decode credit pack JWS failed", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	if clientTransaction.Environment == "Sandbox" {
		api.baseApiHost = settings.Conf.AppleKeys.AppleIapSandboxBaseUrl
	} else {
		api.baseApiHost = settings.Conf.AppleKeys.AppleIapBaseUrl
	}

	if clientTransaction.TransactionID == "" {
		response.Fail(c, response.ApiCode.ParamErr, "transactionId missing")
		return
	}

	transactionInfo, err := api.GetTransactionInfo(clientTransaction.TransactionID)
	if err != nil {
		logger.Logger.Error("GetTransactionInfo failed", zap.Error(err))
		response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
		return
	}

	amount := service.GetCreditsForProduct(transactionInfo.ProductID)
	if amount <= 0 {
		logger.Logger.Warn("unknown credit product", zap.String("productId", transactionInfo.ProductID))
		response.Fail(c, response.ApiCode.ParamErr, "unknown credit product")
		return
	}

	// 幂等：在事务内写 Transaction + GrantCredits
	var creditGranted int
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		var existTx models.Transaction
		findErr := tx.Where("platform = ? AND transaction_id = ?", "apple", transactionInfo.TransactionID).
			First(&existTx).Error

		if findErr == nil {
			// 已处理过，幂等返回
			creditGranted = 0
			return nil
		}
		if !errors.Is(findErr, gorm.ErrRecordNotFound) {
			return findErr
		}

		newTx := models.Transaction{
			UserID:                req.UserID,
			ProductId:             transactionInfo.ProductID,
			Platform:              "apple",
			TransactionID:         transactionInfo.TransactionID,
			OriginalTransactionID: transactionInfo.OriginalTransactionID,
			PurchaseTime:          time.Now().Unix(),
		}
		if err := tx.Create(&newTx).Error; err != nil {
			return err
		}

		if err := service.GrantCredits(tx, req.UserID, newTx.ID, amount); err != nil {
			return err
		}
		creditGranted = amount
		return nil
	})
	if err != nil {
		logger.Logger.Error("credit pack grant failed", zap.Error(err))
		response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
		return
	}

	logger.Logger.Info("credit pack verified",
		zap.Uint("userId", req.UserID),
		zap.String("productId", transactionInfo.ProductID),
		zap.Int("creditGranted", creditGranted))

	response.Success(c, gin.H{
		"credit_granted": creditGranted,
	})
}

func VerisfyAppleIAP(c *gin.Context) {
	var req AppleVerifyRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	if value, ok := c.Get("userId"); ok {
		if userID, ok := value.(uint); ok && userID > 0 {
			req.UserID = userID
		}
	}
	if req.UserID == 0 {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	// 从 JWS payload 解码出 productId，避免依赖客户端传值
	if req.ProductID == "" {
		clientTx, err := decodeJWSPayload[InAppLookupResp](req.TransactionJWS)
		if err == nil && clientTx != nil {
			req.ProductID = clientTx.ProductID
		}
	}

	// 次数包产品走独立验证流程
	if service.IsConsumableProduct(req.ProductID) {
		verifyCreditPackage(c, req)
		return
	}

	result, err := verifyApplePurchase(req)
	if err != nil {
		logger.Logger.Error("Verify Apple IAP failed",
			zap.String("transactionJWS", req.TransactionJWS),
			zap.Error(err))
		response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
		return
	}

	var sub models.Subscription
	err = db.DB.
		Where("original_transaction_id = ?", result.OriginalTransactionID).
		First(&sub).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("query subscription failed", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	if err == nil && sub.UserId == 0 {
		// ✅ 只有查到才处理
		if updateErr := db.DB.Model(&sub).Update("user_id", req.UserID).Error; updateErr != nil {
			logger.Logger.Error("bind subscription user failed", zap.Error(updateErr))
			response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
			return
		}
	}

	now := time.Now().UnixMilli()

	status := 1 // active
	if result.ExpireTime < now {
		status = 2 // expired
	}

	if err == nil {
		// ✅ 防乱序（核心）
		if result.ExpireTime > sub.ExpireTime {
			if updateErr := db.DB.Model(&sub).Updates(models.Subscription{
				ProductId:  result.ProductID,
				ExpireTime: result.ExpireTime,
				Status:     status,
			}).Error; updateErr != nil {
				logger.Logger.Error("update subscription failed", zap.Error(updateErr))
				response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
				return
			}
		}
	} else {
		// 新建
		if createErr := db.DB.Create(&models.Subscription{
			UserId:                req.UserID,
			ProductId:             result.ProductID,
			Platform:              "apple",
			OriginalTransactionId: result.OriginalTransactionID,
			ExpireTime:            result.ExpireTime,
			Status:                status,
			AutoRenew:             result.LatestTransaction.RenewalInfo.AutoRenewStatus,
		}).Error; createErr != nil {
			logger.Logger.Error("create subscription failed", zap.Error(createErr))
			response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
			return
		}
	}

	var tx models.Transaction
	err = db.DB.
		Where("transaction_id = ?", result.TransactionID).
		First(&tx).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("query transaction failed", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if createErr := db.DB.Create(&models.Transaction{
			UserID:                req.UserID,
			ProductId:             result.ProductID,
			Platform:              "apple",
			TransactionID:         result.TransactionID,
			OriginalTransactionID: result.OriginalTransactionID,
			PurchaseTime:          time.Now().Unix(),
			ExpireTime:            result.ExpireTime,
		}).Error; createErr != nil {
			logger.Logger.Error("create transaction failed", zap.Error(createErr))
			response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
			return
		}
	}

	// 5️⃣ 返回
	isPro := result.ExpireTime > time.Now().UnixMilli()
	logger.Logger.Info("Verify Apple IAP success",
		zap.String("transactionJWS", req.TransactionJWS),
		zap.Bool("isPro", isPro),
		zap.Int64("expireTime", result.ExpireTime))

	if isPro {
		go func(userID uint, productID string) {
			_ = service.SendMembershipActivatedPush(userID, productID)
		}(req.UserID, result.ProductID)
	}

	response.Success(c, gin.H{
		"is_pro":      isPro,
		"expire_time": result.ExpireTime,
	})
}

func verifyAppleNotification(signedPayload string) (*AppleNotificationResult, error) {
	payload, err := verifyAndDecodeJWS[AppleNotificationPayload](signedPayload)
	if err != nil {
		return nil, err
	}

	result := &AppleNotificationResult{
		Verified:         true,
		NotificationType: payload.NotificationType,
		Subtype:          payload.Subtype,
		NotificationUUID: payload.NotificationUUID,
		SignedDate:       payload.SignedDate,
		Payload:          payload,
	}

	if payload.Data != nil {
		result.Environment = payload.Data.Environment
		result.BundleID = payload.Data.BundleID

		expectedBundleID := settings.Conf.AppleKeys.AppleBundleId
		if expectedBundleID != "" && payload.Data.BundleID != expectedBundleID {
			return nil, errors.New("apple notification bundleId mismatch")
		}

		// ✅ transaction
		if payload.Data.SignedTransactionInfo != "" {
			transactionInfo, err := verifyAndDecodeJWS[InAppLookupResp](payload.Data.SignedTransactionInfo)
			if err != nil {
				return nil, err
			}

			if expectedBundleID != "" && transactionInfo.BundleID != expectedBundleID {
				return nil, errors.New("transaction bundleId mismatch")
			}

			result.TransactionInfo = transactionInfo
		}

		// ✅ renewal
		if payload.Data.SignedRenewalInfo != "" {
			renewalInfo, err := verifyAndDecodeJWS[JWSRenewalInfoResp](payload.Data.SignedRenewalInfo)
			if err != nil {
				return nil, err
			}
			result.RenewalInfo = renewalInfo
		}
	}

	return result, nil
}

// func NotifyAppleIAP(c *gin.Context) {
// 	var req AppleNotificationRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
// 		return
// 	}

// 	result, err := verifyAppleNotification(req.SignedPayload)
// 	if err != nil {
// 		logger.Logger.Error("Verify Apple notification failed", zap.Error(err))
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"code": response.ApiCode.ParamErr,
// 			"msg":  "invalid apple notification",
// 		})
// 		return
// 	}

// 	txInfo := result.TransactionInfo
// 	if txInfo == nil || txInfo.OriginalTransactionID == "" {
// 		logger.Logger.Info("skip invalid transaction",
// 			zap.String("tx", txInfo.TransactionID))
// 		c.JSON(http.StatusOK, gin.H{"ok": true})
// 		return
// 	}

// 	now := time.Now().UnixMilli()

// 	var subscription models.Subscription
// 	dbTx := db.DB.Where("original_transaction_id = ?", txInfo.OriginalTransactionID).First(&subscription)

// 	// ====== 🧠 防乱序（核心）======
// 	if dbTx.Error == nil && subscription.ExpireTime > 0 {
// 		if txInfo.ExpiresDate > 0 && txInfo.ExpiresDate < subscription.ExpireTime {
// 			// 老通知，直接忽略
// 			logger.Logger.Info("skip outdated notification",
// 				zap.String("tx", txInfo.TransactionID))
// 			c.JSON(http.StatusOK, gin.H{"ok": true})
// 			return
// 		}
// 	}

// 	// ====== 🧾 写 transactions（幂等）======
// 	var tx models.Transaction
// 	err = db.DB.
// 		Where("transaction_id = ?", txInfo.TransactionID).
// 		First(&tx).Error

// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		db.DB.Create(&models.Transaction{
// 			TransactionID:         txInfo.TransactionID,
// 			OriginalTransactionID: txInfo.OriginalTransactionID,
// 			ProductId:             txInfo.ProductID,
// 			ExpireTime:            txInfo.ExpiresDate,
// 		})
// 	}

// 	// ====== 🧠 计算状态（核心逻辑）======
// 	status := 1 // active

// 	switch result.NotificationType {

// 	case "DID_RENEW":
// 		status = 1

// 	case "EXPIRED":
// 		status = 2

// 	case "REFUND", "REVOKE":
// 		status = 2

// 	case "DID_FAIL_TO_RENEW":
// 		// 不直接改状态，用时间判断

// 	case "DID_CHANGE_RENEWAL_STATUS":
// 		// 不改状态

// 	default:
// 		// 保持原状态
// 		if dbTx.Error == nil {
// 			status = subscription.Status
// 		}
// 	}

// 	// ====== ⏰ 时间最终裁决 ======
// 	expireTime := txInfo.ExpiresDate
// 	if expireTime > 0 && expireTime < now {
// 		status = 2
// 	}

// 	// ====== 🔁 autoRenew ======
// 	autoRenew := 0
// 	if dbTx.Error == nil {
// 		autoRenew = subscription.AutoRenew
// 	}

// 	if result.RenewalInfo != nil {
// 		if result.RenewalInfo.AutoRenewStatus == 1 {
// 			autoRenew = 1
// 		} else {
// 			autoRenew = 0
// 		}
// 	}

// 	// ====== 💾 更新 subscription ======
// 	update := models.Subscription{
// 		ProductId:             txInfo.ProductID,
// 		Platform:              "apple",
// 		OriginalTransactionId: txInfo.OriginalTransactionID,
// 		ExpireTime:            expireTime,
// 		Status:                status,
// 		AutoRenew:             autoRenew,
// 	}

// 	if dbTx.Error == nil {
// 		db.DB.Model(&subscription).Updates(update)
// 	} else {
// 		db.DB.Create(&update)
// 	}
// 	logger.Logger.Info("NotifyAppleIAP update subscription success",
// 		zap.String("tx", txInfo.TransactionID),
// 		zap.Int("status", status),
// 		zap.Int("autoRenew", autoRenew),
// 		zap.Int64("expireTime", expireTime))
// 	c.JSON(http.StatusOK, gin.H{"ok": true})
// }

func getUserIDByOriginalTx(originalTx string) uint {
	var sub models.Subscription
	err := db.DB.
		Where("original_transaction_id = ?", originalTx).
		First(&sub).Error

	if err == nil {
		return sub.UserId
	}

	return 0
}

func NotifyAppleIAP(c *gin.Context) {
	var req AppleNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := verifyAppleNotification(req.SignedPayload)
	if err != nil {
		logger.Logger.Error("Verify Apple notification failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": response.ApiCode.ParamErr,
			"msg":  "invalid apple notification",
		})
		return
	}

	// ================== 0️⃣ 幂等：通知去重 ==================
	if result.NotificationUUID != "" {
		var count int64
		db.DB.Model(&models.NotificationLog{}).
			Where("uuid = ?", result.NotificationUUID).
			Count(&count)

		if count > 0 {
			c.JSON(http.StatusOK, gin.H{"ok": true})
			return
		}

		db.DB.Create(&models.NotificationLog{
			UUID: result.NotificationUUID,
		})
	}

	txInfo := result.TransactionInfo
	if txInfo == nil || txInfo.TransactionID == "" {
		c.JSON(http.StatusOK, gin.H{"ok": true})
		return
	}

	now := time.Now().UnixMilli()

	// ================== 1️⃣ Transaction 幂等 ==================
	var tx models.Transaction
	err = db.DB.
		Where("transaction_id = ?", txInfo.TransactionID).
		First(&tx).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		db.DB.Create(&models.Transaction{
			TransactionID:         txInfo.TransactionID,
			OriginalTransactionID: txInfo.OriginalTransactionID,
			ProductId:             txInfo.ProductID,
			ExpireTime:            txInfo.ExpiresDate,
		})
	}

	// ================== 2️⃣ 查 Subscription ==================
	var subscription models.Subscription
	subErr := db.DB.
		Where("original_transaction_id = ?", txInfo.OriginalTransactionID).
		First(&subscription).Error

	// ================== 3️⃣ 防乱序 ==================
	if subErr == nil && subscription.ExpireTime > 0 && txInfo.ExpiresDate > 0 {
		if txInfo.ExpiresDate < subscription.ExpireTime {
			logger.Logger.Info("skip outdated notification",
				zap.String("tx", txInfo.TransactionID))
			c.JSON(http.StatusOK, gin.H{"ok": true})
			return
		}
	}

	// ================== 4️⃣ 状态计算（核心） ==================
	expireTime := txInfo.ExpiresDate

	status := 1 // active
	if expireTime > 0 && expireTime < now {
		status = 2 // expired
	}

	// REFUND / REVOKE 强制失效
	if result.NotificationType == "REFUND" || result.NotificationType == "REVOKE" {
		status = 2
	}

	// ================== 5️⃣ autoRenew ==================
	autoRenew := 0
	if subErr == nil {
		autoRenew = subscription.AutoRenew
	}

	if result.RenewalInfo != nil {
		if result.RenewalInfo.AutoRenewStatus == 1 {
			autoRenew = 1
		} else {
			autoRenew = 0
		}
	}

	// ================== 6️⃣ 更新 Subscription ==================
	update := models.Subscription{
		ProductId:             txInfo.ProductID,
		Platform:              "apple",
		OriginalTransactionId: txInfo.OriginalTransactionID,
		ExpireTime:            expireTime,
		Status:                status,
		AutoRenew:             autoRenew,
	}

	if subErr == nil {
		// ✅ 如果当前 userId = 0，尝试修复
		if subscription.UserId == 0 {
			userID := getUserIDByOriginalTx(txInfo.OriginalTransactionID)
			if userID > 0 {
				update.UserId = userID
			}
		}
		db.DB.Model(&subscription).Updates(update)
	} else {
		// ❗必须补 userId（关键）
		var userID uint

		// 👉 你需要自己实现这个函数
		userID = getUserIDByOriginalTx(txInfo.OriginalTransactionID)
		if userID == 0 {
			logger.Logger.Warn("create subscription without userId",
				zap.String("originalTx", txInfo.OriginalTransactionID))
		}
		update.UserId = userID
		db.DB.Create(&update)
	}

	logger.Logger.Info("NotifyAppleIAP success",
		zap.String("tx", txInfo.TransactionID),
		zap.Int("status", status),
		zap.Int("autoRenew", autoRenew),
		zap.Int64("expireTime", expireTime))

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
