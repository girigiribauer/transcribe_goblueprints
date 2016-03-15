package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/bitly/go-nsq"
	mgo "gopkg.in/mgo.v2"
)

var fatalErr error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

func main() {
	// main の先頭に記述
	// エラーが発生しても必ず処理したいものを内部に書く
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	log.Println("データベースに接続します...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		fatal(err)
		return
	}
	defer func() {
		log.Println("データベース接続を閉じます...")
		db.Close()
	}()
	pollData := db.DB("ballots").C("polls")

	var countsLock sync.Mutex // 複数の goroutine が1つのマップに同時に読み書きしないため
	var counts map[string]int

	log.Println("NSQに接続します...")
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fatal(err)
		return
	}

	// NSQ からメッセージを受け取った際の処理
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))

	// NSQ に接続
	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		fatal(err)
		return
	}
}
