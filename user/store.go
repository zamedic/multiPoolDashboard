package user

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"bytes"
)

type Store interface {
	addUser(email string, password string) error
	findUser(email string, password string) error
	addKey(email string, pool string, key string) error
	ScanUsers(f func([]string)) error
	GetKey(email, pool string) (string, error)
}

type dynamoDB struct {
	db *dynamodb.DynamoDB
}

func (s dynamoDB) GetKey(email, pool string) (string, error) {
	i := dynamodb.QueryInput{
		TableName: aws.String("miningstats"),
		KeyConditions: map[string]*dynamodb.Condition{
			"email": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(email),
					},
				},
			},
			"type": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(pool),
					},
				},
			},
		},
		AttributesToGet: aws.StringSlice([]string{"apikey"}),
	}

	out, err := s.db.Query(&i)
	if err != nil {
		return "", err
	}
	if *out.Count == 0 {
		return "", err
	}

	var r []key
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &r)
	if err != nil {
		return "", err
	}

	return r[0].Apikey, nil

}

func (s dynamoDB) ScanUsers(f func([]string)) error {
	i := dynamodb.ScanInput{
		TableName:       aws.String("miningstats"),
		AttributesToGet: aws.StringSlice([]string{"email"}),
		ScanFilter: map[string]*dynamodb.Condition{
			"type": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("user"),
					},
				},
			},
		},
	}

	err := s.db.ScanPages(&i,
		func(output *dynamodb.ScanOutput, b bool) bool {
			var r []emails
			err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &r)
			if err != nil {
				return false
			}
			a := make([]string, len(r))
			for x, y := range r {
				a[x] = y.Email
			}
			f(a)
			return b
		},
	)
	return err
}

func NewDynamoStore(db *dynamodb.DynamoDB) Store {
	store := dynamoDB{db: db}
	return store
}

func (s dynamoDB) addUser(email string, password string) error {
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Unable to register user, encryption errir: %v", err.Error())
	}
	i := dynamodb.PutItemInput{
		TableName: aws.String("miningstats"),
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
			"type": {
				S: aws.String("user"),
			},
			"password": {
				S: aws.String(string(p)),
			},
		},
	}

	_, err = s.db.PutItem(&i)
	return err
}

func (s dynamoDB) findUser(email string, password string) error {
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Unable to register user, encryption errir: %v", err.Error())
	}

	u, err := s.getUser(email)
	if err != nil {
		return errors.New("invalid username or password")
	}

	if bytes.Equal([]byte(u.Password), p) {
		return errors.New("invalid username or password")
	}

	return nil

}
func (s dynamoDB) addKey(email string, pool string, key string) error {
	i := dynamodb.UpdateItemInput{
		TableName: aws.String("miningstats"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
			"type": {
				S: aws.String(pool),
			},
		},
		UpdateExpression: aws.String("SET apikey = :k"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":k": {
				S: aws.String(key),
			},
		},
	}

	_, err := s.db.UpdateItem(&i)
	return err
}

func (s dynamoDB) getUser(email string) (*user, error) {
	q := dynamodb.QueryInput{
		TableName: aws.String("miningstats"),
		KeyConditions: map[string]*dynamodb.Condition{
			"email": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(email),
					},
				},
			},
			"type": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("user"),
					},
				},
			},
		},
	}
	out, err := s.db.Query(&q)
	if err != nil {
		return nil, err
	}

	if (*out.Count == 0) {
		return nil, errors.New("invalid username or password")
	}

	u := []user{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &u)
	return &u[0], err

}

type user struct {
	Email    string
	Password string
}

type emails struct {
	Email string
}

type key struct {
	Apikey string
}
