package repo

import (
	"github.com/lucasjacques/modb"
	"github.com/lucasjacques/modb/model"
)

func New[M any, PK comparable, C any](dbtx stddbtx, model model.TypedModelCols[M, PK, C]) *modb.ModelRepository[M, PK, C] {
	return modb.NewRepository(wrapDBTX(dbtx), model)
}
