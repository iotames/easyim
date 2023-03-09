package database

type FriendAddLog struct {
	BaseModel `xorm:"extends"`
	UserID    int64  `xorm:"notnull 'user_id'"`
	FriendID  int64  `xorm:"notnull 'friend_id'"`
	Msg       string `xorm:"varchar(64) notnull"`
	Status    uint8  `xorm:"SMALLINT notnull default(0)"` // 0未同意，1已同意，2已拒绝
}

func (f FriendAddLog) TableName() string {
	return "friend_add_logs"
}
