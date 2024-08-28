package blacklistRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *BlacklistRepository) CreateBlackListUser(ctx context.Context, blackListUser entity.BlackListUser) error {
	row := pgEntity.NewUsersBlackListRow().FromBlackListUserEntity(blackListUser)

	if err := r.Create(ctx, row); err != nil {
		return errors.Wrapf(err, "failed to insert user with uid %s to black list", blackListUser.User.Uid)
	}

	return nil
}

func (r *BlacklistRepository) GetBlackListUser(ctx context.Context, userUid uuid.UUID) (entity.BlackListElement, error) {
	row := pgEntity.NewUsersBlackListRow().FromBlackListEntity(entity.BlackListElement{Uid: userUid})

	if err := r.GetOne(ctx, row, row.ConditionUserUidEqual()); err != nil {
		return entity.BlackListElement{}, errors.Wrapf(err, "failed to get")
	}

	return row.ToBlackListEntity(), nil
}

func (r *BlacklistRepository) UpdateBlackListUser(ctx context.Context, elem entity.BlackListElement) error {
	row := pgEntity.NewUsersBlackListRow().FromBlackListEntity(elem)

	if err := r.Update(ctx, row, row.ConditionUserUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to update black list user with uid %s", elem.Uid)
	}

	return nil
}

func (r *BlacklistRepository) DeleteUserFromBlackList(ctx context.Context, userUid uuid.UUID) error {
	row := pgEntity.NewUsersBlackListRow().FromBlackListEntity(entity.BlackListElement{Uid: userUid})

	if err := r.Delete(ctx, row, row.ConditionUserUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to delete user with uid %s from black list", userUid)
	}

	return nil
}

func (r *BlacklistRepository) GetBlackListUsers(ctx context.Context, limit, offset int64) ([]entity.BlackListUser, error) {
	userRow := &pgEntity.UserRow{}
	blackListUserRow := &pgEntity.UsersBlackListRow{}

	userPreffix := "u"
	blackListUserPreffix := "b"

	stmt := fmt.Sprintf(
		"SELECT %s, %s FROM %s AS %s JOIN %s AS %s ON %s=%s LIMIT %d OFFSET %d",
		strings.Join(makePreffix(userPreffix, userRow.Columns()), ","),
		strings.Join(makePreffix(blackListUserPreffix, blackListUserRow.Columns()), ","),
		userRow.Table(), userPreffix,
		blackListUserRow.Table(), blackListUserPreffix,
		userPreffix+"."+userRow.IdColumnName(), blackListUserPreffix+"."+blackListUserRow.IdColumnName(),
		limit, offset,
	)

	rows, err := r.DB().Query(ctx, stmt, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]entity.BlackListUser, 0, limit)

	for rows.Next() {
		if err := rows.Scan(append(userRow.ValuesForScan(), blackListUserRow.ValuesForScan()...)...); err != nil {
			return nil, errors.Wrap(err, "failed to scan")
		}

		res = append(res, entity.BlackListUser{
			User:   userRow.ToEntity(),
			Reason: blackListUserRow.ToBlackListEntity().Reason,
		})
	}

	return res, nil
}

func makePreffix(preffix string, args []string) []string {
	res := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		res[i] = preffix + args[i]
	}
	return res
}
