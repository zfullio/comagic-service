package bq

import (
	"cloud.google.com/go/bigquery"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/googleapi"
	"net/http"
	"strings"
	"time"
)

func CreateTable(ctx context.Context, schemaDTO any, table *bigquery.Table, fieldPartition *string,
	fieldClustering *[]string) error {
	schema, err := bigquery.InferSchema(schemaDTO)
	if err != nil {
		return fmt.Errorf("bigquery.InferSchema: %w", err)
	}

	metadata := &bigquery.TableMetadata{
		Schema: schema,
	}

	if fieldPartition != nil {
		partition := bigquery.TimePartitioning{
			Type:  bigquery.DayPartitioningType,
			Field: *fieldPartition,
		}
		metadata.TimePartitioning = &partition
	}

	if fieldClustering != nil {
		clustering := bigquery.Clustering{
			Fields: *fieldClustering,
		}
		metadata.Clustering = &clustering
	}

	if err := table.Create(ctx, metadata); err != nil {
		return fmt.Errorf("bigquery error: %w", err)
	}

	return nil
}

func DeleteByDateColumn(ctx context.Context, bqClient *bigquery.Client, table *bigquery.Table, dateColumn string, dateStart time.Time, dateFinish time.Time) error {
	q := bqClient.Query(fmt.Sprintf("DELETE %s.%s ", table.DatasetID, table.TableID) + fmt.Sprintf("WHERE %s >= '%s' AND %s <= '%s'", dateColumn, dateStart.Format(time.DateOnly),
		dateColumn, dateFinish.Format(time.DateOnly)))

	job, err := q.Run(ctx)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("job error: %w", err)
	}

	if err := status.Err(); err != nil {
		return fmt.Errorf("status error: %w", err)
	}

	return nil
}

func SendFromCS(ctx context.Context, schemaDTO any, table *bigquery.Table, bucket string, object string) error {
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

	job, err := loader.Run(ctx)
	if err != nil {
		return fmt.Errorf("loader error: %w", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("job error: %w", err)
	}

	if err := status.Err(); err != nil {
		return fmt.Errorf("status error: %w", err)
	}

	return nil
}

func TableExists(ctx context.Context, table *bigquery.Table) error {
	if _, err := table.Metadata(ctx); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code == http.StatusNotFound {
				return errors.New("dataset or table not found")
			}
		}
	}

	return nil
}
