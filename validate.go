package main

import (
	"errors"
)

func authenticateBlog(blog AuthenticatedBlog) error {

	if blog.Passphrase == "" {
		return errors.New("Passphrase cannot be empty")
	}

	if blog.Passphrase != passphrase {
		return errors.New("Passphrase incorrect")
	}
	return nil
}

func validateBlog(blog Blog) error {

	if blog.Author == "" {
		return errors.New("Author cannot be empty")
	}

	if blog.Title == "" {
		return errors.New("Title cannot be empty")
	}

	if err := validateMulilineString(blog.Msg); err != nil {
		return errors.New("Message cannot be empty")
	}

	return nil
}

func validateComment(comment Comment) error {

	if comment.Author == "" {
		return errors.New("Author cannot be empty")
	}

	if err := validateMulilineString(comment.Msg); err != nil {
		return errors.New("Message cannot be empty")
	}

	return nil
}

func validateMulilineString(msg []string) error {
	if len(msg) == 0 {
		return errors.New("Message cannot be empty")
	}

	isEmpty := true
	for _, l := range msg {
		if l != "" {
			isEmpty = false
		}
	}

	if isEmpty {
		return errors.New("Message cannot be empty")
	}

	return nil
}
