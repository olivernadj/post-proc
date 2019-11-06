package api

import (
	"github.com/olivernadj/post-proc/internal/sqlconn"
	"github.com/olivernadj/post-proc/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

var  summaryVec  *prometheus.SummaryVec

func init() {
	summaryVec = BuildSummaryVec("processing_time_milliseconds", "Latency Percentiles in Milliseconds", "statemachine")
}

func HandleProcess(interval int) {
	withMonitoring("check_process", func() {
		handleProcess(interval)
	}, summaryVec)
}
func HandleDelete(interval int) {
	withMonitoring("check_delete", func() {
		handleDelete(interval)
	}, summaryVec)
}

type handleEvent func()

func withMonitoring(name string, handler handleEvent, summary *prometheus.SummaryVec) {
	start := time.Now()
	handler()
	duration := time.Since(start)
	// Store duration of request
	summary.WithLabelValues(name, "OK").Observe(duration.Seconds() * 1000)
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
		withMonitoring("process_action", func() {
			var a model.Action
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
		}, summaryVec)
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
		withMonitoring("delete_action", func() {
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
		}, summaryVec)
	}
	defer results.Close()
}

