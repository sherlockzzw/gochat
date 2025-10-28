// MongoDB初始化脚本
db = db.getSiblingDB('gochat');

// 创建用户
db.createUser({
  user: 'gochat',
  pwd: 'gochat123',
  roles: [
    {
      role: 'readWrite',
      db: 'gochat'
    }
  ]
});

// 创建集合和索引
db.createCollection('messages');
db.createCollection('chatrooms');
db.createCollection('friends');
db.createCollection('messagereads');
db.createCollection('unreadcounts');

// 创建消息集合的索引
db.messages.createIndex({ "from_user_id": 1 });
db.messages.createIndex({ "to_user_id": 1 });
db.messages.createIndex({ "room_id": 1 });
db.messages.createIndex({ "created_at": -1 });
db.messages.createIndex({ "from_user_id": 1, "to_user_id": 1, "created_at": -1 });

// 创建聊天室集合的索引
db.chatrooms.createIndex({ "owner_id": 1 });
db.chatrooms.createIndex({ "members": 1 });
db.chatrooms.createIndex({ "type": 1 });

// 创建好友关系集合的索引
db.friends.createIndex({ "user_id": 1 });
db.friends.createIndex({ "friend_id": 1 });
db.friends.createIndex({ "user_id": 1, "friend_id": 1 }, { unique: true });

// 创建消息已读记录集合的索引
db.messagereads.createIndex({ "message_id": 1 });
db.messagereads.createIndex({ "user_id": 1 });

// 创建未读计数集合的索引
db.unreadcounts.createIndex({ "user_id": 1 });
db.unreadcounts.createIndex({ "from_user_id": 1 });
db.unreadcounts.createIndex({ "room_id": 1 });

print('MongoDB初始化完成！');
