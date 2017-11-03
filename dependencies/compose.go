package dependencies

// provided api
type dbclient struct{}

func dbupdate(client dbclient, field string, value interface{}) error {
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

func BusinessLogicV1(dbhost, query, updatevalue string) (last string, err error) {
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
		err = dbupdate(client, result, updatevalue)
		if err != nil {
			return "", err
		}
		last = result
	}
	return last, nil
}
