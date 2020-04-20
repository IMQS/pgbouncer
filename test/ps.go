package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Test Prepared Statements from lib/pq

func check(e error) {
	if e != nil {
		msg := fmt.Sprintf("Error: %v", e)
		panic(msg)
	}
}

func onceOff(db *sql.DB) {
	key := ""
	sanity := 0
	db.QueryRow("SELECT sessionkey,$1 FROM authsession LIMIT 1", 5).Scan(&key, &sanity)
	fmt.Printf("%v, %v\n", key, sanity)
}

func main() {
	db, err := sql.Open("postgres", "host=localhost port=6432 user=auth password=auth dbname=auth sslmode=disable")
	check(err)

	onceOff(db)
	os.Exit(0)

	poll := func(threadIdx int) {
		ticker := 0
		for {
			fmt.Printf("%v:%v\n", threadIdx, ticker)
			key := ""
			userid := int64(0)
			sanity := 0
			db.QueryRow("SELECT sessionkey,userid,$1 FROM authsession ORDER BY sessionkey LIMIT 1", 5).Scan(&key, &userid, &sanity)
			if key != "6TwgcIJKbp7nM7dGm2XdSyk8izTttO" {
				panic(fmt.Errorf("invalid sessionkey '%v'", key))
			}

			time.Sleep(100 * time.Microsecond)

			email := ""
			username := ""
			modifiedBy := int64(0)
			db.QueryRow("SELECT email,username,modifiedby,$1 FROM authuserstore ORDER BY userid LIMIT 1", 6).Scan(&email, &username, &modifiedBy)
			if email != "dev@dev.com" {
				panic(fmt.Errorf("invalid email '%v'", email))
			}

			time.Sleep(time.Second)
			ticker++
		}
	}

	for i := 0; i < 3; i++ {
		go poll(i)
	}

	for {
		time.Sleep(time.Second)
	}
}
