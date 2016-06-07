package bot

import (
	"fmt"
	"strings"
	"time"

	"strconv"

	"github.com/nlopes/slack"
)

var (
	prefix        string = "party_"
	sprefix       string = "s_" + prefix
	partyIndexKey string = sprefix + "keys"
	meuName       string = "파티 모집원 메우"
)

func rescheduleParty(bot *Meu) {
	items, err := bot.rc.Keys(prefix + "*")
	if err != nil {
		return
	}
	now := time.Now()
	for _, key := range items {
		saved_time, keyword := parseKey(key)
		if now.After(saved_time) {
			members, _ := bot.rc.SetList(key)
			for _, member := range members {
				bot.rc.SetRemove(sprefix+member, key)
			}
			bot.rc.Erase(key)
		} else {
			scheduleParty(bot, &saved_time, keyword)
			registerToIndex(bot, &saved_time, key)
		}
	}

	bot.rc.SortedSetRemoveRange(partyIndexKey, 0, now.Unix())
}

func parseKey(key string) (time.Time, string) {
	parts := strings.Split(key, "_")
	sec, _ := strconv.Atoi(parts[1])
	return time.Unix(int64(sec), 0), parts[2]
}

func registerToIndex(bot *Meu, date *time.Time, key string) {
	if _, err := bot.rc.SortedSetRank(partyIndexKey, key).Result(); err != nil {
		bot.rc.SortedSetAdd(partyIndexKey, int(date.Unix()), key)
	}
}

func scheduleParty(bot *Meu, date *time.Time, keyword string) {
	bot.cr.AddFunc(fmt.Sprintf("0 %d %d %d %d *", date.Minute(), date.Hour(), date.Day(), date.Month()), alarmFuncGenerator(bot, keyword, partyKey(date, keyword)))
}

func partyKey(date *time.Time, keyword string) string {
	return prefix + strconv.FormatInt(date.Unix(), 10) + "_" + keyword
}

func alarmFuncGenerator(bot *Meu, keyword string, key string) func() {
	return func() {
		list, err := bot.rc.SetList(key)
		bot.rc.Erase(key)
		if err == nil {
			members := make([]string, len(list))
			for i, item := range list {
				members[i] = fmt.Sprintf("<@%s>", item)
				bot.rc.SetRemove(sprefix+item, key)
			}
			bot.PostMessage("#random", fmt.Sprintf("'%s' 파티 10분 전이다 메우. %s", keyword, strings.Join(members, " ")), slack.PostMessageParameters{
				AsUser:    false,
				Username:  meuName,
				IconEmoji: ":meu:",
			})
		}
	}
}

func correctDate(matched []string) *time.Time {
	now := time.Now()
	month, err := strconv.Atoi(matched[0])
	var not_set struct {
		month bool
		day   bool
	}
	if err != nil {
		month = int(now.Month())
		not_set.month = true
	}
	day, err := strconv.Atoi(matched[1])
	if err != nil {
		day = now.Day()
		not_set.day = true
	}
	hour, err := strconv.Atoi(matched[2])
	if err != nil {
		return nil
	}
	min, err := strconv.Atoi(matched[3])
	if err != nil {
		min = 0
	}

	date := time.Date(now.Year(), time.Month(month), day, hour, min, 0, 0, now.Location())
	if date.Before(now) {
		corrected := false
		// first try after 12 hour
		if not_set.day {
			if date.Hour() < 12 {
				date = date.Add(time.Hour * 12)
				if corrected = !date.Before(now); !corrected {
					// reset
					date = date.Add(time.Hour * -12)
				}
			}
			if !corrected {
				date = date.AddDate(0, 0, 1)
			}
		} else if not_set.month {
			date = date.AddDate(0, 1, 0)
		} else {
			date = date.AddDate(1, 0, 0)
		}
		if date.Before(now) {
			return nil
		}
	}

	return &date
}

func event_key_to_slack_attach(key string) slack.Attachment {
	date, keyword := parseKey(key)
	return event_to_slack_attach(key, keyword, &date)
}

