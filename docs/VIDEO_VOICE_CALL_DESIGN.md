# 语音/视频通话功能设计文档

## 📋 功能概述

实现单聊和群聊的实时语音通话、视频通话功能，支持：
- 一对一语音通话
- 一对一视频通话
- 群聊语音通话（多人）
- 群聊视频通话（多人）

## 🏗️ 技术架构

### 核心技术栈
- **WebRTC**: 用于P2P音视频通信
- **WebSocket**: 用于信令传输（已实现）
- **STUN/TURN服务器**: 用于NAT穿透（可选，使用第三方服务如Google STUN）

### 架构设计

```
前端 (WebRTC)  <--信令-->  后端 (WebSocket信令服务器)  <--信令-->  前端 (WebRTC)
     |                                                              |
     |                                                              |
     +---------------------- P2P音视频流 ---------------------------+
```

## 📊 数据模型设计

### 1. 通话房间表 (call_rooms)

```go
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
```

### 2. 通话参与者表 (call_participants)

```go
type CallParticipant struct {
    ID          int64          `gorm:"primaryKey;comment:参与者ID"`
    RoomID      int64          `gorm:"column:room_id;not null;index;comment:房间ID"`
    UserID      int64          `gorm:"column:user_id;not null;index;comment:用户ID"`
    Status      string         `gorm:"column:status;type:varchar(20);not null;default:'invited';index;comment:状态(invited已邀请,ringing响铃中,joined已加入,rejected已拒绝,left已离开)"`
    JoinedAt    int64          `gorm:"column:joined_at;type:bigint;not null;default:0;comment:加入时间(时间戳)"`
    LeftAt      int64          `gorm:"column:left_at;type:bigint;not null;default:0;comment:离开时间(时间戳)"`
    Duration    int64          `gorm:"column:duration;type:bigint;not null;default:0;comment:参与时长(秒)"`
    CreatedAt   int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
    UpdatedAt   int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
}
```

### 3. 通话记录表 (call_records)

```go
type CallRecord struct {
    ID          int64          `gorm:"primaryKey;comment:记录ID"`
    RoomID      int64          `gorm:"column:room_id;not null;index;comment:房间ID"`
    UserID      int64          `gorm:"column:user_id;not null;index;comment:用户ID"`
    Type        string         `gorm:"column:type;type:varchar(20);not null;index;comment:类型(voice语音,video视频)"`
    CallType    string         `gorm:"column:call_type;type:varchar(20);not null;index;comment:通话类型(private私聊,group群聊)"`
    OtherUserID int64          `gorm:"column:other_user_id;index;comment:对方用户ID(私聊时使用)"`
    GroupID     int64          `gorm:"column:group_id;index;comment:群组ID(群聊时使用)"`
    Direction   string         `gorm:"column:direction;type:varchar(20);not null;index;comment:方向(incoming来电,outgoing去电)"`
    Status      string         `gorm:"column:status;type:varchar(20);not null;index;comment:状态(missed未接,answered已接,rejected已拒绝)"`
    Duration    int64          `gorm:"column:duration;type:bigint;not null;default:0;comment:通话时长(秒)"`
    StartedAt   int64          `gorm:"column:started_at;type:bigint;not null;default:0;comment:开始时间(时间戳)"`
    EndedAt     int64          `gorm:"column:ended_at;type:bigint;not null;default:0;comment:结束时间(时间戳)"`
    CreatedAt   int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
}
```

## 🔄 业务流程设计

### 1. 发起通话流程

```
1. 用户A发起通话请求
   ↓
2. 后端创建通话房间 (CallRoom)
   ↓
3. 后端创建参与者记录 (CallParticipant)
   ↓
4. 后端通过WebSocket推送通话邀请给用户B
   ↓
5. 用户B收到邀请，显示响铃界面
   ↓
6. 用户B接受/拒绝
   ↓
7. 如果接受，建立WebRTC连接
   ↓
8. 开始通话
```

### 2. WebRTC信令流程

```
1. 发起方创建Offer
   ↓
2. 通过WebSocket发送Offer给接收方
   ↓
3. 接收方创建Answer
   ↓
4. 通过WebSocket发送Answer给发起方
   ↓
5. 双方交换ICE Candidate
   ↓
6. 建立P2P连接
   ↓
7. 开始音视频传输
```

### 3. 群聊通话流程

```
1. 用户A发起群聊通话
   ↓
2. 后端创建通话房间
   ↓
3. 后端为所有群成员创建参与者记录
   ↓
4. 后端通过WebSocket推送邀请给所有在线成员
   ↓
5. 成员逐个加入房间
   ↓
6. 使用SFU (Selective Forwarding Unit) 模式
   - 每个成员发送流到服务器
   - 服务器转发给其他成员
   ↓
7. 开始多人通话
```

## 📡 WebSocket消息类型

### 信令消息类型

