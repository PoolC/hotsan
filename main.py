import config
from butterfield import Bot
import butterfield
import redis
import meu, anzu

def main():
    r = redis.StrictRedis(host='localhost', port=6379, db=0)

    exit = None

    meu_bot = Bot(config.token_meu)
    meu_bot.listen(meu.process)

    anzu_bot = Bot(config.token_anzu)
    anzu_bot.listen(anzu.process)

    butterfield.run(meu_bot, anzu_bot)

if __name__ == "__main__":
    from daemonize import Daemonize
    from sys import argv
    daemon = Daemonize(app="test_app", pid=argv[1], action=main)
    daemon.start()
