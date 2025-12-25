package models

import (
	"gorm.io/gorm"
)

// UserBalance 用户余额表
type UserBalance struct {
	UserID    int64 `gorm:"primaryKey;column:user_id;comment:用户ID"`
	Balance   int64 `gorm:"column:balance;type:bigint;not null;default:0;comment:余额(单位:分)"`
	UpdatedAt int64 `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
}

// TableName 指定表名
func (UserBalance) TableName() string {
	return "user_balances"
}

// RechargeRequest 充值申请表
type RechargeRequest struct {
	ID        int64          `gorm:"primaryKey;comment:充值申请ID"`
	UserID    int64          `gorm:"column:user_id;not null;index;comment:用户ID"`
	Amount    int64          `gorm:"column:amount;type:bigint;not null;comment:充值金额(单位:分)"`
	Status    string         `gorm:"column:status;type:varchar(20);not null;default:'pending';index;comment:状态(pending待审核,approved已通过,rejected已驳回)"`
	Remark    string         `gorm:"column:remark;type:varchar(500);comment:备注说明"`
	AuditorID int64          `gorm:"column:auditor_id;comment:审核人ID(管理员)"`
	AuditTime int64          `gorm:"column:audit_time;type:bigint;not null;default:0;comment:审核时间(时间戳)"`
	AuditNote string         `gorm:"column:audit_note;type:varchar(500);comment:审核备注"`
	CreatedAt int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 指定表名
func (RechargeRequest) TableName() string {
	return "recharge_requests"
}

// WithdrawRequest 提现申请表
type WithdrawRequest struct {
	ID          int64          `gorm:"primaryKey;comment:提现申请ID"`
	UserID      int64          `gorm:"column:user_id;not null;index;comment:用户ID"`
	Amount      int64          `gorm:"column:amount;type:bigint;not null;comment:提现金额(单位:分)"`
	Fee         int64          `gorm:"column:fee;type:bigint;not null;default:0;comment:手续费(单位:分)"`
	AccountInfo string         `gorm:"column:account_info;type:varchar(500);not null;comment:收款账户信息(加密存储)"`
	Status      string         `gorm:"column:status;type:varchar(20);not null;default:'pending';index;comment:状态(pending待审核,approved已通过,rejected已驳回)"`
	AuditorID   int64          `gorm:"column:auditor_id;comment:审核人ID(管理员)"`
	AuditTime   int64          `gorm:"column:audit_time;type:bigint;not null;default:0;comment:审核时间(时间戳)"`
	AuditNote   string         `gorm:"column:audit_note;type:varchar(500);comment:审核备注"`
	CreatedAt   int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt   int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 指定表名
func (WithdrawRequest) TableName() string {
	return "withdraw_requests"
}

// RedPacket 红包记录表
type RedPacket struct {
	ID              int64          `gorm:"primaryKey;comment:红包ID"`
	Type            string         `gorm:"column:type;type:varchar(20);not null;index;comment:类型(private私聊,group群聊)"`
	RedPacketType   string         `gorm:"column:red_packet_type;type:varchar(20);not null;default:'lucky';comment:红包类型(lucky拼手气,normal普通)"`
	SenderID        int64          `gorm:"column:sender_id;not null;index;comment:发送者ID"`
	ReceiverID      int64          `gorm:"column:receiver_id;index;comment:接收者ID(私聊时使用,群聊时为0)"`
	GroupID         int64          `gorm:"column:group_id;index;comment:群组ID(群聊时使用,私聊时为0)"`
	Amount          int64          `gorm:"column:amount;type:bigint;not null;comment:红包金额(单位:分,普通红包时为单个金额,拼手气红包时为总金额)"`
	TotalAmount     int64          `gorm:"column:total_amount;type:bigint;not null;comment:总金额(单位:分,群聊时为总金额,私聊时等于amount)"`
	Count           int            `gorm:"column:count;type:int;not null;comment:红包个数"`
	RemainingAmount int64          `gorm:"column:remaining_amount;type:bigint;not null;comment:剩余金额(单位:分,用于并发控制)"`
	RemainingCount  int            `gorm:"column:remaining_count;type:int;not null;comment:剩余个数(用于并发控制)"`
	Status          string         `gorm:"column:status;type:varchar(20);not null;default:'sent';index;comment:状态(sent已发送,received已领取,expired已过期,refunded已退回)"`
	Message         string         `gorm:"column:message;type:varchar(200);comment:祝福语"`
	ExpiredAt       int64          `gorm:"column:expired_at;type:bigint;not null;default:0;comment:过期时间(时间戳,群聊24小时,私聊可配置)"`
	CreatedAt       int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt       int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 指定表名
func (RedPacket) TableName() string {
	return "red_packets"
}

// RedPacketReceive 红包领取记录表(群聊时使用,记录每个成员的领取情况)
type RedPacketReceive struct {
	ID          int64 `gorm:"primaryKey;comment:领取记录ID"`
	RedPacketID int64 `gorm:"column:red_packet_id;not null;index:idx_red_packet_user,unique;comment:红包ID"`
	UserID      int64 `gorm:"column:user_id;not null;index:idx_red_packet_user,unique;comment:领取者ID"`
	Amount      int64 `gorm:"column:amount;type:bigint;not null;comment:领取金额(单位:分)"`
	ReceivedAt  int64 `gorm:"column:received_at;type:bigint;not null;default:0;comment:领取时间(时间戳)"`
}

// TableName 指定表名
func (RedPacketReceive) TableName() string {
	return "red_packet_receives"
}

// Transfer 转账记录表
type Transfer struct {
	ID         int64          `gorm:"primaryKey;comment:转账ID"`
	Type       string         `gorm:"column:type;type:varchar(20);not null;index;comment:类型(private私聊,group群聊)"`
	SenderID   int64          `gorm:"column:sender_id;not null;index;comment:发送者ID"`
	ReceiverID int64          `gorm:"column:receiver_id;not null;index;comment:接收者ID"`
	GroupID    int64          `gorm:"column:group_id;index;comment:群组ID(群聊时使用,私聊时为0)"`
	Amount     int64          `gorm:"column:amount;type:bigint;not null;comment:转账金额(单位:分)"`
	Status     string         `gorm:"column:status;type:varchar(20);not null;default:'pending';index;comment:状态(pending待收款,received已收款,cancelled已撤销,refunded已退回)"`
	Remark     string         `gorm:"column:remark;type:varchar(200);comment:转账备注"`
	CreatedAt  int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt  int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 指定表名
func (Transfer) TableName() string {
	return "transfers"
}

// BalanceFlow 资金流水表
type BalanceFlow struct {
	ID        int64  `gorm:"primaryKey;comment:流水ID"`
	UserID    int64  `gorm:"column:user_id;not null;index;comment:用户ID"`
	Type      string `gorm:"column:type;type:varchar(30);not null;index;comment:流水类型(recharge充值,withdraw提现,redpacket_send发红包,redpacket_receive收红包,transfer_send转账支出,transfer_receive转账收入,fee手续费,refund退款)"`
	Amount    int64  `gorm:"column:amount;type:bigint;not null;comment:金额(单位:分,正数表示收入,负数表示支出)"`
	Balance   int64  `gorm:"column:balance;type:bigint;not null;comment:操作后余额(单位:分)"`
	RelatedID int64  `gorm:"column:related_id;index;comment:关联记录ID(充值/提现/红包/转账ID)"`
	Remark    string `gorm:"column:remark;type:varchar(500);comment:备注说明"`
	CreatedAt int64  `gorm:"column:created_at;type:bigint;not null;default:0;index;comment:创建时间(时间戳)"`
}

// TableName 指定表名
func (BalanceFlow) TableName() string {
	return "balance_flows"
}

// 状态常量
const (
	// 充值/提现状态
	StatusPending  = "pending"  // 待审核
	StatusApproved = "approved" // 已通过
	StatusRejected = "rejected" // 已驳回

	// 红包状态
	RedPacketStatusSent     = "sent"     // 已发送
	RedPacketStatusReceived = "received" // 已领取
	RedPacketStatusExpired  = "expired"  // 已过期
	RedPacketStatusRefunded = "refunded" // 已退回

	// 转账状态
	TransferStatusPending   = "pending"   // 待收款
	TransferStatusReceived  = "received"  // 已收款
	TransferStatusCancelled = "cancelled" // 已撤销
	TransferStatusRefunded  = "refunded"  // 已退回

	// 红包类型
	RedPacketTypePrivate = "private" // 私聊红包
	RedPacketTypeGroup   = "group"   // 群聊红包

	// 红包分配类型
	RedPacketAllocTypeLucky  = "lucky"  // 拼手气红包（随机金额）
	RedPacketAllocTypeNormal = "normal" // 普通红包（固定金额）

	// 转账类型
	TransferTypePrivate = "private" // 私聊转账
	TransferTypeGroup   = "group"   // 群聊转账

	// 流水类型
	FlowTypeRecharge         = "recharge"          // 充值
	FlowTypeWithdraw         = "withdraw"          // 提现
	FlowTypeRedPacketSend    = "redpacket_send"    // 发红包
	FlowTypeRedPacketReceive = "redpacket_receive" // 收红包
	FlowTypeTransferSend     = "transfer_send"     // 转账支出
	FlowTypeTransferReceive  = "transfer_receive"  // 转账收入
	FlowTypeFee              = "fee"               // 手续费
	FlowTypeRefund           = "refund"            // 退款
)
