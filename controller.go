package nyobain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/badoux/checkmail"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/argon2"
)

var (
	Response     nyobain.Response
	user         nyobain.User
	pengguna     nyobain.Pengguna
	seller       nyobain.Seller
	product      nyobain.Product
	Orderproduct nyobain.Orderproduct
	password     nyobain.Password
)

func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

// crud
func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error GetAllDocs %s: %s", col, err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return err
	}
	return docs
}

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func UpdateOneDoc(id primitive.ObjectID, db *mongo.Database, col string, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error update: %v", err)
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("tidak ada data yang diubah")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// signup
func SignUpPengguna(db *mongo.Database, insertedDoc nyobain.Pengguna) error {
	objectId := primitive.NewObjectID()
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("Dimohon untuk melengkapi data")
	}
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	}
	userExists, _ := GetUserFromEmail(insertedDoc.Akun.Email, db)
	if insertedDoc.Akun.Email == userExists.Email {
		return fmt.Errorf("Email sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Akun.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Akun.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Akun.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"_id":      objectId,
		"email":    insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt":     hex.EncodeToString(salt),
		"role":     "pengguna",
	}
	pengguna := bson.M{
		"namalengkap":  insertedDoc.NamaLengkap,
		"tanggallahir": insertedDoc.TanggalLahir,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp":      insertedDoc.NomorHP,
		"alamat":       insertedDoc.Alamat,
		"akun": nyobain.User{
			ID: objectId,
		},
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	_, err = InsertOneDoc(db, "pengguna", pengguna)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	return nil
}

func SignUpSeller(db *mongo.Database, insertedDoc nyobain.Seller) error {
	objectId := primitive.NewObjectID()
	if insertedDoc.NamaLengkap == "" || insertedDoc.NamaToko == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("dimohon untuk melengkapi data")
	}
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	userExists, _ := GetUserFromEmail(insertedDoc.Akun.Email, db)
	if insertedDoc.Akun.Email == userExists.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Akun.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Akun.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Akun.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"_id":      objectId,
		"email":    insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt":     hex.EncodeToString(salt),
		"role":     "seller",
	}
	seller := bson.M{
		"namalengkap":  insertedDoc.NamaLengkap,
		"jeniskelamin": insertedDoc.NamaToko,
		"nomorhp":      insertedDoc.NomorHP,
		"alamat":       insertedDoc.Alamat,
		"akun": nyobain.User{
			ID: objectId,
		},
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return err
	}
	_, err = InsertOneDoc(db, "seller", seller)
	if err != nil {
		return err
	}
	return nil
}

// login
func LogIn(db *mongo.Database, insertedDoc nyobain.User) (user nyobain.User, err error) {
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return user, fmt.Errorf("Dimohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return user, fmt.Errorf("Email tidak valid")
	}
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return
	}
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		return user, fmt.Errorf("kesalahan server : salt")
	}
	hash := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		return user, fmt.Errorf("password salah")
	}
	return existsDoc, nil
}

// user
func UpdateEmailUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc nyobain.User) error {
	dataUser, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	if insertedDoc.Email == "" {
		return fmt.Errorf("Dimohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	}
	existsDoc, _ := GetUserFromEmail(insertedDoc.Email, db)
	if existsDoc.Email == insertedDoc.Email {
		return fmt.Errorf("Email sudah terdaftar")
	}
	user := bson.M{
		"email":    insertedDoc.Email,
		"password": dataUser.Password,
		"salt":     dataUser.Salt,
		"role":     dataUser.Role,
	}
	err = UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePasswordUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc nyobain.Password) error {
	dataUser, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	salt, err := hex.DecodeString(dataUser.Salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hash := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != dataUser.Password {
		return fmt.Errorf("password lama salah")
	}
	if insertedDoc.Newpassword == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if strings.Contains(insertedDoc.Newpassword, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Newpassword) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt = make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Newpassword), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"email":    dataUser.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt":     hex.EncodeToString(salt),
		"role":     dataUser.Role,
	}
	err = UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc nyobain.User) error {
	dataUser, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	existsDoc, _ := GetUserFromEmail(insertedDoc.Email, db)
	if existsDoc.Email == insertedDoc.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"email":    insertedDoc.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt":     hex.EncodeToString(salt),
		"role":     dataUser.Role,
	}
	err = UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUser(db *mongo.Database) (user []nyobain.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return user, fmt.Errorf("error GetAllUser mongo: %s", err)
	}
	err = cursor.All(context.Background(), &user)
	if err != nil {
		return user, fmt.Errorf("error GetAllUser context: %s", err)
	}
	return user, nil
}

