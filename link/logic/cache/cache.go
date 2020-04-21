package cache

import "github.com/aijie/michat/datas/model"

var (
	AppCache appCache
)

type appCache struct {
}

func (a *appCache)Get(appId int64) (*model.App, error) {

}

