//go:generate wire
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/asynccnu/be-calendar/cron"
	"github.com/asynccnu/be-calendar/grpc"
	"github.com/asynccnu/be-calendar/ioc"
	"github.com/asynccnu/be-calendar/repository/cache"
	"github.com/asynccnu/be-calendar/repository/dao"
	"github.com/asynccnu/be-calendar/service"
	"github.com/google/wire"
)

func InitApp() App {
	wire.Build(
		// 第三方
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitFeedClient,
		ioc.InitQiniu,
		ioc.InitGRPCxKratosServer,
		grpc.NewCalendarServiceServer,
		service.NewCachedCalendarService,
		cache.NewRedisCalendarCache,
		dao.NewMysqlCalendarDAO,
		cron.NewHolidayController,
		cron.NewCalendarController,
		cron.NewCron,
		NewApp,
	)
	return App{}
}
