package database

type ChatGroupUser struct {
	BaseModel `xorm:"extends"`
	GroupID   int64 `xorm:"notnull 'group_id'"`          // 聊天群ID
	UserID    int64 `xorm:"notnull 'user_id'"`           // 聊天群成员用户ID
	Status    uint8 `xorm:"SMALLINT notnull default(0)"` // 0正常 1禁言
	Sort      int   `xorm:"notnull default(0)"`          // 按从大到小排序。置顶时, sort值+1
	BeQuiet   bool  `xorm:"notnull 'be_quiet'"`          // 消息免打扰
}

func (c ChatGroupUser) TableName() string {
	return "chat_group_users"
}
