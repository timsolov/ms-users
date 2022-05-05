package repository

//go:generate mockgen -destination=../../infrastructure/repository/mockrepo/mockrepo.go -package=mockrepo ms-users/app/domain/repository Repository

// Repository is an API for work with database
type Repository interface {
	UserRepository
}
