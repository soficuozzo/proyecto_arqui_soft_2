package controllers
import (
	"proyecto_arqui_soft_2/user-api/domain"
	"proyecto_arqui_soft_2/user-api/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// función para el login
func Login(c *gin.Context) {
	// se crea variable de tipo dto.LoginRequest
	var request domain.LoginRequest

	// chequeo de que el parametro q se manda por url no tenga errores
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.Resultado{
			Mensaje: fmt.Sprintf("Request invalido: %s", err.Error()),
		})
		return
	}

	// se llama al servicio
	token, err := services.UsuarioServicio.Login(request.Email, request.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.Resultado{
			Mensaje: fmt.Sprintf("Login no autorizado: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, domain.LoginResponse{
		Token: token,
	})
}

// función de obtener el usuario por el email
func GetUsuariobyEmail(c *gin.Context) {

	// se almacena el email que se manda por URL
	email := c.Param("email")

	// creamos una variable de tipo dto.UsuarioData
	var usuarioData domain.UsuarioData

	// llamamos al servicio
	usuarioData, err := services.UsuarioServicio.GetUsuariobyEmail(email)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, usuarioData)

}

// función de obtener el usuario por el id
func GetUsuariobyID(c *gin.Context) {

	// se almacena el email que se manda por URL
	id := c.Param("id")

	// hacemos que id se vuelva int64
	usuarioID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// llamamos al servicio
	usuarioData, er := services.UsuarioServicio.GetUsuariobyID(usuarioID)

	if er != nil {
		c.JSON(er.Status(), er)
		return
	}

	c.JSON(http.StatusOK, usuarioData)

}

// función para crear el usuario
func CrearUsuario(c *gin.Context) {

	//creamos una variable de tipo dto.UsuarioData
	var newusuario domain.UsuarioData

	// chequeo y almacenamiento del parametro q se manda por URL
	err := c.BindJSON(&newusuario)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// llamamos al servicio
	newusuario, er := services.UsuarioServicio.CrearUsuario(newusuario)

	if er != nil {
		c.JSON(er.Status(), er)
		return
	}

	c.JSON(http.StatusCreated, newusuario)
}
