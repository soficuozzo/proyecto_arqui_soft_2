package services

import (
	client "arqui_soft_mio/backend/clientes/usuario"
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	e "arqui_soft_mio/backend/utils"
)

// dependencias
type usuarioServicio struct {
	usuarioCliente client.UsuarioClienteInterface
}

// funciones
type usuarioServiceInterface interface {
	Login(email string, password string) (string, error)
	GetUsuariobyEmail(email string) (dto.UsuarioData, e.ApiError)
	GetUsuariobyID(id int64) (dto.UsuarioData, e.ApiError)
	CrearUsuario(newusuario dto.UsuarioData) (dto.UsuarioData, e.ApiError)
}

var (
	UsuarioServicio usuarioServiceInterface
)

// inicializamos el servicio y sus dependencias
func initUsuarioService(usuarioCliente client.UsuarioClienteInterface) usuarioServiceInterface {
	service := new(usuarioServicio)
	service.usuarioCliente = usuarioCliente
	return service
}

func init() {
	UsuarioServicio = initUsuarioService(client.UsuarioCliente)

}

// función para el hash de la contraseña
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

// función de login
func (s *usuarioServicio) Login(email string, password string) (string, error) {

	// chequeamos que esten completos los espacios del login
	if strings.TrimSpace(email) == "" {
		return "", errors.New("Debe ingresar un email.")
	}

	if strings.TrimSpace(password) == "" {
		return "", errors.New("Debe ingresar una contraseña.")
	}

	// creamos el hash
	hash := generateHash(password)

	// llamamos a la función en clientes
	usuarioo, err := s.usuarioCliente.GetUsuariobyEmail(email)

	if err != nil {
		return "", fmt.Errorf("Hubo un error al buscar el usuario en la Base de Datos.")

	}

	// el hash debe ser el mismo que el que esta guardado en la base de datos
	if hash != usuarioo.Passwordhash {
		return "", fmt.Errorf("Contraseña incorrecta.")
	}

	token, err := generarJWT(email)

	if err != nil {
		return "", fmt.Errorf("Error al generar JWT token: %w", err)
	}
	return token, nil
}

// función de obtener el usuario con su email
func (s *usuarioServicio) GetUsuariobyEmail(email string) (dto.UsuarioData, e.ApiError) {

	// creamos una variable de tipo model
	var usuarioo model.Usuario

	// llamamos a la función en cliente
	usuarioo, err := s.usuarioCliente.GetUsuariobyEmail(email)

	//creamos una variable de tipo dto
	var us dto.UsuarioData

	if err != nil {
		return us, e.NewBadRequestApiError("Usuario no encontrado")
	}

	// igualamos
	us.Nombre = usuarioo.Nombre
	us.Apellido = usuarioo.Apellido
	us.Tipo = usuarioo.Tipo
	us.Email = usuarioo.Email
	us.Passwordhash = usuarioo.Passwordhash
	us.UsuarioID = usuarioo.UsuarioID

	return us, nil
}

// función para obtener el usuario con su id
func (s *usuarioServicio) GetUsuariobyID(id int64) (dto.UsuarioData, e.ApiError) {

	// creamos una variable de tipo model
	var usuarioo model.Usuario

	// llamamos a la función en clientes
	usuarioo, err := s.usuarioCliente.GetUsuariobyID(id)

	// creamos una variable de tipo dto
	var us dto.UsuarioData

	if err != nil {
		return us, e.NewBadRequestApiError("Usuario no encontrado")
	}

	// igualamos
	us.Nombre = usuarioo.Nombre
	us.Apellido = usuarioo.Apellido
	us.Tipo = usuarioo.Tipo
	us.Email = usuarioo.Email
	us.Passwordhash = usuarioo.Passwordhash
	us.UsuarioID = usuarioo.UsuarioID

	return us, nil
}

// función de crear usuario
func (s *usuarioServicio) CrearUsuario(newusuario dto.UsuarioData) (dto.UsuarioData, e.ApiError) {

	// creamos una variable tipo model
	var user model.Usuario

	// igualamos
	user.Nombre = newusuario.Nombre
	user.Apellido = newusuario.Apellido
	user.Email = newusuario.Email

	// generamos el hash
	hash := generateHash(newusuario.Passwordhash)

	// igualamos
	user.Passwordhash = hash
	newusuario.Passwordhash = user.Passwordhash
	user.Tipo = newusuario.Tipo

	// llamamos a la función en clientes
	user = s.usuarioCliente.CrearUsuario(user)

	newusuario.UsuarioID = user.UsuarioID

	return newusuario, nil

}
