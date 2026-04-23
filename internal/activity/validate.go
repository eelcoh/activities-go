package activity

import "errors"

func ValidateBlog(blog Blog) error {
	if blog.Author == "" {
		return errors.New("author cannot be empty")
	}
	if blog.Title == "" {
		return errors.New("title cannot be empty")
	}
	if err := validateMultilineString(blog.Msg); err != nil {
		return err
	}

	return nil
}

func ValidateComment(comment Comment) error {
	if comment.Author == "" {
		return errors.New("author cannot be empty")
	}
	if err := validateMultilineString(comment.Msg); err != nil {
		return err
	}

	return nil
}

func validateMultilineString(msg []string) error {
	if len(msg) == 0 {
		return errors.New("message cannot be empty")
	}

	for _, line := range msg {
		if line != "" {
			return nil
		}
	}

	return errors.New("message cannot be empty")
}
