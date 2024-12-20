package repositories

import (
	"context"
	// "errors"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "gorm.io/gorm"
)

// Paginated Result
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	TotalPages int64 `json:"total_pages"`
}

// Handle DB Error
func HandleDBError(err error) error {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
	}
	return nil
}

// Find One without fixed struct result
func FindOne[T any](collection *mongo.Collection, filter interface{}, projection interface{}, dest *T) (*T, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	err := collection.FindOne(context.TODO(), filter, opts).Decode(dest)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return dest, nil
}

// FindMany with pagination support
func FindMany[T any](
	collection *mongo.Collection,
	filter interface{},
	sort interface{},
	projection interface{},
	page, pageSize int64,
	dest *[]T,
) (*PaginatedResult[T], error) {

	if filter == nil {
		filter = bson.D{}
	}

	// 设置选项
	findOptions := options.Find()
	if projection != nil {
		findOptions.SetProjection(projection)
	}
	if sort != nil {
		findOptions.SetSort(sort)
	}

	// 计算总数
	total, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int64(math.Ceil(float64(total) / float64(pageSize)))

	// 验证页码
	if page < 1 {
		page = 1
	}
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	// 设置分页
	skip := (page - 1) * pageSize
	findOptions.SetSkip(skip)
	findOptions.SetLimit(pageSize)

	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// 解码结果
	var items []T
	if err = cursor.All(context.TODO(), &items); err != nil {
		return nil, err
	}
	*dest = items
	// 返回分页结果
	return &PaginatedResult[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// FindManyWithoutPagination for cases where pagination is not needed
func FindManyWithoutPagination[T any](
	collection *mongo.Collection,
	filter interface{},
	projection interface{},
	sort interface{},
	dest *[]T,
) ([]T, error) {

	if filter == nil {
		filter = bson.D{}
	}

	findOptions := options.Find()
	if projection != nil {
		findOptions.SetProjection(projection)
	}
	if sort != nil {
		findOptions.SetSort(sort)
	}

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var items []T
	if err = cursor.All(context.TODO(), &items); err != nil {
		return nil, err
	}
	*dest = items
	return items, nil
}

// InsertOne
func InsertOne[T any](collection *mongo.Collection, data *T) error {
	_, err := collection.InsertOne(context.TODO(), data)
	return err
}

// UpdateOne
func UpdateOne(collection *mongo.Collection, filter interface{}, update interface{}) error {
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// UpdateMany
func UpdateMany(collection *mongo.Collection, filter interface{}, update interface{}) error {
	_, err := collection.UpdateMany(context.TODO(), filter, update)
	return err
}

// DeleteOne
func DeleteOne(collection *mongo.Collection, filter interface{}) error {
	_, err := collection.DeleteOne(context.TODO(), filter)
	return err
}

// DeleteMany
func DeleteMany(collection *mongo.Collection, filter interface{}) error {
	_, err := collection.DeleteMany(context.TODO(), filter)
	return err
}
