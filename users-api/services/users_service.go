package services

import (
	//"crypto/md5"
	//"encoding/hex"
	//"fmt"
	dao "proyecto_arqui_soft_2/users-api/dao"
	//domain "proyecto_arqui_soft_2/users-api/domain"
)

type Repository interface {
	GetAll() ([]dao.Usuario, error)
	GetByID(id int64) (dao.Usuario, error)
	GetByUsername(username string) (dao.Usuario, error)
	Create(user dao.Usuario) (int64, error)
	Update(user dao.Usuario) error
	Delete(id int64) error
}

type Tokenizer interface {
	GenerateToken(username string, userID int64) (string, error)
}

type Service struct {
	mainRepository      Repository
	cacheRepository     Repository
	memcachedRepository Repository
	tokenizer           Tokenizer
}

func NewService(mainRepository, cacheRepository, memcachedRepository Repository, tokenizer Tokenizer) Service {
	return Service{
		mainRepository:      mainRepository,
		cacheRepository:     cacheRepository,
		memcachedRepository: memcachedRepository,
		tokenizer:           tokenizer,
	}
}
