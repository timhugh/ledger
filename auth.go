package ledger

type User struct {
    UserUUID string
    Login string
    PasswordHash string
}

type Session struct {
    SessionUUID string
    UserUUID string
    CreatedAt string

    User *User
}

type SessionGetter interface {
    GetSession(uuid string) (*Session, error)
}
