package database

type Member struct {
	ID     string `db:"id"`
	Guild  string `db:"guild"`
	Points int64  `db:"points"`
	Level  int    `db:"level"`
}

type GuildConfig struct {
	ID          string `db:"id"`
	QRCode      bool   `db:"qrcode"`
	CaseChannel int64  `db:"CaseChannel"`
}

type Config struct {
	ID          string `db:"id"`
	Prefix      string `db:"prefix"`
	CaseChannel int64  `db:"case"`
}

type Case struct {
	GuildID               int64  `db:"gid"`
	ID                    int    `db:"id"`
	UserID                string `db:"uid"`
	ModID                 int64  `db:"mid"`
	AuthorUsernameDiscrim string `db:"aud"`
	Message               string `db:"reason"`
	Type                  string `db:"type"`
}
