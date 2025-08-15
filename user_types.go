package main

// User 定义了用户的基本信息模型
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar,omitempty"`
	Plan     string `json:"plan"` // "免费版" | "专业版"
	Phone    string `json:"phone,omitempty"`
	Region   string `json:"region,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	Bio      string `json:"bio,omitempty"`
}

// LoginCredentials 定义了登录时所需的认证信息
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegistrationData 定义了用户注册时所需的数据
type RegistrationData struct {
	LoginCredentials
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	TermsAccepted bool   `json:"termsAccepted"`
}

// NotificationSetting 定义了单条通知设置
type NotificationSetting struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// SecuritySetting 定义了单条安全设置
type SecuritySetting struct {
	Title       string `json:"title"`
	Status      string `json:"status"`
	StatusClass string `json:"statusClass"`
	Action      string `json:"action"`
}

// UserSettings 定义了用户账户的完整设置信息
type UserSettings struct {
	User             User                  `json:"user"`
	Notifications    []NotificationSetting `json:"notifications"`
	SecuritySettings []SecuritySetting     `json:"securitySettings"`
	ProPlanFeatures  []string              `json:"proPlanFeatures"`
}
