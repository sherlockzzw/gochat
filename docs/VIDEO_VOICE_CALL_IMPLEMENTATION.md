# 语音/视频通话功能实现计划

## 📁 文件结构

```
gochat/
├── api/api/call/
│   └── call.proto                    # 通话相关API定义
├── internal/
│   ├── infrastructure/
│   │   ├── models/
│   │   │   └── call.go              # 通话相关数据模型
│   │   └── dao/
│   │       └── call_dao.go         # 通话数据访问层
│   ├── application/
│   │   └── handler/
│   │       └── call/
│   │           ├── call.go          # 通话Handler
│   │           ├── room.go          # 房间管理
│   │           ├── signal.go        # 信令处理
│   │           └── record.go        # 通话记录
│   └── infrastructure/
│       └── websocket/
│           └── call_handler.go      # WebSocket通话信令处理
└── cmd/
    └── setup.go                     # 添加通话表迁移
```

## 🗄️ 数据库表设计

### 1. call_rooms 表

```sql
CREATE TABLE `call_rooms` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '房间ID',
  `type` varchar(20) NOT NULL COMMENT '类型(voice语音,video视频)',
  `call_type` varchar(20) NOT NULL COMMENT '通话类型(private私聊,group群聊)',
  `creator_id` bigint NOT NULL COMMENT '创建者ID',
  `group_id` bigint DEFAULT '0' COMMENT '群组ID(群聊时使用,私聊时为0)',
  `room_token` varchar(100) NOT NULL COMMENT '房间令牌',
  `status` varchar(20) NOT NULL DEFAULT 'calling' COMMENT '状态(calling呼叫中,ringing响铃中,connected已连接,ended已结束)',
  `started_at` bigint NOT NULL DEFAULT '0' COMMENT '开始时间(时间戳)',
  `ended_at` bigint NOT NULL DEFAULT '0' COMMENT '结束时间(时间戳)',
  `duration` bigint NOT NULL DEFAULT '0' COMMENT '通话时长(秒)',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间(时间戳)',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '更新时间(时间戳)',
  `deleted_at` bigint DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_call_type` (`call_type`),
  KEY `idx_creator_id` (`creator_id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_status` (`status`),
  UNIQUE KEY `idx_room_token` (`room_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通话房间表';
```

### 2. call_participants 表

```sql
CREATE TABLE `call_participants` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '参与者ID',
  `room_id` bigint NOT NULL COMMENT '房间ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `status` varchar(20) NOT NULL DEFAULT 'invited' COMMENT '状态(invited已邀请,ringing响铃中,joined已加入,rejected已拒绝,left已离开)',
  `joined_at` bigint NOT NULL DEFAULT '0' COMMENT '加入时间(时间戳)',
  `left_at` bigint NOT NULL DEFAULT '0' COMMENT '离开时间(时间戳)',
  `duration` bigint NOT NULL DEFAULT '0' COMMENT '参与时长(秒)',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间(时间戳)',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '更新时间(时间戳)',
  PRIMARY KEY (`id`),
  KEY `idx_room_id` (`room_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  UNIQUE KEY `idx_room_user` (`room_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通话参与者表';
```

### 3. call_records 表

```sql
CREATE TABLE `call_records` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `room_id` bigint NOT NULL COMMENT '房间ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `type` varchar(20) NOT NULL COMMENT '类型(voice语音,video视频)',
  `call_type` varchar(20) NOT NULL COMMENT '通话类型(private私聊,group群聊)',
  `other_user_id` bigint DEFAULT '0' COMMENT '对方用户ID(私聊时使用)',
  `group_id` bigint DEFAULT '0' COMMENT '群组ID(群聊时使用)',
  `direction` varchar(20) NOT NULL COMMENT '方向(incoming来电,outgoing去电)',
  `status` varchar(20) NOT NULL COMMENT '状态(missed未接,answered已接,rejected已拒绝)',
  `duration` bigint NOT NULL DEFAULT '0' COMMENT '通话时长(秒)',
  `started_at` bigint NOT NULL DEFAULT '0' COMMENT '开始时间(时间戳)',
  `ended_at` bigint NOT NULL DEFAULT '0' COMMENT '结束时间(时间戳)',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间(时间戳)',
  PRIMARY KEY (`id`),
  KEY `idx_room_id` (`room_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_call_type` (`call_type`),
  KEY `idx_other_user_id` (`other_user_id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_direction` (`direction`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通话记录表';
```

## 🔧 核心实现

### 1. 数据模型 (models/call.go)

```go
package models

import "gorm.io/gorm"

// CallRoom 通话房间
type CallRoom struct {
    ID          int64          `gorm:"primaryKey;comment:房间ID"`
    Type        string         `gorm:"column:type;type:varchar(20);not null;index;comment:类型(voice语音,video视频)"`
    CallType    string         `gorm:"column:call_type;type:varchar(20);not null;index;comment:通话类型(private私聊,group群聊)"`
    CreatorID   int64          `gorm:"column:creator_id;not null;index;comment:创建者ID"`
    GroupID     int64          `gorm:"column:group_id;index;comment:群组ID(群聊时使用,私聊时为0)"`
    RoomToken   string         `gorm:"column:room_token;type:varchar(100);uniqueIndex;not null;comment:房间令牌"`
    Status      string         `gorm:"column:status;type:varchar(20);not null;default:'calling';index;comment:状态(calling呼叫中,ringing响铃中,connected已连接,ended已结束)"`
    StartedAt   int64          `gorm:"column:started_at;type:bigint;not null;default:0;comment:开始时间(时间戳)"`
    EndedAt     int64          `gorm:"column:ended_at;type:bigint;not null;default:0;comment:结束时间(时间戳)"`
    Duration    int64          `gorm:"column:duration;type:bigint;not null;default:0;comment:通话时长(秒)"`
    CreatedAt   int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
    UpdatedAt   int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
    DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

func (CallRoom) TableName() string {
    return "call_rooms"
}

// CallParticipant 通话参与者
type CallParticipant struct {
    ID        int64  `gorm:"primaryKey;comment:参与者ID"`
    RoomID    int64  `gorm:"column:room_id;not null;index;comment:房间ID"`
    UserID    int64  `gorm:"column:user_id;not null;index;comment:用户ID"`
    Status    string `gorm:"column:status;type:varchar(20);not null;default:'invited';index;comment:状态(invited已邀请,ringing响铃中,joined已加入,rejected已拒绝,left已离开)"`
    JoinedAt  int64  `gorm:"column:joined_at;type:bigint;not null;default:0;comment:加入时间(时间戳)"`
    LeftAt    int64  `gorm:"column:left_at;type:bigint;not null;default:0;comment:离开时间(时间戳)"`
    Duration  int64  `gorm:"column:duration;type:bigint;not null;default:0;comment:参与时长(秒)"`
    CreatedAt int64  `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
    UpdatedAt int64  `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
}

func (CallParticipant) TableName() string {
    return "call_participants"
}

// CallRecord 通话记录
type CallRecord struct {
    ID          int64  `gorm:"primaryKey;comment:记录ID"`
    RoomID      int64  `gorm:"column:room_id;not null;index;comment:房间ID"`
    UserID      int64  `gorm:"column:user_id;not null;index;comment:用户ID"`
    Type        string `gorm:"column:type;type:varchar(20);not null;index;comment:类型(voice语音,video视频)"`
    CallType    string `gorm:"column:call_type;type:varchar(20);not null;index;comment:通话类型(private私聊,group群聊)"`
    OtherUserID int64  `gorm:"column:other_user_id;index;comment:对方用户ID(私聊时使用)"`
    GroupID     int64  `gorm:"column:group_id;index;comment:群组ID(群聊时使用)"`
    Direction   string `gorm:"column:direction;type:varchar(20);not null;index;comment:方向(incoming来电,outgoing去电)"`
    Status      string `gorm:"column:status;type:varchar(20);not null;index;comment:状态(missed未接,answered已接,rejected已拒绝)"`
    Duration    int64  `gorm:"column:duration;type:bigint;not null;default:0;comment:通话时长(秒)"`
    StartedAt   int64  `gorm:"column:started_at;type:bigint;not null;default:0;comment:开始时间(时间戳)"`
    EndedAt     int64  `gorm:"column:ended_at;type:bigint;not null;default:0;comment:结束时间(时间戳)"`
    CreatedAt   int64  `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
}

func (CallRecord) TableName() string {
    return "call_records"
}

// 常量定义
const (
    // 通话类型
    CallTypeVoice = "voice"
    CallTypeVideo = "video"
    
    // 通话场景
    CallScenePrivate = "private"
    CallSceneGroup   = "group"
    
    // 房间状态
    RoomStatusCalling   = "calling"   // 呼叫中
    RoomStatusRinging   = "ringing"    // 响铃中
    RoomStatusConnected = "connected" // 已连接
    RoomStatusEnded     = "ended"     // 已结束
    
    // 参与者状态
    ParticipantStatusInvited = "invited"  // 已邀请
    ParticipantStatusRinging = "ringing"  // 响铃中
    ParticipantStatusJoined  = "joined"   // 已加入
    ParticipantStatusRejected = "rejected" // 已拒绝
    ParticipantStatusLeft    = "left"     // 已离开
    
    // 通话方向
    CallDirectionIncoming = "incoming" // 来电
    CallDirectionOutgoing = "outgoing" // 去电
    
    // 通话记录状态
    RecordStatusMissed   = "missed"   // 未接
    RecordStatusAnswered = "answered" // 已接
    RecordStatusRejected = "rejected" // 已拒绝
)
```

### 2. WebSocket信令消息格式

```go
// CallSignalMessage 通话信令消息
type CallSignalMessage struct {
    Type      string      `json:"type"`       // 信令类型
    RoomID    int64       `json:"room_id"`    // 房间ID
    RoomToken string      `json:"room_token"` // 房间令牌
    FromUserID int64      `json:"from_user_id"` // 发送者ID
    ToUserID  int64       `json:"to_user_id"`   // 接收者ID（可选）
    Data      interface{} `json:"data"`       // 信令数据
}

// WebRTC信令数据
type WebRTCSignalData struct {
    SDP        string `json:"sdp"`         // SDP信息
    Type       string `json:"type"`        // offer/answer
    Candidate  string `json:"candidate"`   // ICE Candidate
    CandidateType string `json:"candidate_type"` // host/srflx/relay
}
```

## 🚀 实现优先级

### Phase 1: 单聊语音/视频通话（核心功能）
1. 数据模型和数据库迁移
2. 发起私聊通话API
3. WebSocket信令处理（邀请/接受/拒绝）
4. WebRTC信令交换（Offer/Answer/ICE）
5. 结束通话和记录保存

### Phase 2: 群聊通话
1. 发起群聊通话API
2. 多人加入/离开处理
3. 群聊信令广播

### Phase 3: 优化和增强
1. 通话录音
2. 通话转文字
3. 美颜滤镜
4. 网络质量检测

---

**建议**: 先从单聊语音通话开始实现，验证整个流程后再扩展到视频和群聊。

