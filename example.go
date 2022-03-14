package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	BooksCollection   *mongo.Collection
	AuthorsCollection *mongo.Collection
	Ctx               = context.TODO()
)

/*Setup opens a database connection to mongodb*/
func Setup() {
	host := "127.0.0.1"
	port := "27017"
	connectionURI := "mongodb://" + host + ":" + port + "/"
	clientOptions := options.Client().ApplyURI(connectionURI)
	client, err := mongo.Connect(Ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(Ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("example")
	BooksCollection = db.Collection("books")
	AuthorsCollection = db.Collection("authors")
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

func main() {
	Setup()

	//author := Author{ "GeeksforGeeks" }

	//newBook1 := Book{ "How to Use Go With MongoDB", "GeeksforGeeks", 100 } // TODO: 622f6e3fc134f9ed331b4b33
	//newBook2 := Book{ "How to Do CRUD Transactions in MongoDB with Go", "hackajob Staff", 99 } // TODO: 622f6e3fc134f9ed331b4b34
	//newBook3 := Book{ "MongoDB Go Driver туториал", "pocoZ", 101 } // TODO: 622f6e3fc134f9ed331b4b35
	//CreateBook(newBook1)
	//CreateBook(newBook2)
	//CreateBook(newBook3)

	//result, err := GetBook("622f6e3fc134f9ed331b4b33")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Name='%v'; Author='%v'; PageCount='%v'; \n", result.Name, result.Author, result.PageCount)

	results, err := GetBooks()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(results)

	/*
	 * https://www.geeksforgeeks.org/loops-in-go-language
	 */
	//allRec := 0
	for count, result := range results {
		//allRec = count
		//fmt.Printf("Name='%v'; Author='%v'; PageCount='%v'; \n", result.Name, result.Author, result.PageCount)
		fmt.Printf("%v. Name='%v'; Author='%v'; PageCount='%v'; \n", count, result.Name, result.Author, result.PageCount)
	}
	//allRec = allRec * allRec

	/*
	 * @see https://serveanswer.com/questions/how-to-convert-string-to-primitive-objectid-in-golang
	 *      https://stackoverflow.com/questions/60864873/primitive-objectid-to-string-in-golang
	 */
	bookId, err := primitive.ObjectIDFromHex("622f6e3fc134f9ed331b4b34")
	UpdateBook(bookId, 199)

	DeleteBook(bookId)

	fullName := "GeeksforGeeks"
	resultsAll, errAll := FindAuthorBooks(fullName)
	if errAll != nil {
		log.Fatal(errAll)
	}
	for count, result := range resultsAll {
		fmt.Printf("%v. FullName='%v'; > Name='%v'; Author='%v'; PageCount='%v'; \n", count, fullName, result.Name, result.Author, result.PageCount)
	}

}
