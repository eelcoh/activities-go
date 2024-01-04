package main

// # type Activity
// # = NewBet MetaData Name UUID
// # | Comment MetaData Author Message
// # | Blog MetaData Author Title Message
// # | NewRanking MetaData
type Activity interface {
	SetMetaData(meta MetaData) Activity
	GetMetaData() MetaData
	SetUUID(string) Activity
	GetUUID() string
	GetType() string
	IsActive() bool
	ToggleActive() Activity
}

type MetaData struct {
	Date   int64  `json:"date"`
	Active bool   `json:"active"`
	UUID   string `json:"uuid"`
}

type CommentMsg struct {
	Comment PlainComment `json:"comment"`
}

type PlainComment struct {
	Author string   `json:"author"`
	Msg    []string `json:"msg"`
}

type Comment struct {
	Author       string   `json:"author"`
	Msg          []string `json:"msg"`
	Meta         MetaData `json:"meta"`
	ActivityType string   `json:"type"`
}

type Activities struct {
	Activities []Activity `json:"activities"`
}

func (c Comment) SetMetaData(meta MetaData) Activity {
	c.Meta = meta
	return c
}

func (c Comment) GetMetaData() MetaData {
	return c.Meta
}

func (c Comment) GetUUID() string {
	return c.Meta.UUID
}

func (c Comment) SetUUID(uu string) Activity {
	c.Meta.UUID = uu
	return c
}

func (c Comment) GetType() string {
	return "comment"
}

func (c Comment) ToggleActive() Activity {
	c.Meta.Active = !c.Meta.Active
	return c
}

func (c Comment) IsActive() bool {
	return c.Meta.Active
}

type BlogMsg struct {
	Blog AuthenticatedBlog `json:"blog"`
}

type AuthenticatedBlog struct {
	Passphrase string   `json:"passphrase"`
	Author     string   `json:"author"`
	Title      string   `json:"title"`
	Msg        []string `json:"msg"`
}

type Blog struct {
	Author       string   `json:"author"`
	Title        string   `json:"title"`
	Msg          []string `json:"msg"`
	Meta         MetaData `json:"meta"`
	ActivityType string   `json:"type"`
}

func (b Blog) SetMetaData(meta MetaData) Activity {
	b.Meta = meta
	return b
}

func (b Blog) GetMetaData() MetaData {
	return b.Meta
}

func (b Blog) GetUUID() string {
	return b.Meta.UUID
}

func (b Blog) SetUUID(uu string) Activity {
	b.Meta.UUID = uu
	return b
}

func (b Blog) GetType() string {
	return "blog"
}

func (b Blog) ToggleActive() Activity {
	b.Meta.Active = !b.Meta.Active
	return b
}

func (b Blog) IsActive() bool {
	return b.Meta.Active
}

type NewBet struct {
	Name         string   `json:"name"`
	UUID         string   `json:"uuid"`
	Meta         MetaData `json:"meta"`
	ActivityType string   `json:"type"`
}

func (b NewBet) SetMetaData(meta MetaData) Activity {
	b.Meta = meta
	return b
}

func (b NewBet) GetMetaData() MetaData {
	return b.Meta
}

func (b NewBet) GetUUID() string {
	return b.Meta.UUID
}

func (b NewBet) SetUUID(uu string) Activity {
	b.Meta.UUID = uu
	return b
}

func (b NewBet) GetType() string {
	return "new bet"
}
func (b NewBet) ToggleActive() Activity {
	b.Meta.Active = !b.Meta.Active
	return b
}

func (b NewBet) IsActive() bool {
	return b.Meta.Active
}

type NewRanking struct {
	Meta         MetaData `json:"meta"`
	ActivityType string   `json:"type"`
}

func (r NewRanking) SetMetaData(meta MetaData) Activity {
	r.Meta = meta
	return r
}

func (r NewRanking) GetMetaData() MetaData {
	return r.Meta
}

func (r NewRanking) GetUUID() string {
	return r.Meta.UUID
}

func (r NewRanking) SetUUID(uu string) Activity {
	r.Meta.UUID = uu
	return r
}

func (r NewRanking) GetType() string {
	return "new ranking"
}

func (r NewRanking) ToggleActive() Activity {
	r.Meta.Active = !r.Meta.Active
	return r
}

func (r NewRanking) IsActive() bool {
	return r.Meta.Active
}

type StoredActivity struct {
	ID           string   `firestore:"id" json:"ActivityId"`
	Activity     Activity `firestore:"activity" json:"activity"`
	ActivityType string   `firestore:"activityType" json:"activityType"`
}
