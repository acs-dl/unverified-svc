package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer

	Amqp() *AmqpData
	JwtParams() *JwtCfg
	Registrator() RegistratorConfig
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	getter kv.Getter

	amqp        comfig.Once
	jwt         comfig.Once
	registrator comfig.Once
}

func New(getter kv.Getter) Config {
	return &config{
		getter:     getter,
		Databaser:  pgdb.NewDatabaser(getter),
		Copuser:    copus.NewCopuser(getter),
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
	}
}
