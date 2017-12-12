package dependencies

// Composition, Testing and Dependencies
// With a lot of dynamic languages, it is easy to replace a given function or method call with a mocked version
// By changing the behaviour, you make testing existing code much easier.
// Another way to make testing easier, especially with static languages is with dependency injection.
// A different approach is to use function composition which has different trade offs to dependency injection.
// This is very much inspired by https://medium.com/javascript-scene/mocking-is-a-code-smell-944a70c90a6a

// The assumption I'm making is that there's no active record style library available
// The way to implement and write tests for the following functionality would be completely
// different based on the use of that library.

// Part of the goal is to maintain the public method's signature
// I don't want the caller to need to know about the function's dependencies
// or be affected by any refactoring.

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
func BasicVersion(
	dbhost, entityID, delta string,
	updateOptions map[string]string,
) (lastUpdated string, err error) {
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
func ExtractedFunctionVersion(
	dbhost, entityID, delta string,
	updateOptions map[string]string,
) (lastUpdated string, err error) {
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

func InterfaceClientVersion(
	dbhost, entityID, delta string,
	updateOptions map[string]string,
) (lastUpdated string, err error) {
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

func FunctionDependenciesVersion(
	dbhost, entityID, delta string,
	updateOptions map[string]string,
) (lastUpdated string, err error) {
	return functionDependenciesVersion(
		clientDep(NewDbClient, dbhost),
		queryDep,
		updateDep(delta, updateOptions),
	)(entityID)
}

// adapters, change the interface from the existing library to a more useful interface for our specific use case
// extracting the dependencies to functions allows them to be tested separately
// This can be tested using mocks and pure functions to confirm it behaves correctly

func clientDep(newClient func(config string) (dbclient, error), dbhost string) func() (dbclient, error) {
	return func() (dbclient, error) {
		return newClient(dbconnectionstring(dbhost))
	}
}

func queryDep(client dbclient, entityID string) ([]string, error) {
	return client.DbQuery(makeQuery(entityID))
}

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

// Level 1 of this would be to pass in the function directly. To refactor existing code, this is probably the first step.
// Level 2 is to adapt the dependency function to our specific use case. By creating adapter functions, the resulting code only includes the required details, without irrelevant details.

func functionDependenciesVersion(
	dbclient func() (dbclient, error),
	query func(dbclient, string) ([]string, error),
	update func(dbclient, string) error,
) func(entityId string) (lastUpdated string, err error) {
	return func(entityId string) (lastUpdated string, err error) {
		client, err := dbclient()
		if err != nil {
			return "", err
		}
		results, err := query(client, entityId)
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
}

// to make sure the query and updates use the same connection, the client isn't abstracted in the previous version

// when memoization is used for the client function, it is pulled out into the composition
func MemoizedVersion(
	dbhost, entityID, delta string,
	updateOptions map[string]string,
) (lastUpdated string, err error) {
	clientFunc := memoize(clientDep(NewDbClient, dbhost))
	return memoizedVersion(
		queryDepWithClient(clientFunc, entityID),
		updateDepWithClient(clientFunc, delta, updateOptions),
	)
}

func memoize(makeClient func() (dbclient, error)) func() (dbclient, error) {
	var client dbclient
	var err error
	exists := false
	return func() (dbclient, error) {
		if exists {
			return client, err
		}
		client, err = makeClient()
		return client, err
	}
}

// This can be tested using mocks and pure functions to confirm it behaves correctly
func queryDepWithClient(makeClient func() (dbclient, error), entityID string) func() ([]string, error) {
	return func() ([]string, error) {
		client, err := makeClient()
		if err != nil {
			return nil, err
		}
		return client.DbQuery(makeQuery(entityID))
	}
}

func updateDepWithClient(makeClient func() (dbclient, error), delta string, updateOptions map[string]string) func(result string) error {
	return func(field string) error {
		client, err := makeClient()
		if err != nil {
			return err
		}
		return client.DbUpdate(updateOptions, field, delta)
	}
}

// the client is now pulled out of this function
// the only issue with this is if makeClient returns an error, it will now be handled by the query error handling
func memoizedVersion(
	query func() ([]string, error),
	update func(string) error,
) (lastUpdated string, err error) {
	results, err := query()
	if err != nil {
		return "", err
	}
	for _, field := range results {
		err = update(field)
		if err != nil {
			return lastUpdated, err
		}
		lastUpdated = field
	}
	return lastUpdated, nil
}
