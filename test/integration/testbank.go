package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hthunberg/course-golang-postgres-grpc-api/dbtest"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestBank struct {
	container tc.Container
	URI       string
}

func setupTestBank(ctx context.Context, testBankContainerRequest tc.GenericContainerRequest) (*TestBank, error) {
	container, err := tc.GenericContainer(ctx, testBankContainerRequest)
	if err != nil {
		return nil, fmt.Errorf("setup test bank:create container: %v", err)
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("setup test bank:container host: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, fmt.Errorf("setup test bank:container port: %v", err)
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	// Follow application logs
	lc := newTestLogConsumer([]string{}, make(chan bool))
	container.FollowOutput(&lc)

	_ = container.StartLogProducer(ctx)

	return &TestBank{container: container, URI: uri}, nil
}

// TearDown tears down the running bank container
func (tdb *TestBank) TearDown() {
	_ = tdb.container.Terminate(context.Background())
}

func TestBankContainerRequest(dbAddr string) tc.GenericContainerRequest {
	dbSource := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbtest.DBUser, dbtest.DBPass, dbAddr, dbtest.DBName)

	// TODO: Hans fixa
	// file:./<path> to be relative to working directory
	// hostPathMigrationsURL := os.Getenv("MIGRATION_URL")
	hostPathMigrationsURL := "/Users/hansthunberg/git-views/golang/course-golang-postgres-grpc-api/build/db/migrations"

	env := map[string]string{
		"ENVIRONMENT":   "integrationtest",
		"DB_SOURCE":     dbSource,
		"MIGRATION_URL": "file://migrations",
		"LOG_LEVEL":     "DEBUG",
	}
	port := "8080/tcp"

	req := tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Name:         "testbank",
			Image:        "bank:latest",
			ExposedPorts: []string{port},
			Env:          env,
			Mounts: tc.ContainerMounts{
				tc.ContainerMount{
					Source: tc.GenericBindMountSource{
						HostPath: hostPathMigrationsURL,
					},
					Target: tc.ContainerMountTarget("/app/bin/migrations"),
				},
			},
			WaitingFor: wait.ForAll(
				wait.ForLog("initializing: starting application").WithStartupTimeout(5 * time.Second),
			),
		},
		Started: true,
	}

	return req
}

type TestBankClient struct {
	httpClient http.Client
	baseURL    string
}

func newTestBankCLient(baseURL string) (*TestBankClient, error) {
	return &TestBankClient{httpClient: *http.DefaultClient, baseURL: baseURL}, nil
}

func (t *TestBankClient) createUser(reqBody io.Reader) (res *http.Response, body []byte, err error) {
	req, err := http.NewRequest(
		"POST",
		t.baseURL+"/users",
		reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("create user:new request: %v", err)
	}

	if res, err = t.httpClient.Do(req); err != nil {
		return nil, nil, fmt.Errorf("create user:do request: %v", err)
	}

	if body, err = io.ReadAll(res.Body); err != nil {
		return nil, nil, fmt.Errorf("create user:read response: %v", err)
	}

	_ = res.Body.Close()

	return res, body, nil
}

func (t *TestBankClient) createAccount(reqBody io.Reader) (res *http.Response, body []byte, err error) {
	req, err := http.NewRequest(
		"POST",
		t.baseURL+"/accounts",
		reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("create account:new request: %v", err)
	}

	if res, err = t.httpClient.Do(req); err != nil {
		return nil, nil, fmt.Errorf("create account:do request: %v", err)
	}

	if body, err = io.ReadAll(res.Body); err != nil {
		return nil, nil, fmt.Errorf("create account:read response: %v", err)
	}

	_ = res.Body.Close()

	return res, body, nil
}

type UserRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type AccountRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
}
