package processor

import (
	"github.com/acs-dl/unverified-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) handleDeleteUsersAction(msg data.ModulePayload) error {
	p.log.Infof("start handle message action with id `%s`", msg.RequestId)

	for _, user := range msg.Users {
		err := p.validateUser(user)
		if err != nil {
			p.log.WithError(err).Errorf("failed to validate fields for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to validate fields")
		}

		if err = p.usersQ.FilterByModules(user.Module).FilterByModuleIds(user.ModuleId).Delete(); err != nil {
			p.log.WithError(err).Errorf("failed to delete user in db for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to delete user in user db")
		}
	}

	p.log.Infof("finish handle message action with id `%s`", msg.RequestId)
	return nil
}
