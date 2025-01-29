package mailer

import "embed"

const (
	fromName            = "Golang Media"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitations.tmpl"
)

// const userWelcomeTemplate string = "user_invitations"

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}