func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc nyobain.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetUserFromEmail(email string, db *mongo.Database) (doc nyobain.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// pengguna
func UpdatePengguna(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc nyobain.Pengguna) error {
	pengguna, err := GetPenggunaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if pengguna.ID != idparam {
		return fmt.Errorf("Anda bukan pemilik data ini")
	}
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" {
		return fmt.Errorf("Dimohon untuk melengkapi data")
	}
	pgn := bson.M{
		"namalengkap":  insertedDoc.NamaLengkap,
		"tanggallahir": insertedDoc.TanggalLahir,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp":      insertedDoc.NomorHP,
		"alamat":       insertedDoc.Alamat,
		"akun": nyobain.User{
			ID: pengguna.Akun.ID,
		},
	}
	err = UpdateOneDoc(idparam, db, "pengguna", pgn)
	if err != nil {
		return err
	}
	return nil
}

func GetAllPengguna(db *mongo.Database) (pengguna []nyobain.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna mongo: %s", err)
	}
	err = cursor.All(context.Background(), &pengguna)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna context: %s", err)
	}
	return pengguna, nil
}

func GetPenggunaFromID(_id primitive.ObjectID, db *mongo.Database) (doc nyobain.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetPenggunaFromAkun(akun primitive.ObjectID, db *mongo.Database) (doc nyobain.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{"akun._id": akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("pengguna tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// by admin
func GetPenggunaFromIDByAdmin(idparam primitive.ObjectID, db *mongo.Database) (pengguna nyobain.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{
		"_id": idparam,
	}
	err = collection.FindOne(context.Background(), filter).Decode(&pengguna)
	if err != nil {
		return pengguna, fmt.Errorf("error GetPenggunaFromID mongo: %s", err)
	}
	user, err := GetUserFromID(pengguna.Akun.ID, db)
	if err != nil {
		return pengguna, fmt.Errorf("error GetPenggunaFromID mongo: %s", err)
	}
	akun := nyobain.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	pengguna.Akun = akun
	return pengguna, nil
}

func GetAllPenggunaByAdmin(db *mongo.Database) (pengguna []nyobain.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna mongo: %s", err)
	}
	err = cursor.All(context.Background(), &pengguna)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna context: %s", err)
	}
	return pengguna, nil
}

func GetSellerFromIDByAdmin(idparam primitive.ObjectID, db *mongo.Database) (seller nyobain.Seller, err error) {
	collection := db.Collection("seller")
	filter := bson.M{
		"_id": idparam,
	}
	err = collection.FindOne(context.Background(), filter).Decode(&seller)
	if err != nil {
		return seller, err
	}
	user, err := GetUserFromID(seller.Akun.ID, db)
	if err != nil {
		return seller, err
	}
	akun := nyobain.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	seller.Akun = akun
	return seller, nil
}

func UpdateSeller(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc nyobain.Seller) error {
	seller, err := GetSellerFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if seller.ID != idparam {
		return fmt.Errorf("anda bukan pemilik data ini")
	}
	if insertedDoc.NamaLengkap == "" || insertedDoc.NamaToko == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" {
		return fmt.Errorf("dimohon untuk melengkapi data")
	}
	drv := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"namatoko":    insertedDoc.NamaToko,
		"nomorhp":     insertedDoc.NomorHP,
		"alamat":      insertedDoc.Alamat,
		"akun": nyobain.User{
			ID: seller.Akun.ID,
		},
	}
	err = UpdateOneDoc(idparam, db, "seller", drv)
	if err != nil {
		return err
	}
	return nil
}

func GetAllSeller(db *mongo.Database) (seller []nyobain.Seller, err error) {
	collection := db.Collection("seller")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return seller, fmt.Errorf("error GetAllSeller mongo: %s", err)
	}
	err = cursor.All(context.Background(), &seller)
	if err != nil {
		return seller, fmt.Errorf("error GetAllSeller context: %s", err)
	}
	return seller, nil
}

