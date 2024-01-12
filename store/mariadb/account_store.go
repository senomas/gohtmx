package mariadb

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/senomas/gohtmx/store"
)

type MariadbAccountStore struct {
	db       *sqlx.DB
	maxLimit int
}

func init() {
	store.AddAccountStore("mariadb", func() store.AccountStore {
		v := MariadbAccountStore{}
		return v.init()
	})
}

var err_duplicate_rx = regexp.MustCompile(`^Error 1062 \((?P<code>\d+)\): Duplicate entry '(?P<value>[^']+)' for key '(?P<field>[^']+)'$`)

func (s *MariadbAccountStore) init() store.AccountStore {
	url := os.Getenv("DB_URL")
	if url == "" {
		url = "root:dodol123@tcp(localhost:13306)/test"
	}
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		panic(fmt.Errorf("error opening database [%s]: %v", url, err))
	}

	err = db.Ping()
	for i := 0; err != nil && i < 2; i++ {
		if strings.Contains(err.Error(), "bad connection") ||
			strings.Contains(err.Error(), "connection refused") {
			time.Sleep(1 * time.Second)
			err = db.Ping()
		} else {
			break
		}
	}
	if err != nil {
		panic(fmt.Errorf("error ping database [%s]: %v", url, err))
	}

	qry := `CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    UNIQUE(name),
    UNIQUE(email)
  )`
	_, err = db.Exec(qry)
	if err != nil {
		panic(fmt.Errorf("error creating table: %v\n\n%s", err, qry))
	}

	qry = `CREATE TABLE IF NOT EXISTS privilege (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    UNIQUE(name)
  )`
	_, err = db.Exec(qry)
	if err != nil {
		panic(fmt.Errorf("error creating table: %v\n\n%s", err, qry))
	}

	qry = `CREATE TABLE IF NOT EXISTS user_privilege (
    user INTEGER NOT NULL,
    privilege INTEGER NOT NULL,
    UNIQUE(user, privilege),
    FOREIGN KEY(user) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY(privilege) REFERENCES privilege(id)
  )`
	_, err = db.Exec(qry)
	if err != nil {
		panic(fmt.Errorf("error creating table: %v\n\n%s", err, qry))
	}
	ctx := MariadbAccountStore{
		db: db,
	}

	maxLimit := os.Getenv("DB_MAX_LIMIT")
	if maxLimit == "" {
		ctx.maxLimit = 100
	} else {
		v, err := strconv.ParseInt(maxLimit, 10, 32)
		if err != nil {
			panic(fmt.Errorf("invalid DB_MAX_LIMIT '%s': %v", maxLimit, err))
		}
		ctx.maxLimit = int(v)
	}
	return &ctx
}

func (s *MariadbAccountStore) Close() error {
	return s.db.Close()
}

func (s *MariadbAccountStore) ValidLimit(limit int) bool {
	return limit > 0 && limit <= s.maxLimit
}

func (s *MariadbAccountStore) ValueString(v interface{}) string {
	bstr, _ := json.Marshal(v)
	str := string(bstr)
	if strings.HasPrefix(str, "{") {
		return fmt.Sprintf("(%s)", str[1:len(str)-1])
	}
	return str
}

type filter struct {
	filters []string
	args    []interface{}
}

func (ctx *filter) Int64(field string, f store.FilterInt64) {
	switch f.Op {
	case store.OP_NOP:
	case store.OP_EQ:
		ctx.filters = append(ctx.filters, field+" = ?")
		ctx.args = append(ctx.args, f.Value)
	default:
		panic(fmt.Errorf("invalid op %s: %+v", field, f))
	}
}

func (ctx *filter) String(field string, f store.FilterString) {
	switch f.Op {
	case store.OP_NOP:
	case store.OP_EQ:
		ctx.filters = append(ctx.filters, field+" = ?")
		ctx.args = append(ctx.args, f.Value)
	case store.OP_LIKE:
		ctx.filters = append(ctx.filters, field+" like ?")
		ctx.args = append(ctx.args, f.Value)
	default:
		panic(fmt.Errorf("invalid op %s: %+v", field, f))
	}
}

func (ctx *filter) AppendWhere(query string) string {
	if len(ctx.filters) > 0 {
		return query + " WHERE " + strings.Join(ctx.filters, " AND ")
	}
	return query
}
