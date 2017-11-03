package dependencies

// provided api
type dbclient struct{}

func dbupdate(client dbclient, options, field string, value interface{}) error {
	return nil
}

func dbquery(client dbclient, query string) (result []string, err error) {
	return []string{}, nil
}

func getdbclient(config string) (dbclient, error) {
	return dbclient{}, nil
}

// code under source control
func dbconnectionstring(host string) string {
	return "url:" + host
}

func BusinessLogicV1(dbhost, query, updatevalue, updateOptions string) (last string, err error) {
	connectiondetails := dbconnectionstring(dbhost)
	client, err := getdbclient(connectiondetails)
	if err != nil {
		return "", err
	}
	results, err := dbquery(client, query)
	if err != nil {
		return "", err
	}
	for _, result := range results {
		err = dbupdate(client, updateOptions, result, updatevalue)
		if err != nil {
			return "", err
		}
		last = result
	}
	return last, nil
}

func BusinessLogicV2(dbhost, query, updatevalue, updateOptions string) (last string, err error) {
	client, err := makeClient(dbhost)
	if err != nil {
		return "", err
	}
	return businessLogicV2(
		func() ([]string, error) {
			return dbquery(client, query)
		},
		func(result string) error {
			return dbupdate(client, updateOptions, result, updatevalue)
		},
	)
}

func makeClient(dbhost string) (dbclient, error) {
	connectiondetails := dbconnectionstring(dbhost)
	client, err := getdbclient(connectiondetails)
	if err != nil {
		return dbclient{}, err
	}
	return client, nil
}

func businessLogicV2(
	query func() ([]string, error),
	update func(string) error,
) (last string, err error) {
	results, err := query()
	if err != nil {
		return "", err
	}
	for _, result := range results {
		err = update(result)
		if err != nil {
			return "", err
		}
		last = result
	}
	return last, nil
}
