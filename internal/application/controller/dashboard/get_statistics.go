package dashboard

import (
	"gochat/internal/infrastructure/models"
	"gochat/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardStatisticsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type DashboardStatisticsData struct {
	TotalUsers             int64                    `json:"total_users"`
	TodayActiveUsers       int64                    `json:"today_active_users"`
	TodayTransactionAmount float64                  `json:"today_transaction_amount"`
	PendingReviews         int64                    `json:"pending_reviews"`
	UserGrowth             []DashboardChartPoint    `json:"user_growth"`
	TransactionBreakdown   []DashboardPieChartPoint `json:"transaction_breakdown"`
}

type DashboardChartPoint struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

type DashboardPieChartPoint struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (c *DashboardController) GetStatistics(ctx *gin.Context) {
	// 1) 总用户数
	var totalUsers int64
	_ = utils.DB.Model(&models.UserBasic{}).Count(&totalUsers).Error

	// 2) 今日活跃：使用 last_active_time(秒) 落在今日区间的用户数；若字段不可用则降级为 0
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	end := start + 24*3600

	var todayActive int64
	_ = utils.DB.Model(&models.UserBasic{}).Where("last_active_time >= ? AND last_active_time < ?", start, end).Count(&todayActive).Error

	// 3) 今日交易额：根据充值申请(approved) + 提现申请(approved) + 红包发送(按 total_amount) + 转账(按 amount)
	// 注意：金额字段单位均为“分”，这里返回“元”用于前端展示。
	var rechargeApprovedSum int64
	_ = utils.DB.Model(&models.RechargeRequest{}).
		Where("status = ? AND created_at >= ? AND created_at < ?", models.StatusApproved, start, end).
		Select("COALESCE(SUM(amount),0)").
		Scan(&rechargeApprovedSum).Error

	var withdrawApprovedSum int64
	_ = utils.DB.Model(&models.WithdrawRequest{}).
		Where("status = ? AND created_at >= ? AND created_at < ?", models.StatusApproved, start, end).
		Select("COALESCE(SUM(amount),0)").
		Scan(&withdrawApprovedSum).Error

	var redPacketSum int64
	_ = utils.DB.Model(&models.RedPacket{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Select("COALESCE(SUM(total_amount),0)").
		Scan(&redPacketSum).Error

	var transferSum int64
	_ = utils.DB.Model(&models.Transfer{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Select("COALESCE(SUM(amount),0)").
		Scan(&transferSum).Error

	todayAmount := float64(rechargeApprovedSum+withdrawApprovedSum+redPacketSum+transferSum) / 100.0

	// 4) 待审核：充值申请 pending + 提现申请 pending
	var rechargePending int64
	_ = utils.DB.Model(&models.RechargeRequest{}).Where("status = ?", models.StatusPending).Count(&rechargePending).Error
	var withdrawPending int64
	_ = utils.DB.Model(&models.WithdrawRequest{}).Where("status = ?", models.StatusPending).Count(&withdrawPending).Error
	pending := rechargePending + withdrawPending

	// 4) 用户增长趋势：最近 7 天新增用户（按 created_at 秒）
	userGrowth := make([]DashboardChartPoint, 0, 7)
	for i := 6; i >= 0; i-- {
		d := now.AddDate(0, 0, -i)
		dStart := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()).Unix()
		dEnd := dStart + 24*3600

		var cnt int64
		_ = utils.DB.Model(&models.UserBasic{}).Where("created_at >= ? AND created_at < ?", dStart, dEnd).Count(&cnt).Error

		userGrowth = append(userGrowth, DashboardChartPoint{
			Label: d.Format("01-02"),
			Value: cnt,
		})
	}

	// 5) 交易统计饼图：按今日金额汇总（单位：元）
	transaction := []DashboardPieChartPoint{
		{Name: "充值", Value: float64(rechargeApprovedSum) / 100.0},
		{Name: "提现", Value: float64(withdrawApprovedSum) / 100.0},
		{Name: "红包", Value: float64(redPacketSum) / 100.0},
		{Name: "转账", Value: float64(transferSum) / 100.0},
	}

	// 如果前端希望金额展示为 ¥xx,xxx，可在前端格式化；这里保持数值
	data := DashboardStatisticsData{
		TotalUsers:             totalUsers,
		TodayActiveUsers:       todayActive,
		TodayTransactionAmount: todayAmount,
		PendingReviews:         pending,
		UserGrowth:             userGrowth,
		TransactionBreakdown:   transaction,
	}

	c.response.JsonSuccess(ctx, data)
}
