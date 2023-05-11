package processor

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/acs/unverified-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) validateUser(user data.User) error {
	return validation.Errors{
		"module": validation.Validate(user.Module, validation.Required),
	}.Filter() //other fields are optional and depends on module
}

func (p *processor) handleSetUsersAction(msg data.ModulePayload) error {
	p.log.Infof("start handle message action with id `%s`", msg.RequestId)

	for _, user := range msg.Users {
		err := p.validateUser(user)
		if err != nil {
			p.log.WithError(err).Errorf("failed to validate fields for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to validate fields")
		}

		if err = p.usersQ.Upsert(user); err != nil {
			p.log.WithError(err).Errorf("failed to create user in db for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to create user in user db")
		}
	}

	p.log.Infof("finish handle message action with id `%s`", msg.RequestId)
	return nil
}