func event_to_slack_attach(key string, keyword string, date *time.Time) slack.Attachment {
	return slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "일시",
				Value: date.String(),
			},
			slack.AttachmentField{
				Title: "이름",
				Value: keyword,
			},
			slack.AttachmentField{
				Title: "파티 ID",
				Value: key,
			},
		},
	}
}

func register_party(bot *Meu, e *slack.MessageEvent, matched []string) {
	keyword := strings.TrimSpace(matched[5])

	date := correctDate(matched[1:])
	if date == nil {
		bot.replySimple(e, "과거에 대해서 파티를 모집할 수 없다 메우")
		return
	}

	key := partyKey(date, keyword)
	inserted := bot.rc.SetAdd(key, e.User)
	registerToIndex(bot, date, key)
	responseData := slack.PostMessageParameters{
		AsUser:    false,
		IconEmoji: ":meu:",
		Username:  meuName,
		Attachments: []slack.Attachment{
			event_to_slack_attach(key, keyword, date),
		},
	}
	if inserted.Val() == 1 {
		bot.PostMessage(e.Channel, fmt.Sprintf("<%s> 파티 대기에 들어갔다 메우", e.User), responseData)
		cardinal := bot.rc.SetCard(key)
		bot.rc.SetAdd(sprefix+e.User, key)
		if cardinal.Val() == 1 {
			scheduleParty(bot, date, keyword)
		}
	} else {
		bot.PostMessage(e.Channel, fmt.Sprintf("<%s> 이미 들어가있는 파티다 메우.", e.User), responseData)
	}
}

func list_party(bot *Meu, e *slack.MessageEvent, matched []string) {
	b_t := correctDate(matched[1:5])
	e_t := correctDate(matched[5:9])
	var (
		end   time.Time
		begin time.Time
	)
	if b_t == nil {
		keys, err := bot.rc.SetList(sprefix + e.User)
		if err != nil || len(keys) == 0 {
			bot.replySimple(e, "대기중인 파티가 없다 메우.")
		} else {
			attachments := make([]slack.Attachment, len(keys))
			for i, key := range keys {
				attachments[i] = event_key_to_slack_attach(key)
			}
			bot.PostMessage(e.Channel, fmt.Sprintf("<%s> 지금 대기중인 파티는 다음과 같다 메우.", e.User), slack.PostMessageParameters{
				AsUser:      false,
				IconEmoji:   ":meu:",
				Username:    meuName,
				Attachments: attachments,
			})
		}
		return
	} else if e_t == nil {
		d, _ := time.ParseDuration("1h")
		end = begin.Add(d)
		d, _ = time.ParseDuration("-1h")
		begin = begin.Add(d)
	} else {
		end = *e_t
	}
	begin = *b_t

	keys, _ := bot.rc.SortedSetRange(partyIndexKey, begin.Unix(), end.Unix())
	attachments := make([]slack.AttachmentField, len(keys))
	for i, key := range keys {
		t, k := parseKey(key)
		attachments[i] = slack.AttachmentField{
			Title: k,
			Value: t.String(),
		}
	}

	bot.PostMessage(e.Channel,
		fmt.Sprintf("%s ~ %s 사이에 있는 파티 목록은 다음과 같다 메우.", begin.String(), end.String()),
		slack.PostMessageParameters{
			AsUser:    false,
			IconEmoji: ":meu:",
			Username:  meuName,
			Attachments: []slack.Attachment{
				slack.Attachment{
					Fields: attachments,
				},
			},
		})
}

func exit_party(bot *Meu, e *slack.MessageEvent, matched []string) {
	key := strings.TrimSpace(matched[1])
	if bot.rc.SetRemove(key, e.User).Val() == 1 {
		bot.replySimple(e, "성공적으로 파티 대기에서 빠졌다 메우")
		bot.rc.SetRemove(sprefix+e.User, key)
		if bot.rc.SetCard(key).Val() == 0 {
			bot.rc.Erase(key)
			bot.rc.SortedSetRemove(partyIndexKey, key)
		}
	} else {
		bot.replySimple(e, "잘못된 파티 이름이거나 대기중이 아닌 파티이다 메우")
	}
}
