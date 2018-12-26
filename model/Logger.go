package model

import (
	"encoding/json"
	"log"
	"net"
)

type LoggerFluent struct {
	Conn net.Conn
}
type MsgFluent struct {
	Tag  string
	Data interface{}
}

var Logger *LoggerFluent

func InitLogger(source string) *LoggerFluent {
	logger, err := net.Dial("tcp", source)
	if err != nil {
		log.Fatal(err)
	}
	Logger = &LoggerFluent{
		Conn: logger,
	}
	return Logger
}

func (l *LoggerFluent) Write(tag string, st interface{}) error {
	msg := MsgFluent{
		Tag:  tag,
		Data: st,
	}
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = l.Conn.Write(raw)
	if err != nil {
		return err
	}
	return nil
}

func (l *LoggerFluent) WriteBytes(bytes []byte) error {
	_, err := l.Conn.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (l *LoggerFluent) Close() {
	l.Conn.Close()
}
