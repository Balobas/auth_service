package pgEntity

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const sessionTableName = "sessions"

var sessionTableColumns = []string{
	"uid",
	"user_uid",
	"created_at",
	"updated_at",
}

type SessionRow struct {
	Uid       pgtype.UUID
	UserUid   pgtype.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func NewSessionRow() *SessionRow {
	return &SessionRow{}
}

func (s *SessionRow) FromEntity(session entity.Session) *SessionRow {
	s.Uid = pgtype.UUID{
		Bytes:  session.Uid,
		Status: pgtype.Present,
	}
	s.UserUid = pgtype.UUID{
		Bytes:  session.UserUid,
		Status: pgtype.Present,
	}
	if session.CreatedAt.Unix() == 0 {
		s.CreatedAt = pgtype.Timestamp{
			Status: pgtype.Null,
		}
	} else {
		s.CreatedAt = pgtype.Timestamp{
			Time:   session.CreatedAt,
			Status: pgtype.Present,
		}
	}
	if session.UpdatedAt.Unix() == 0 {
		s.UpdatedAt = pgtype.Timestamp{
			Status: pgtype.Null,
		}
	} else {
		s.UpdatedAt = pgtype.Timestamp{
			Time:   session.UpdatedAt,
			Status: pgtype.Present,
		}
	}
	return s
}

func (s *SessionRow) ToEntity() entity.Session {
	return entity.Session{
		Uid:       s.Uid.Bytes,
		UserUid:   s.UserUid.Bytes,
		CreatedAt: s.CreatedAt.Time,
		UpdatedAt: s.UpdatedAt.Time,
	}
}

func (s *SessionRow) IdColumnName() string {
	return "uid"
}

func (s *SessionRow) Values() []interface{} {
	return []interface{}{
		s.Uid,
		s.UserUid,
		s.CreatedAt,
		s.UpdatedAt,
	}
}

func (s *SessionRow) Columns() []string {
	return sessionTableColumns
}

func (s *SessionRow) Table() string {
	return sessionTableName
}

func (s *SessionRow) Scan(row pgx.Row) error {
	return row.Scan(
		&s.Uid,
		&s.UserUid,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
}
func (s *SessionRow) ColumnsForUpdate() []string {
	return nil
}

func (s *SessionRow) ValuesForUpdate() []interface{} {
	return nil
}

func (s *SessionRow) ConditionUidEqual() sq.Eq {
	return sq.Eq{
		"uid": s.Uid,
	}
}

func (s *SessionRow) ConditionUserUidEqual() sq.Eq {
	return sq.Eq{
		"user_uid": s.UserUid,
	}
}

type SessionRows struct {
	sessions []*SessionRow
}

func NewSessionRows() *SessionRows {
	return &SessionRows{}
}

func (s *SessionRows) ScanAll(rows pgx.Rows) error {
	for rows.Next() {
		newRow := &SessionRow{}

		if err := newRow.Scan(rows); err != nil {
			return err
		}
		s.sessions = append(s.sessions, newRow)
	}

	return nil
}

func (s *SessionRows) ToEntities() []entity.Session {
	if len(s.sessions) == 0 {
		return nil
	}

	res := make([]entity.Session, len(s.sessions))

	for i := 0; i < len(s.sessions); i++ {
		res[i] = s.sessions[i].ToEntity()
	}

	return res
}
