package handler

import "database/sql"

type CommentVM struct {
	AvatarURL string
	UserName  string
	CreatedAt string
	CommentID int32
	Content   string
	Image     string
	PostID    int
	ParentID  int
}

type PostVM struct {
	UserAvatar string
	UserName   string
	CreatedAt  string
	PostID     int32
	Title      string
	Content    string
	Comments   []CommentVM
	Image      string
	Deleted    sql.NullTime
}
