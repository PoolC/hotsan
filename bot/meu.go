package bot

import (
	"fmt"
	"log"
	"regexp"

	"strings"

	. "github.com/PoolC/slack_bot/util"
	"github.com/marcmak/calc/calc"
	"github.com/nlopes/slack"
	"github.com/robfig/cron"
)

type Meu struct {
	*BaseBot
	rc        RedisClient
	timetable map[string]*TimeTable
	cr        *cron.Cron
}

var (
	calc_re        *regexp.Regexp = regexp.MustCompile("^계산하라 메우 (.+)")
	et_register    *regexp.Regexp = regexp.MustCompile("에타 등록 ([^ ]+)")
	et_call        *regexp.Regexp = regexp.MustCompile("다음 (시간|수업)")
	party_register *regexp.Regexp = regexp.MustCompile("(?:(?:(\\d+)월 *)?(\\d+)일 *)?(\\d+)시(?: *(\\d+)분)?(.+)")
	party_list     *regexp.Regexp = regexp.MustCompile("파티 목록( (?:(?:(\\d+)월 *)?(\\d+)일 *)?(\\d+)시(?: *(\\d+)분)( ~ (?:(?:(\\d+)월 *)?(\\d+)일 *)?(\\d+)시(?: *(\\d+)분))?)?")
	party_exit     *regexp.Regexp = regexp.MustCompile("파티 탈퇴 (.+)")
)

func NewMeu(token string, stop *chan struct{}, redisClient RedisClient) *Meu {
	c := cron.New()
	c.Start()
	return &Meu{NewBot(token, stop), redisClient, map[string]*TimeTable{}, c}
}

func (bot *Meu) onConnectedEvent(e *slack.ConnectedEvent) {
	bot.BaseBot.onConnectedEvent(e)
	rescheduleParty(bot)
}

func (bot *Meu) onMessageEvent(e *slack.MessageEvent) {
	message := meuMessageProcess(bot, e)
	switch message.(type) {
	case string:
		bot.sendSimple(e, message.(string))
	}
}

func meuMessageProcess(bot *Meu, e *slack.MessageEvent) interface{} {
	switch {
	case e.Text == "메우, 멱살":
		return "사람은 일을 하고 살아야한다. 메우"
	case e.Text == "메우메우 펫탄탄":
		return `메메메 메메메메 메우메우
메메메 메우메우
펫땅펫땅펫땅펫땅 다이스키`
	}

	text := strings.TrimSpace(e.Text)
	if matched, ok := MatchRE(text, calc_re); ok {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered : %g\n", r)
				bot.replySimple(e, "에러났다 메우")
			}
		}()
		return fmt.Sprintf("%f", calc.Solve(matched[1]))
	} else {
		specialResponses(bot.getBase(), e)
	}
	if bot.IsBeginWithMention(e) {
		if matched, ok := MatchRE(text, et_register); ok {
			register_et(bot, e, matched)
		} else if _, ok := MatchRE(text, et_call); ok {
			next_et(bot, e, matched)
		} else if matched, ok := MatchRE(text, party_register); ok {
			register_party(bot, e, matched)
		} else if matched, ok := MatchRE(text, party_list); ok {
			list_party(bot, e, matched)
		} else if matched, ok := MatchRE(text, party_exit); ok {
			exit_party(bot, e, matched)
		}
	}

	return nil
}
