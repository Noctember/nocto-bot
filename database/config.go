package database

func GetPrefix(id int64) string {
	var g Config
	err := Get(&g, "SELECT prefix FROM guilds WHERE id = $1", id)
	if err != nil {
		return "="
	}
	return g.Prefix
}

func IsPlonked(id int64) bool {
	var plonked bool
	err := Get(&plonked, "SELECT plonked FROM plonked WHERE id = $1", id)
	if err != nil {
		return false
	}
	return plonked
}

func CreateCase(p *Case) error {
	count := 0
	err := Get(&count, "SELECT count(*) FROM modlogs")
	p.ID = count
	if err == nil {
		_, err := Exec(`INSERT INTO modlogs ("id","gid", "uid", "mid", "aud", "reason", "type") VALUES ($1, $2, $3, $4, $5, $6, $7)`, count+1, p.GuildID, p.UserID, p.ModID, p.AuthorUsernameDiscrim, p.Message, p.Type)
		return err
	}
	return err
}

func GetCases(id, gid int64) []Case {
	cases := []Case{}
	err := Select(&cases, "SELECT * FROM modlogs WHERE uid = $1 AND gid = $2", id, gid)
	if err != nil {
		return nil
	}
	//return cases
	return cases
}

func GetConfig(id int64) (*Config, error) {
	var c Config
	err := Get(&c, "SELECT * FROM guilds WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetGConfig(id int64) (*GuildConfig, error) {
	var c GuildConfig
	err := Get(&c, "SELECT * FROM config WHERE id = $1", id)
	if err != nil {
		_, err := Exec(`INSERT INTO config ("id", "qrcode", "case") VALUES ($1, $2, $3)`, id, false, 0)
		err = Get(&c, "SELECT * FROM config WHERE id = $1", id)
		return &c, err
	}
	return &c, nil
}

func SetValue(id int64, key, val string) error {
	_, err := Exec("UPDATE config SET "+key+" = $1 WHERE id = $2", val, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateConfig(id int64) error {
	_, err := Exec(`INSERT INTO guilds ("id", "prefix", "case") VALUES ($1, $2, "")`, id, "=")
	if err != nil {
		return err
	}
	return nil
}

func SetPrefix(guild int64, prefix string) error {
	_, err := Exec("UPDATE guilds SET PREFIX = $1 WHERE id = $2", prefix, guild)
	return err
}
