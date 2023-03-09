package database

func getAllTables() []interface{} {
	return []interface{}{
		new(User),
		new(ChatGroup),
		new(ChatGroupUser),
		new(ChatGroupUserAddLog),
		new(Friend),
		new(FriendAddLog),
		// Code generated Begin; DO NOT EDIT.
		// Code generated End; DO NOT EDIT.
	}
}
