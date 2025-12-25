# 用户相关扩展设计总结

## ✅ 已完成的工作

### 1. 扩展UserBasic模型 (`models/userBasic.go`)

**新增字段**:
- `LoginAccount` - 登录账号(5-32位，字母开头，不可修改，唯一索引)
- `PhoneVerified` - 手机号是否已验证
- `Email1` - 邮箱2
- `Email2` - 邮箱3
- `PrivacySettings` - 隐私设置(JSON格式)
- `LoginFailCount` - 登录失败次数
- `LockedUntil` - 锁定到期时间(时间戳，0表示未锁定)
- `LastActiveTime` - 最后活跃时间(时间戳)

### 2. 创建资金系统模型 (`models/balance.go`)

**所有金额使用int64存储，单位：分**
- 1元 = 100分
- 0.01元 = 1分
- 5000元 = 500000分

**数据模型**:
1. **UserBalance** - 用户余额表
   - `user_id` (主键)
   - `balance` (余额，单位：分)
   - `updated_at`

2. **RechargeRequest** - 充值申请表
   - 状态: pending(待审核), approved(已通过), rejected(已驳回)
   - 支持审核流程

3. **WithdrawRequest** - 提现申请表
   - 状态: pending(待审核), approved(已通过), rejected(已驳回)
   - 支持手续费、收款账户信息(加密存储)

4. **RedPacket** - 红包记录表
   - 类型: private(私聊), group(群聊)
   - 状态: sent(已发送), received(已领取), expired(已过期), refunded(已退回)
   - 支持过期时间

5. **RedPacketReceive** - 红包领取记录表(群聊时使用)
   - 记录每个成员的领取情况

6. **Transfer** - 转账记录表
   - 类型: private(私聊), group(群聊)
   - 状态: pending(待收款), received(已收款), cancelled(已撤销), refunded(已退回)

7. **BalanceFlow** - 资金流水表
   - 记录所有资金操作流水
   - 类型: recharge, withdraw, redpacket_send, redpacket_receive, transfer_send, transfer_receive, fee, refund

### 3. 创建设备管理模型 (`models/device.go`)

**UserDevice** - 用户设备表(多端登录管理)
- `device_type`: ios, android, web_mobile, pc_desktop, pc_web
- `device_id`: 设备唯一标识
- `token`: 登录Token
- `last_active_at`: 最后活跃时间

### 4. 更新Protocol Buffers定义 (`api/api/user/user.proto`)

**扩展UserInfo消息**:
- 添加所有新字段

**新增API**:
- `ChangePassword` - 修改密码
- `ResetPassword` - 重置密码
- `BindPhone` - 绑定手机号
- `BindEmail` - 绑定邮箱
- `UpdatePrivacySettings` - 更新隐私设置
- `GetDeviceList` - 获取设备列表
- `LogoutDevice` - 下线设备

### 5. 创建DAO层

**BalanceDao** (`internal/infrastructure/dao/balance_dao.go`):
- `GetBalance` - 获取用户余额
- `UpdateBalance` - 更新余额(事务中使用)
- `GetBalanceForUpdate` - 获取余额并加锁(SELECT FOR UPDATE)
- `CreateRechargeRequest` - 创建充值申请
- `GetRechargeRequests` - 获取充值申请列表
- `CreateWithdrawRequest` - 创建提现申请
- `GetWithdrawRequests` - 获取提现申请列表
- `GetTodayWithdrawCount` - 获取今日提现次数
- `CreateRedPacket` - 创建红包
- `GetRedPacket` - 获取红包
- `CreateRedPacketReceive` - 创建红包领取记录
- `CreateTransfer` - 创建转账
- `CreateBalanceFlow` - 创建资金流水
- `GetBalanceFlows` - 获取资金流水列表
- `GetExpiredRedPackets` - 获取过期红包(定时任务用)

**DeviceDao** (`internal/infrastructure/dao/device_dao.go`):
- `CreateDevice` - 创建设备
- `GetUserDevices` - 获取用户所有设备
- `UpdateDeviceToken` - 更新设备Token
- `UpdateDeviceActiveTime` - 更新活跃时间
- `DeleteDevice` - 删除设备(下线)
- `DeleteUserDevices` - 删除用户所有设备(除指定设备外)

### 6. 创建金额工具函数 (`utils/money.go`)

**函数**:
- `YuanToFen` - 元转分
- `YuanToFenFromString` - 从字符串元转分(推荐，避免精度问题)
- `FenToYuan` - 分转元
- `FormatMoney` - 格式化金额显示(保留2位小数)
- `FormatMoneyYuan` - 格式化金额显示(带"元"单位)

### 7. 更新数据库迁移 (`cmd/setup.go`)

添加了所有新表的AutoMigrate:
- user_balances (用户余额表)
- recharge_requests (充值申请表)
- withdraw_requests (提现申请表)
- red_packets (红包记录表)
- red_packet_receives (红包领取记录表)
- transfers (转账记录表)
- balance_flows (资金流水表)
- user_devices (用户设备表)

## 📋 使用示例

### 金额转换示例

```go
import "gochat/utils"

// 元转分
fen := utils.YuanToFen(1.23)  // 123分
fen, _ := utils.YuanToFenFromString("1.23")  // 123分(推荐)

// 分转元
yuan := utils.FenToYuan(123)  // 1.23

// 格式化显示
str := utils.FormatMoney(123)  // "1.23"
str := utils.FormatMoneyYuan(123)  // "1.23元"
```

### 余额操作示例(需要在事务中)

```go
// 在事务中操作余额
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// 获取余额并加锁
balance, err := balanceDao.GetBalanceForUpdate(tx, userID)
if err != nil {
    tx.Rollback()
    return err
}

// 检查余额
if balance.Balance < amount {
    tx.Rollback()
    return errors.New("余额不足")
}

// 更新余额
if err := tx.Model(&models.UserBalance{}).
    Where("user_id = ?", userID).
    Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
    tx.Rollback()
    return err
}

// 记录流水
flow := &models.BalanceFlow{
    UserID:  userID,
    Type:    models.FlowTypeRedPacketSend,
    Amount:  -amount,  // 负数表示支出
    Balance: balance.Balance - amount,
    Remark:  "发送红包",
}
if err := tx.Create(flow).Error; err != nil {
    tx.Rollback()
    return err
}

// 提交事务
return tx.Commit().Error
```

## ⚠️ 注意事项

1. **金额存储**: 所有金额必须使用int64存储，单位是"分"，不要使用float或decimal
2. **事务保证**: 所有余额变更操作必须在数据库事务中完成
3. **并发控制**: 使用`GetBalanceForUpdate`获取余额时加锁，防止并发问题
4. **精度问题**: 金额转换时，推荐使用`YuanToFenFromString`避免浮点数精度问题
5. **流水记录**: 所有资金操作必须记录到BalanceFlow表，便于审计

## 🔄 下一步工作

1. 实现验证码服务(Redis存储)
2. 实现用户注册/登录逻辑(支持验证码、登录账号规则)
3. 实现资金服务(充值、提现、红包、转账)
4. 实现设备管理服务(多端登录)
5. 实现Handler层业务逻辑
6. 更新路由注册

---

**创建时间**: 2025-01-18  
**最后更新**: 2025-01-18



