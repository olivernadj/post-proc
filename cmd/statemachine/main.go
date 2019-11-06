package main

import (
	"github.com/olivernadj/post-proc/internal/sqlconn"
	"github.com/olivernadj/post-proc/pkg/model"
	"log"
	"os"
	"os/signal"
	"time"
)

//SELECT *
//FROM `action`
//WHERE created < NOW() - INTERVAL 1 MINUTE
//LIMIT 50


func main() {
	//check for action timeouts in every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	// trap Ctrl+C and call cancel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		log.Printf("signal stop")
		signal.Stop(c)
	}()

	for {
		select {
		case <- ticker.C:
			// do stuff
			handleProcess(1)
			handleDelete(1)
		case <-c:
			log.Printf("gracefully stop statemachinge")
			ticker.Stop()
			return
		}
	}
}

func handleProcess(interval int) {
	db, err := sqlconn.GetConnection()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	results, err := db.Query("" +
		"SELECT `id`, `state`, `processed` " +
		"	FROM `action` " +
		"	WHERE `state` = 'new' AND created < NOW() - INTERVAL ? MINUTE", interval)
	if err != nil {
		log.Println(err)
		return
	}
	for results.Next() {
		var a model.Action
		//var p string
		err = results.Scan(&a.Id, &a.State, &a.Processed)
		if err != nil {
			log.Println(err)
			return
		}
		aws := a.NewWithState()
		err = aws.FSM.Event("process")
		if err != nil {
			log.Println(err)
			return
		}
		err = aws.Save()
		if err != nil {
			log.Println(err)
			return
		}
	}
	defer results.Close()
}

func handleDelete(interval int) {
	db, err := sqlconn.GetConnection()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	results, err := db.Query("" +
		"SELECT `id`, `state`, `processed` " +
		"	FROM `action` " +
		"	WHERE `state` = 'processed' AND processed < NOW() - INTERVAL ? MINUTE", interval)
	if err != nil {
		log.Println(err)
		return
	}
	for results.Next() {
		var a model.Action
		//var p string
		err = results.Scan(&a.Id, &a.State, &a.Processed)
		if err != nil {
			log.Println(err)
			return
		}
		aws := a.NewWithState()
		err = aws.FSM.Event("delete")
		if err != nil {
			log.Println(err)
			return
		}
		err = aws.Save()
		if err != nil {
			log.Println(err)
			return
		}
	}
	defer results.Close()
}