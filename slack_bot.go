package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"

	bot "github.com/PoolC/slack_bot/bot"
	command "github.com/PoolC/slack_bot/command"
	. "github.com/PoolC/slack_bot/util"
	daemon "github.com/sevlyar/go-daemon"
	gcfg "gopkg.in/gcfg.v1"
	redis "gopkg.in/redis.v3"
)

var (
	signal = flag.String("s", "", `send signal to the daemon
		quit — graceful shutdown
		stop — fast shutdown`)

	cfg struct {
		General struct {
			PidFileName string
			LogFileName string
		}
		Bot map[string]*struct {
			Token string
		}
		Redis struct {
			Host string
			Port int
		}
		Command  command.CommandInfo
		Commands command.CommandsInfo
	}

	redisClient RedisClient
)

func main() {
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)

	error := gcfg.ReadFileInto(&cfg, "config")
	if error != nil {
		fmt.Println(error)
		return
	}

	cntxt := &daemon.Context{
		PidFileName: cfg.General.PidFileName,
		PidFilePerm: 0644,
		LogFileName: cfg.General.LogFileName,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func worker() {
	var wg sync.WaitGroup
	wg.Add(3)

	redisClient = &RedisClientWrap{redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})}

	anzu := bot.NewAnzu(cfg.Bot["Anzu"].Token, &stop, redisClient)
	meu := bot.NewMeu(cfg.Bot["Meu"].Token, &stop, redisClient)

	go bot.StartBot(anzu, &wg)
	go bot.StartBot(meu, &wg)

	server := command.NewServer(cfg.Commands, cfg.Command)

	go server.Start(&wg)

	wg.Wait()

	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}
