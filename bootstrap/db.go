package bootstrap

import (
	"os"

	"github.com/facebookgo/inject"
	_ "github.com/lib/pq"

	"golang.ysitd.cloud/db"
)

func initDB() *db.GeneralOpener {
	return db.NewOpener("postges", os.Getenv("DB_URL"))
}

func InjectDB(graph *inject.Graph) {
	graph.Provide(
		&inject.Object{Value: initDB()},
	)
}
