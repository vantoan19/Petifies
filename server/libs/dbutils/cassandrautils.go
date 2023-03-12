package dbutils

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var cassandraLogger = logging.New("Libs.DBUtils.Cassandra")

func ConnectToCassandraDB(dbUrl, dbUser, dbPassword, keyspace string) (*gocql.Session, error) {
	attempt := 0
	cassandraLogger.InfoData("Start ConnectToCassandraDB", logging.Data{"dbUrl": dbUrl})

	for {
		client, err := openCassandraDB(dbUrl, dbUser, dbPassword, keyspace)
		if err != nil {
			attempt++
			cassandraLogger.WarningData("Executing ConnectToCassandraDB: Connect to the database fails, attempt again...", logging.Data{"attempt": attempt})
		} else {
			cassandraLogger.InfoData("Finished ConnectToCassandraDB: SUCCESSFUL", logging.Data{"dbUrl": dbUrl})
			return client, nil
		}

		if attempt > 10 {
			cassandraLogger.ErrorData("Finished ConnectToCassandraDB: FAILED", logging.Data{"dbUrl": dbUrl, "error": err.Error()})
			return nil, err
		}

		cassandraLogger.Info("Executing ConnectToCassandraDB: Wait for 2 seconds before retrying to connect to the database")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openCassandraDB(dbUrl, dbUser, dbPassword, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(dbUrl)
	cluster.ConnectTimeout = time.Second * 10
	cluster.Keyspace = keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: dbUser,
		Password: dbPassword,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}
