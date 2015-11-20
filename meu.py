import asyncio
import logging

@asyncio.coroutine
def process(bot, message: 'message'):
    if message.get('text', '').startswith('계산해라 메우'):
        expr = message['text'][8:]
        if expr in ('bot', 'message', 'expr', 'process', 'meu', 'anzu'):
            yield from bot.post(message['channel'], '이상한 짓 하지 말고 그냥 코드를 봐라 메우')
        elif '안즈쨩' in expr or '안즈쨩' in eval(expr):
            yield from bot.post(message['channel'], '일도 안하는놈을 부르게 하려고 하다니 메우 거절이다 메우')
        else:
            try:
                yield from bot.post(message['channel'], str(eval(expr)))
            except Exception as e:
                logging.error(repr(e))
                yield from bot.post(message['channel'], '에러났다 메우. 이상한것 좀 시키지 마라 메우')
    elif message.get('text', '') == '메우, 멱살':
        yield from bot.post(message['channel'], '사람은 일을 하고 살아야한다. 메우')
    elif message.get('text', '') == '메우메우 펫탄탄':
        yield from bot.post(message['channel'], '메메메 메메메메 메우메우\n메메메 메우메우\n펫땅펫땅펫땅펫땅 다이스키')
