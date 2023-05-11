package processor

import (
	"fmt"
	"github.com/acs-dl/unverified-svc/internal/config"
	"github.com/acs-dl/unverified-svc/internal/data"
	"github.com/acs-dl/unverified-svc/internal/data/postgres"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	serviceName = data.ModuleName + "-processor"

	//add needed actions for module
	SetUsersAction    = "set_users"
	DeleteUsersAction = "delete_users"
)

type Processor interface {
	HandleNewMessage(msg data.ModulePayload) error
}

type processor struct {
	log    *logan.Entry
	usersQ data.Users
}

var handleActions = map[string]func(proc *processor, msg data.ModulePayload) error{
	SetUsersAction:    (*processor).handleSetUsersAction,
	DeleteUsersAction: (*processor).handleDeleteUsersAction,
}

func NewProcessor(cfg config.Config) Processor {
	return &processor{
		log:    cfg.Log().WithField("service", serviceName),
		usersQ: postgres.NewUsersQ(cfg.DB()),
	}
}

func (p *processor) HandleNewMessage(msg data.ModulePayload) error {
	p.log.Infof("handling message with id `%s`", msg.RequestId)

	err := validation.Errors{
		"action": validation.Validate(msg.Action, validation.Required, validation.In(SetUsersAction, DeleteUsersAction)),
	}.Filter()
	if err != nil {
		p.log.WithError(err).Errorf("no such action to handle for message with id `%s`", msg.RequestId)
		return errors.Wrap(err, fmt.Sprintf("no such action `%s` to handle for message with id `%s`", msg.Action, msg.RequestId))
	}

	requestHandler := handleActions[msg.Action]
	if err = requestHandler(p, msg); err != nil {
		p.log.WithError(err).Errorf("failed to handle message with id `%s`", msg.RequestId)
		return err
	}

	p.log.Infof("finish handling message with id `%s`", msg.RequestId)
	return nil
}
