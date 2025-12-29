package pgEntity

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	basePgEntity "github.com/balobas/sport_city_common/repository/postgres/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		Bytes: mqMessage.Uid,
		Valid: true,
	}
	m.SubjectName = mqMessage.SubjectName
	m.Payload = string(mqMessage.Payload)
	m.LastErrorMessage = mqMessage.LastErrorMessage

	m.CreatedAt = pgtype.Timestamp{
		Time:  mqMessage.CreatedAt,
		Valid: true,
	}
	m.UpdatedAt = pgtype.Timestamp{Time: mqMessage.UpdatedAt, Valid: true}
	m.SendAt = pgtype.Timestamp{Time: mqMessage.SendAt, Valid: true}
	if mqMessage.SendAt.Equal(time.Time{}) {
		m.SendAt.Valid = false
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
	return sq.Eq{"send_at": pgtype.Timestamp{Valid: false}}
}

func NewMqMessageRows() *basePgEntity.Rows[*MqMessageRow, entity.MqMessage] {
	return &basePgEntity.Rows[*MqMessageRow, entity.MqMessage]{}
}
