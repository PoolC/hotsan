package bot

import (
	"fmt"
	"strings"
	"time"

	"strconv"

	"github.com/nlopes/slack"
)

var prefix string = "party_"

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
		}
	}
}

func scheduleParty(bot *Meu, date *time.Time, keyword string) {
	dur, _ := time.ParseDuration("-10m")
	noti_date := date.Add(dur)
	bot.cr.AddFunc(fmt.Sprintf("0 %d %d %d %d *", noti_date.Minute(), noti_date.Hour(), noti_date.Day(), noti_date.Month()), alarmFuncGenerator(bot, keyword, partyKey(date, keyword)))
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
