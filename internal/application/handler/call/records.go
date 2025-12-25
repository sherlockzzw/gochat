package call

import (
	"gochat/api/api/call"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetCallRecords 获取通话记录
func (h *CallHandler) GetCallRecords(ctx *gin.Context) {
	req, err := analysis.BindQuery[call.GetCallRecordsRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getCallRecordsLogic(ctx, &req)
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

func (h *CallHandler) getCallRecordsLogic(ctx *gin.Context, req *call.GetCallRecordsRequest) (resp *call.GetCallRecordsResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	page := int(req.GetPage())
	if page <= 0 {
		page = 1
	}

	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	callType := req.GetType()
	if callType == "" {
		callType = "all"
	}

	// 获取通话记录
	records, total, err := h.dao.GetUserCallRecords(userID, page, pageSize, callType)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 转换为响应格式
	var callRecords []*call.CallRecord
	for _, r := range records {
		callRecord := &call.CallRecord{
			Id:          uint64(r.ID),
			RoomId:      uint64(r.RoomID),
			UserId:      uint32(r.UserID),
			Type:        r.Type,
			CallType:    r.CallType,
			OtherUserId: uint32(r.OtherUserID),
			GroupId:     uint32(r.GroupID),
			Direction:   r.Direction,
			Status:      r.Status,
			Duration:    r.Duration,
		}

		if r.StartedAt > 0 {
			callRecord.StartedAt = timestamppb.New(time.Unix(r.StartedAt, 0))
		}
		if r.EndedAt > 0 {
			callRecord.EndedAt = timestamppb.New(time.Unix(r.EndedAt, 0))
		}
		if r.CreatedAt > 0 {
			callRecord.CreatedAt = timestamppb.New(time.Unix(r.CreatedAt, 0))
		}

		callRecords = append(callRecords, callRecord)
	}

	resp = &call.GetCallRecordsResponse{
		Records:     callRecords,
		TotalCount:  int32(total),
		CurrentPage: int32(page),
		PageSize:    int32(pageSize),
	}

	return resp, 0, nil
}

