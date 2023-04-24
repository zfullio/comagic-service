package bq

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"strings"
	"time"
)

func CreateTable(ctx context.Context, schemaDTO any, dataset *bigquery.Dataset, tableID string) (err error) {
	schema, err := bigquery.InferSchema(schemaDTO)
	if err != nil {
		return fmt.Errorf("bigquery.InferSchema: %w", err)
	}

	err = dataset.Table(tableID).Create(ctx, &bigquery.TableMetadata{
		Schema: schema})
	if err != nil {
		fmt.Printf("bigquery.Create: %s", err)
	}

	return nil
}

func DeleteByDateColumn(ctx context.Context, bqClient bigquery.Client, datasetID string, tableID string, dateColumn string, dateStart time.Time, dateFinish time.Time) (err error) {
	q := bqClient.Query(`DELETE FROM ` + fmt.Sprintf(" %s.%s ", datasetID, tableID) +
		`WHERE ` + fmt.Sprintf("%s >= '%s' AND %s <= '%s'", dateColumn, dateStart.Format(time.DateOnly),
		dateColumn, dateFinish.Format(time.DateOnly)))

	job, err := q.Run(ctx)
	if err != nil {
		return err
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}

	if err := status.Err(); err != nil {
		return err
	}

	return err
}

func SendFromCS(ctx context.Context, schemaDTO any, table *bigquery.Table, bucket string, object string) (err error) {
	schema, err := bigquery.InferSchema(schemaDTO)
	if err != nil {
		return fmt.Errorf("bigquery.InferSchema: %w", err)
	}

	filePath := strings.Split(object, "/")
	gcsRef := bigquery.NewGCSReference(fmt.Sprintf("gs://%s/%s", bucket, filePath[len(filePath)-1]))

	gcsRef.SourceFormat = bigquery.CSV
	gcsRef.FieldDelimiter = "|"
	gcsRef.SkipLeadingRows = 1
	gcsRef.AllowJaggedRows = true
	gcsRef.AllowQuotedNewlines = true

	gcsRef.Schema = schema
	loader := table.LoaderFrom(gcsRef)
	loader.CreateDisposition = bigquery.CreateNever
	loader.WriteDisposition = bigquery.WriteAppend

	_, err = loader.Run(ctx)
	if err != nil {
		return fmt.Errorf("loader error: %w", err)
	}

	return err
}
