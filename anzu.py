import asyncio

@asyncio.coroutine
def process(bot, message: 'message'):
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
