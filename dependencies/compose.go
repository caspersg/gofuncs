package dependencies

// Testing, Dependencies and Composition
// With a lot of dynamic languages, it is easy to replace a given function or method call with a mocked version
// By changing the behaviour, you make testing existing code much easier.
// Another way to make testing easier, especially with static languages is with dependency injection.
// A different approach is to use function composition which has different trade offs to dependency injection.
// This is very much inspired by https://medium.com/javascript-scene/mocking-is-a-code-smell-944a70c90a6a

// The assumption I'm making is that there's no active record style library available
// The way to implement and write tests for the following functionality would be completely
// different based on the use of that library.

//
// provided library api
type dbclient struct{}

func (client dbclient) DbUpdate(options map[string]string, field string, value interface{}) error {
	// IO side effect, updating the DB
	return nil
}

func (client dbclient) DbQuery(query string) (result []string, err error) {
	// IO side effect, reading data from DB
	return []string{}, nil
}

func NewDbClient(config string) (dbclient, error) {
	// IO side effect, creating new connection to DB
	return dbclient{}, nil
}

// code under source control

//

// Basic straight forward code, using db libraries directly.
// there's no separation of side effects from the logic
// no dependency is being injected, so nothing can be mocked
// we basically need a test instance of the DB in order to test this
// As a result tests are slower and more prone to fragility,
// so they'll often end up only testing the happy path.
func BasicVersion(dbhost, entityID, delta string, updateOptions map[string]string) (lastUpdated string, err error) {
	connectionString := "url:" + dbhost
	client, err := NewDbClient(connectionString)
	if err != nil {
		return "", err
	}
	query := "id=" + entityID
	results, err := client.DbQuery(query)
	if err != nil {
		return "", err
	}
	for _, field := range results {
		err = client.DbUpdate(updateOptions, field, delta)
		if err != nil {
			return "", err
		}
		lastUpdated = field
	}
	return lastUpdated, nil
}

//

// One approach to refactoring is to extract functions
// it may be a bit easier to read,
// but this has not improved our ease of testing at all
// If the error handling was more extensive, then that could also be partially abstracted in the sub functions
func ExtractedFunctionVersion(dbhost, entityID, delta string, updateOptions map[string]string) (lastUpdated string, err error) {
	client, err := makeClient(dbhost)
	if err != nil {
		return "", err
	}
	results, err := query(client, entityID)
	if err != nil {
		return "", err
	}
	for _, field := range results {
		err := update(client, field, delta, updateOptions)
		if err != nil {
			return "", err
		}
		lastUpdated = field
	}
	return lastUpdated, nil
}

func dbconnectionstring(host string) string {
	return "url:" + host
}

func makeClient(dbhost string) (dbclient, error) {
	return NewDbClient(dbconnectionstring(dbhost))
}

func makeQuery(entityID string) string {
	return "id=" + entityID
}

func query(client dbclient, entityID string) ([]string, error) {
	return client.DbQuery(makeQuery(entityID))
}

func update(client dbclient, field, delta string, updateOptions map[string]string) error {
	return client.DbUpdate(updateOptions, field, delta)
}

//

// Basic interface for dependencies being injected.
// to test this, we'd need to pass in a mock version of the client
// but this then locks down the implementation to a specific query and update implementation,
// which means we still need an integration or functional test to confirm that it actually works correctly.
type iClient interface {
	DbUpdate(options map[string]string, field string, value interface{}) error
	DbQuery(query string) (result []string, err error)
}

func InterfaceClientVersion(dbhost, entityID, delta string, updateOptions map[string]string) (lastUpdated string, err error) {
	client, err := NewDbClient(dbconnectionstring(dbhost))
	if err != nil {
		return "", err
	}
	return interfaceClientVersion(client, entityID, delta, updateOptions)
}

func interfaceClientVersion(client iClient, entityID, delta string, updateOptions map[string]string) (lastUpdated string, err error) {
	results, err := client.DbQuery(makeQuery(entityID))
	if err != nil {
		return "", err
	}
	for _, field := range results {
		err = client.DbUpdate(updateOptions, field, delta)
		if err != nil {
			return "", err
		}
		lastUpdated = field
	}
	return lastUpdated, nil
}

// this function needs to be integration tested to make sure it's doing the right thing
// a unit test wouldn't be able to confirm that our query or update is actually correct.
// This is because the dependencies being injected (which is the entire contents of this function), produce side effects.
// There are multiple possible queries/updates that would acheive the same result so only an integration test is sufficient.
// For example, should the test fail if an additional value in options is passed in?

// If Go had generics, then this composition could be written more generically.
// but it would still need integration testing, so not a lot is lost
func FunctionDependenciesVersion(dbhost, entityID, delta string, updateOptions map[string]string) (lastUpdated string, err error) {
	return functionDependenciesVersion(
		func() (dbclient, error) {
			return NewDbClient(dbconnectionstring(dbhost))
		},
		queryDep(entityID),
		updateDep(delta, updateOptions),
	)
}

func queryDep(entityID string) func(client dbclient) ([]string, error) {
	return func(client dbclient) ([]string, error) {
		return client.DbQuery(makeQuery(entityID))
	}
}

// If the dependency is more complicated then it can be extracted to a separate function like this
// This can be tested using mocks and pure functions to confirm it behaves correctly
func updateDep(delta string, updateOptions map[string]string) func(client dbclient, result string) error {
	return func(client dbclient, field string) error {
		return client.DbUpdate(updateOptions, field, delta)
	}
}

// testing this function is as simple as passing in pure functions or closures
// instead of making the test fit whatever interface your library api implementation has,
// you create the perfect interface for your purpose, and make your dependencies match that
// this can be perfectly unit tested, because the logic can be isolated from the side effects
// the specific details of the library api are abstracted from the logic as well.
// I think this function is far more readable now.
func functionDependenciesVersion(
	dbclient func() (dbclient, error),
	query func(dbclient) ([]string, error),
	update func(dbclient, string) error,
) (lastUpdated string, err error) {
	client, err := dbclient()
	if err != nil {
		return "", err
	}
	results, err := query(client)
	if err != nil {
		return "", err
	}
	for _, field := range results {
		err = update(client, field)
		if err != nil {
			return lastUpdated, err
		}
		lastUpdated = field
	}
	return lastUpdated, nil
}
