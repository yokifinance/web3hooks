package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"

	"yoki.finance/common/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type MapStringInterface map[string]interface{}

func (m MapStringInterface) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MapStringInterface) Scan(v interface{}) error {
	vBytes, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("expected %v to be []byte, not %T", v, v)
	}

	var vValue interface{}
	if err := json.Unmarshal(vBytes, &vValue); err != nil {
		return err
	}

	if *m, ok = vValue.(map[string]interface{}); !ok {
		return fmt.Errorf("expected %v to be map[string]interface{}, not %T", vValue, vValue)
	}

	return nil
}

var Conn *sql.DB
var ORM *bun.DB

func init() {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DBhost, config.DBport, config.DBuser, config.DBpass, config.DBname)

	Conn, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = Conn.Ping()
	if err != nil {
		panic(err)
	}

	ORM = bun.NewDB(Conn, pgdialect.New())

	// enables ORM queries logging – can be customized by environment variable BUNDEBUG
	ORM.AddQueryHook(bundebug.NewQueryHook(

		bundebug.FromEnv("BUNDEBUG"),
	))

	fmt.Printf("DB %s successfully connected (%s:%s)!\n", config.DBname, config.DBhost, config.DBport)
}
