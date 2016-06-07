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
	partyIndexKey string = prefix + "keys"
)

func rescheduleParty(bot *Meu) {
	items, err := bot.rc.Keys(prefix + "*")
	if err != nil {
		return
	}
	now := time.Now()
	for _, key := range items {
		parts := strings.Split(key, "_")
		sec, _ := strconv.Atoi(parts[1])
		saved_time := time.Unix(int64(sec), 0)
		if now.After(saved_time) {
			bot.rc.Erase(key)
		} else {
			scheduleParty(bot, &saved_time, parts[2])
			registerToIndex(bot, &saved_time, key)
		}
	}

	bot.rc.SortedSetRemoveRange(partyIndexKey, 0, now.Unix())
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
			}
			bot.PostMessage("#random", fmt.Sprintf("'%s' 파티 10분 전이다 메우. %s", keyword, strings.Join(members, " ")), slack.PostMessageParameters{
				AsUser:    false,
				Username:  "파티 안내원 메우",
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

func register_party(bot *Meu, e *slack.MessageEvent, matched []string) {
	keyword := matched[5]

	date := correctDate(matched[1:])
	if date == nil {
		bot.replySimple(e, "과거에 대해서 파티를 모집할 수 없다 메우")
		return
	}

	key := partyKey(date, keyword)
	inserted := bot.rc.SetAdd(key, e.User)
	registerToIndex(bot, date, key)
	if inserted.Val() == 1 {
		bot.replySimple(e, fmt.Sprintf("파티 대기에 들어갔다 메우. - %s %s", date.String(), keyword))
		cardinal := bot.rc.SetCard(key)
		if cardinal.Val() == 1 {
			scheduleParty(bot, date, keyword)
		}
	} else {
		bot.replySimple(e, fmt.Sprintf("이미 들어가있는 파티다 메우. - %s %s", date.String(), keyword))
	}
}
