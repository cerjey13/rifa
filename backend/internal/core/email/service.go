package email

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"rifa/backend/internal/types"
)

// Mailer defines the contract for sending emails.
type Mailer interface {
	SendPurchaseConfirmation(purchase types.Purchase) error
}

type mailerooClient struct {
	apiKey   string
	from     string
	to       string
	emailURL string
	client   *http.Client
}

// NewMailerooClient creates a new Maileroo API client.
func NewMailerooClient(apiKey, from, to, url string) Mailer {
	return &mailerooClient{
		apiKey:   apiKey,
		from:     from,
		to:       to,
		emailURL: url,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// SendPurchaseConfirmation sends a confirmation email using Maileroo.
func (m *mailerooClient) SendPurchaseConfirmation(purchase types.Purchase) error {
	htmlBody := fmt.Sprintf(
		EmailHTMLTemplate,
		purchase.UserID,
		purchase.CreatedAt.Format("02 Jan 2006 15:04"),
		purchase.Quantity,
		purchase.MontoUSD,
		purchase.MontoBs,
		purchase.PaymentMethod,
		purchase.TransactionDigits,
	)
	payload := EmailPayload{
		FromEmail: EmailObject{Address: m.from, DisplayName: "Compras"},
		ToEmail:   []EmailObject{{Address: m.to}},
		Subject:   "Compra recibida",
		HtmlBody:  htmlBody,
		Attachments: []File{{
			Name:        "Capture.jpg",
			ContentType: "image/jpeg",
			Content: base64.StdEncoding.EncodeToString(
				purchase.PaymentScreenshot,
			),
		}},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		m.emailURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("failed to create email request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("email send failed: status %s", resp.Status)
	}

	return nil
}
