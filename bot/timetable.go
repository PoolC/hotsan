package bot

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"sort"

	"encoding/json"

	"log"

	"github.com/nlopes/slack"
	"golang.org/x/net/html"
)

type TimeTableResponse struct {
	XMLName  xml.Name `xml:"response"`
	Year     int      `xml:"year,attr"`
	Semester int      `xml:"semester,attr"`
	Subjects []struct {
		Name struct {
			Value string `xml:"value,attr"`
		} `xml:"name"`
		Professor struct {
			Value string `xml:"value,attr"`
		} `xml:"professor"`
		Times []struct {
			Place     string `xml:"place,attr"`
			EndTime   int    `xml:"endtime,attr"`
			StartTime int    `xml:"starttime,attr"`
			Day       int    `xml:"day,attr"`
		} `xml:"time>data"`
	} `xml:"subject"`
}

type Subject struct {
	Name      string
	Professor string
}
type TimeTableEvent struct {
	Id        int
	Place     string
	StartTime int
	EndTime   int
}

type ByStartTime []TimeTableEvent

func (a ByStartTime) Len() int           { return len(a) }
func (a ByStartTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStartTime) Less(i, j int) bool { return a[i].StartTime < a[j].StartTime }

type TimeTableDay struct {
	Events []TimeTableEvent
}

type TimeTable struct {
	Year     int
	Semester int
	Days     []TimeTableDay
	Subjects []Subject
}

func (resp *TimeTableResponse) toTimeTable() *TimeTable {
	ret := &TimeTable{}

	ret.Year = resp.Year
	ret.Semester = resp.Semester

	ret.Days = make([]TimeTableDay, 5)
	ret.Subjects = make([]Subject, len(resp.Subjects))
	for i, subject := range resp.Subjects {
		ret.Subjects[i] = Subject{
			subject.Name.Value,
			subject.Professor.Value,
		}
		for _, day := range subject.Times {
			ret.Days[day.Day].Events = append(ret.Days[day.Day].Events, TimeTableEvent{
				i,
				day.Place,
				day.StartTime,
				day.EndTime,
			})
		}
	}

	for _, day := range ret.Days {
		sort.Sort(ByStartTime(day.Events))
	}

	return ret
}

// return value is dummy...
func getEveryTimeTable(bot *Meu, userid string, username string, et_nick string) error {
	ret := &slack.OutgoingMessage{}
	quit := func(msg string, args ...interface{}) error {
		ret.Text = fmt.Sprintf("<@%s> "+msg, userid, args)
		log.Print(ret.Text)
		bot.SendMessage(ret)
		return nil
	}

	url_str := fmt.Sprintf("http://everytime.kr/@%s", et_nick)
	resp, err := http.Get(url_str)
	if err != nil {
		return quit("에러: 시간표 목록을 가져오는데 실패했습니다.")
	}

	latest_info, err := parseTimeTableList(resp.Body)
	if err != nil {
		return quit("에러: 시간표 목록을 파싱하는데 실패했습니다. %q", err)
	}

	now := time.Now()
	if now.Year() != latest_info.Year || (int(now.Month())-1)/6+1 != latest_info.Semester {
		return quit("에러: 현재 학기의 시간표가 없습니다.")
	}

	url_str = "http://everytime.kr/ajax/timetable/wizard/tableload"
	resp, err = http.PostForm(url_str, url.Values{
		"id": {latest_info.Id},
	})
	if err != nil {
		return quit("에러: 시간표 정보를 가져오는데 실패했습니다.")
	}

	schedule_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return quit("에러: 시간표 읽기 실패")
	}

	timetableResp, err := parseTimeTable(schedule_body)
	if err != nil {
		return quit("에러: 시간표 파싱 에러 - %q", err)
	}

	timetable := timetableResp.toTimeTable()
	bot.timetable[userid] = timetable
	serialized, err := json.Marshal(timetable)
	if err != nil {
		return quit("에러: 시간표 저장 실패")
	}

	bot.rc.Set(TimeTableKeyName(userid), serialized, 0)

	return nil
}

func TimeTableKeyName(userid string) string {
	return fmt.Sprintf("et_tt_%s", userid)
}

