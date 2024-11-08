package repositories

import (
	//"encoding/json"
	//"errors"
	"encoding/json"
	"errors"
	"fmt"
	"proyecto_arqui_soft_2/users-api/dao"
	"proyecto_arqui_soft_2/users-api/domain"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedConfig struct {
	Host string
	Port string
}

type Memcached struct {
	client *memcache.Client
}

func usernameKey(username string) string {
	return fmt.Sprintf("username:%s", username)
}

func NewMemcached(config MemcachedConfig) Memcached {
	// Connect to Memcached
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	client := memcache.New(address)

	return Memcached{client: client}
}

func (repository Memcached) GetUsuariobyEmail(email string) (dao.Usuario, error) {

	eKey := fmt.Sprintf("user:email:%s", email)

	item, err := repository.client.Get(eKey)

	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return dao.Usuario{}, fmt.Errorf("Usuario no encontrado")
		}
		return dao.Usuario{}, fmt.Errorf("error fetching user by username from memcached: %w", err)
	}

	var user dao.Usuario
	// check de que esten bien lo encontrado

	if err := json.Unmarshal(item.Value, &user); err != nil {
		return dao.Usuario{}, fmt.Errorf("error unmarshaling user: %w", err)
	}

	return user, nil

}

func (repository Memcached) GetUsuariobyID(id int64) (dao.Usuario, error) {

	idKey := fmt.Sprintf("user:id:%d", id)

	item, err := repository.client.Get(idKey)

	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return dao.Usuario{}, fmt.Errorf("Usuario no encontrado")
		}
		return dao.Usuario{}, fmt.Errorf("error fetching user by username from memcached: %w", err)
	}

	var user dao.Usuario
	// check de que esten bien lo encontrado

	if err := json.Unmarshal(item.Value, &user); err != nil {
		return dao.Usuario{}, fmt.Errorf("error unmarshaling user: %w", err)
	}

	repository.client.Set(&memcache.Item{Key: "user:id:1", Value: item.Value})

	return user, nil

}

func (repository Memcached) Actualizar(usuario domain.UsuarioData) (int64, error) {

	// Serialize user data
	data, err := json.Marshal(usuario)

	fmt.Println("JSON serializado:", string(data))

	if err != nil {
		return 0, fmt.Errorf("error marshaling user: %w", err)
	}

	// Store user with ID as key and username as an alternate key
	idKey := fmt.Sprintf("user:id:%d", usuario.UsuarioID)

	if err := repository.client.Set(&memcache.Item{Key: idKey, Value: data}); err != nil {
		return 0, fmt.Errorf("error storing user in memcached: %w", err)
	}

	// Set key for username as well for easier lookup by username
	emailKey := fmt.Sprintf("user:email:%s", usuario.Email)

	if err := repository.client.Set(&memcache.Item{Key: emailKey, Value: data}); err != nil {
		return 0, fmt.Errorf("error storing username in memcached: %w", err)
	}

	return usuario.UsuarioID, nil

}

func (repository Memcached) EliminarUsuario(id int64) error {
	// Si el id es 0, la clave sería "user:id:0"
	idKey := fmt.Sprintf("user:id:%d", id)

	// Eliminar el usuario del cache
	err := repository.client.Delete(idKey)
	if err != nil {
		return fmt.Errorf("error eliminando usuario de Memcached: %w", err)
	}

	return nil
}

func (repository Memcached) CrearUsuario(newusuario dao.Usuario) (dao.Usuario, error) {

	// Serialize user data
	data, err := json.Marshal(newusuario)
	if err != nil {
		return newusuario, fmt.Errorf("error marshaling user: %w", err)
	}

	idKey := fmt.Sprintf("user:id:%d", newusuario.UsuarioID)
	if err := repository.client.Set(&memcache.Item{Key: idKey, Value: data}); err != nil {
		return newusuario, fmt.Errorf("error storing user in memcached: %w", err)
	}

	emailKey := fmt.Sprintf("user:email:%s", newusuario.Email)

	if err := repository.client.Set(&memcache.Item{Key: emailKey, Value: data}); err != nil {
		return newusuario, fmt.Errorf("error storing username in memcached: %w", err)
	}

	return newusuario, nil

}