func GetSellerFromID(_id primitive.ObjectID, db *mongo.Database) (doc nyobain.Seller, err error) {
	collection := db.Collection("seller")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("id seller tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func GetSellerFromAkun(akun primitive.ObjectID, db *mongo.Database) (doc nyobain.Seller, err error) {
	collection := db.Collection("seller")
	filter := bson.M{"akun._id": akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("Seller tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

//obat

func InsertProduct(iduser primitive.ObjectID, db *mongo.Database, insertedDoc nyobain.Product) error {
	if insertedDoc.NamaProduct == "" || insertedDoc.Deskripsi == "" || insertedDoc.Kategori == "" || insertedDoc.Harga == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}

	tkt := bson.M{
		"namaproduct": insertedDoc.NamaProduct,
		"deskripsi":   insertedDoc.Deskripsi,
		"kategori":    insertedDoc.Kategori,
		"harga":       insertedDoc.Harga,
	}

	_, err := InsertOneDoc(db, "produk", tkt)
	if err != nil {
		return fmt.Errorf("error saat menyimpan data produk: %s", err)
	}
	return nil
}

func Updateproduct(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc nyobain.Product) error {
	_, err := GetProductFromID(idparam, db)
	if err != nil {
		return err
	}
	if insertedDoc.NamaProduct == "" || insertedDoc.Deskripsi == "" || insertedDoc.Kategori == "" || insertedDoc.Harga == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	tkt := bson.M{
		"namaproduct": insertedDoc.NamaProduct,
		"deskripsi":   insertedDoc.Deskripsi,
		"kategori":    insertedDoc.Kategori,
		"harga":       insertedDoc.Harga,
	}

	err = UpdateOneDoc(idparam, db, "produk", tkt)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProduct(idparam, iduser primitive.ObjectID, db *mongo.Database) error {
	_, err := GetProductFromID(idparam, db)
	if err != nil {
		return err
	}
	err = DeleteOneDoc(idparam, db, "produk")
	if err != nil {
		return err
	}
	return nil
}

func GetAllProduct(db *mongo.Database) (product []nyobain.Product, err error) {
	collection := db.Collection("produk")
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return product, fmt.Errorf("error GetAllProduct mongo: %s", err)
	}
	err = cursor.All(context.TODO(), &product)
	if err != nil {
		return product, fmt.Errorf("error GetAllProduct context: %s", err)
	}
	return product, nil
}

func GetProductFromID(_id primitive.ObjectID, db *mongo.Database) (doc nyobain.Product, err error) {
	collection := db.Collection("product")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("id product tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

//order

func InsertOrderProduct(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc nyobain.Orderproduct) error {

	if insertedDoc.NamaProduct == "" || insertedDoc.Quantity == "" || insertedDoc.TotalCost == "" || insertedDoc.Status == "" {
		return fmt.Errorf("harap lengkapi semua data order")
	}

	ord := bson.M{
		"pengguna": bson.M{
			"_id": iduser,
		},
		"seller": bson.M{
			"_id": insertedDoc.Seller.ID,
		},
		"order": bson.M{
			"_id": idparam,
		},
		"namaproduct": insertedDoc.NamaProduct,
		"quantity":    insertedDoc.Quantity,
		"total_cost":  insertedDoc.TotalCost,
		"status":      insertedDoc.Status,
	}

	_, err := InsertOneDoc(db, "order", ord)
	if err != nil {
		return fmt.Errorf("error saat menyimpan data order produk: %s", err)
	}
	return nil
}

// update status pengiriman
func UpdateStatusOrderProduct(idorder primitive.ObjectID, db *mongo.Database, insertedDoc nyobain.Orderproduct) error {
	order, err := GetOrderFromID(idorder, db)
	if err != nil {
		return err
	}

	data := bson.M{
		"namaproduct": order.NamaProduct,
		"quantity":    order.Quantity,
		"total_cost":  order.TotalCost,
		"status":      insertedDoc.Status,
	}

	err = UpdateOneDoc(idorder, db, "order", data)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOrder(idparam, iduser primitive.ObjectID, db *mongo.Database) error {
	_, err := GetOrderFromID(idparam, db)
	if err != nil {
		return err
	}
	err = DeleteOneDoc(idparam, db, "order")
	if err != nil {
		return err
	}
	return nil
}

func GetOrderFromID(_id primitive.ObjectID, db *mongo.Database) (doc nyobain.Orderproduct, err error) {
	collection := db.Collection("order")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("id transaksi tidak ditemukan")
		}
		return doc, err
	}
	return doc, nil
}

func GetAllOrder(db *mongo.Database) (order []nyobain.Orderproduct, err error) {
	collection := db.Collection("order")
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return order, fmt.Errorf("error GetAllOrder mongo: %s", err)
	}
	err = cursor.All(context.TODO(), &order)
	if err != nil {
		return order, fmt.Errorf("error GetAllOrder context: %s", err)
	}
	return order, nil
}
