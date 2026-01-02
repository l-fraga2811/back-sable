package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/l-fraga2811/back-sable/internal/config"
)

type Client struct {
	projectURL string
	anonKey    string
	client     *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		projectURL: strings.TrimRight(cfg.SupabaseURL, "/"),
		anonKey:    cfg.SupabaseKey,
		client:     &http.Client{Timeout: 15 * time.Second},
	}
}

// ItemRow represents the database row structure
type ItemRow struct {
	ID          string  `json:"itemId"`
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"` // Changed to float64 to simplify vs Numeric
	Completed   bool    `json:"completed"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CreateItemPayload struct {
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Completed   bool    `json:"completed"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateItemPayload struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Completed   *bool    `json:"completed,omitempty"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type SignInCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpCredentials struct {
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Data     map[string]interface{} `json:"data"`
	Phone    string                 `json:"phone,omitempty"`
}

func (c *Client) SignIn(ctx context.Context, creds SignInCredentials) (AuthResponse, error) {
	var response AuthResponse
	// Supabase Auth endpoint: /auth/v1/token?grant_type=password
	q := url.Values{}
	q.Set("grant_type", "password")

	err := c.doJSON(ctx, http.MethodPost, "/auth/v1/token", "", q, creds, &response, nil)
	if err != nil {
		return AuthResponse{}, err
	}
	return response, nil
}

func (c *Client) SignUp(ctx context.Context, creds SignUpCredentials) (AuthResponse, error) {
	var response AuthResponse
	// Supabase Auth endpoint: /auth/v1/signup
	err := c.doJSON(ctx, http.MethodPost, "/auth/v1/signup", "", nil, creds, &response, nil)
	if err != nil {
		return AuthResponse{}, err
	}
	return response, nil
}

func (c *Client) ListItems(ctx context.Context, accessToken string) ([]ItemRow, error) {
	q := url.Values{}
	q.Set("select", "*")
	q.Set("order", "created_at.desc")
	return c.getItems(ctx, accessToken, q)
}

func (c *Client) GetItemByID(ctx context.Context, accessToken string, id string) (ItemRow, bool, error) {
	q := url.Values{}
	q.Set("select", "*")
	q.Set("itemId", "eq."+id)

	items, err := c.getItems(ctx, accessToken, q)
	if err != nil {
		return ItemRow{}, false, err
	}
	if len(items) == 0 {
		return ItemRow{}, false, nil
	}
	return items[0], true, nil
}

func (c *Client) CreateItem(ctx context.Context, accessToken string, payload CreateItemPayload) (ItemRow, error) {
	var created []ItemRow
	// Assuming /rest/v1/items is the correct endpoint path on Supabase
	if err := c.doJSON(ctx, http.MethodPost, "/rest/v1/items", accessToken, nil, payload, &created, map[string]string{"Prefer": "return=representation"}); err != nil {
		return ItemRow{}, err
	}
	if len(created) == 0 {
		return ItemRow{}, errors.New("supabase did not return created item")
	}
	return created[0], nil
}

func (c *Client) UpdateItem(ctx context.Context, accessToken string, id string, payload UpdateItemPayload) (ItemRow, bool, error) {
	q := url.Values{}
	q.Set("itemId", "eq."+id)

	var updated []ItemRow
	err := c.doJSON(ctx, http.MethodPatch, "/rest/v1/items", accessToken, q, payload, &updated, map[string]string{"Prefer": "return=representation"})
	if err != nil {
		return ItemRow{}, false, err
	}
	if len(updated) == 0 {
		return ItemRow{}, false, nil
	}
	return updated[0], true, nil
}

func (c *Client) DeleteItem(ctx context.Context, accessToken string, id string) (bool, error) {
	q := url.Values{}
	q.Set("itemId", "eq."+id)

	var deleted []ItemRow
	if err := c.doJSON(ctx, http.MethodDelete, "/rest/v1/items", accessToken, q, nil, &deleted, map[string]string{"Prefer": "return=representation"}); err != nil {
		return false, err
	}
	return len(deleted) > 0, nil
}

func (c *Client) getItems(ctx context.Context, accessToken string, q url.Values) ([]ItemRow, error) {
	var out []ItemRow
	if err := c.doJSON(ctx, http.MethodGet, "/rest/v1/items", accessToken, q, nil, &out, nil); err != nil {
		return nil, err
	}
	if out == nil {
		out = []ItemRow{}
	}
	return out, nil
}

func (c *Client) doJSON(ctx context.Context, method string, path string, accessToken string, q url.Values, payload any, out any, extraHeaders map[string]string) error {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	resp, err := c.do(ctx, method, path, accessToken, q, body, extraHeaders)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		if len(b) == 0 {
			return errors.New("supabase request failed")
		}
		// Log the error for debugging
		fmt.Printf("Supabase Error [Status %d]: %s\n", resp.StatusCode, string(b))
		return errors.New(string(b))
	}

	if out == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) do(ctx context.Context, method string, path string, accessToken string, q url.Values, body io.Reader, extraHeaders map[string]string) (*http.Response, error) {
	base := strings.TrimRight(c.projectURL, "/")
	fullURL := base + path
	if q != nil {
		fullURL = fullURL + "?" + q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.anonKey)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("empty response")
	}
	return resp, nil
}
