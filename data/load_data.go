package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db *dynamodb.DynamoDB

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

func main() {
	InitializeDB()
	CreateTable()
	StoreData(ParseFile())
}

// InitializeDB create a DynamoDB reference
func InitializeDB() {
	fmt.Println("Initializing DynamoDB session ...")

	db = dynamodb.New(session.New(&aws.Config{
		Endpoint: aws.String(os.Getenv("DYNAMODB_ENDPOINT")),
	}))
}

// CreateTable creates the dynamodb table
func CreateTable() {
	fmt.Println("Creating table 'Movies' ...")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("movie-id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("movie-id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Movies"),
	}

	_, err := db.CreateTable(input)

	if err != nil {
		fmt.Println("Failed to create table 'Movies'. Got error:", err)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Created the table Movies in us-east-1.")
}

// ParseFile parses the movies.csv file return an array of Movies
func ParseFile() []Movie {
	fmt.Println("Parsing file 'movies.csv' ...")

	csvFile, _ := os.Open("movies.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var movies []Movie
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}
		movies = append(movies, Movie{
			MovieID: line[3],
			Title:   line[17],
			Budget: func() int64 {
				value, _ := strconv.ParseInt(line[0], 10, 64)
				return value
			}(),
			ReleaseDate: func() time.Time {
				tokens := strings.Split(line[11], "-")
				year, _ := strconv.Atoi(tokens[0])
				month, _ := strconv.Atoi(tokens[1])
				day, _ := strconv.Atoi(tokens[2])
				return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			}(),
			Revenue: func() int64 {
				value, _ := strconv.ParseInt(line[12], 10, 64)
				return value
			}(),
			Runtime: func() int64 {
				value, _ := strconv.ParseInt(line[13], 10, 64)
				return value
			}(),
			VoteAverage: func() float32 {
				value, _ := strconv.ParseFloat(line[18], 32)
				return float32(value)
			}(),
			VoteCount: func() int64 {
				value, _ := strconv.ParseInt(line[19], 10, 64)
				return value
			}(),
		})
	}
	return movies
}

// StoreData creates the items in DynamoDB
func StoreData(movies []Movie) {
	fmt.Println("Storing movies into DynamoDB ...")

	if movies == nil || len(movies) == 0 {
		return
	}

	for i := 0; i < len(movies); i++ {
		av, err := dynamodbattribute.MarshalMap(movies[i])

		if err != nil {
			fmt.Println("Failed to store movie:", movies[i], "Got error:", err)
			continue
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("Movies"),
			ExpressionAttributeNames: map[string]*string{
				"#m": aws.String("movie-id"),
			},
			ConditionExpression: aws.String("attribute_not_exists(#m)"),
		}

		_, err = db.PutItem(input)

		if err != nil {
			fmt.Println("Failed to store movie:", movies[i], "Got error:", err)
			continue
		}

		fmt.Println("Stored movie:", movies[i])
	}
}
