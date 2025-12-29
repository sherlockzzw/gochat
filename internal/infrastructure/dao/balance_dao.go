package dao

import (
	"gochat/internal/infrastructure/models"
	"time"

	"gorm.io/gorm"
)

type BalanceDao struct {
	db *gorm.DB
}

func NewBalanceDao(db *gorm.DB) *BalanceDao {
	return &BalanceDao{db: db}
}

// GetDB 获取数据库连接（用于事务）
func (d *BalanceDao) GetDB() *gorm.DB {
	return d.db
}

// GetBalance 获取用户余额
func (d *BalanceDao) GetBalance(userID int64) (*models.UserBalance, error) {
	var balance models.UserBalance
	err := d.db.Where("user_id = ?", userID).First(&balance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，创建初始余额记录
			balance = models.UserBalance{
				UserID:    userID,
				Balance:   0,
				UpdatedAt: time.Now().Unix(),
			}
			if err := d.db.Create(&balance).Error; err != nil {
				return nil, err
			}
			return &balance, nil
		}
		return nil, err
	}
	return &balance, nil
}

// UpdateBalance 更新余额(必须在事务中使用)
func (d *BalanceDao) UpdateBalance(userID int64, amount int64) error {
	return d.db.Model(&models.UserBalance{}).
		Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

// GetBalanceForUpdate 获取余额并加锁(用于事务中的SELECT FOR UPDATE)
func (d *BalanceDao) GetBalanceForUpdate(tx *gorm.DB, userID int64) (*models.UserBalance, error) {
	var balance models.UserBalance
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("user_id = ?", userID).
		First(&balance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，创建初始余额记录
			balance = models.UserBalance{
				UserID:    userID,
				Balance:   0,
				UpdatedAt: time.Now().Unix(),
			}
			if err := tx.Create(&balance).Error; err != nil {
				return nil, err
			}
			return &balance, nil
		}
		return nil, err
	}
	return &balance, nil
}

// CreateRechargeRequest 创建充值申请
func (d *BalanceDao) CreateRechargeRequest(req *models.RechargeRequest) error {
	return d.db.Create(req).Error
}

