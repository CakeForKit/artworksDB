package api

import (
	"errors"
	"net/http"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/adminserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminRouter struct {
	adminServ adminserv.AdminService
	authe     auth.AuthEmployee
	authZ     auth.AuthZ
}

func (r *AdminRouter) Init(
	router *gin.RouterGroup, adminserv adminserv.AdminService,
	authe auth.AuthEmployee, authZ auth.AuthZ) {
	r.adminServ = adminserv
	r.authe = authe
	r.authZ = authZ
	gr := router.Group("employeelist")
	gr.GET("/", r.GetAllEmployees)
	gr.POST("/register-employee", r.Register)
	gr.PUT("/change-rights", r.ChangeRights)
	gru := router.Group("userlist")
	gru.GET("/", r.GetAllUsers)
}

// GetAllEmployees godoc
// @Summary Получить всех сотрудников
// @Description Возвращает список всех сотрудников
// @Tags Администратор
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {array} jsonreqresp.EmployeeResponse
// @Failure 401 "Не авторизован"
// @Router /admin/employeelist/ [get]
func (r *AdminRouter) GetAllEmployees(c *gin.Context) {
	ctx := c.Request.Context()
	employees, err := r.adminServ.GetAllEmployees(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	employeesResp := make([]jsonreqresp.EmployeeResponse, len(employees))
	for i, e := range employees {
		employeesResp[i] = e.ToEmployeeResponse()
	}
	c.JSON(http.StatusOK, employeesResp)
}

// GetAllUsers godoc
// @Summary Получить всех пользователей
// @Description Возвращает список всех пользователей
// @Tags Администратор
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {array} jsonreqresp.EmployeeResponse
// @Failure 401 "Не авторизован"
// @Router /admin/userlist/ [get]
func (r *AdminRouter) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := r.adminServ.GetAllUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	usersResp := make([]jsonreqresp.UserResponse, len(users))
	for i, u := range users {
		usersResp[i] = u.ToUserResponse()
	}
	c.JSON(http.StatusOK, usersResp)
}

// Register employee Handler
// @Summary Регистрация сотрудника
// @Description Регистрирует нового сотрудника
// @Tags Администратор
// @Accept json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body auth.RegisterEmployeeRequest true "Данные для регистрации"
// @Success 200 "Сотрудник зарегистрирован"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка авторизации"
// @Failure 409 "Попытка повторной регистрации"
// @Router /admin/employeelist/register-employee [post]
func (r *AdminRouter) Register(c *gin.Context) {
	ctx := c.Request.Context()

	adminID, err := r.authZ.AdminIDFromContext(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNotAuthZ) || errors.Is(err, auth.ErrHasNoRights) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var req auth.RegisterEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.authe.RegisterEmployee(ctx, req, adminID); err != nil {
		if errors.Is(err, employeerep.ErrDuplicateLoginEmp) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Change employee rights
// @Summary Изменить права сотрудника
// @Description Изменяет поле valid сотрудника
// @Tags Администратор
// @Accept json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.UpdateValidEmployeeRequest true "Данные для обновления"
// @Success 200 "Успешное обновление"
// @Failure 404 "Сотрудник не найден"
// @Router /admin/employeelist/change-rights [put]
func (r *AdminRouter) ChangeRights(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.UpdateValidEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.adminServ.ChangeEmployeeRights(ctx, uuid.MustParse(req.ID), req.Valid)
	if err != nil {
		if errors.Is(err, employeerep.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
