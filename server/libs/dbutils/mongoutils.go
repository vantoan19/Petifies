package dbutils

import (
	"context"
	"fmt"
	"time"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var mongoLogger = logging.New("Libs.DBUtils.Mongo")

var (
	wc        = writeconcern.New(writeconcern.WMajority())
	rc        = readconcern.Snapshot()
	transOpts = options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
)

func ConnectToMongoDB(dbUrl string) (*mongo.Client, error) {
	attempt := 0
	mongoLogger.InfoData("Start ConnectToMongoDB", logging.Data{"dbUrl": dbUrl})

	for {
		client, err := openMongoDB(dbUrl)
		if err != nil {
			attempt++
			mongoLogger.WarningData("Executing ConnectToMongoDB: Connect to the database fails, attempt again...", logging.Data{"attempt": attempt})
		} else {
			mongoLogger.InfoData("Finished ConnectToMongoDB: SUCCESSFUL", logging.Data{"dbUrl": dbUrl})
			return client, nil
		}

		if attempt > 10 {
			mongoLogger.ErrorData("Finished ConnectToMongoDB: FAILED", logging.Data{"dbUrl": dbUrl, "error": err.Error()})
			return nil, err
		}

		mongoLogger.Info("Executing ConnectToMongoDB: Wait for 2 seconds before retrying to connect to the database")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openMongoDB(dbUrl string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI((options.ServerAPIVersion1))
	opts := options.Client().ApplyURI(dbUrl).SetServerAPIOptions(serverAPI).SetTimeout(time.Second * 2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

type MongoOperation func(ssCtx mongo.SessionContext) error

func ExecWithSession(ctx context.Context, client *mongo.Client, operations ...MongoOperation) error {
	session, err := client.StartSession()
	defer session.EndSession(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	err = session.StartTransaction(transOpts)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	sessionContext := mongo.NewSessionContext(ctx, session)
	for _, ope := range operations {
		if err = ope(sessionContext); err != nil {
			if abErr := session.AbortTransaction(ctx); abErr != nil {
				return status.Errorf(codes.Internal, fmt.Sprintf("session err: %v, abort err: %v", err, abErr))
			}
			if status.FromContextError(err).Code() == codes.NotFound {
				return err
			}
			return status.Errorf(codes.Internal, err.Error())
		}
	}

	return session.CommitTransaction(ctx)
}
