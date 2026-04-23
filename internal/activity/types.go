package activity

import "errors"

const (
	TypeComment    = "comment"
	TypeBlog       = "blog"
	TypeNewBet     = "new bet"
	TypeNewRanking = "new ranking"
)

var (
	ErrNotFound     = errors.New("activity not found")
	ErrUnauthorized = errors.New("passphrase incorrect")
)

type Activity interface {
	activityKind()
	Type() string
	GetMeta() Meta
}

type Meta struct {
	Date   int64  `firestore:"date" json:"date"`
	Active bool   `firestore:"active" json:"active"`
	UUID   string `firestore:"uuid" json:"uuid"`
}

type CommentEnvelope struct {
	Comment CommentInput `json:"comment"`
}

type CommentInput struct {
	Author string   `json:"author"`
	Msg    []string `json:"msg"`
}

type Comment struct {
	Author       string   `firestore:"author" json:"author"`
	Msg          []string `firestore:"msg" json:"msg"`
	Meta         Meta     `firestore:"meta" json:"meta"`
	ActivityType string   `firestore:"type" json:"type"`
}

func (Comment) activityKind() {}

func (c Comment) Type() string {
	return TypeComment
}

func (c Comment) GetMeta() Meta {
	return c.Meta
}

type BlogEnvelope struct {
	Blog BlogInput `json:"blog"`
}

type BlogInput struct {
	Passphrase string   `json:"passphrase"`
	Author     string   `json:"author"`
	Title      string   `json:"title"`
	Msg        []string `json:"msg"`
}

type Blog struct {
	Author       string   `firestore:"author" json:"author"`
	Title        string   `firestore:"title" json:"title"`
	Msg          []string `firestore:"msg" json:"msg"`
	Meta         Meta     `firestore:"meta" json:"meta"`
	ActivityType string   `firestore:"type" json:"type"`
}

func (Blog) activityKind() {}

func (b Blog) Type() string {
	return TypeBlog
}

func (b Blog) GetMeta() Meta {
	return b.Meta
}

type NewBet struct {
	Name         string `firestore:"name" json:"name"`
	Meta         Meta   `firestore:"meta" json:"meta"`
	ActivityType string `firestore:"type" json:"type"`
}

func (NewBet) activityKind() {}

func (b NewBet) Type() string {
	return TypeNewBet
}

func (b NewBet) GetMeta() Meta {
	return b.Meta
}

type NewRanking struct {
	Meta         Meta   `firestore:"meta" json:"meta"`
	ActivityType string `firestore:"type" json:"type"`
}

func (NewRanking) activityKind() {}

func (r NewRanking) Type() string {
	return TypeNewRanking
}

func (r NewRanking) GetMeta() Meta {
	return r.Meta
}

type ListResponse struct {
	Activities []Activity `json:"activities"`
}
