package movie

import (
	"errors"
	"time"

	"../../db"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Movie represents a movie
type Movie struct {
	MovieID     string    `json:"movie-id"`
	Title       string    `json:"title" binding:"required"`
	Budget      int64     `json:"budget" binding:"required"`
	ReleaseDate time.Time `json:"release-date" binding:"required"`
	Revenue     int64     `json:"revenue" binding:"required"`
	Runtime     int64     `json:"runtime" binding:"required"`
	VoteAverage float32   `json:"vote-average" binding:"required"`
	VoteCount   int64     `json:"vote-count" binding:"required"`
}

// CreateMovie creates a new movie item and store it in the database
func CreateMovie(movie *Movie) (bool, error) {
	if movie.MovieID == "" {
		return false, errors.New("Missing movie id property")
	}

	av, err := dynamodbattribute.MarshalMap(movie)

	if err != nil {
		return false, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Movies"),
		ExpressionAttributeNames: map[string]*string{
			"#m": aws.String("movie-id"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#m)"),
	}

	_, err = db.GetDynamoDB().PutItem(input)

	if err != nil {
		return false, err
	}

	return true, nil
}

// RetrieveMovie retrieves the movie from the database given the movie id
func RetrieveMovie(id string) (*Movie, error) {
	result, err := db.GetDynamoDB().GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Movies"),
		Key: map[string]*dynamodb.AttributeValue{
			"movie-id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	movie := Movie{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &movie)

	if err != nil {
		return nil, err
	}

	if movie.MovieID == "" {
		return nil, errors.New("Movie not found")
	}

	return &movie, nil
}

// UpdateMovie updates an existing movie in the database
func UpdateMovie(movie *Movie) (bool, error) {
	av, err := dynamodbattribute.MarshalMap(movie)

	if err != nil {
		return false, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Movies"),
		ExpressionAttributeNames: map[string]*string{
			"#m": aws.String("movie-id"),
		},
		ConditionExpression: aws.String("attribute_exists(#m)"),
	}

	_, err = db.GetDynamoDB().PutItem(input)

	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteMovie removes a existing move from the database
func DeleteMovie(id string) (bool, error) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"movie-id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String("Movies"),
		ExpressionAttributeNames: map[string]*string{
			"#m": aws.String("movie-id"),
		},
		ConditionExpression: aws.String("attribute_exists(#m)"),
	}

	_, err := db.GetDynamoDB().DeleteItem(input)

	if err != nil {
		return false, err
	}

	return true, nil
}
