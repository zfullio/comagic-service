package main

import (
	"Comagic/internal/config"
	"Comagic/pb"
	"context"
	"flag"
	"fmt"
	"github.com/go-co-op/gocron"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"
)

func main() {
	var fileConfig = flag.String("f", "schedule_config.yml", "configuration file")
	flag.Parse()

	cfg, err := config.NewScheduleConfig(*fileConfig)
	if err != nil {
		log.Fatalf("could not read config: %v", err)
	}

	fmt.Printf("Количество отчетов: %d\n", len(cfg.Reports))
	fmt.Printf("Запуск ежедневно в: %s\n", cfg.Time)

	s := gocron.NewScheduler(time.UTC)
	s.WaitForScheduleAll()
	location, err := time.LoadLocation("Local")
	s.ChangeLocation(location)

	_, err = s.Every(1).Day().At(cfg.Time).Do(scheduleRun, *cfg)
	if err != nil {
		log.Fatalf("Ошибка планировщика %v", err)
	}

	s.StartBlocking()

	recover()
}

func scheduleRun(cfg config.ScheduleConfig) {
	defer GracefulShutdown()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", cfg.GRPC.IP, cfg.GRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	c := pb.NewComagicServiceClient(conn)

	ctx := context.Background()
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -89).Format("2006-01-02")
	dateTill := now.Format("2006-01-02")

	for _, report := range cfg.Reports {
		log.Printf("Сбор отчета для: %s\n", report.ObjectName)
		callsReq, err := c.PushCallsToBQ(ctx, &pb.PushCallsToBQRequest{
			ComagicToken: report.ComagicToken,
			BqConfig: &pb.BqConfig{
				ProjectID:  report.ProjectID,
				DatasetID:  report.DatasetID,
				TableID:    report.CallsTable,
				ServiceKey: report.GoogleServiceKey,
			},
			CsConfig: &pb.CsConfig{
				ServiceKey: report.GoogleServiceKey,
				BucketName: report.BucketName,
			},
			DateFrom: dateFrom,
			DateTill: dateTill,
		})
		if err != nil {
			log.Println(err)
		}
		log.Printf("Статус отчета по звонкам: %v ", callsReq.IsOK)

		messagesReq, err := c.PushOfflineMessagesToBQ(ctx, &pb.PushOfflineMessagesToBQRequest{
			ComagicToken: report.ComagicToken,
			BqConfig: &pb.BqConfig{
				ProjectID:  report.ProjectID,
				DatasetID:  report.DatasetID,
				TableID:    report.OfflineMessageTable,
				ServiceKey: report.GoogleServiceKey,
			},
			CsConfig: &pb.CsConfig{
				ServiceKey: report.GoogleServiceKey,
				BucketName: report.BucketName,
			},
			DateFrom: dateFrom,
			DateTill: dateTill,
		})
		if err != nil {
			log.Println(err)
		}
		log.Printf("Статус отчета по заявкам: %v ", messagesReq.IsOK)

	}
}

func GracefulShutdown() {
	if err := recover(); err != nil {
		fmt.Println("Критическая ошибка:", err)
	}
	os.Exit(0)
}
