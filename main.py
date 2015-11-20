import config
from butterfield import Bot
import butterfield
import redis, logging
import meu, anzu

logging.basicConfig(filename='bot.log')

def main():
    r = redis.StrictRedis(host='localhost', port=6379, db=0)

    exit = None

    meu_bot = Bot(config.token_meu)
    meu_bot.listen(meu.process)

    anzu_bot = Bot(config.token_anzu)
    anzu_bot.listen(anzu.process(r))

    butterfield.run(meu_bot, anzu_bot)

if __name__ == "__main__":
    from daemonize import Daemonize
    from sys import argv
    if argv[1] == '-i':
        main()
    else:
        daemon = Daemonize(app="test_app", pid=argv[1], action=main)
        daemon.start()
