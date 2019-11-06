package model

import (
	"github.com/looplab/fsm"
	"github.com/olivernadj/post-proc/internal/sqlconn"
	"log"
	"time"
)

type ActionWithState struct {
	id int
	Action string
	SourceType string
	Processed time.Time
	FSM *fsm.FSM
}

func NewActionWithState (id int, a string, st string, s string, p time.Time) *ActionWithState {
	aws := &ActionWithState{
		id:			id,
		Action:     a,
		SourceType: st,
		Processed:  p,
	}
	aws.FSM = fsm.NewFSM(
		s,
		fsm.Events{
			{Name: "process", Src: []string{"new"}, Dst: "processed"},
			{Name: "delete", Src: []string{"processed"}, Dst: "deleted"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { aws.enterState(e) },
		},
	)
	return aws
}

func (aws *ActionWithState) enterState(e *fsm.Event) {
	log.Printf("Enter state for id:%d is %s\n", aws.id, e.Dst)
	if e.Dst == "processed" {
		aws.Processed = time.Now()
	}
}

// saves state and processed date only as business requirements
func (aws *ActionWithState) Save() error {
	db, err := sqlconn.GetConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	update, err := db.Query("" +
		"UPDATE `action` " +
		"	SET `state` = ?, " +
		"		`processed` = ? " +
		"	WHERE `id` = ?", aws.FSM.Current(), aws.Processed, aws.id)
	if err != nil {
		return err
	}
	defer update.Close()
	return nil
}