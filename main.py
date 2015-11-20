import config
import asyncio
from butterfield import Bot
import butterfield
import redis
r = redis.StrictRedis(host='localhost', port=6379, db=0)

exit = None

@asyncio.coroutine
def process_meu(bot, message: 'message'):
    if message.get('text', '').startswith('계산해라 메우'):
        expr = message['text'][8:]
        if expr in ('bot', 'message', 'expr', 'process_meu', 'process_anzu', 'anzu_dict', 'meu', 'anzu'):
            yield from bot.post(message['channel'], '이상한 짓 하지 말고 그냥 코드를 봐라 메우')
        elif '안즈쨩' in expr or '안즈쨩' in eval(expr):
            yield from bot.post(message['channel'], '일도 안하는놈을 부르게 하려고 하다니 메우 거절이다 메우')
        else:
            try:
                yield from bot.post(message['channel'], str(eval(expr)))
            except Exception:
                yield from bot.post(message['channel'], '에러났다 메우. 이상한것 좀 시키지 마라 메우')
    elif message.get('text', '') == '메우, 멱살':
        yield from bot.post(message['channel'], '사람은 일을 하고 살아야한다. 메우')
    elif message.get('text', '') == '메우메우 펫탄탄':
        yield from bot.post(message['channel'], '메메메 메메메메 메우메우\n메메메 메우메우\n펫땅펫땅펫땅펫땅 다이스키')

@asyncio.coroutine
def process_anzu(bot, message: 'message'):
    if message.get('text', '') == '사람은 일을 하고 살아야한다. 메우':
        yield from bot.post(message['channel'], '이거 놔라 이 퇴근도 못하는 놈이')
    elif message.get('text', '').startswith('안즈쨩 카와이'):
        yield from bot.post(message['channel'], "뭐... 뭐라는거야\n기억해 [key]/[val], 알려줘 [key] 만 시키라구")
    elif message.get('text', '').startswith('안즈쨩 기억해'):
        key,val = message['text'][8:].split('/',1)
        key = key.strip()
        val = val.strip()
        if not key or not val:
            yield from bot.post(message['channel'], '에...?')
        elif '안즈쨩 알려줘' in val:
            yield from bot.post(message['channel'], '에... 귀찮아')
        else:
            r.set(key, val)
            yield from bot.post(message['channel'], '에... 귀찮지만 기억했어')
    elif message.get('text', '').startswith('안즈쨩 알려줘'):
        key = message['text'][8:].strip()
        val = r.get(key)
        if not val:
            val = '그런거 몰라'
        if isinstance(val, bytes):
            val = val.decode('utf-8')
        yield from bot.post(message['channel'], val)
    elif message.get('text', '') == '안즈쨩 뭐해?':
        yield from bot.post(message['channel'], '숨셔')

meu = Bot(config.token_meu)
meu.listen(process_meu)

anzu = Bot(config.token_anzu)
anzu.listen(process_anzu)

butterfield.run(meu, anzu)
