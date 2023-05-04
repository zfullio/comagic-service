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
	"sync"
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
	if err != nil {
		log.Fatalf("Ошибка планировщика %v", err)
	}

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
	dateFrom := now.AddDate(0, 0, -60).Format("2006-01-02")
	dateTill := now.Format("2006-01-02")

	wg := &sync.WaitGroup{}
	wg.Add(len(cfg.Reports))

	for _, report := range cfg.Reports {
		report := report
		go getReport(ctx, wg, c, &report, dateFrom, dateTill)
	}

	wg.Wait()
}

func getReport(ctx context.Context, wg *sync.WaitGroup, c pb.ComagicServiceClient, report *config.Report, dateFrom string, dateTill string) {
	defer wg.Done()
	log.Printf("%s // Сбор отчета\n", report.ObjectName)

	if report.CallsTable != "" {
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
			log.Println(fmt.Errorf("%s // %w", report.ObjectName, err))
		} else {
			log.Printf("%s // Статус отчета по звонкам: %v ", report.ObjectName, callsReq.IsOK)
		}
	}

	if report.OfflineMessageTable != "" {
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
			log.Println(fmt.Errorf("%s // %w", report.ObjectName, err))
		} else {
			log.Printf("%s // Статус отчета по заявкам: %v ", report.ObjectName, messagesReq.IsOK)
		}
	}
}

func GracefulShutdown() {
	if err := recover(); err != nil {
		fmt.Println("Критическая ошибка:", err)
	}

	os.Exit(0)
}
