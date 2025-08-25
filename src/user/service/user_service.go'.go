package service

import (
	"errors"
	"gorm.io/gorm"
	"st-novel-go/src/user/dao"
	"st-novel-go/src/user/model"
	"st-novel-go/src/utils"
)

// --- Payloads ---
type RegisterPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserPayload struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
	Phone  string `json:"phone"`
}

type NotificationSettingPayload struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Enabled bool   `json:"enabled"`
}

type UpdateUserSettingsPayload struct {
	User          UpdateUserPayload            `json:"user"`
	Notifications []NotificationSettingPayload `json:"notifications"`
}

// --- Responses ---
type LoginResponse struct {
	User  model.User `json:"user"`
	Token string     `json:"token"`
}

type NotificationSetting struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type SecuritySetting struct {
	Title       string `json:"title"`
	Status      string `json:"status"`
	StatusClass string `json:"statusClass"`
	Action      string `json:"action"`
}

type UserSettingsResponse struct {
	User             *model.User           `json:"user"`
	Notifications    []NotificationSetting `json:"notifications"`
	SecuritySettings []SecuritySetting     `json:"securitySettings"`
	ProPlanFeatures  []string              `json:"proPlanFeatures"`
}

// --- Services ---

func Register(payload RegisterPayload) (*model.User, error) {
	_, err := dao.FindUserByEmail(payload.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &model.User{
		Email:    payload.Email,
		Password: payload.Password,
		Name:     payload.Name,
	}

	err = dao.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func Login(payload LoginPayload) (*LoginResponse, error) {
	user, err := dao.FindUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found or password incorrect")
		}
		return nil, err
	}

	if !user.CheckPassword(payload.Password) {
		return nil, errors.New("user not found or password incorrect")
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	response := &LoginResponse{
		User:  *user,
		Token: token,
	}
	// Do not expose password hash
	response.User.Password = ""

	return response, nil
}

func GetUserSettings(userID uint) (*UserSettingsResponse, error) {
	user, err := dao.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	user.Password = "" // Ensure password is not sent

	// MOCK DATA: These would be loaded from their own tables in a real app.
	notifications := []NotificationSetting{
		{ID: 1, Title: "产品更新", Description: "获取新功能和产品更新的通知", Enabled: true},
		{ID: 2, Title: "账户动态", Description: "关于账户安全和账单的重要通知", Enabled: true},
		{ID: 3, Title: "每周摘要", Description: "每周接收您的创作活动摘要", Enabled: false},
	}
	securitySettings := []SecuritySetting{
		{Title: "密码", Status: "上次更新于 3 个月前", StatusClass: "text-yellow-600 dark:text-yellow-400", Action: "修改"},
		{Title: "双重验证 (2FA)", Status: "您尚未启用双重验证", StatusClass: "text-gray-500", Action: "启用"},
		{Title: "登录设备", Status: "2 个活跃会话", StatusClass: "text-gray-500", Action: "管理"},
		{Title: "登录历史", Status: "查看最近的登录活动", StatusClass: "text-gray-500", Action: "查看"},
	}
	proPlanFeatures := []string{
		"无限次创作",
		"高级 AI 模型访问",
		"优先技术支持",
		"云端同步与备份",
	}

	return &UserSettingsResponse{
		User:             user,
		Notifications:    notifications,
		SecuritySettings: securitySettings,
		ProPlanFeatures:  proPlanFeatures,
	}, nil
}

func UpdateUserSettings(userID uint, payload UpdateUserSettingsPayload) error {
	user, err := dao.FindUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Update user fields
	user.Name = payload.User.Name
	user.Avatar = payload.User.Avatar
	user.Bio = payload.User.Bio
	user.Phone = payload.User.Phone

	// In a real application, you would also update the notification settings in the database.
	// For now, we can just log it.
	// log.Printf("Updating notification settings for user %d: %+v", userID, payload.Notifications)

	return dao.UpdateUser(user)
}
