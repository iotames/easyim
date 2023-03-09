package database

type ChatGroup struct {
	BaseModel `xorm:"extends"`
	UserID    int64  `xorm:"notnull 'user_id'"` // 创建者user_id
	Name      string `xorm:"varchar(64) notnull"`
	// 进群是否审核.0不需要1需要.默认0 需要审核时使用GroupUserAddLog
	JoinCheck bool  `xorm:"notnull default(0) 'join_check'"`
	Status    uint8 `xorm:"SMALLINT notnull default(0)"` // 0正常 1全体禁言
}

func (c ChatGroup) TableName() string {
	return "chat_groups"
}
