package request

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"time"
)

type Time time.Time

func (rt Time) MarshalJSON() ([]byte, error) {
	return []byte("\"" + rt.String() + "\""), nil
}

func (rt *Time) UnmarshalJSON(b []byte) (err error) {
	return json.Unmarshal(b, rt)
}

func (rt *Time) String() string {
	t := time.Time(*rt)
	return t.Format(time.RFC822)
}

func (rt *Time) Scan(src interface{}) error {
	if t, ok := src.(Time); ok {
		rt = &t
	}
	return nil
}

func (rt Time) Value() (driver.Value, error) {
	return time.Time(rt), nil
}

type CaughtRequest struct {
	Id            int64
	Time          Time   `db:"created_at" json:"created_at"`
	Method        string `db:"method" json:"method"`
	ContentLength int64  `db:"content_length" json:"content_length"`
	RemoteAddr    string `db:"remote_addr" json:"remote_addr"`
	Url           string `db:"url" json:"url"`
	Headers       string `db:"headers" json:"headers"`
	Body          string `db:"body" json:"body"`
}

func (v *CaughtRequest) ParsedHeaders() map[string][]string {
	var headers map[string][]string
	err := json.Unmarshal([]byte(v.Headers), &headers)
	if err != nil {
		log.Print("Failed to unmarshal headers json", err)
	}
	return headers
}
