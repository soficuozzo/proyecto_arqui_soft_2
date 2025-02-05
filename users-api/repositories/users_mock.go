package repositories

import (
	"proyecto_arqui_soft_2/users-api/dao"
	"proyecto_arqui_soft_2/users-api/domain"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) GetUsuariobyEmail(email string) (dao.Usuario, error) {
	args := m.Called(email)
	if err := args.Error(1); err != nil {
		return dao.Usuario{}, err
	}

	return args.Get(0).(dao.Usuario), nil
}

func (m *Mock) GetUsuariobyID(id int64) (dao.Usuario, error) {
	args := m.Called(id)
	if err := args.Error(1); err != nil {
		return dao.Usuario{}, err
	}

	return args.Get(0).(dao.Usuario), nil
}

func (m *Mock) CrearUsuario(usuario dao.Usuario) (dao.Usuario, error) {
	args := m.Called(usuario)
	if err := args.Error(1); err != nil {
		return dao.Usuario{}, err
	}

	return args.Get(0).(dao.Usuario), nil

}

func (m *Mock) Actualizar(usuario domain.UsuarioData) error {

	args := m.Called(usuario)
	if err := args.Error(1); err != nil {
		return err
	}

	return nil
}

func (m *Mock) GenerarJWT(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}
