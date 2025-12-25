package balance

import (
	"fmt"
	"gochat/api/api/balance"
	"gochat/internal/application/handler/common"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SendPrivateRedPacket 发送私聊红包
func (h *BalanceHandler) SendPrivateRedPacket(ctx *gin.Context) {
	req, err := analysis.BindParameter[balance.SendPrivateRedPacketRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.sendPrivateRedPacketLogic(ctx, &req)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *BalanceHandler) sendPrivateRedPacketLogic(ctx *gin.Context, req *balance.SendPrivateRedPacketRequest) (resp *balance.SendPrivateRedPacketResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	amount := req.GetAmount()
	receiverID := req.GetReceiverId()

	// 使用事务处理红包发送
	err = h.dao.GetDB().Transaction(func(tx *gorm.DB) error {
		balanceDao := dao.NewBalanceDao(tx)

		// 1. 检查并扣减发送者余额（加锁）
		userBalance, err := balanceDao.GetBalanceForUpdate(tx, userID)
		if err != nil {
			return err
		}
		if userBalance.Balance < amount {
			return gorm.ErrRecordNotFound // TODO: 返回余额不足错误
		}

		// 2. 扣减余额
		err = balanceDao.UpdateBalance(userID, -amount)
		if err != nil {
			return err
		}

		// 3. 创建红包记录
		now := time.Now().Unix()
		redPacket := &models.RedPacket{
			Type:        models.RedPacketTypePrivate,
			SenderID:    userID,
			ReceiverID:  receiverID,
			GroupID:     0,
			Amount:      amount,
			TotalAmount: amount,
			Status:      models.RedPacketStatusSent,
			Message:     req.GetMessage(),
			ExpiredAt:   now + 86400*7, // 私聊红包7天过期（可配置）
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		err = balanceDao.CreateRedPacket(redPacket)
		if err != nil {
			return err
		}

		// 4. 记录资金流水
		newBalance := userBalance.Balance - amount
		flow := &models.BalanceFlow{
			UserID:    userID,
			Type:      models.FlowTypeRedPacketSend,
			Amount:    -amount, // 负数表示支出
			Balance:   newBalance,
			RelatedID: redPacket.ID,
			Remark:    "发送私聊红包",
			CreatedAt: now,
		}
		err = balanceDao.CreateBalanceFlow(flow)
		if err != nil {
			return err
		}

		resp = &balance.SendPrivateRedPacketResponse{
			RedPacketId: redPacket.ID,
			Success:     true,
		}
		return nil
	})

	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return resp, 0, nil
}

// SendGroupRedPacket 发送群聊红包（手气红包）
func (h *BalanceHandler) SendGroupRedPacket(ctx *gin.Context) {
	req, err := analysis.BindParameter[balance.SendGroupRedPacketRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.sendGroupRedPacketLogic(ctx, &req)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *BalanceHandler) sendGroupRedPacketLogic(ctx *gin.Context, req *balance.SendGroupRedPacketRequest) (resp *balance.SendGroupRedPacketResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	redPacketType := req.GetType()
	count := int(req.GetCount())
	groupID := req.GetGroupId()

	// 验证红包个数
	if count <= 0 || count > 100 {
		return nil, code_msg.BadRequest, nil
	}

	var totalAmount int64
	var singleAmount int64
	var redPacketAllocType string

	// 根据红包类型验证金额
	if redPacketType == balance.RedPacketType_LUCKY {
		// 拼手气红包：总金额
		totalAmount = req.GetTotalAmount()
		if totalAmount < int64(count) {
			return nil, code_msg.BadRequest, nil // 总金额不能小于红包个数（每个至少1分）
		}
		redPacketAllocType = models.RedPacketAllocTypeLucky
	} else {
		// 普通红包：单个金额
		singleAmount = req.GetSingleAmount()
		if singleAmount <= 0 {
			return nil, code_msg.BadRequest, nil
		}
		totalAmount = singleAmount * int64(count)
		redPacketAllocType = models.RedPacketAllocTypeNormal
	}

	// 使用事务处理红包发送
	err = h.dao.GetDB().Transaction(func(tx *gorm.DB) error {
		balanceDao := dao.NewBalanceDao(tx)

		// 1. 检查并扣减发送者余额（加锁）
		userBalance, err := balanceDao.GetBalanceForUpdate(tx, userID)
		if err != nil {
			return err
		}
		if userBalance.Balance < totalAmount {
			return gorm.ErrRecordNotFound // TODO: 返回余额不足错误
		}

		// 2. 扣减余额
		err = balanceDao.UpdateBalance(userID, -totalAmount)
		if err != nil {
			return err
		}

		// 3. 创建红包记录
		now := time.Now().Unix()
		redPacket := &models.RedPacket{
			Type:            models.RedPacketTypeGroup,
			RedPacketType:   redPacketAllocType,
			SenderID:        userID,
			ReceiverID:      0,
			GroupID:         groupID,
			Amount:          singleAmount, // 普通红包时存储单个金额，拼手气红包时存储总金额
			TotalAmount:     totalAmount,
			Count:           count,
			RemainingAmount: totalAmount, // 剩余金额，用于并发控制
			RemainingCount:  count,       // 剩余个数，用于并发控制
			Status:          models.RedPacketStatusSent,
			Message:         req.GetMessage(),
			ExpiredAt:       now + 86400, // 群聊红包24小时过期
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		err = balanceDao.CreateRedPacket(redPacket)
		if err != nil {
			return err
		}

		// 4. 记录资金流水
		newBalance := userBalance.Balance - totalAmount
		flow := &models.BalanceFlow{
			UserID:    userID,
			Type:      models.FlowTypeRedPacketSend,
			Amount:    -totalAmount,
			Balance:   newBalance,
			RelatedID: redPacket.ID,
			Remark:    "发送群聊红包",
			CreatedAt: now,
		}
		err = balanceDao.CreateBalanceFlow(flow)
		if err != nil {
			return err
		}

		resp = &balance.SendGroupRedPacketResponse{
			RedPacketId: redPacket.ID,
			Success:     true,
		}
		return nil
	})

	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return resp, 0, nil
}

// ReceiveRedPacket 领取红包
func (h *BalanceHandler) ReceiveRedPacket(ctx *gin.Context) {
	req, err := analysis.BindParameter[balance.ReceiveRedPacketRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.receiveRedPacketLogic(ctx, &req)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *BalanceHandler) receiveRedPacketLogic(ctx *gin.Context, req *balance.ReceiveRedPacketRequest) (resp *balance.ReceiveRedPacketResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	redPacketID := req.GetRedPacketId()

	// 使用事务处理红包领取
	var receiveAmount int64
	err = h.dao.GetDB().Transaction(func(tx *gorm.DB) error {
		balanceDao := dao.NewBalanceDao(tx)
		now := time.Now().Unix()

		// 1. 检查是否已领取（先检查，避免不必要的加锁）
		receive, _ := balanceDao.GetRedPacketReceive(redPacketID, userID)
		if receive != nil {
			return gorm.ErrRecordNotFound // TODO: 返回已领取错误
		}

		// 2. 获取红包信息并加锁（SELECT FOR UPDATE，确保并发安全）
		redPacket, err := balanceDao.GetRedPacketForUpdate(tx, redPacketID)
		if err != nil {
			return err
		}
		if redPacket == nil {
			return gorm.ErrRecordNotFound
		}

		// 3. 检查红包状态
		if redPacket.Status != models.RedPacketStatusSent {
			return gorm.ErrRecordNotFound // TODO: 返回红包已领取或过期错误
		}

		// 4. 检查是否过期
		if redPacket.ExpiredAt > 0 && redPacket.ExpiredAt < now {
			// 红包已过期，退回余额
			balanceDao.UpdateRedPacket(redPacketID, map[string]interface{}{
				"status":     models.RedPacketStatusExpired,
				"updated_at": now,
			})
			// TODO: 退回余额给发送者
			return gorm.ErrRecordNotFound
		}

		// 5. 检查剩余金额和个数（并发控制的关键）
		if redPacket.RemainingAmount <= 0 || redPacket.RemainingCount <= 0 {
			// 红包已领完
			balanceDao.UpdateRedPacket(redPacketID, map[string]interface{}{
				"status":     models.RedPacketStatusReceived,
				"updated_at": now,
			})
			return gorm.ErrRecordNotFound
		}

		// 6. 根据红包类型分配金额
		if redPacket.Type == models.RedPacketTypePrivate {
			// 私聊红包：固定金额
			if redPacket.ReceiverID != userID {
				return gorm.ErrRecordNotFound // TODO: 返回权限错误
			}
			receiveAmount = redPacket.Amount
		} else {
			// 群聊红包
			if redPacket.RedPacketType == models.RedPacketAllocTypeNormal {
				// 普通红包：固定金额
				receiveAmount = redPacket.Amount
			} else {
				// 拼手气红包：随机分配金额
				if redPacket.RemainingCount == 1 {
					// 最后一个红包，剩余金额全部分配
					receiveAmount = redPacket.RemainingAmount
				} else {
					// 随机分配：在1分到剩余金额的2倍/剩余个数之间随机
					// 确保至少1分，最多不超过剩余金额
					maxAmount := redPacket.RemainingAmount / int64(redPacket.RemainingCount) * 2
					if maxAmount < 1 {
						maxAmount = 1
					}
					if maxAmount > redPacket.RemainingAmount {
						maxAmount = redPacket.RemainingAmount
					}
					receiveAmount = int64(rand.Intn(int(maxAmount))) + 1
					if receiveAmount > redPacket.RemainingAmount {
						receiveAmount = redPacket.RemainingAmount
					}
					// 确保至少1分
					if receiveAmount < 1 {
						receiveAmount = 1
					}
				}
			}

			// 再次检查：确保分配的金额不超过剩余金额（双重检查）
			if receiveAmount > redPacket.RemainingAmount {
				receiveAmount = redPacket.RemainingAmount
			}
		}

		// 7. 原子更新红包的剩余金额和剩余个数（关键：使用WHERE条件确保并发安全）
		newRemainingAmount := redPacket.RemainingAmount - receiveAmount
		newRemainingCount := redPacket.RemainingCount - 1

		updateData := map[string]interface{}{
			"remaining_amount": newRemainingAmount,
			"remaining_count":  newRemainingCount,
			"updated_at":       now,
		}

		// 如果领完了，更新状态
		if newRemainingAmount <= 0 || newRemainingCount <= 0 {
			updateData["status"] = models.RedPacketStatusReceived
		}

		// 使用原子更新，WHERE条件确保剩余金额和个数足够
		result := tx.Model(&models.RedPacket{}).
			Where("id = ? AND remaining_amount >= ? AND remaining_count > 0", redPacketID, receiveAmount).
			Updates(updateData)

		if result.Error != nil {
			return result.Error
		}

		// 检查是否更新成功（如果没有更新，说明并发冲突）
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound // 并发冲突，红包已被其他用户领取完或金额不足
		}

		// 8. 增加接收者余额
		receiverBalance, err := balanceDao.GetBalanceForUpdate(tx, userID)
		if err != nil {
			return err
		}
		err = balanceDao.UpdateBalance(userID, receiveAmount)
		if err != nil {
			return err
		}

		// 9. 创建领取记录
		receiveRecord := &models.RedPacketReceive{
			RedPacketID: redPacketID,
			UserID:      userID,
			Amount:      receiveAmount,
			ReceivedAt:  now,
		}
		err = balanceDao.CreateRedPacketReceive(receiveRecord)
		if err != nil {
			return err
		}

		// 10. 记录资金流水
		newBalance := receiverBalance.Balance + receiveAmount
		flow := &models.BalanceFlow{
			UserID:    userID,
			Type:      models.FlowTypeRedPacketReceive,
			Amount:    receiveAmount,
			Balance:   newBalance,
			RelatedID: redPacketID,
			Remark:    "领取红包",
			CreatedAt: now,
		}
		err = balanceDao.CreateBalanceFlow(flow)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 创建通知：给发送者发送红包被领取通知
	// 获取红包信息
	redPacket, _ := h.dao.GetRedPacket(redPacketID)
	if redPacket != nil {
		// 获取接收者信息
		userDao := dao.NewUserDao(globalUtils.DB)
		receiver, _ := userDao.GetUserByID(userID)
		receiverName := "用户"
		if receiver != nil {
			receiverName = receiver.Name
		}

		title := "红包被领取"
		content := fmt.Sprintf("%s 领取了你的红包，金额：%.2f元", receiverName, float64(receiveAmount)/100.0)

		// 检查红包是否已领完
		if redPacket.RemainingAmount <= 0 || redPacket.RemainingCount <= 0 {
			content = fmt.Sprintf("%s 领取了你的红包，金额：%.2f元，红包已领完", receiverName, float64(receiveAmount)/100.0)
		}

		// 异步创建通知
		go func() {
			if err := common.CreateRedPacketNotification(
				redPacket.SenderID,
				title,
				content,
				receiveAmount,
				redPacketID,
			); err != nil {
				fmt.Printf("Failed to create red packet receive notification: %v\n", err)
			}
		}()
	}

	return &balance.ReceiveRedPacketResponse{
		Amount:  receiveAmount,
		Success: true,
	}, 0, nil
}
