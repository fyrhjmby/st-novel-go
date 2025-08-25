package model

// ApiProvider represents detailed information about a supported API provider for the main list view.
type ApiProvider struct {
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	Description string `json:"description"`
	StatusText  string `json:"statusText"`
	ActiveKeys  int    `json:"activeKeys"`
	TotalCalls  string `json:"totalCalls"`
}

// ModalProvider represents simplified provider information for the selection modal.
type ModalProvider struct {
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	Description string `json:"description"`
}

// APIKeyResponse is the DTO for sending API key data to the frontend.
type APIKeyResponse struct {
	ID            uint      `json:"id"`
	Provider      string    `json:"provider"`
	ProviderShort string    `json:"providerShort"`
	Name          string    `json:"name"`
	KeyFragment   string    `json:"keyFragment"`
	Model         string    `json:"model"`
	Calls         string    `json:"calls"`
	Status        KeyStatus `json:"status"`
	Created       string    `json:"created"`
	BaseURL       string    `json:"baseUrl"`
}
