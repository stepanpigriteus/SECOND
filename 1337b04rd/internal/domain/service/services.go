package service

type AllServices struct {
	Post    PostService
	User    UserService
	Comment CommentService
	Session SessionService
	Storage StorageService
}
