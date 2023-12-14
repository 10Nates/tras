package db

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-pg/pg/v10"
)

type Connection struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	DB       *pg.DB
}

func (c *Connection) Connect() error {
	// open database
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		User:     c.User,
		Password: c.Password,
		Database: c.DBName,
	})

	err := createSchema(db)
	if err != nil {
		return err
	}

	c.DB = db

	// check db
	err = c.DB.Ping(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) CloseOnInterrupt() {
	s := make(chan os.Signal, 3)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	go func() {
		<-s
		err := c.DB.Close()
		if err != nil {
			os.Exit(99) // incorrect exit
		}
		os.Exit(0)
	}()
}
