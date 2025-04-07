package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	guardianEntity "github.com/winartodev/apollo/modules/guardian/entities"
)

var (
	ErrorInvalidLengthPermissionID = errors.New("different length between permission ids and and slugs")
)

type GuardianPermissionRepositoryItf interface {
	GetApplicationServicePermissionsDB(ctx context.Context, userID int64, roleID int64, applicationID int64, applicationServiceID int64) (res []guardianEntity.GuardianPermission, err error)
}

type GuardianPermissionRepository struct {
	DB *sql.DB
}

func NewGuardianPermissionRepository(repository GuardianPermissionRepository) GuardianPermissionRepositoryItf {
	return &GuardianPermissionRepository{
		DB: repository.DB,
	}
}

func (r *GuardianPermissionRepository) GetApplicationServicePermissionsDB(ctx context.Context, userID int64, roleID int64, applicationID int64, applicationServiceID int64) (res []guardianEntity.GuardianPermission, err error) {
	stmt, err := r.DB.PrepareContext(ctx, GetServicePermissionByUserID)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res = []guardianEntity.GuardianPermission{}
	var serviceID int64
	var permissionIDs []int64
	var permissionSlugs []string
	err = stmt.QueryRowContext(ctx,
		userID,
		roleID,
		applicationID,
		applicationServiceID,
	).Scan(
		&serviceID,
		pq.Array(&permissionIDs),
		pq.Array(&permissionSlugs),
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	res, err = r.buildPermissionToList(permissionIDs, permissionSlugs)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (r *GuardianPermissionRepository) buildPermissionToList(ids []int64, slugs []string) (res []guardianEntity.GuardianPermission, err error) {

	if len(ids) != len(slugs) {
		return nil, ErrorInvalidLengthPermissionID
	}

	res = make([]guardianEntity.GuardianPermission, 0)
	for i, permissionID := range ids {
		res = append(res, guardianEntity.GuardianPermission{
			ID:   permissionID,
			Slug: slugs[i],
		})
	}

	return res, nil
}