type TimeTableInfo struct {
	Year     int
	Semester int
	Id       string
}

func parseTimeTableList(reader io.Reader) (*TimeTableInfo, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}
	ret := &TimeTableInfo{}

	var parser func(*html.Node) bool
	parser = func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			switch {
			case n.Data == "input":
				var input_id, input_val string
				for _, attr := range n.Attr {
					if attr.Key == "id" {
						input_id = attr.Val
					} else if attr.Key == "value" {
						input_val = attr.Val
					}
				}
				switch input_id {
				case "year":
					ret.Year, _ = strconv.Atoi(input_val)
				case "semester":
					ret.Semester, _ = strconv.Atoi(input_val)
				case "tableId":
					ret.Id = input_val
				}
				break
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if !parser(c) {
				return false
			}
		}

		return true
	}
	parser(doc)

	return ret, nil
}

func parseTimeTable(body []byte) (*TimeTableResponse, error) {
	schedule := &TimeTableResponse{}
	err := xml.Unmarshal(body, &schedule)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (nextEvent *TimeTableEvent) toSlackAttachment(subject *Subject) slack.Attachment {
	now := time.Now()
	begin := time.Date(now.Year(), now.Month(), now.Day(), nextEvent.StartTime/12, (nextEvent.StartTime%12)*5, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), nextEvent.EndTime/12, (nextEvent.EndTime%12)*5, 0, 0, now.Location())

	return slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "수업명",
				Value: subject.Name,
			},
			slack.AttachmentField{
				Title: "교수",
				Value: subject.Professor,
			},
			slack.AttachmentField{
				Title: "장소",
				Value: nextEvent.Place,
			},
			slack.AttachmentField{
				Title: "수업 시각",
				Value: fmt.Sprintf("%s ~ %s", begin.Format(time.Kitchen), end.Format(time.Kitchen)),
			},
		},
	}
}

func register_et(bot *Meu, e *slack.MessageEvent, matched []string) {
	// 에브리타임 기록
	et_nick := matched[1]
	bot.rc.Set(fmt.Sprintf("et_nick_%s", e.User), et_nick, 0)
	bot.replySimple(e, "기억했다 메우. 시간표를 가져오겠다 메우.\n시간이 좀 걸릴거다 메우.")

	user, _ := bot.GetUserInfo(e.User)
	go getEveryTimeTable(bot, e.User, user.Name, et_nick)
}

func next_et(bot *Meu, e *slack.MessageEvent, matched []string) {
	// 에브리타임 다음 시간
	log.Printf("%q", bot.timetable)
	timetable, exists := bot.timetable[e.User]
	if !exists {
		log.Print("Get from redis")
		result := bot.rc.Get(TimeTableKeyName(e.User))
		if result == nil {
			bot.replySimple(e, "시간표 정보가 없다 메우. 등록부터 해달라 메우.")
			return
		}

		bytes, _ := result.Bytes()
		timetable = &TimeTable{}
		if json.Unmarshal(bytes, timetable) != nil {
			bot.replySimple(e, "저장된 시간표가 이상하다 메우. 새로 등록해달라 메우.")
			return
		}

		bot.timetable[e.User] = timetable
	}

	now := time.Now()
	weekDay := int(now.Weekday() - 1)
	curHour := now.Hour()*12 + now.Minute()/5
	var nextEvent *TimeTableEvent
	if weekDay < 0 || weekDay >= 5 {
		nextEvent = nil
	} else {
		log.Printf("%q", timetable)
		todayEvents := timetable.Days[weekDay].Events
		for _, event := range todayEvents {
			if event.StartTime > curHour {
				nextEvent = &event
				break
			}
		}
	}
	if nextEvent == nil {
		bot.replySimple(e, "오늘은 더 이상 수업이 없다 메우.")
	} else {
		subject := &timetable.Subjects[nextEvent.Id]

		bot.PostMessage(e.Channel, fmt.Sprintf("다음 수업은 \"%s\"다 메우.", subject.Name), slack.PostMessageParameters{
			AsUser:    false,
			IconEmoji: ":meu:",
			Username:  "시간표 알려주는 메우",
			Attachments: []slack.Attachment{
				nextEvent.toSlackAttachment(subject),
			},
		})
	}
}
