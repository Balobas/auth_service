package pgEntity

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type MqMessageRow struct {
	Uid              pgtype.UUID
	SubjectName      string
	Payload          string
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
	LastErrorMessage string
	SendAt           pgtype.Timestamp
}

func NewMqMessageRow() *MqMessageRow {
	return &MqMessageRow{}
}

func (m *MqMessageRow) New() *MqMessageRow {
	return &MqMessageRow{}
}

func (m *MqMessageRow) FromEntity(mqMessage entity.MqMessage) *MqMessageRow {
	m.Uid = pgtype.UUID{
		Bytes:  mqMessage.Uid,
		Status: pgtype.Present,
	}
	m.SubjectName = mqMessage.SubjectName
	m.Payload = string(mqMessage.Payload)
	m.LastErrorMessage = mqMessage.LastErrorMessage

	m.CreatedAt = pgtype.Timestamp{
		Time:   mqMessage.CreatedAt,
		Status: pgtype.Present,
	}
	m.UpdatedAt = pgtype.Timestamp{Time: mqMessage.UpdatedAt, Status: pgtype.Present}
	if mqMessage.UpdatedAt.Equal(time.Time{}) {
		m.UpdatedAt.Status = pgtype.Null
	}
	m.SendAt = pgtype.Timestamp{Time: mqMessage.SendAt, Status: pgtype.Present}
	if mqMessage.SendAt.Equal(time.Time{}) {
		m.SendAt.Status = pgtype.Null
	}

	return m
}

func (m *MqMessageRow) ToEntity() entity.MqMessage {
	return entity.MqMessage{
		Uid:              m.Uid.Bytes,
		SubjectName:      m.SubjectName,
		Payload:          []byte(m.Payload),
		LastErrorMessage: m.LastErrorMessage,
		CreatedAt:        m.CreatedAt.Time,
		UpdatedAt:        m.UpdatedAt.Time,
		SendAt:           m.SendAt.Time,
	}
}

func (m *MqMessageRow) IdColumnName() string {
	return "uid"
}

func (m *MqMessageRow) Values() []interface{} {
	return []interface{}{
		m.Uid, m.SubjectName, m.Payload, m.CreatedAt, m.UpdatedAt, m.LastErrorMessage, m.SendAt,
	}
}

func (m *MqMessageRow) Columns() []string {
	return []string{
		"uid", "subject_name", "payload", "created_at", "updated_at", "last_error_msg", "send_at",
	}
}

func (m *MqMessageRow) Table() string {
	return "outbox_messages"
}

func (m *MqMessageRow) Scan(row pgx.Row) error {
	return row.Scan(&m.Uid, &m.SubjectName, &m.Payload, &m.CreatedAt, &m.UpdatedAt, &m.LastErrorMessage, &m.SendAt)
}

func (m *MqMessageRow) ColumnsForUpdate() []string {
	return []string{
		"updated_at", "last_error_msg", "send_at",
	}
}

func (m *MqMessageRow) ValuesForUpdate() []interface{} {
	return []interface{}{
		m.UpdatedAt, m.LastErrorMessage, m.SendAt,
	}
}

func (m *MqMessageRow) ConditionUidEqual() sq.Eq {
	return sq.Eq{"uid": m.Uid}
}

func (m *MqMessageRow) ConditionSendAtIsNull() sq.Eq {
	return sq.Eq{"send_at": pgtype.Timestamp{Status: pgtype.Null}}
}

func NewMqMessageRows() *Rows[*MqMessageRow, entity.MqMessage] {
	return &Rows[*MqMessageRow, entity.MqMessage]{}
}
