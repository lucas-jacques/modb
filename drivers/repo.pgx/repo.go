package repo

import (
	"github.com/lucasjacques/modb"
	"github.com/lucasjacques/modb/model"
)

func wrapDBTX(conn pgxdbtx) modb.DBTX {
	return &dbtx{
		conn: conn,
	}
}

func Repo[M any, PK comparable, C model.Schema[M]](dbtx pgxdbtx, model model.TypedModelCols[M, PK, C]) *modb.ModelRepository[M, PK, C] {
	return modb.NewRepository(wrapDBTX(dbtx), model)
}
