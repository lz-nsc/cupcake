package session

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
)

type IBeforeQuery interface {
	BeforeQuery(s *Session) error
}
type IAfterQuery interface {
	AfterQuery(s *Session) error
}
type IBeforeInsert interface {
	BeforeInsert(s *Session) error
}
type IAfterInsert interface {
	AfterInsert(s *Session) error
}
type IBeforeUpdate interface {
	BeforeUpdate(s *Session) error
}
type IAfterUpdate interface {
	AfterUpdate(s *Session) error
}
type IBeforeDelete interface {
	BeforeDelete(s *Session) error
}
type IAfterDelete interface {
	AfterDelete(s *Session) error
}

func (s *Session) TriggerHook(method string, value interface{}) {
	if s.Schema() == nil {
		return
	}
	target := s.Schema().Model
	if value != nil {
		target = value
	}

	switch method {
	case BeforeQuery:
		if v, ok := target.(IBeforeQuery); ok {
			v.BeforeQuery(s)
		}
	case AfterQuery:
		if v, ok := target.(IAfterQuery); ok {
			v.AfterQuery(s)
		}
	case BeforeInsert:
		if v, ok := target.(IBeforeInsert); ok {
			v.BeforeInsert(s)
		}
	case AfterInsert:
		if v, ok := target.(IAfterInsert); ok {
			v.AfterInsert(s)
		}
	case BeforeUpdate:
		if v, ok := target.(IBeforeUpdate); ok {
			v.BeforeUpdate(s)
		}
	case AfterUpdate:
		if v, ok := target.(IAfterUpdate); ok {
			v.AfterUpdate(s)
		}
	case BeforeDelete:
		if v, ok := target.(IBeforeDelete); ok {
			v.BeforeDelete(s)
		}
	case AfterDelete:
		if v, ok := target.(IAfterDelete); ok {
			v.AfterDelete(s)
		}
	}
}
