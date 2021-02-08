package tasks

import (
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v2"

	redisbackend "github.com/RichardKnop/machinery/v1/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v1/brokers/redis"
	eagerlock "github.com/RichardKnop/machinery/v1/locks/eager"
)

func StartServer() (*machinery.Server, error) {
	cnf := &config.Config{
		DefaultQueue:    "machinery_tasks",
		ResultsExpireIn: 3600,
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}

	// Create server instance
	broker := redisbroker.NewGR(cnf, []string{"localhost:6379"}, 0)
	backend := redisbackend.NewGR(cnf, []string{"localhost:6379"}, 0)
	lock := eagerlock.New()
	server := machinery.NewServer(cnf, broker, backend, lock)

	tasks := map[string]interface{}{
		"add":               exampletasks.Add,
		"multiply":          exampletasks.Multiply,
		"sum_ints":          exampletasks.SumInts,
		"sum_floats":        exampletasks.SumFloats,
		"concat":            exampletasks.Concat,
		"split":             exampletasks.Split,
		"panic_task":        exampletasks.PanicTask,
		"long_running_task": exampletasks.LongRunningTask,
	}

	return server, server.RegisterTasks(tasks)
}