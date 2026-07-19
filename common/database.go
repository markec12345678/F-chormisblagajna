package common

import (
	"fmt"
	"sync"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var singleDBInstance *mongo.Client
var lock = &sync.Mutex{}

func GetDatabaseClient(logger logger.ILogger, conf *config.Config) (*mongo.Client, error) {

	if singleDBInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleDBInstance == nil {
			logger.Info("Creating DB single instance now.")

			uri := fmt.Sprintf("mongodb://%s:%v", conf.Databases[0].Host, conf.Databases[0].Port)
			if conf.Databases[0].Username != "" {
				uri = fmt.Sprintf("mongodb://%s:%s@%s:%v",
					conf.Databases[0].Username,
					conf.Databases[0].Password,
					conf.Databases[0].Host,
					conf.Databases[0].Port,
				)
			}

			opts := options.Client().ApplyURI(uri)

			client, err := mongo.Connect(opts)
			if err != nil {
				return nil, err
			}
			singleDBInstance = client
		}
	}

	return singleDBInstance, nil
}
