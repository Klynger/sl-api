package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
)

func baseURL() string {
	if u := os.Getenv("SEED_BASE_URL"); u != "" {
		return u
	}

	return "http://localhost:8080"
}

// postJSON marshals body to JSON and POSTs it. The caller checks the status.
func postJSON(ctx context.Context, client *http.Client, base, path string, body any) (*http.Response, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func getJSON(ctx context.Context, client *http.Client, base, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	return client.Do(req)
}

// registerUser creates a user. created is false (with no error) when the user
// already exists (HTTP 409), which the caller uses as an idempotency signal.
func registerUser(ctx context.Context, client *http.Client, base, name, lastName, username, password string) (created bool, err error) {
	resp, err := postJSON(ctx, client, base, "/v1/users", map[string]any{
		"name":     name,
		"lastName": lastName,
		"username": username,
		"password": password,
	})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		return true, nil
	case http.StatusConflict:
		return false, nil
	default:
		return false, fmt.Errorf("register %s: unexpected status %d", username, resp.StatusCode)
	}
}

// login authenticates as the given user. The client's cookie jar captures the
// session, so later requests from the same client act as this user.
func login(ctx context.Context, client *http.Client, base, username, password string) error {
	resp, err := postJSON(ctx, client, base, "/v1/users/login", map[string]any{
		"username": username,
		"password": password,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login %s: unexpected status %d", username, resp.StatusCode)
	}

	return nil
}

// createProducts creates each named product, then returns the catalog's ids.
// Create returns no body, so the ids come from a follow-up GET /v1/products.
func createProducts(ctx context.Context, client *http.Client, base string, names []string) ([]string, error) {
	// Created sequentially on purpose: these are authenticated requests, and this
	// API re-issues the session cookie on every response, so firing them in
	// parallel races on the shared cookie jar and can drop the session.
	for _, name := range names {
		resp, err := postJSON(ctx, client, base, "/v1/products", map[string]any{"name": name})
		if err != nil {
			return nil, err
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return nil, fmt.Errorf("create product %q: unexpected status %d", name, resp.StatusCode)
		}
	}

	resp, err := getJSON(ctx, client, base, "/v1/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list products: unexpected status %d", resp.StatusCode)
	}

	var products []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(products))
	for _, p := range products {
		ids = append(ids, p.ID)
	}

	return ids, nil
}

// createGroup creates a group and returns its id (from the JSON response).
func createGroup(ctx context.Context, client *http.Client, base, name string) (string, error) {
	resp, err := postJSON(ctx, client, base, "/v1/groups", map[string]any{"name": name})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("create group %q: unexpected status %d", name, resp.StatusCode)
	}

	var out struct {
		GroupID string `json:"groupId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	return out.GroupID, nil
}

// createList creates a list inside the group and returns its id.
func createList(ctx context.Context, client *http.Client, base, groupID, name string) (string, error) {
	resp, err := postJSON(ctx, client, base, fmt.Sprintf("/v1/groups/%s/lists", groupID),
		map[string]any{
			"groupId": groupID,
			"name":    name,
		})

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("create list %q: unexpected status %d", name, resp.StatusCode)
	}

	var out struct {
		ListID string `json:"listId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	return out.ListID, nil
}

// addItem adds a product to a list.
func addItem(ctx context.Context, client *http.Client, base, listID, productID string) error {
	resp, err := postJSON(ctx, client, base, fmt.Sprintf("/v1/lists/%s/items", listID),
		map[string]any{
			"productId": productID,
		})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("add product %q into list %q: unexpected status %d", productID, listID, resp.StatusCode)
	}

	return nil
}

// inviteByUsername invites a user to the group by their username.
func inviteByUsername(ctx context.Context, client *http.Client, base, groupID, username string) error {
	resp, err := postJSON(ctx, client, base, fmt.Sprintf("/v1/groups/%s/invites", groupID),
		map[string]any{"invitedUsername": username})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("invite %q to group %q: unexpected status %d", username, groupID, resp.StatusCode)
	}

	return nil
}

// acceptInvite accepts the pending invite for whoever the client is logged in as.
func acceptInvite(ctx context.Context, client *http.Client, base, groupID string) error {
	resp, err := postJSON(ctx, client, base, fmt.Sprintf("/v1/groups/%s/invites/accept", groupID),
		map[string]any{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("accept invite for group %q: unexpected status %d", groupID, resp.StatusCode)
	}

	return nil
}

func main() {
	base := baseURL()

	// The cookie jar is what makes the session work: after login it stores the
	// auth cookie and attaches it to every later request from this client.
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{Jar: jar}
	ctx := context.Background()

	// Registering dev is the idempotency gate: if it already exists, we assume
	// the database is seeded and stop.
	created, err := registerUser(ctx, client, base, "Dev", "User", "dev", "pass")
	if err != nil {
		log.Fatalf("register dev: %v", err)
	}
	if !created {
		log.Println("dev already exists; the database looks seeded, nothing to do")
		return
	}
	log.Println("registered dev")

	if err := login(ctx, client, base, "dev", "pass"); err != nil {
		log.Fatalf("login dev: %v", err)
	}
	log.Println("logged in as dev")

	productIDs, err := createProducts(ctx, client, base, []string{"Milk", "Eggs", "Bread", "Coffee", "Bananas"})
	if err != nil {
		log.Fatalf("create products: %v", err)
	}

	groupID, err := createGroup(ctx, client, base, "Home")
	if err != nil {
		log.Fatalf("create group: %v", err)
	}

	listID, err := createList(ctx, client, base, groupID, "Weekly groceries")
	if err != nil {
		log.Fatalf("create list: %v", err)
	}

	for _, productID := range productIDs {
		if err := addItem(ctx, client, base, listID, productID); err != nil {
			log.Fatalf("add item: %v", err)
		}
	}

	// Register two more users and invite both WHILE still logged in as dev (the
	// owner). friend will accept; pending will not, leaving a live pending invite.
	if _, err := registerUser(ctx, client, base, "Fri", "End", "friend", "pass"); err != nil {
		log.Fatalf("register friend: %v", err)
	}
	if _, err := registerUser(ctx, client, base, "Pen", "Ding", "pending", "pass"); err != nil {
		log.Fatalf("register pending: %v", err)
	}

	if err := inviteByUsername(ctx, client, base, groupID, "friend"); err != nil {
		log.Fatalf("invite friend: %v", err)
	}
	if err := inviteByUsername(ctx, client, base, groupID, "pending"); err != nil {
		log.Fatalf("invite pending: %v", err)
	}

	// Logging in as friend swaps the session in the jar: from here the client acts
	// as friend, who accepts. pending never logs in, so their invite stays pending.
	if err := login(ctx, client, base, "friend", "pass"); err != nil {
		log.Fatalf("login friend: %v", err)
	}

	if err := acceptInvite(ctx, client, base, groupID); err != nil {
		log.Fatalf("accept invite: %v", err)
	}

	log.Println("seed complete")
}
