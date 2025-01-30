package services_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	dao "proyecto_arqui_soft_2/users-api/dao"
	domain "proyecto_arqui_soft_2/users-api/domain"
	repositories "proyecto_arqui_soft_2/users-api/repositories"
	servicio "proyecto_arqui_soft_2/users-api/services"
)

var (
	// Creamos mocks
	mainRepo      = repositories.NewMock()
	cacheRepo     = repositories.NewMock()
	memcachedRepo = repositories.NewMock()
	usersService  = servicio.NewService(mainRepo, cacheRepo, memcachedRepo)
)

func TestService(t *testing.T) {

	t.Run("GetUsuariobyID - Success from Cache", func(t *testing.T) {
		hashedPassword := servicio.GenerateHash("password123")
		mockUser := dao.Usuario{UsuarioID: 1, Nombre: "nombre1", Apellido: "apellido1", Email: "email1", Tipo: "estudiante", Passwordhash: hashedPassword}

		// nombre de la función + el parametro q necesita
		cacheRepo.On("GetUsuariobyID", int64(1)).Return(mockUser, nil).Once()

		result, err := usersService.GetUsuariobyID(1)

		assert.NoError(t, err)
		assert.Equal(t, "email1", result.Email)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetUsuariobyID - Not Found in Cache, Found in Memcached", func(t *testing.T) {
		hashedPassword := servicio.GenerateHash("password123")
		mockUser := dao.Usuario{UsuarioID: 1, Nombre: "nombre1", Apellido: "apellido1", Email: "email1", Tipo: "estudiante", Passwordhash: hashedPassword}

		cacheRepo.On("GetUsuariobyID", int64(1)).Return(dao.Usuario{}, errors.New("not found")).Once()
		memcachedRepo.On("GetUsuariobyID", int64(1)).Return(mockUser, nil).Once()
		cacheRepo.On("CrearUsuario", mockUser).Return(int64(1), nil).Once()

		result, err := usersService.GetUsuariobyID(1)

		assert.NoError(t, err)
		assert.Equal(t, "email1", result.Email)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetUsuariobyID - Not Found in Cache or Memcached, Found in Main Repo", func(t *testing.T) {
		hashedPassword := servicio.GenerateHash("password123")
		mockUser := dao.Usuario{UsuarioID: 1, Nombre: "nombre1", Apellido: "apellido1", Email: "email1", Tipo: "estudiante", Passwordhash: hashedPassword}

		cacheRepo.On("GetUsuariobyID", int64(1)).Return(dao.Usuario{}, errors.New("not found")).Once()
		memcachedRepo.On("GetUsuariobyID", int64(1)).Return(dao.Usuario{}, errors.New("not found")).Once()
		mainRepo.On("GetUsuariobyID", int64(1)).Return(mockUser, nil).Once()
		cacheRepo.On("CrearUsuario", mockUser).Return(int64(1), nil).Once()
		memcachedRepo.On("CrearUsuario", mockUser).Return(int64(1), nil).Once()

		result, err := usersService.GetUsuariobyID(1)

		assert.NoError(t, err)
		assert.Equal(t, "email", result.Email)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("GetUsuariobyID - Error in Main Repo", func(t *testing.T) {
		cacheRepo.On("GetUsuariobyID", int64(1)).Return(dao.Usuario{}, errors.New("not found")).Once()
		memcachedRepo.On("GetUsuariobyID", int64(1)).Return(dao.Usuario{}, errors.New("not found")).Once()
		mainRepo.On("GetUsuariobyID", int64(1)).Return(dao.Usuario{}, errors.New("db error")).Once()

		result, err := usersService.GetUsuariobyID(1)

		assert.Error(t, err)
		assert.Equal(t, "error getting user by ID: db error", err.Error())
		assert.Equal(t, domain.UsuarioData{}, result)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("CrearUsuario - Success", func(t *testing.T) {
		hashedPassword := servicio.GenerateHash("password123")
		newUser := dao.Usuario{Nombre: "nombre1", Apellido: "apellido1", Email: "email1", Tipo: "estudiante", Passwordhash: hashedPassword}
		mainRepo.On("CrearUsuario", newUser).Return(dao.Usuario(newUser), nil).Once()

		newUser.UsuarioID = 1

		cacheRepo.On("CrearUsuario", newUser).Return(dao.Usuario(newUser), nil).Once()
		memcachedRepo.On("CrearUsuario", newUser).Return(dao.Usuario(newUser), nil).Once()

		usuario, err := usersService.CrearUsuario(domain.UsuarioData{Nombre: "nombre1", Apellido: "apellido1", Email: "email1", Tipo: "estudiante", Passwordhash: "password123"})

		assert.NoError(t, err)
		assert.Equal(t, dao.Usuario(newUser), usuario)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("CrearUsuario - Error", func(t *testing.T) {

		hashedPassword := servicio.GenerateHash("password123")
		newUser := dao.Usuario{Nombre: "nombre1", Apellido: "apellido1", Email: "email1", Tipo: "estudiante", Passwordhash: hashedPassword}

		mainRepo.On("CrearUsuario", newUser).Return(int64(0), errors.New("db error")).Once()

		id, err := usersService.CrearUsuario(domain.UsuarioData{Nombre: "nombre1", Apellido: "apellido1", Email: "email1", Tipo: "estudiante", Passwordhash: hashedPassword})

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Equal(t, "error creating user: db error", err.Error())

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - Success", func(t *testing.T) {
		email := "email1"
		password := "password"
		hashedPassword := servicio.GenerateHash(password)

		mockUser := dao.Usuario{UsuarioID: 1, Email: email, Passwordhash: hashedPassword}
		cacheRepo.On("GetUsuariobyEmail", email).Return(mockUser, nil).Once()

		response, err := usersService.Login(email, password)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), mockUser.UsuarioID)
		assert.Equal(t, "token", response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - Invalid Credentials", func(t *testing.T) {
		email := "email1"
		password := "passwordmala"
		hashedPassword := servicio.GenerateHash(password)

		mockUser := dao.Usuario{UsuarioID: 1, Email: email, Passwordhash: hashedPassword}
		cacheRepo.On("GetUsuariobyEmail", email).Return(mockUser, nil).Once()

		response, err := usersService.Login(email, password)

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - User Not Found", func(t *testing.T) {
		email := "email1"
		password := "password"

		cacheRepo.On("GetUsuariobyEmail", email).Return(dao.Usuario{}, errors.New("not found")).Once()
		memcachedRepo.On("GetUsuariobyEmail", email).Return(dao.Usuario{}, errors.New("not found")).Once()
		mainRepo.On("GetUsuariobyEmail", email).Return(dao.Usuario{}, errors.New("not found")).Once()

		response, err := usersService.Login(email, password)

		assert.Error(t, err)
		assert.Equal(t, "error getting user by username from main repository: not found", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})

	t.Run("Login - Token Generation Error", func(t *testing.T) {
		email := "email1"
		password := "password"
		hashedPassword := servicio.GenerateHash(password)

		mockUser := dao.Usuario{UsuarioID: 1, Email: email, Passwordhash: hashedPassword}
		cacheRepo.On("GetUsuariobyEmail", email).Return(mockUser, nil).Once()

		response, err := usersService.Login(email, password)

		assert.Error(t, err)
		assert.Equal(t, "error generating token: token error", err.Error())
		assert.Equal(t, domain.LoginResponse{}, response)

		mainRepo.AssertExpectations(t)
		cacheRepo.AssertExpectations(t)
		memcachedRepo.AssertExpectations(t)
	})
}