// GetRechargeRequest 获取充值申请
func (d *BalanceDao) GetRechargeRequest(id int64) (*models.RechargeRequest, error) {
	var req models.RechargeRequest
	err := d.db.Where("id = ?", id).First(&req).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

// GetRechargeRequests 获取用户的充值申请列表
func (d *BalanceDao) GetRechargeRequests(userID int64, page, pageSize int) ([]*models.RechargeRequest, int64, error) {
	var requests []*models.RechargeRequest
	var total int64

	query := d.db.Model(&models.RechargeRequest{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error

	return requests, total, err
}

// UpdateRechargeRequest 更新充值申请状态
func (d *BalanceDao) UpdateRechargeRequest(id int64, updates map[string]interface{}) error {
	return d.db.Model(&models.RechargeRequest{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// CreateWithdrawRequest 创建提现申请
func (d *BalanceDao) CreateWithdrawRequest(req *models.WithdrawRequest) error {
	return d.db.Create(req).Error
}

// GetWithdrawRequest 获取提现申请
func (d *BalanceDao) GetWithdrawRequest(id int64) (*models.WithdrawRequest, error) {
	var req models.WithdrawRequest
	err := d.db.Where("id = ?", id).First(&req).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

// GetWithdrawRequests 获取用户的提现申请列表
func (d *BalanceDao) GetWithdrawRequests(userID int64, page, pageSize int) ([]*models.WithdrawRequest, int64, error) {
	var requests []*models.WithdrawRequest
	var total int64

	query := d.db.Model(&models.WithdrawRequest{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error

	return requests, total, err
}

// UpdateWithdrawRequest 更新提现申请状态
func (d *BalanceDao) UpdateWithdrawRequest(id int64, updates map[string]interface{}) error {
	return d.db.Model(&models.WithdrawRequest{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// GetTodayWithdrawCount 获取用户今日提现次数
func (d *BalanceDao) GetTodayWithdrawCount(userID int64) (int64, error) {
	var count int64
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	endTime := startTime + 86400 // 24小时后

	err := d.db.Model(&models.WithdrawRequest{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startTime, endTime).
		Count(&count).Error

	return count, err
}

// CreateRedPacket 创建红包记录
func (d *BalanceDao) CreateRedPacket(packet *models.RedPacket) error {
	return d.db.Create(packet).Error
}

// GetRedPacket 获取红包记录
func (d *BalanceDao) GetRedPacket(id int64) (*models.RedPacket, error) {
	var packet models.RedPacket
	err := d.db.Where("id = ?", id).First(&packet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &packet, nil
}

// GetRedPacketForUpdate 获取红包记录并加锁(用于事务中的SELECT FOR UPDATE)
func (d *BalanceDao) GetRedPacketForUpdate(tx *gorm.DB, id int64) (*models.RedPacket, error) {
	var packet models.RedPacket
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", id).
		First(&packet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &packet, nil
}

// UpdateRedPacket 更新红包状态
func (d *BalanceDao) UpdateRedPacket(id int64, updates map[string]interface{}) error {
	return d.db.Model(&models.RedPacket{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// CreateRedPacketReceive 创建红包领取记录
func (d *BalanceDao) CreateRedPacketReceive(receive *models.RedPacketReceive) error {
	return d.db.Create(receive).Error
}

// GetRedPacketReceive 检查用户是否已领取红包
func (d *BalanceDao) GetRedPacketReceive(redPacketID, userID int64) (*models.RedPacketReceive, error) {
	var receive models.RedPacketReceive
	err := d.db.Where("red_packet_id = ? AND user_id = ?", redPacketID, userID).First(&receive).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &receive, nil
}

// GetRedPacketReceives 获取红包的所有领取记录
func (d *BalanceDao) GetRedPacketReceives(redPacketID int64) ([]*models.RedPacketReceive, error) {
	var receives []*models.RedPacketReceive
	err := d.db.Where("red_packet_id = ?", redPacketID).
		Order("received_at ASC").
		Find(&receives).Error
	return receives, err
}

// CreateTransfer 创建转账记录
func (d *BalanceDao) CreateTransfer(transfer *models.Transfer) error {
	return d.db.Create(transfer).Error
}

// GetTransfer 获取转账记录
func (d *BalanceDao) GetTransfer(id int64) (*models.Transfer, error) {
	var transfer models.Transfer
	err := d.db.Where("id = ?", id).First(&transfer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &transfer, nil
}

// UpdateTransfer 更新转账状态
func (d *BalanceDao) UpdateTransfer(id int64, updates map[string]interface{}) error {
	return d.db.Model(&models.Transfer{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// GetUserTransfers 获取用户的转账记录(发送或接收)
func (d *BalanceDao) GetUserTransfers(userID int64, page, pageSize int) ([]*models.Transfer, int64, error) {
	var transfers []*models.Transfer
	var total int64

	query := d.db.Model(&models.Transfer{}).
		Where("sender_id = ? OR receiver_id = ?", userID, userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&transfers).Error

	return transfers, total, err
}

// CreateBalanceFlow 创建资金流水
func (d *BalanceDao) CreateBalanceFlow(flow *models.BalanceFlow) error {
	return d.db.Create(flow).Error
}

// GetPendingRecharges 获取所有待审核的充值申请（管理员用）
func (d *BalanceDao) GetPendingRecharges(page, pageSize int) ([]*models.RechargeRequest, int64, error) {
	var requests []*models.RechargeRequest
	var total int64

	query := d.db.Model(&models.RechargeRequest{}).Where("status = ?", models.StatusPending)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetPendingWithdraws 获取所有待审核的提现申请（管理员用）
func (d *BalanceDao) GetPendingWithdraws(page, pageSize int) ([]*models.WithdrawRequest, int64, error) {
	var requests []*models.WithdrawRequest
	var total int64

	query := d.db.Model(&models.WithdrawRequest{}).Where("status = ?", models.StatusPending)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetAllBalanceFlows 获取所有资金流水（管理员用）
func (d *BalanceDao) GetAllBalanceFlows(userID int64, flowType string, startTime, endTime int64, page, pageSize int) ([]*models.BalanceFlow, int64, error) {
	var flows []*models.BalanceFlow
	var total int64

	query := d.db.Model(&models.BalanceFlow{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if flowType != "" {
		query = query.Where("type = ?", flowType)
	}
	if startTime > 0 {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("created_at <= ?", endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&flows).Error; err != nil {
		return nil, 0, err
	}

	return flows, total, nil
}

// GetBalanceFlows 获取用户的资金流水
func (d *BalanceDao) GetBalanceFlows(userID int64, flowType string, page, pageSize int) ([]*models.BalanceFlow, int64, error) {
	var flows []*models.BalanceFlow
	var total int64

	query := d.db.Model(&models.BalanceFlow{}).Where("user_id = ?", userID)
	if flowType != "" {
		query = query.Where("type = ?", flowType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&flows).Error

	return flows, total, err
}

// GetExpiredRedPackets 获取过期的红包(用于定时任务退回)
func (d *BalanceDao) GetExpiredRedPackets(beforeTime int64) ([]*models.RedPacket, error) {
	var packets []*models.RedPacket
	err := d.db.Where("status = ? AND expired_at > 0 AND expired_at < ?", models.RedPacketStatusSent, beforeTime).
		Find(&packets).Error
	return packets, err
}
