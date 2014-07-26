package mixpanel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	host       = "http://api.mixpanel.com"
	trackPath  = "track"
	engagePath = "engage"
)

type mixpanel struct {
	token string
}

/*
	example track event data:
  {
    "event": "Signed Up",
    "properties": {
      "distinct_id": "13793",
      "token": "e3bc4100330c35722740fb8c6f5abddc",
      "Referred By": "Friend"
    }
  }
	example engage event data:
	{
		"$token": "36ada5b10da39a1347559321baf13063",
		"$distinct_id": "13793",
		"$ip": "123.123.123.123",
		"$set": {
			"Address": "1313 Mockingbird Lane"
		}
	}
*/

type eventData struct {
	Event string                 `json:"event"`
	Props map[string]interface{} `json:"properties"`
}

func New(token string) *mixpanel {
	return &mixpanel{
		token: token,
	}
}

func (mp *mixpanel) Track(uid int64, e string, p map[string]interface{}) bool {
	data := &eventData{
		Event: e,
		Props: map[string]interface{}{
			"time":  time.Now().Unix(),
			"token": mp.token,
		},
	}
	if uid != 0 {
		data.Props["distinct_id"] = strconv.Itoa(int(uid))
	}
	for k, v := range p {
		data.Props[k] = v
	}

	marshaledData, err := json.Marshal(data)
	if err != nil {
		return false
	}

	url := fmt.Sprintf("%s/%s/?data=%s", host, trackPath, base64.StdEncoding.EncodeToString(marshaledData))

	_, err = http.Get(url)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	return true
}

func (mp *mixpanel) Engage(uid int64, e string, p map[string]interface{}) bool {
	return false
}
