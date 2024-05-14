package credentials

import "github.com/Dafaque/sshaman/internal/credentials"

type Credentials struct {
	ID     int64  `db:"id"`
	UserID int64  `db:"user_id"`
	Space  string `db:"space"`
	credentials.Credentials
}
