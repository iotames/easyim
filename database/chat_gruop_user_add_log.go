package database

type ChatGroupUserAddLog struct {
	BaseModel `xorm:"extends"`
	GroupID   int64 `xorm:"notnull 'group_id'"`
	UserID    int64 `xorm:"notnull 'user_id'"`
	Status    uint8 `xorm:"SMALLINT notnull default(0)"` // 0未同意，1已同意，2已拒绝
}

func (g ChatGroupUserAddLog) TableName() string {
	return "chat_group_user_add_logs"
}
