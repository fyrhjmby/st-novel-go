package mock

// --- API Management ---
type ApiProvider struct {
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	Description string `json:"description"`
	StatusText  string `json:"statusText"`
	ActiveKeys  int    `json:"activeKeys"`
	TotalCalls  string `json:"totalCalls"`
}

type ApiKey struct {
	ID            int     `json:"id"`
	Provider      string  `json:"provider"`
	ProviderShort string  `json:"providerShort"`
	Name          string  `json:"name"`
	KeyFragment   string  `json:"keyFragment"`
	Model         string  `json:"model"`
	Calls         string  `json:"calls"`
	Status        string  `json:"status"` // "启用" | "暂停"
	Created       string  `json:"created"`
	BaseURL       string  `json:"baseUrl,omitempty"`
	Description   string  `json:"description"`
	Temperature   float64 `json:"temperature"`
	MaxTokens     int     `json:"maxTokens"`
}

type ModalProvider struct {
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	Description string `json:"description"`
}

// --- Data Privacy ---
type DataCollectionSetting struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type DataUsageItem struct {
	Title    string `json:"title"`
	Tag      string `json:"tag"`
	Includes string `json:"includes"`
	Purpose  string `json:"purpose"`
}

type DataPermission struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

type PrivacySettingsData struct {
	CollectionSettings []DataCollectionSetting `json:"collectionSettings"`
	Usage              []DataUsageItem         `json:"usage"`
	Permissions        []DataPermission        `json:"permissions"`
	Promises           []string                `json:"promises"`
}

// --- System Settings ---
type Theme struct {
	Name string `json:"name"`
}

type SettingItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type SystemSettings struct {
	ActiveTheme          string        `json:"activeTheme"`
	ZoomLevel            float64       `json:"zoomLevel"` // Use float64 for better precision
	Language             string        `json:"language"`
	DateFormat           string        `json:"dateFormat"`
	NotificationSettings []SettingItem `json:"notificationSettings"`
	AppSettings          []SettingItem `json:"appSettings"`
}

// --- Usage Logs ---
type UsageStat struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Trend string `json:"trend"`
}

type ApiLog struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Endpoint  string `json:"endpoint"`
	Model     string `json:"model"`
	Tokens    string `json:"tokens"`
	Status    string `json:"status"` // "成功" | "失败"
	Duration  string `json:"duration"`
}

type ChartDataPoint struct {
	Label    string `json:"label"`
	Requests int    `json:"requests"`
	Tokens   int    `json:"tokens"`
}

type UsageData struct {
	Stats      []UsageStat      `json:"stats"`
	Logs       []ApiLog         `json:"logs"`
	ChartData  []ChartDataPoint `json:"chartData"`
	TotalLogs  int              `json:"totalLogs"`
	TotalPages int              `json:"totalPages"`
}