```go
const (
    // 通话相关
    SignalTypeCallInvite    = "call_invite"     // 通话邀请
    SignalTypeCallAccept   = "call_accept"     // 接受通话
    SignalTypeCallReject   = "call_reject"     // 拒绝通话
    SignalTypeCallCancel   = "call_cancel"     // 取消通话
    SignalTypeCallEnd      = "call_end"        // 结束通话
    SignalTypeCallBusy     = "call_busy"       // 通话中（占线）
    
    // WebRTC信令
    SignalTypeOffer        = "offer"           // WebRTC Offer
    SignalTypeAnswer       = "answer"          // WebRTC Answer
    SignalTypeIceCandidate = "ice_candidate"   // ICE Candidate
    SignalTypeIceComplete  = "ice_complete"    // ICE完成
    
    // 通话状态
    SignalTypeCallJoined   = "call_joined"     // 加入通话
    SignalTypeCallLeft    = "call_left"       // 离开通话
    SignalTypeCallMute     = "call_mute"       // 静音
    SignalTypeCallUnmute   = "call_unmute"     // 取消静音
    SignalTypeCallVideoOn  = "call_video_on"   // 开启视频
    SignalTypeCallVideoOff = "call_video_off"  // 关闭视频
)
```

### WebSocket消息格式

```json
{
  "type": "call_invite",
  "data": {
    "room_id": 123,
    "room_token": "abc123",
    "caller_id": 1,
    "caller_name": "用户A",
    "call_type": "private",
    "media_type": "video",
    "group_id": 0
  }
}
```

## 🛠️ API接口设计

### 1. 发起通话

```protobuf
// 发起私聊通话
rpc StartPrivateCall(StartPrivateCallRequest) returns (StartPrivateCallResponse);

message StartPrivateCallRequest {
  uint32 receiver_id = 1;  // 接收者ID
  string type = 2;         // voice/video
}

message StartPrivateCallResponse {
  uint64 room_id = 1;
  string room_token = 2;
  bool success = 3;
}
```

### 2. 发起群聊通话

```protobuf
// 发起群聊通话
rpc StartGroupCall(StartGroupCallRequest) returns (StartGroupCallResponse);

message StartGroupCallRequest {
  uint32 group_id = 1;     // 群组ID
  string type = 2;         // voice/video
}

message StartGroupCallResponse {
  uint64 room_id = 1;
  string room_token = 2;
  bool success = 3;
}
```

### 3. 接受/拒绝通话

```protobuf
// 接受通话
rpc AcceptCall(AcceptCallRequest) returns (AcceptCallResponse);

// 拒绝通话
rpc RejectCall(RejectCallRequest) returns (RejectCallResponse);

message AcceptCallRequest {
  uint64 room_id = 1;
}

message RejectCallRequest {
  uint64 room_id = 1;
  string reason = 2;  // 拒绝原因（可选）
}
```

### 4. 结束通话

```protobuf
// 结束通话
rpc EndCall(EndCallRequest) returns (EndCallResponse);

message EndCallRequest {
  uint64 room_id = 1;
}
```

### 5. 获取通话记录

```protobuf
// 获取通话记录
rpc GetCallRecords(GetCallRecordsRequest) returns (GetCallRecordsResponse);

message GetCallRecordsRequest {
  int32 page = 1;
  int32 page_size = 2;
  string type = 3;  // voice/video/all
}

message GetCallRecordsResponse {
  repeated CallRecord records = 1;
  int32 total_count = 2;
}
```

## 🎯 实现步骤

### 第一阶段：基础功能（单聊）

1. **数据模型**
   - [ ] 创建 CallRoom 模型
   - [ ] 创建 CallParticipant 模型
   - [ ] 创建 CallRecord 模型
   - [ ] 数据库迁移

2. **API接口**
   - [ ] 发起私聊通话接口
   - [ ] 接受/拒绝通话接口
   - [ ] 结束通话接口
   - [ ] 获取通话记录接口

3. **WebSocket信令**
   - [ ] 通话邀请信令
   - [ ] 接受/拒绝信令
   - [ ] WebRTC Offer/Answer/ICE信令
   - [ ] 通话状态更新信令

4. **业务逻辑**
   - [ ] 创建通话房间
   - [ ] 邀请参与者
   - [ ] 通话状态管理
   - [ ] 通话记录保存

### 第二阶段：群聊功能

1. **群聊通话**
   - [ ] 发起群聊通话
   - [ ] 多人加入/离开
   - [ ] SFU模式实现（可选，或使用P2P Mesh）

2. **权限管理**
   - [ ] 群主/管理员权限
   - [ ] 成员权限控制

### 第三阶段：优化功能

1. **通话质量**
   - [ ] 网络质量检测
   - [ ] 自适应码率
   - [ ] 回声消除

2. **用户体验**
   - [ ] 通话录音（可选）
   - [ ] 通话转文字（可选）
   - [ ] 美颜滤镜（视频）

## 🔐 安全考虑

1. **房间令牌验证**
   - 每个房间生成唯一token
   - 加入房间时验证token

2. **权限验证**
   - 验证用户是否有权限发起/加入通话
   - 群聊验证群成员身份

3. **防骚扰**
   - 限制通话频率
   - 黑名单功能

## 📝 注意事项

1. **STUN/TURN服务器**
   - 开发环境可以使用Google STUN服务器
   - 生产环境建议自建TURN服务器

2. **性能优化**
   - 群聊使用SFU模式（服务器转发）而非Mesh模式
   - 限制群聊最大人数（建议10-20人）

3. **前端实现**
   - 使用WebRTC API (getUserMedia, RTCPeerConnection)
   - 处理设备权限请求
   - 处理网络断开重连

4. **移动端**
   - iOS/Android需要使用原生WebRTC SDK
   - 处理后台运行限制

---

**最后更新**: 2025-01-18

