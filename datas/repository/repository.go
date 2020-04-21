package repository

import (
	"github.com/aijie/michat/datas/model"
)

type IMRepository interface {
	SetAppInfo(app model.App) error
	GetAppInfo(appId int64) (*model.App, error)
}
