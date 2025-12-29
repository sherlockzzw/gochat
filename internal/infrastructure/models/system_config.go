package models

// SystemConfig 系统配置
type SystemConfig struct {
	ID          int64  `gorm:"primaryKey;comment:配置ID"`
	ConfigKey   string `gorm:"column:config_key;type:varchar(100);not null;uniqueIndex;comment:配置键"`
	ConfigValue string `gorm:"column:config_value;type:json;comment:配置值JSON"`
	Description string `gorm:"column:description;type:text;comment:描述"`
	UpdatedAt   int64  `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}

// 配置键常量
const (
	// 表情包配置
	ConfigKeyEmojiAllowCustom = "emoji_allow_custom"    // 允许用户自定义
	ConfigKeyEmojiMaxSize     = "emoji_max_size"        // 文件大小限制(MB)
	ConfigKeyEmojiMaxCount    = "emoji_max_count"       // 用户最大表情包数量
	
	// 易翻译配置
	ConfigKeyTranslateEnabled     = "translate_enabled"      // 全局启用
	ConfigKeyTranslateDailyLimit  = "translate_daily_limit"  // 每日调用次数限制
	ConfigKeyTranslateApiKey       = "translate_api_key"      // API密钥
	ConfigKeyTranslateLanguages    = "translate_languages"    // 支持的语言列表
	
	// 群组配置
	ConfigKeyGroupAllowCreate      = "group_allow_create"     // 允许创建群组
	ConfigKeyGroupRedPacketEnabled = "group_redpacket_enabled" // 群红包启用
	ConfigKeyGroupRedPacketLimit   = "group_redpacket_limit"  // 群红包每日次数上限
	
	// 红包配置
	ConfigKeyRedPacketPrivateMin = "redpacket_private_min"    // 私聊红包最小金额
	ConfigKeyRedPacketPrivateMax = "redpacket_private_max"    // 私聊红包最大金额
	ConfigKeyRedPacketPrivateCount = "redpacket_private_count" // 私聊红包个数上限
	ConfigKeyRedPacketGroupMin   = "redpacket_group_min"      // 群聊红包最小金额
	ConfigKeyRedPacketGroupMax   = "redpacket_group_max"      // 群聊红包最大金额
	ConfigKeyRedPacketGroupCount = "redpacket_group_count"    // 群聊红包个数上限
	
	// 资金配置
	ConfigKeyFinanceRechargeEnabled = "finance_recharge_enabled" // 充值启用
	ConfigKeyFinanceWithdrawEnabled = "finance_withdraw_enabled" // 提现启用
	ConfigKeyFinanceRechargeDailyLimit = "finance_recharge_daily_limit" // 单用户单日充值上限
	ConfigKeyFinanceWithdrawDailyLimit = "finance_withdraw_daily_limit" // 单用户单日提现上限
	ConfigKeyFinanceWithdrawDailyCount = "finance_withdraw_daily_count" // 每日提现次数限制
	ConfigKeyFinanceWithdrawFeeRate = "finance_withdraw_fee_rate" // 提现手续费率
	
	// 功能开关
	ConfigKeyFeatureDeleteMessage = "feature_delete_message" // 双向删除消息功能
	ConfigKeyFeatureRecallMessage = "feature_recall_message" // 撤回消息功能
	ConfigKeyFeatureCreateGroup   = "feature_create_group"   // 创建群组功能
)

