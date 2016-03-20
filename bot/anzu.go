package bot

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"fmt"

	. "github.com/PoolC/slack_bot/util"
	slack "github.com/nlopes/slack"
)

var (
	remember_re   *regexp.Regexp = regexp.MustCompile("^안즈쨩? 기억해? ([^/]+)/(.+)")
	tell_re       *regexp.Regexp = regexp.MustCompile("^안즈쨩? 알려줘 (.+)")
	kawaii_re     *regexp.Regexp = regexp.MustCompile("^안즈쨩? 카와이")
	give_candy_re *regexp.Regexp = regexp.MustCompile("^안즈쨩? 사탕줄게")
	dice          *regexp.Regexp = regexp.MustCompile("^안즈쨩? 주사위 ([0-9]+)")
)

type Anzu struct {
	*BaseBot
	rc RedisClient
}

func NewAnzu(token string, stop *chan struct{}, redisClient RedisClient) *Anzu {
	return &Anzu{NewBot(token, stop), redisClient}
}

func anzuMessageProcess(bot *Anzu, e *slack.MessageEvent) interface{} {
	force_accept := false
	switch {
	case e.Text == "사람은 일을 하고 살아야한다. 메우":
		return "이거 놔라 이 퇴근도 못하는 놈이"
	case e.Text == "안즈쨩 뭐해?":
		return "숨셔"
	default:
		if AcceptRE(e.Text, give_candy_re) {
			last := bot.rc.Get(fmt.Sprintf("%s_lastfail", e.User)).Val()
			if last == "" {
				return ""
			}
			force_accept = true
			e.Text = last
		}

		if matched, ok := MatchRE(e.Text, remember_re); ok {
			key, val := strings.TrimSpace(matched[1]), strings.TrimSpace(matched[2])
			var ret string
			switch {
			case key == "" || val == "":
				ret = "에...?"
			case AcceptRE(val, tell_re):
				ret = "에... 귀찮아..."
			case force_accept:
				ret = "응응 기억했어"
				fallthrough
			case rand.Float32() < 0.4:
				bot.rc.Set(key, val, 0)
				if len(ret) == 0 {
					ret = "에... 귀찮지만 기억했어"
				}
			default:
				ret = "귀찮아..."
				bot.rc.Set(fmt.Sprintf("%s_lastfail", e.User), e.Text, time.Duration(300*time.Second))
			}

			return ret
		} else if matched, ok := MatchRE(e.Text, tell_re); ok {
			key := strings.TrimSpace(matched[1])
			val := bot.rc.Get(key).Val()
			var ret string
			switch {
			case val == "":
				ret = "그런거 몰라"
			case force_accept:
				ret = fmt.Sprintf("%s 물어봤지? %s 야", key, val)
			case rand.Float32() < 0.4:
				ret = val
			default:
				bot.rc.Set(fmt.Sprintf("%s_lastfail", e.User), e.Text, time.Duration(300*time.Second))
				ret = "Zzz..."
			}
			return ret
		} else if _, ok := MatchRE(e.Text, kawaii_re); ok {
			return "뭐... 뭐라는거야"
		} else if matched, ok := MatchRE(e.Text, dice); ok {
			num := matched[1]
			return fmt.Sprintf("%d", num * rand.Float32())
		}
	}
	return nil
}

func (bot *Anzu) onMessageEvent(e *slack.MessageEvent) {
	message := anzuMessageProcess(bot, e)
	switch message.(type) {
	case string:
		bot.sendSimple(e, message.(string))
	}
}
