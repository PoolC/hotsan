import config
from butterfield import Bot
import butterfield
import redis
import meu, anzu
r = redis.StrictRedis(host='localhost', port=6379, db=0)

exit = None

meu_bot = Bot(config.token_meu)
meu_bot.listen(meu.process)

anzu_bot = Bot(config.token_anzu)
anzu_bot.listen(anzu.process)

butterfield.run(meu_bot, anzu_bot)
