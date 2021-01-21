package database

import (
	"database/sql"
	"github.com/AkameMoe/BilibiliFilter/define"
	"github.com/AkameMoe/BilibiliFilter/utils"
	_ "github.com/mattn/go-sqlite3"
)

var (
	MainDatabase *sql.DB
	err          error
)

func StartDatabaseModule() {
	MainDatabase, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		utils.Logger.Panic().Msg(err.Error())
	} else {
		utils.Logger.Info().Msg("Database Started Successfully")
	}
}

func SaveUser(user *define.User) {
	statement, err := MainDatabase.Prepare("INSERT or replace INTO user (uid,level,name,silence) VALUES (?,?,?,(CASE WHEN (SELECT silence FROM user WHERE uid = ?)=1 THEN 1 ELSE ? END))")
	if err != nil {
		utils.Logger.Error().Msg(err.Error())
		return
	}
	_, err = statement.Exec(user.Uid, user.Level, user.Name, user.Uid, user.Silence)
	if err != nil {
		utils.Logger.Error().Msg(err.Error())
		return
	}
	defer statement.Close()
}
