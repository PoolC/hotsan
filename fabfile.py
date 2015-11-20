from fabric.api import run
from fabric.context_managers import cd
from fabric.contrib.files import exists
from fabric.operations import put

COPY_FILES = (
    'main.py',
    'anzu.py',
    'meu.py',
    'requirements.txt'
)

def deploy(path):
    with cd(path):
        for copy_file in COPY_FILES:
            put(copy_file)
        run('./env/bin/pip install -r requirements.txt', warn_only=True)
        if exists('bot.pid'):
            pid = run('cat bot.pid')
            run('kill -9 ' + pid, quiet=True)
        run('./env/bin/python main.py bot.pid')
