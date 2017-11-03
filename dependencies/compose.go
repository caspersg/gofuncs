package dependencies

// Testing, Dependencies and Composition
// With a lot of dynamic languages, it is easy to replace a given function or method call with a mocked version
// By changing the behaviour, you make testing existing code much easier.
// This is very much inspired by https://medium.com/javascript-scene/mocking-is-a-code-smell-944a70c90a6a
// Another way to make testing easier, especially with static languages is with dependency injection.

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

func newDbClient(config string) (dbclient, error) {
	// IO side effect, creating new connection to DB
	return dbclient{}, nil
}

// code under source control
// There's always some amount of processing required to take app config and to use it with a given library.
func dbconnectionstring(host string) string {
	return "url:" + host
}

// Basic straight forward code, using db libraries directly.
// there's no separation of side effects from the logic
// no dependency is being injected, so nothing can be mocked
// we basically need a test instance of the DB in order to test this
// The tests are slower and more prone to fragility,
// so they'll often end up only testing the happy path.
func BasicVersion(dbhost, query, updatevalue string, updateOptions map[string]string) (last string, err error) {
	client, err := newDbClient(dbconnectionstring(dbhost))
	if err != nil {
		return "", err
	}
	results, err := client.DbQuery(query)
	if err != nil {
		return "", err
	}
	for _, result := range results {
		err = client.DbUpdate(updateOptions, result, updatevalue)
		if err != nil {
			return "", err
		}
		last = result
	}
	return last, nil
}

// the first refactoring is usually to extract functions
// it may be a bit easier to read,
// but this has not improved our ease of testing at all
// If the error handling was more extensive, then that could be partially abstracted in the sub functions
func ExtractedFunctionVersion(dbhost, query, updatevalue string, updateOptions map[string]string) (last string, err error) {
	client, err := makeClient(dbhost)
	if err != nil {
		return "", err
	}
	results, err := client.DbQuery(query)
	if err != nil {
		return "", err
	}
	return updateAndGetLast(client, results, updatevalue, updateOptions)
}

func makeClient(dbhost string) (dbclient, error) {
	return newDbClient(dbconnectionstring(dbhost))
}

func updateAndGetLast(client dbclient, results []string, updatevalue string, updateOptions map[string]string) (string, error) {
	last := ""
	for _, result := range results {
		err := client.DbUpdate(updateOptions, result, updatevalue)
		if err != nil {
			return "", err
		}
		last = result
	}
	return last, nil
}

// Basic interface for dependencies being injected.
// to test this, we'd need to pass in a mock version of the client
// but this then locks down the implementation to an exact query/update
// which means we still need an integration or functional test to confirm that it actually works correctly.
// For example, should the test fail if an additional value in options is passed in?
type iClient interface {
	DbUpdate(options map[string]string, field string, value interface{}) error
	DbQuery(query string) (result []string, err error)
}

func InterfaceClient(dbhost, query, updatevalue string, updateOptions map[string]string) (last string, err error) {
	client, err := newDbClient(dbconnectionstring(dbhost))
	if err != nil {
		return "", err
	}
	return interfaceClient(client, query, updatevalue, updateOptions)
}

func interfaceClient(client iClient, query, updatevalue string, updateOptions map[string]string) (last string, err error) {
	results, err := client.DbQuery(query)
	if err != nil {
		return "", err
	}
	for _, result := range results {
		err = client.DbUpdate(updateOptions, result, updatevalue)
		if err != nil {
			return "", err
		}
		last = result
	}
	return last, nil
}

// this function needs to be integration tested to make sure it's doing the right thing
// a unit test wouldn't be able to confirm that our query or update is actually correct.
// There are multiple possible queries/updates that would acheive the same result so only an integration test is sufficient.
// This is because the dependencies being injected (which is the entire contents of this function), produce side effects.
// If Go had generics, then this composition could be written more generically.
// but it would still need integration testing, so not a lot is lost
func FunctionDependencies(dbhost, query, updatevalue string, updateOptions map[string]string) (last string, err error) {
	return functionDependencies(
		func() (dbclient, error) {
			return newDbClient(dbconnectionstring(dbhost))
		},
		func(client dbclient) ([]string, error) {
			return client.DbQuery(query)
		},
		updateDep(updatevalue, updateOptions),
	)
}

// If the dependency is more complicated then it can be extracted to a separate function like this
// This can be tested using mocks and pure functions to confirm it behaves correctly
func updateDep(updatevalue string, updateOptions map[string]string) func(client dbclient, result string) error {
	return func(client dbclient, result string) error {
		return client.DbUpdate(updateOptions, result, updatevalue)
	}
}

// testing this function is as simple as passing in pure functions or closures
// instead of making the test fit whatever interface your library api implementation has,
// you create the perfect interface for your purpose, and make your dependencies match that
// this can be perfectly unit tested, because the logic can be isolated from the side effects
// the specific details of the library api are abstracted from the logic as well.
// I think this function is far more readable now.
func functionDependencies(
	dbclient func() (dbclient, error),
	query func(dbclient) ([]string, error),
	update func(dbclient, string) error,
) (last string, err error) {
	client, err := dbclient()
	if err != nil {
		return "", err
	}
	results, err := query(client)
	if err != nil {
		return "", err
	}
	for _, result := range results {
		err = update(client, result)
		if err != nil {
			return last, err
		}
		last = result
	}
	return last, nil
}
