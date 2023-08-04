package postgresql

import (
	"context"
	"fmt"
	"github.com/ArkamFahry/uploadnexus/server/rest/errors"
	"github.com/ArkamFahry/uploadnexus/server/rest/models"
	"github.com/lib/pq"
)

func (c *DatabaseClient) CreateBucket(ctx context.Context, bucket models.Bucket) error {
	const Op errors.Op = "postgresql.CreateBucket"

	query := fmt.Sprintf(`INSERT INTO %s (id, name, description, allowed_mime_types, file_size_limit, is_public, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, models.BucketCollection)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.NewError(Op, "failed to prepare statement", err)
	}

	_, err = stmt.ExecContext(ctx, bucket.Id, bucket.Name, bucket.Description, pq.Array(bucket.AllowedMimeTypes), bucket.FileSizeLimit, bucket.IsPublic, bucket.CreatedAt)
	if err != nil {
		return errors.NewError(Op, "failed to create bucket", err)
	}

	return nil
}

func (c *DatabaseClient) UpdateBucket(ctx context.Context, bucket models.Bucket) error {
	const Op errors.Op = "postgresql.UpdateBucket"

	query := fmt.Sprintf(`UPDATE %s SET name = $1, description = $2, allowed_mime_types = $3, file_size_limit = $4, is_public = $5, updated_at = $6 WHERE id = $7`, models.BucketCollection)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.NewError(Op, "failed to prepare statement", err)
	}

	_, err = stmt.ExecContext(ctx, bucket.Name, bucket.Description, pq.Array(bucket.AllowedMimeTypes), bucket.FileSizeLimit, bucket.IsPublic, bucket.UpdatedAt, bucket.Id)
	if err != nil {
		return errors.NewError(Op, "failed to update bucket", err)
	}

	return nil
}

func (c *DatabaseClient) DeleteBucket(ctx context.Context, id string) error {
	const Op errors.Op = "postgresql.DeleteBucket"

	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, models.BucketCollection)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.NewError(Op, "failed to prepare statement", err)
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return errors.NewError(Op, "failed to delete bucket", err)
	}

	return nil
}

func (c *DatabaseClient) CheckIfBucketExistsById(ctx context.Context, id string) (bool, error) {
	const Op errors.Op = "postgresql.CheckIfBucketExistsById"
	var exists bool

	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE id = $1)`, models.BucketCollection)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return false, errors.NewError(Op, "failed to prepare statement", err)
	}

	result := stmt.QueryRowContext(ctx, id)
	err = result.Scan(&exists)
	if err != nil {
		return false, errors.NewError(Op, "failed to check if bucket exists by id", err)
	}
	if exists {
		return true, nil
	}

	return false, nil
}

func (c *DatabaseClient) GetBucketById(ctx context.Context, id string) (*models.Bucket, error) {
	const Op errors.Op = "postgresql.GetBucketById"
	var bucket models.Bucket
	var allowedMimeTypes pq.StringArray

	query := fmt.Sprintf(`SELECT id, name, description, allowed_mime_types, file_size_limit, is_public, created_at, updated_at FROM %s WHERE id = $1`, models.BucketCollection)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.NewError(Op, "failed to prepare statement", err)
	}

	result := stmt.QueryRowContext(ctx, id)
	err = result.Scan(&bucket.Id, &bucket.Name, &bucket.Description, &allowedMimeTypes, &bucket.FileSizeLimit, &bucket.IsPublic, &bucket.CreatedAt, &bucket.UpdatedAt)
	if err != nil {
		return nil, errors.NewError(Op, "failed to get bucket by id", err)
	}

	bucket.AllowedMimeTypes = allowedMimeTypes

	return &bucket, nil
}

func (c *DatabaseClient) CheckIfBucketExistsByName(ctx context.Context, name string) (bool, error) {
	const Op errors.Op = "postgresql.CheckIfBucketExistsByName"
	var exists bool

	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE name = $1)`, models.BucketCollection)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return false, errors.NewError(Op, "failed to prepare statement", err)
	}

	result := stmt.QueryRowContext(ctx, name)
	err = result.Scan(&exists)
	if err != nil {
		return false, errors.NewError(Op, "failed to check if bucket exists by name", err)
	}
	if exists {
		return true, nil
	}

	return false, nil
}

func (c *DatabaseClient) GetBucketByName(ctx context.Context, name string) (*models.Bucket, error) {
	const Op errors.Op = "postgresql.GetBucketByName"
	var bucket models.Bucket
	var allowedMimeTypes pq.StringArray

	query := fmt.Sprintf(`SELECT id, name, description, allowed_mime_types, file_size_limit, is_public, created_at, updated_at FROM %s WHERE name = $1`, models.BucketCollection)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.NewError(Op, "failed to prepare statement", err)
	}

	result := stmt.QueryRowContext(ctx, name)
	err = result.Scan(&bucket.Id, &bucket.Name, &bucket.Description, &allowedMimeTypes, &bucket.FileSizeLimit, &bucket.IsPublic, &bucket.CreatedAt, &bucket.UpdatedAt)
	if err != nil {
		return nil, errors.NewError(Op, "failed to get bucket by name", err)
	}

	bucket.AllowedMimeTypes = allowedMimeTypes

	return &bucket, nil
}

func (c *DatabaseClient) GetBuckets(ctx context.Context) (*[]models.Bucket, error) {
	const Op errors.Op = "postgresql.GetBuckets"
	var buckets []models.Bucket

	query := fmt.Sprintf(`SELECT id, name, description, allowed_mime_types, file_size_limit, is_public, created_at, updated_at FROM %s`, models.BucketCollection)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.NewError(Op, "failed to prepare statement", err)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, errors.NewError(Op, "failed to get buckets", err)
	}

	for rows.Next() {
		var bucket models.Bucket
		var allowedMimeTypes pq.StringArray
		err := rows.Scan(&bucket.Id, &bucket.Name, &bucket.Description, &allowedMimeTypes, &bucket.FileSizeLimit, &bucket.IsPublic, &bucket.CreatedAt, &bucket.UpdatedAt)
		if err != nil {
			return nil, errors.NewError(Op, "failed to get bucket", err)
		}
		bucket.AllowedMimeTypes = allowedMimeTypes
		buckets = append(buckets, bucket)
	}

	return &buckets, nil
}
