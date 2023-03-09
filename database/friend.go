package database

type Friend struct {
	BaseModel `xorm:"extends"`
	UserID    int64  `xorm:"notnull 'user_id'"`
	FriendID  int64  `xorm:"notnull 'friend_id'"`
	Remark    string `xorm:"varchar(64) notnull"`
	Sort      int    `xorm:"notnull default(0)"` // 按从大到小排序。置顶时, sort值+1
	// Status    uint8  `xorm:"SMALLINT notnull default(0)"`
}

func (f Friend) TableName() string {
	return "friends"
}
