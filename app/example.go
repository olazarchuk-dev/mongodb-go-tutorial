package app

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mongodb-go-tutorial/models"
	"os"
)

var (
	BooksCollection   *mongo.Collection
	AuthorsCollection *mongo.Collection
	UsersCollection   *mongo.Collection
	Ctx               = context.TODO()
)

/*Setup opens a database connection to mongodb*/
func Setup() {

	/**
	 * @see https://www.loginradius.com/blog/async/environment-variables-in-golang
	 *      https://github.com/LoginRadius/engineering-blog-samples/tree/master/GoLang/EnvironmentVariables/godotenvtest
	 * https://www.geeksforgeeks.org/golang-environment-variables
	 *
	 * load .env file
	 */
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connectionURI := "mongodb://" + os.Getenv("mongo_host") + ":" + os.Getenv("mongo_port") + "/"
	clientOptions := options.Client().ApplyURI(connectionURI)
	client, err := mongo.Connect(Ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(Ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(os.Getenv("mongo_database"))
	BooksCollection = db.Collection("books")
	AuthorsCollection = db.Collection("authors")
	UsersCollection = db.Collection("users")
}

/**
 * Определим наши модели
 */
type Author struct {
	FullName string `bson:"full_name"`
}

type Book struct {
	Name      string `bson:"name"`
	Author    string `bson:"author"`
	PageCount int    `bson:"page_count"`
}

/**
 * Напишем 6 функций для наших операций и позволим им делать следующее:
 *
 * 1. CreateBook() -> Добавляет новую запись в коллекцию книг
 * 2. GetBook() -> Возвращает запись из коллекции книг по параметру id
 * 3. GetBooks() -> Возвращает все записи
 * 4. UpdateBooks() -> Обновляет запись в соответствии с параметром id
 * 5. FindAuthorBooks() -> Извлекает все книги автора
 * 6. DeleteBook() -> Удаляет запись с параметром id
 */

// 1. Создать книгу:
func CreateBook(b Book) (string, error) {
	result, err := BooksCollection.InsertOne(Ctx, b)
	if err != nil {
		return "0", err
	}
	return fmt.Sprintf("%v", result.InsertedID), err
}

func CreateUser(u models.User) (string, error) {
	result, err := UsersCollection.InsertOne(Ctx, u)
	if err != nil {
		return "0", err
	}
	return fmt.Sprintf("%v", result.InsertedID), err
}

// 2. Получить книгу:
func GetBook(id string) (Book, error) {
	var b Book
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return b, err
	}

	err = BooksCollection.
		FindOne(Ctx, bson.D{{"_id", objectId}}).
		Decode(&b)
	if err != nil {
		return b, err
	}
	return b, nil
}

// 3. Получить книги:
func GetBooks() ([]Book, error) {
	var book Book
	var books []Book

	cursor, err := BooksCollection.Find(Ctx, bson.D{})
	if err != nil {
		defer cursor.Close(Ctx)
		return books, err
	}

	for cursor.Next(Ctx) {
		err := cursor.Decode(&book)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}

	return books, nil
}

func GetUsers() ([]models.User, error) {
	var user models.User
	var users []models.User

	cursor, err := UsersCollection.Find(Ctx, bson.D{})
	if err != nil {
		defer cursor.Close(Ctx)
		return users, err
	}

	for cursor.Next(Ctx) {
		err := cursor.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

// 4. Обновить книгу:
func UpdateBook(id primitive.ObjectID, pageCount int) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"page_count", pageCount}}}}
	_, err := BooksCollection.UpdateOne(
		Ctx,
		filter,
		update,
	)
	return err
}

// 5. Найти автора книги:
type AuthorBooks struct {
	FullName string `bson:"full_name"`
	Books    []Book
}

/**
 * @see https://jira.mongodb.org/browse/GODRIVER-1129
 *      https://www.mongodb.com/blog/post/quick-start-golang--mongodb--data-aggregation-pipeline
 *      https://www.digitalocean.com/community/tutorials/how-to-use-aggregations-in-mongodb
 */
func FindAuthorBooks(fullName string) ([]Book, error) {
	matchStage := bson.D{{"$match", bson.D{{"full_name", fullName}}}}

	lookupStage := bson.D{{"$lookup",
		bson.D{{"from", "books"},
			{"localField", "full_name"},
			{"foreignField", "author"},
			{"as", "books"}}}}

	showLoadedCursor, err := AuthorsCollection.Aggregate(Ctx, mongo.Pipeline{matchStage, lookupStage})
	if err != nil {
		return nil, err
	}

	var a []AuthorBooks
	if err = showLoadedCursor.All(Ctx, &a); err != nil { // https://jira.mongodb.org/browse/GODRIVER-1129
		return nil, err
	}
	return a[0].Books, err
}

// 6. Удалить книгу:
func DeleteBook(id primitive.ObjectID) error {
	_, err := BooksCollection.DeleteOne(Ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}
