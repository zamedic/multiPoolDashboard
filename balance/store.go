package balance

import (
	"time"
	"fmt"
	"strconv"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Store interface {
	SaveBalance(email string, value map[string][]CoinValue)
	getBalanceByPool(email string, start time.Time) map[int64]map[string][]CoinValue
	getBalanceByCoinName(email, coin string, start time.Time) map[int64]float64
}

type dynamoDB struct {
	db *dynamodb.DynamoDB
}

func NewDynamoStore(db *dynamodb.DynamoDB) Store {
	return &dynamoDB{db: db}
}

func (s dynamoDB) getBalanceByCoinName(email, coin string, start time.Time) map[int64]float64 {
	i := getQueryInput(email, start)
	i.FilterExpression = aws.String("#coin = :name")
	i.ExpressionAttributeNames = map[string]*string{
		"#coin": aws.String("coinname"),
	}
	i.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
		":name": {
			S:aws.String(coin),
		},
	}

	result := map[int64]float64{}
	s.db.QueryPages(&i, func(output *dynamodb.QueryOutput, b bool) bool {
		var r []balancedb
		dynamodbattribute.UnmarshalListOfMaps(output.Items, &r)
		for _, item := range r {
			_, ok := result[item.Timestamp]
			if !ok {
				result[item.Timestamp] = item.Coinvalue
			} else {
				result[item.Timestamp] = result[item.Timestamp] + item.Coinvalue
			}

		}
		return !b
	})

	return result
}

func (s dynamoDB) getBalanceByPool(email string, start time.Time) map[int64]map[string][]CoinValue {

	i := getQueryInput(email, start)
	result := map[int64]map[string][]CoinValue{}

	s.db.QueryPages(&i, func(output *dynamodb.QueryOutput, b bool) bool {
		var r []balancedb
		dynamodbattribute.UnmarshalListOfMaps(output.Items, &r)
		for _, item := range r {
			timevalue, ok := result[item.Timestamp]
			if !ok {
				timevalue = map[string][]CoinValue{}
			}

			poolvalue, ok := timevalue[item.Pool]
			if !ok {
				poolvalue = []CoinValue{}
			}
			poolvalue = append(poolvalue, CoinValue{Coins: item.Coinvalue, Name: item.Coinname})
			timevalue[item.Pool] = poolvalue
			result[item.Timestamp] = timevalue
		}

		return !b

	})

	return result
}

func (s dynamoDB) SaveBalance(email string, value map[string][]CoinValue) {
	b := dynamodb.BatchWriteItemInput{}
	var w []*dynamodb.WriteRequest
	t := time.Now().Unix()
	for poolname, poolcoins := range value {
		for _, coins := range poolcoins {
			a := &dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{
					Item: map[string]*dynamodb.AttributeValue{
						"email": {
							S: aws.String(email),
						},
						"type": {
							S: aws.String(fmt.Sprintf("%v|%v|%v", poolname, coins.Name, strconv.FormatInt(t, 10))),
						},
						"pool": {
							S: aws.String(poolname),
						},
						"timestamp": {
							N: aws.String(strconv.FormatInt(t, 10)),
						},
						"coinname": {
							S: aws.String(coins.Name),
						},
						"coinvalue": {
							N: aws.String(strconv.FormatFloat(coins.Coins, 'f', 6, 64)),
						},
					},
				},
			}
			w = append(w, a)
		}
	}

	b.RequestItems = map[string][]*dynamodb.WriteRequest{
		"miningstats": w,
	}

	s.db.BatchWriteItem(&b)
}

func getQueryInput(email string, start time.Time) dynamodb.QueryInput {
	i := dynamodb.QueryInput{
		TableName: aws.String("miningstats"),
		IndexName: aws.String("email-timestamp-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"email": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(email),
					},
				},
			},
			"timestamp": {
				ComparisonOperator: aws.String("GE"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						N: aws.String(strconv.FormatInt(start.Unix(), 10)),
					},
				},
			},
		},
	}
	return i
}

type CoinValue struct {
	Name  string
	Coins float64
}

type balancedb struct {
	Pool      string
	Timestamp int64
	Coinname  string
	Coinvalue float64
}
