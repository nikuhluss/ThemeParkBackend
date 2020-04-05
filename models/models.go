package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// NullString wraps sql.NullString for correct JSON marshalling.
type NullString struct {
	sql.NullString
}

func FromSQLNullString(nullString sql.NullString) NullString {
	return NullString{nullString}
}

func (ns *NullString) UnmarshalJSON(data []byte) error {

	var x *string

	err := json.Unmarshal(data, x)
	if err != nil {
		return err
	}

	if x != nil {
		ns.String = *x
		ns.Valid = true
	} else {
		ns.Valid = false
	}

	return nil
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ns.String)
}

// NullTime wraps sql.NullTime for correct JSON marshalling.
type NullTime struct {
	sql.NullTime
}

func FromSQLNullTime(nullTime sql.NullTime) NullTime {
	return NullTime{nullTime}
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {

	var x *time.Time

	err := json.Unmarshal(data, x)
	if err != nil {
		return err
	}

	if x != nil {
		nt.Time = *x
		nt.Valid = true
	} else {
		nt.Valid = false
	}

	return nil
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.Time)
}
