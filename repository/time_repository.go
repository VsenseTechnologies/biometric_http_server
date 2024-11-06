package repository

import (
	"database/sql"
	"encoding/json"
	"io"
	"sync"

	"vsensetech.in/go_fingerprint_server/models"
)

type timeRepository struct {
	db *sql.DB
	mut *sync.Mutex
}

func NewTimeRepository(db *sql.DB , mut *sync.Mutex) *timeRepository{
	return &timeRepository{
		db,
		mut,
	}
}

func(tr *timeRepository) SetTime(reader *io.ReadCloser) error {
	tr.mut.Lock()
	defer tr.mut.Unlock()
	var newTime models.TimeModel

	if err := json.NewDecoder(*reader).Decode(&newTime); err != nil {
		return err
	}

	if _ , err := tr.db.Exec("UPDATE times SET morning_start=$1 , morning_end=$2 , afternoon_start=$3 , afternoon_end=$4 , evening_start=$5 , evening_end=$6 WHERE user_id=$7",newTime.MorningStart , newTime.MorningEnd , newTime.AfternoonStart , newTime.AfternoonEnd , newTime.EveningStart , newTime.EveningEnd , newTime.UserID); err != nil {
		return err
	}
	return nil
}