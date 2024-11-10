package services

import (
	"crypto/md5"
	"errors"
	"fmt"
	dao "proyecto_arqui_soft_2/users-api/dao"
	domain "proyecto_arqui_soft_2/users-api/domain"

	//e "proyecto_arqui_soft_2/users-api/utils"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Repository interface {
	GetUsuariobyEmail(email string) (dao.Usuario, error)
	GetUsuariobyID(id int64) (dao.Usuario, error)

	// actualizar cache y memcache
	Actualizar(usuario domain.UsuarioData) error
	CrearUsuario(newusuario dao.Usuario) (dao.Usuario, error)
}

type Service struct {
	mainRepository      Repository
	cacheRepository     Repository
	memcachedRepository Repository
}

func NewService(mainRepo, cacheRepo, memcachedRepo Repository) Service {
	return Service{
		mainRepository:      mainRepo,
		cacheRepository:     cacheRepo,
		memcachedRepository: memcachedRepo,
	}
}

func Usuario(user dao.Usuario) domain.UsuarioData {

	var us domain.UsuarioData

	us.Nombre = user.Nombre
	us.Apellido = user.Apellido
	us.Tipo = user.Tipo
	us.Email = user.Email
	us.Passwordhash = user.Passwordhash
	us.UsuarioID = user.UsuarioID

	return us
}

func generateHash(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

var jwtSecreto = []byte("llave_secreta")

func generarJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecreto)
}

func gettoken(password string, usuario domain.UsuarioData) (string, error) {

	hash := generateHash(password)

	if hash != usuario.Passwordhash {
		return "", fmt.Errorf("Contraseña incorrecta.")
	}

	token, err := generarJWT(usuario.Email)

	if err != nil {
		return "", fmt.Errorf("Error al generar JWT token: %w", err)
	}
	return token, nil

}

func (service Service) Login(email string, password string) (string, error) {

	if strings.TrimSpace(email) == "" {
		return "", errors.New("Debe ingresar un email.")
	}

	if strings.TrimSpace(password) == "" {
		return "", errors.New("Debe ingresar una contraseña.")
	}

	var token string
	var error error

	usuarioo, err := service.cacheRepository.GetUsuariobyEmail(email)

	if err == nil {

		result := Usuario(usuarioo)

		token, error = gettoken(password, result)

		return token, error

	}

	usuarioo, err = service.memcachedRepository.GetUsuariobyEmail(email)

	if err == nil {

		result := Usuario(usuarioo)

		token, error = gettoken(password, result)

		service.cacheRepository.Actualizar(result)

		return token, error

	}

	usuarioo, err = service.mainRepository.GetUsuariobyEmail(email)

	if err == nil {

		result := Usuario(usuarioo)

		token, error = gettoken(password, result)

		// actualizar cache
		service.cacheRepository.Actualizar(result)

		// actualizar memcache
		service.memcachedRepository.Actualizar(result)

		return token, error

	}

	return token, error
}

func (service Service) GetUsuariobyEmail(email string) (domain.UsuarioData, error) {

	// primero me fijo en la cache
	usuarioo, err := service.cacheRepository.GetUsuariobyEmail(email)

	if err == nil {

		result := Usuario(usuarioo)

		return result, nil

	}

	usuarioo, err = service.memcachedRepository.GetUsuariobyEmail(email)

	if err == nil {

		result := Usuario(usuarioo)

		service.cacheRepository.Actualizar(result)

		return result, nil

	}

	usuarioo, err = service.mainRepository.GetUsuariobyEmail(email)

	if err == nil {

		result := Usuario(usuarioo)

		// actualizar cache
		service.cacheRepository.Actualizar(result)

		// actualizar memcache
		service.memcachedRepository.Actualizar(result)

		return result, nil

	}

	usuarioencontrado := Usuario(usuarioo)

	return usuarioencontrado, nil
}

func (service Service) CrearUsuario(newusuario domain.UsuarioData) (domain.UsuarioData, error) {

	var user dao.Usuario

	user.Nombre = newusuario.Nombre
	user.Apellido = newusuario.Apellido
	user.Email = newusuario.Email

	hash := generateHash(newusuario.Passwordhash)

	user.Passwordhash = hash
	newusuario.Passwordhash = user.Passwordhash
	user.Tipo = newusuario.Tipo

	user, err := service.mainRepository.CrearUsuario(user)

	if err != nil {
		return newusuario, fmt.Errorf("error caching new user: %w", err)
	}

	newusuario.UsuarioID = user.UsuarioID

	service.memcachedRepository.CrearUsuario(user)
	service.cacheRepository.CrearUsuario(user)

	return newusuario, nil
}

func (service Service) GetUsuariobyID(id int64) (domain.UsuarioData, error) {

	// primero me fijo en la cache
	usuarioo, err := service.cacheRepository.GetUsuariobyID(id)

	if err == nil {
		log.Println("Datos obtenidos desde cache")

		result := Usuario(usuarioo)

		return result, nil

	}

	usuarioo, err = service.memcachedRepository.GetUsuariobyID(id)

	if err == nil {
		log.Println("Datos obtenidos desde Memcached")

		result := Usuario(usuarioo)

		service.cacheRepository.Actualizar(result)

		return result, nil

	}

	usuarioo, err = service.mainRepository.GetUsuariobyID(id)

	if err == nil {
		log.Println("Datos obtenidos desde base de datos")

		result := Usuario(usuarioo)

		// actualizar cache
		service.cacheRepository.Actualizar(result)

		// actualizar memcache
		service.memcachedRepository.Actualizar(result)

		return result, nil

	}

	usuarioencontrado := Usuario(usuarioo)

	return usuarioencontrado, nil
}
