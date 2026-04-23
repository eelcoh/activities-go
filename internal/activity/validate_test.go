package activity

import "testing"

func TestValidateComment(t *testing.T) {
	tests := []struct {
		name    string
		comment Comment
		wantErr bool
	}{
		{
			name: "valid comment",
			comment: Comment{
				Author: "Eelco",
				Msg:    []string{"hello"},
			},
		},
		{
			name: "missing author",
			comment: Comment{
				Msg: []string{"hello"},
			},
			wantErr: true,
		},
		{
			name: "missing message",
			comment: Comment{
				Author: "Eelco",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateComment(tt.comment)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateBlog(t *testing.T) {
	tests := []struct {
		name    string
		blog    Blog
		wantErr bool
	}{
		{
			name: "valid blog",
			blog: Blog{
				Author: "Eelco",
				Title:  "Title",
				Msg:    []string{"hello"},
			},
		},
		{
			name: "missing title",
			blog: Blog{
				Author: "Eelco",
				Msg:    []string{"hello"},
			},
			wantErr: true,
		},
		{
			name: "missing author",
			blog: Blog{
				Title: "Title",
				Msg:   []string{"hello"},
			},
			wantErr: true,
		},
		{
			name: "missing message",
			blog: Blog{
				Author: "Eelco",
				Title:  "Title",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBlog(tt.blog)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateBlog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceNewComment(t *testing.T) {
	svc := NewService("secret", func() int64 { return 42 }, func() string { return "uuid-1" })

	comment, err := svc.NewComment(CommentInput{
		Author: "Eelco",
		Msg:    []string{"hello"},
	})
	if err != nil {
		t.Fatalf("NewComment() error = %v", err)
	}

	if comment.Meta.UUID != "uuid-1" {
		t.Fatalf("comment.Meta.UUID = %q, want %q", comment.Meta.UUID, "uuid-1")
	}
	if comment.Meta.Date != 42 {
		t.Fatalf("comment.Meta.Date = %d, want %d", comment.Meta.Date, 42)
	}
	if !comment.Meta.Active {
		t.Fatal("comment.Meta.Active = false, want true")
	}
}
