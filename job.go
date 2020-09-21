package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/hcl"
	"github.com/jmoiron/sqlx"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Job struct {
	Driver        string `hcl:"driver"`
	DSN           string `hcl:"dsn"`
	Query         string `hcl:"query"`
	Channel       string `hcl:"channel"`
	Schedule      string `hcl:"schedule"`
	MessageString string `hcl:"message"`

	messageCompiled *JSVM      `hcl:"-"`
	conn            *sqlx.DB   `hcl:"-"`
	stmnt           *sqlx.Stmt `hcl:"-"`
}

func ParseJobs(jobsdir string) (map[string]*Job, error) {
	files, err := filepath.Glob(filepath.Join(jobsdir, "*.s2s.hcl"))
	if err != nil {
		return nil, err
	}

	result := map[string]*Job{}

	for _, filename := range files {
		var fileJobs struct {
			Jobs map[string]*Job `hcl:"job"`
		}

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		if err := hcl.Decode(&fileJobs, string(data)); err != nil {
			return nil, errors.New("#hcl: " + err.Error())
		}

		for k, job := range fileJobs.Jobs {
			job.messageCompiled, err = NewJSVM(k, fmt.Sprintf("(function(){%s})()", job.MessageString))
			if err != nil {
				return nil, errors.New("#javascript: " + err.Error())
			}

			job.conn, err = sqlx.Connect(job.Driver, job.DSN)
			if err != nil {
				return nil, errors.New("#sql:" + k + ": " + err.Error())
			}

			job.stmnt, err = job.conn.Preparex(job.Query)
			if err != nil {
				return nil, errors.New("#sql:" + k + ": " + err.Error())
			}

			if job.Channel == "" {
				return nil, errors.New("#channel:" + k + ": channel is required")
			}

			if err := (func(job *Job) error {
				_, err := cronhub.AddFunc(job.Schedule, func() {
					if err := job.Exec(); err != nil {
						panic(err)
					}
				})
				return err
			})(job); err != nil {
				return nil, errors.New("#cron:" + k + ":" + err.Error())
			}

			result[k] = job
		}
	}

	return result, nil
}

func (j *Job) Exec() error {
	rows, err := j.stmnt.Queryx()
	if err != nil {
		return err
	}
	defer rows.Close()
	var res []map[string]interface{}
	for rows.Next() {
		o := map[string]interface{}{}
		if err := rows.MapScan(o); err != nil {
			return err
		}
		for k, v := range o {
			if nil == v {
				continue
			}

			switch v.(type) {
			case []uint8:
				v = []byte(v.([]uint8))
			default:
				v, _ = json.Marshal(v)
			}

			var d interface{}
			if nil == json.Unmarshal(v.([]byte), &d) {
				o[k] = d
			} else {
				o[k] = string(v.([]byte))
			}
		}
		res = append(res, o)
	}
	msg := ""
	ctx := map[string]interface{}{
		"$rows": res,
		"log":   log.Println,
		"say": func(in ...interface{}) {
			msg += fmt.Sprint(in...) + "\n"
		},
	}

	if err := j.messageCompiled.Exec(ctx); err != nil {
		return err
	}

	_, err = resty.New().R().SetDoNotParseResponse(true).SetHeader("content-type", "application/json").SetBody(map[string]interface{}{
		"text": msg,
	}).Post(j.Channel)

	return err
}
