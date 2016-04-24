package bot

import "github.com/nlopes/slack"
import "strings"
import "math/rand"

func postResponse(bot *BaseBot, channel string, emoji string, name string, response string) {
	bot.PostMessage(channel, response, slack.PostMessageParameters{
		AsUser:    false,
		IconEmoji: emoji,
		Username:  name,
		Parse:     "full",
	})
}

func randomResponse(bot *BaseBot, channel string, emoji string, name string, responses []string) {
	response := responses[rand.Intn(len(responses))]
	postResponse(bot, channel, emoji, name, response)
}

var (
	anduinresp []string = []string{
		"안녕하세요!",
		"감사합니다!",
		"이게... 아닌데...",
		"빛이 당신을 태울 것입니다!",
		"정말 잘하셨어요.",
		"죄송합니다.",
	}
	yayoiresp []string = []string{
		"ζ*'ヮ')ζ 웃우─!",
		"ζ° ͜ʖ ͡° ζ/ 웃우─! 프로듀서 로우탓-치!",
		"ζ*'ヮ')ζ / 프로듀사 하이탓-치! 이예이!",
		"ζ*'ヮ')ζ 오늘은 숙주나물 축제에요!",
		"ζ*'ヮ')ζ 저기, 타카츠키 야요이, 14세입니다. 저, 집이 빈곤해서, 저도 무언가 할 수 있는게 없을까 생각해서 아이돌이 되보려고 했습니다. 잘 부탁드립니다! 이에이!",
		"ζ*'ヮ')ζ타카츠키 야요이, 14세입니다-! 저, 건강이 장점이니까, 아무리 많은 일이어도 걱정 없어요. 그러니까 척척 일해나가서, 같이 톱 아이돌을 목표해주세요. 에헤헤♪",
		"ζ*'ヮ')ζ 타카츠키 야요이에요! 조금이라도 집안에 보탬이 되지 않을까 해서 아이돌을 시작했어요. 저, 건강과 의욕만큼은 분명하니까, 프로듀스 잘 부탁드려요! 에헤헷♪",
		"ζ*'ヮ')ζ 최근 상점가에서 쇼핑을 하고 있으면, 누가 말을 거는 일이 많아졌어요-. 이것도, 아이돌로서 유명해졌다는 걸까나! 웃우-, 프로듀서 감사합니다-!",
		"ζ*'ヮ')ζ 수영복을 입으면, 청소로 젖어도 갈아입지 않아도 괜찮아서 득일까나! 프로듀서도 수영복으로 갈아입어서, 같이 놀지 않을래요? 아, 물론 청소를 끝내는 것이 우선이에요-!",
		"ζ*'ヮ')ζ 이 옷, 팔랑팔랑~ 둥실둥실~ 거려서, 왠지 제가 아닌 것 같아요-. 조금 쑥스러울지도! 에헤헤.... 프로듀서, 저, 팔랑팔랑하고 둥실둥실 거리는 거, 어울리나요-?",
		"ζ*'ヮ')ζ 프, 프로듀서. 미아 오리씨들이 저를 따라버려서, 모두가 뒤를 쫒아 오고 있어요-! 어...어쩌면 좋죠-! 하앗, 저기, 모두? 저는 엄마가 아니에요-!",
		"ζ*'ヮ')ζ 프로듀서, 봐주세요-! 오리씨와 오리 엄마와 개구리씨와 돼지씨와... 그리고... 어쨌든, 잔뜩 잔뜩 있는 모두와 노래할게요-♪ 에헤헷, 귀여워요-!",
		"ζ*'ヮ')ζ 에헤헷, 늑대씨 의상이에요-! 커흠-, 먹어버릴테다-! ...앗, 사실은 먹지 않을테니까 무서워하지 말아주세요~! 상냥한 늑대씨가 되고 싶으니깐요!",
	}
	guguresp []string = []string{
		"구구구",
		"@scarlet9",
	}
)

func specialResponses(bot *BaseBot, e *slack.MessageEvent) {
	// ignore all bot_message
	if e.SubType == "bot_message" {
		return
	}
	
	if strings.Contains(e.Text, "72") || strings.Contains(e.Text, "치하야") || strings.Contains(e.Text, "큿") {
		postResponse(bot, e.Channel, ":kutt:", "치하야", "큿")
	}
	if strings.Contains(e.Text, "크킄") {
		postResponse(bot, e.Channel, ":chuni:", "Dark Flame Master", "흐콰한다")
	}
	if strings.Contains(e.Text, "안두인") {
		randomResponse(bot, e.Channel, ":anduin:", "안두인", anduinresp)
	}
	if strings.Contains(e.Text, "웃우") {
		randomResponse(bot, e.Channel, ":yayoyee:", "타카츠키 야요이", yayoiresp)
	}
	if strings.Contains(e.Text, "혼란하다 혼란해") {
		postResponse(bot, e.Channel, ":honse:", "혼세마왕", "혼세혼세")
	}
	if strings.Contains(e.Text, "비둘기") {
		randomResponse(bot, e.Channel, ":gugu:", "비둘기", guguresp)
	}
	if strings.Contains(e.Text, "신촌 셔틀") {
		processShuttleCommand(bot, e.Channel)
	}
}
