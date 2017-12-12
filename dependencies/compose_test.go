package dependencies

import (
	"errors"
	"testing"
)

var updateOptions = map[string]string{
	"test": "value",
}

func TestBasicVersion(t *testing.T) {
	// TODO setup an actual DB
	result, err := BasicVersion("test_hostname", "query", "value", updateOptions)
	if err != nil {
		t.Errorf("result == %v expected %v", err, nil)
	}
	if result != "success" {
		t.Errorf("result == %v expected %v", result, "success")
	}
}

func TestExtractedFunctionVersion(t *testing.T) {
	// TODO setup an actual DB
	result, err := ExtractedFunctionVersion("test_hostname", "query", "value", updateOptions)
	if err != nil {
		t.Errorf("result == %v expected %v", err, nil)
	}
	if result != "success" {
		t.Errorf("result == %v expected %v", result, "success")
	}
}

type clientMock struct{}

func (client clientMock) DbUpdate(options map[string]string, field string, value interface{}) error {
	// TODO store something to do verifcation
	return nil
}

func (client clientMock) DbQuery(query string) (result []string, err error) {
	// TODO store something to do verifcation
	return []string{}, nil
}

func TestInterfaceClientVersion(t *testing.T) {
	// TODO mock out values, so we get some useful response
	// probably using a mocking library
	client := clientMock{}
	result, err := interfaceClientVersion(client, "query", "value", updateOptions)
	if err != nil {
		t.Errorf("result == %v expected %v", err, nil)
	}
	if result != "success" {
		t.Errorf("result == %v expected %v", result, "success")
	}
	// TODO verify mock functions were called correctly
}

func TestFunctionDependencyVersion(t *testing.T) {
	// very straightforward test setup with no special libaries required
	// no mocking required, as long as some instance of client can be created
	expectedClient := dbclient{}
	expectedUpdates := []string{"a", "b"}
	client := func() (dbclient, error) { return expectedClient, nil }
	query := func(client dbclient, id string) ([]string, error) { return expectedUpdates, nil }
	updateCount := 0
	entityID := "123"
	update := func(client dbclient, value string) error {
		if value != expectedUpdates[updateCount] {
			t.Errorf("result == %v expected %v", value, expectedUpdates[updateCount])
		}
		updateCount++
		return nil
	}
	result, err := functionDependenciesVersion(client, query, update)(entityID)
	if err != nil {
		t.Errorf("result == %v expected %v", err, nil)
	}
	if updateCount != 2 {
		t.Errorf("result == %v expected %v", updateCount, 2)
	}
	if result != expectedUpdates[1] {
		t.Errorf("result == %v expected %v", result, expectedUpdates[1])
	}
	// TODO verify mock was called correctly
}

func TestClientDep(t *testing.T) {
	expectedClient := dbclient{}
	expectedConfig := "url:host"
	expectedErr := errors.New("some error")
	dbhost := "host"
	newClient := func(config string) (dbclient, error) {
		if config != expectedConfig {
			t.Errorf("result == %v expected %v", config, expectedConfig)
		}
		return expectedClient, expectedErr
	}
	client, err := clientDep(newClient, dbhost)()
	if client != expectedClient {
		t.Errorf("result == %v expected %v", client, expectedClient)
	}
	if err != expectedErr {
		t.Errorf("result == %v expected %v", err, expectedErr)
	}
}

func TestMemoizedVersion(t *testing.T) {
	// the test is a bit simpler as nothing about the client is required in the test
	expectedUpdates := []string{"a", "b"}
	query := func() ([]string, error) { return expectedUpdates, nil }
	updateCount := 0
	update := func(value string) error {
		if value != expectedUpdates[updateCount] {
			t.Errorf("result == %v expected %v", value, expectedUpdates[updateCount])
		}
		updateCount++
		return nil
	}
	result, err := memoizedVersion(query, update)
	if err != nil {
		t.Errorf("result == %v expected %v", err, nil)
	}
	if updateCount != 2 {
		t.Errorf("result == %v expected %v", updateCount, 2)
	}
	if result != expectedUpdates[1] {
		t.Errorf("result == %v expected %v", result, expectedUpdates[1])
	}
	// TODO verify mock was called correctly
}
