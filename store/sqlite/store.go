package sqlite

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/senomas/gohtmx/store"
)

type SqliteStoreCtx struct {
	db       *sqlx.DB
	maxLimit int
}

func init() {
	v := SqliteStoreCtx{}
	store.AddImplementation("sqlite", v.init())
}

func (s *SqliteStoreCtx) init() store.StoreCtx {
	url := os.Getenv("DB_URL")
	if url == "" {
		url = ":memory:"
	}
	db, err := sqlx.Open("sqlite3", url)
	if err != nil {
		panic(fmt.Errorf("error opening database [%s]: %v", url, err))
	}

	qry := `PRAGMA foreign_keys = ON`
	_, err = db.Exec(qry)
	if err != nil {
		panic(fmt.Errorf("error setup conn : %v\n\n%s", err, qry))
	}

	qry = `CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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
	ctx := SqliteStoreCtx{
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

func (s *SqliteStoreCtx) Close() error {
	return s.db.Close()
}

func (s *SqliteStoreCtx) ValidLimit(limit int) bool {
	return limit > 0 && limit <= s.maxLimit
}

func (s *SqliteStoreCtx) ValueString(v interface{}) string {
	bstr, _ := json.Marshal(v)
	str := string(bstr)
	if strings.HasPrefix(str, "{") {
		return fmt.Sprintf("(%s)", str[1:len(str)-1])
	}
	return str
}

type filterCtx struct {
	filters []string
	args    []interface{}
}

func (ctx *filterCtx) Int64(field string, f store.FilterInt64) {
	switch f.Op {
	case store.OP_NOP:
	case store.OP_EQ:
		ctx.filters = append(ctx.filters, field+" = ?")
		ctx.args = append(ctx.args, f.Value)
	default:
		panic(fmt.Errorf("invalid op %s: %+v", field, f))
	}
}

func (ctx *filterCtx) String(field string, f store.FilterString) {
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

func (ctx *filterCtx) AppendWhere(query string) string {
	if len(ctx.filters) > 0 {
		return query + " WHERE " + strings.Join(ctx.filters, " AND ")
	}
	return query
}
