package balance

import (
	"gochat/api/api/balance"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SendPrivateTransfer 私聊转账
func (h *BalanceHandler) SendPrivateTransfer(ctx *gin.Context) {
	req, err := analysis.BindParameter[balance.SendPrivateTransferRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.sendPrivateTransferLogic(ctx, &req)
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

func (h *BalanceHandler) sendPrivateTransferLogic(ctx *gin.Context, req *balance.SendPrivateTransferRequest) (resp *balance.SendPrivateTransferResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	amount := req.GetAmount()
	receiverID := req.GetReceiverId()

	// 使用事务处理转账
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

		// 3. 创建转账记录
		now := time.Now().Unix()
		transfer := &models.Transfer{
			Type:        models.TransferTypePrivate,
			SenderID:    userID,
			ReceiverID:  receiverID,
			GroupID:     0,
			Amount:      amount,
			Status:      models.TransferStatusPending,
			Remark:      req.GetRemark(),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		err = balanceDao.CreateTransfer(transfer)
		if err != nil {
			return err
		}

		// 4. 记录资金流水
		newBalance := userBalance.Balance - amount
		flow := &models.BalanceFlow{
			UserID:    userID,
			Type:      models.FlowTypeTransferSend,
			Amount:    -amount,
			Balance:   newBalance,
			RelatedID: transfer.ID,
			Remark:    "私聊转账",
			CreatedAt: now,
		}
		err = balanceDao.CreateBalanceFlow(flow)
		if err != nil {
			return err
		}

		resp = &balance.SendPrivateTransferResponse{
			TransferId: transfer.ID,
			Success:    true,
		}
		return nil
	})

	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return resp, 0, nil
}

// SendGroupTransfer 群聊转账
func (h *BalanceHandler) SendGroupTransfer(ctx *gin.Context) {
	req, err := analysis.BindParameter[balance.SendGroupTransferRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.sendGroupTransferLogic(ctx, &req)
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

func (h *BalanceHandler) sendGroupTransferLogic(ctx *gin.Context, req *balance.SendGroupTransferRequest) (resp *balance.SendGroupTransferResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	amount := req.GetAmount()
	receiverID := req.GetReceiverId()
	groupID := req.GetGroupId()

	// 使用事务处理转账
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

		// 3. 创建转账记录
		now := time.Now().Unix()
		transfer := &models.Transfer{
			Type:       models.TransferTypeGroup,
			SenderID:   userID,
			ReceiverID: receiverID,
			GroupID:    groupID,
			Amount:     amount,
			Status:     models.TransferStatusPending,
			Remark:     req.GetRemark(),
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		err = balanceDao.CreateTransfer(transfer)
		if err != nil {
			return err
		}

		// 4. 记录资金流水
		newBalance := userBalance.Balance - amount
		flow := &models.BalanceFlow{
			UserID:    userID,
			Type:      models.FlowTypeTransferSend,
			Amount:    -amount,
			Balance:   newBalance,
			RelatedID: transfer.ID,
			Remark:    "群聊转账",
			CreatedAt: now,
		}
		err = balanceDao.CreateBalanceFlow(flow)
		if err != nil {
			return err
		}

		resp = &balance.SendGroupTransferResponse{
			TransferId: transfer.ID,
			Success:    true,
		}
		return nil
	})

	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return resp, 0, nil
}

// ReceiveTransfer 接收转账
func (h *BalanceHandler) ReceiveTransfer(ctx *gin.Context) {
	req, err := analysis.BindParameter[balance.ReceiveTransferRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.receiveTransferLogic(ctx, &req)
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

func (h *BalanceHandler) receiveTransferLogic(ctx *gin.Context, req *balance.ReceiveTransferRequest) (resp *balance.ReceiveTransferResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	transferID := req.GetTransferId()

	// 使用事务处理转账接收
	err = h.dao.GetDB().Transaction(func(tx *gorm.DB) error {
		balanceDao := dao.NewBalanceDao(tx)

		// 1. 获取转账记录（加锁）
		transfer, err := balanceDao.GetTransfer(transferID)
		if err != nil {
			return err
		}
		if transfer == nil {
			return gorm.ErrRecordNotFound
		}

		// 2. 检查转账状态
		if transfer.Status != models.TransferStatusPending {
			return gorm.ErrRecordNotFound // TODO: 返回已收款或已撤销错误
		}

		// 3. 检查接收者
		if transfer.ReceiverID != userID {
			return gorm.ErrRecordNotFound // TODO: 返回权限错误
		}

		// 4. 增加接收者余额
		receiverBalance, err := balanceDao.GetBalanceForUpdate(tx, userID)
		if err != nil {
			return err
		}
		err = balanceDao.UpdateBalance(userID, transfer.Amount)
		if err != nil {
			return err
		}

		// 5. 更新转账状态
		now := time.Now().Unix()
		err = balanceDao.UpdateTransfer(transferID, map[string]interface{}{
			"status":     models.TransferStatusReceived,
			"updated_at": now,
		})
		if err != nil {
			return err
		}

		// 6. 记录资金流水
		newBalance := receiverBalance.Balance + transfer.Amount
		flow := &models.BalanceFlow{
			UserID:    userID,
			Type:      models.FlowTypeTransferReceive,
			Amount:    transfer.Amount,
			Balance:   newBalance,
			RelatedID: transferID,
			Remark:    "接收转账",
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

	return &balance.ReceiveTransferResponse{
		Success: true,
	}, 0, nil
}

// CancelTransfer 撤销转账（仅私聊，且未收款前）
func (h *BalanceHandler) CancelTransfer(ctx *gin.Context) {
	req, err := analysis.BindParameter[balance.CancelTransferRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.cancelTransferLogic(ctx, &req)
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

func (h *BalanceHandler) cancelTransferLogic(ctx *gin.Context, req *balance.CancelTransferRequest) (resp *balance.CancelTransferResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	transferID := req.GetTransferId()

	// 使用事务处理转账撤销
	err = h.dao.GetDB().Transaction(func(tx *gorm.DB) error {
		balanceDao := dao.NewBalanceDao(tx)

		// 1. 获取转账记录（加锁）
		transfer, err := balanceDao.GetTransfer(transferID)
		if err != nil {
			return err
		}
		if transfer == nil {
			return gorm.ErrRecordNotFound
		}

		// 2. 检查权限：只有发送者可以撤销
		if transfer.SenderID != userID {
			return gorm.ErrRecordNotFound // TODO: 返回权限错误
		}

		// 3. 检查转账状态：只有待收款状态可以撤销
		if transfer.Status != models.TransferStatusPending {
			return gorm.ErrRecordNotFound // TODO: 返回已收款或已撤销错误
		}

		// 4. 检查转账类型：只有私聊转账可以撤销
		if transfer.Type != models.TransferTypePrivate {
			return gorm.ErrRecordNotFound // TODO: 返回群聊转账不可撤销错误
		}

		// 5. 退回余额给发送者
		senderBalance, err := balanceDao.GetBalanceForUpdate(tx, userID)
		if err != nil {
			return err
		}
		err = balanceDao.UpdateBalance(userID, transfer.Amount)
		if err != nil {
			return err
		}

		// 6. 更新转账状态
		now := time.Now().Unix()
		err = balanceDao.UpdateTransfer(transferID, map[string]interface{}{
			"status":     models.TransferStatusCancelled,
			"updated_at": now,
		})
		if err != nil {
			return err
		}

		// 7. 记录资金流水（退款）
		newBalance := senderBalance.Balance + transfer.Amount
		flow := &models.BalanceFlow{
			UserID:    userID,
			Type:      models.FlowTypeRefund,
			Amount:    transfer.Amount,
			Balance:   newBalance,
			RelatedID: transferID,
			Remark:    "撤销转账退回",
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

	return &balance.CancelTransferResponse{
		Success: true,
	}, 0, nil
}


