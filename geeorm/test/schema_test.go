package test

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"geeorm/session"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

var (
	testDB      *sql.DB
	TestDial, _ = dialect.GetDialect("mysql")
)

func TestMain(m *testing.M) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/geeorm"
	var err error
	testDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("link session")
	}
	code := m.Run()
	_ = testDB.Close()
	os.Exit(code)
}

func NewSession() *session.Session {
	return session.New(testDB, TestDial)
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}
