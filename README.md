[![Build Status](https://travis-ci.org/PoolC/slack_bot.svg)](https://travis-ci.org/PoolC/slack_bot)
[![Coverage Status](https://coveralls.io/repos/PoolC/slack_bot/badge.svg?branch=master&service=github)](https://coveralls.io/github/PoolC/slack_bot?branch=master)
# slack_bot
우리들의 장난감 메우와 안즈

# Prerequirement for develop
- go
- godep (recommended)
  - `go get github.com/tools/godep`

# Prepare build environment for beginner(sample)
`go get github.com/Perlmint/goautoenv`

## Linux/OSX
```bash
git clone https://github.com/PoolC/slack_bot.git
cd slack_bot
goautoenv init
source .goenv/bin/activate
```
## Windows(powershell)
```
git clone https://github.com/PoolC/slack_bot.git
cd slack_bot
goautoenv init
.\.goenv\bin\activate.ps1
```

# Build instruction(sample)
```bash
mkdir -p .workspace/src/github.com/PoolC
ln -s `pwd` .workspace/src/github.com/PoolC/slack_bot
export GOPATH=`pwd`/.workspace
cd $GOPATH/src/github.com/PoolC/slack_bot
# instll dependencies
godep restore
# build it
go build
```

# Run testcode
```bash
go test -cover ./...
```

# Update dependencies
```bash
godep save -r
```

# command
slash command서버.  
sample로 echo server가 구현되어있음.

# bot
RTM bot.  
anzu, meu 두개의 구현체가 존재. 

# config file
암호화된 prod_config를 실제 서비스 실행시 적용함.  
`openssl`의 `enc`명령을 사용, `aes-256-ecb` cipher 사용해서 암/복호화 수행함.  
`openssl`의 `enc`사용법은 https://www.openssl.org/docs/manmaster/apps/enc.html 참조.  
암호화에 사용된 키는 *#bot_test* 채널에 pin되어있으니 참고.
