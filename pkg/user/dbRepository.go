package user

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	_pg "github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type Repository interface {
	CreateUser(dCtx context.Context, u *User) error
	UpdateUser(dCtx context.Context, u *User) error
	Fetch(dCtx context.Context, rID int) (user *User, err error)
	FetchByMobileNumber(dCtx context.Context, mobile string) (user *User, err error)
	FetchAllUsers(dCtx context.Context, req *UserRequest) (users []User, pagination Pagination, err error)
}

// NewRepositoryIn is function param struct of func `NewRepository`
type NewRepositoryIn struct {
	fx.In

	Log *logrus.Logger
	DB  *pg.DB `name:"gouserDB"`
}

// PGRepo is postgres implementation
type PGRepo struct {
	log *logrus.Logger
	db  *pg.DB
}

// NewDBRepository returns a new persistence layer object which can be used for
// CRUD on db
func NewDBRepository(i NewRepositoryIn) (Repo Repository, err error) {

	Repo = &PGRepo{
		log: i.Log,
		db:  i.DB,
	}

	return
}

func (r *PGRepo) CreateUser(ctx context.Context, u *User) (err error) {
	_, err = r.db.ModelContext(ctx, u).Insert()
	return
}

func (r *PGRepo) UpdateUser(ctx context.Context, u *User) (err error) {
	// var err error
	query := r.db.ModelContext(ctx, u) //.WherePK().Update()

	if u.Mobile != "" {
		query.Set("mobile=?", u.Mobile)
	}
	if u.ProfilePicture != "" {
		query.Set("profile_picture=?", u.ProfilePicture)
	}
	if u.Metadata != nil {
		query.Set("metadata=?", u.Metadata)
	}
	if !u.DOB.IsZero() {
		query.Set("dob=?", u.DOB)
	}
	if u.FirstName != "" {
		query.Set("first_name=?first_name,last_name=?last_name")
	}
	query.Set("updated_at=?", time.Now())
	k, err := query.WherePK().Update()
	if err != nil {
		r.log.Error(err.Error())
	}
	r.log.Info(k)

	return err
}

func (r *PGRepo) FetchAllUsers(dCtx context.Context, req *UserRequest) (users []User, pagination Pagination, err error) {
	users = []User{}
	query := r.db.ModelContext(dCtx, &User{}).Returning("*")
	if req.Mobile != nil {
		query.Where(`mobile ILIKE '%` + *req.Mobile + `%'`)
	}
	if req.Name != nil {
		nameString := strings.Split(*req.Name, " ")

		if len(nameString) == 1 {
			query.Where(`first_name ILIKE '%` + nameString[0] + `%'`)
		}
		if len(nameString) >= 2 {
			query.Where(`first_name ILIKE '%` + nameString[0] + `%'`)
			query.Where(`last_name ILIKE '%` + nameString[1] + `%'`)
		}
	}
	var count int
	if req.Limit != -1 {
		count, err = query.Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Order("user.id desc").SelectAndCount(&users)
		if err != nil {
			r.log.WithContext(dCtx).Info(dCtx, "unable to do pagination error :", err.Error())
			return
		}
	} else {
		err = query.Order("user.id desc").Select(&users)
		if err != nil {
			r.log.WithContext(dCtx).Info(dCtx, "unable to fetch all data :", err.Error())
			return
		}
	}
	pagination.TotalDataCount = count
	pagination.CurrentPage = req.Page
	d := float64(count) / float64(req.Limit)
	pagination.TotalPages = int(math.Ceil(d))
	if err != nil {
		if err == _pg.ErrNoRows {
			r.log.WithContext(dCtx).Info(dCtx, err.Error())
			return users, pagination, nil
		}
		return
	}
	return users, pagination, nil
}

func (r *PGRepo) Fetch(dCtx context.Context, rID int) (user *User, err error) {
	user = &User{
		ID: rID,
	}
	err = r.db.ModelContext(dCtx, user).Column("user.*").WherePK().Select()
	if err != nil {
		return
	}
	return
}

func (r *PGRepo) FetchByMobileNumber(dCtx context.Context, mobile string) (user *User, err error) {
	user = &User{
		Mobile: mobile,
	}
	err = r.db.ModelContext(dCtx, user).Column("user.*").Where("mobile=?mobile").Select()
	if err != nil {
		return
	}
	return
}
