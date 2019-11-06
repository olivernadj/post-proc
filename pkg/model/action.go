package model

import (
	"github.com/olivernadj/post-proc/internal/sqlconn"
	"time"
)

type Action struct {
	Id int
	Action string `json:"action,omitempty"`
	// Status of the item
	State string `json:"state,omitempty"`
	SourceType string
	Processed time.Time
}

func (a Action) Insert() error {
	db, err := sqlconn.GetConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	insert, err := db.Query("INSERT INTO `action` (`action`, `state`, `source_type`) VALUES (?, ?, ?)", a.Action, a.State, a.SourceType)
	if err != nil {
		return err
	}
	defer insert.Close()
	return nil
}

func (a Action) NewWithState () *ActionWithState {
	 return NewActionWithState(a.Id, a.Action, a.SourceType, a.State, a.Processed)
}