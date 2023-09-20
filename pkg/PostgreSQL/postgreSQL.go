package PostgreSQL

import "database/sql"
import _ "github.com/lib/pq"

type DataBase struct {
	DB *sql.DB
}

//func NewDataBase() (*DataBase, error) {
//	db := &DataBase{}
//	err := db.Open()
//	if err != nil {
//
//	}
//}

func (d *DataBase) Open() error {
	var err error = nil

	//d.DB, err = sql.Open(cfg.Storage.DbDriver, cfg.Storage.DbDriver+"://"+cfg.Storage.Username+":"+
	//	""+cfg.Storage.Password+"@"+cfg.Storage.Host+":"+cfg.Storage.Port+"/"+cfg.Storage.Database+""+
	//	"?sslmode="+cfg.Storage.SSLMode)
	d.DB, err = sql.Open("postgres", "postgres://dmilyano:qwerty@localhost:5432/messenger?sslmode=disable")

	if err != nil {
		return err
	}

	if err = d.DB.Ping(); err != nil {
		return err
	}

	return nil
}
