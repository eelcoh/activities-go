package activity

import "errors"

type Service struct {
	passphrase string
	now        func() int64
	newID      func() string
}

func NewService(passphrase string, now func() int64, newID func() string) Service {
	return Service{
		passphrase: passphrase,
		now:        now,
		newID:      newID,
	}
}

func (s Service) NewComment(input CommentInput) (Comment, error) {
	comment := Comment{
		Author:       input.Author,
		Msg:          cloneLines(input.Msg),
		Meta:         s.newMeta(),
		ActivityType: TypeComment,
	}

	if err := ValidateComment(comment); err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func (s Service) UpdateComment(existing Comment, input CommentInput) (Comment, error) {
	existing.Author = input.Author
	existing.Msg = cloneLines(input.Msg)
	existing.ActivityType = TypeComment

	if err := ValidateComment(existing); err != nil {
		return Comment{}, err
	}

	return existing, nil
}

func (s Service) NewBlog(input BlogInput) (Blog, error) {
	if err := s.authorize(input.Passphrase); err != nil {
		return Blog{}, err
	}

	blog := Blog{
		Author:       input.Author,
		Title:        input.Title,
		Msg:          cloneLines(input.Msg),
		Meta:         s.newMeta(),
		ActivityType: TypeBlog,
	}

	if err := ValidateBlog(blog); err != nil {
		return Blog{}, err
	}

	return blog, nil
}

func (s Service) UpdateBlog(existing Blog, input BlogInput) (Blog, error) {
	if err := s.authorize(input.Passphrase); err != nil {
		return Blog{}, err
	}

	existing.Author = input.Author
	existing.Title = input.Title
	existing.Msg = cloneLines(input.Msg)
	existing.ActivityType = TypeBlog

	if err := ValidateBlog(existing); err != nil {
		return Blog{}, err
	}

	return existing, nil
}

func (s Service) authorize(passphrase string) error {
	if passphrase == "" {
		return errors.New("passphrase cannot be empty")
	}
	if passphrase != s.passphrase {
		return ErrUnauthorized
	}

	return nil
}

func (s Service) newMeta() Meta {
	return Meta{
		Active: true,
		UUID:   s.newID(),
		Date:   s.now(),
	}
}

func cloneLines(lines []string) []string {
	return append([]string(nil), lines...)
}
