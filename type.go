package nyobain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email    string             `bson:"email,omitempty" json:"email,omitempty"`
	Password string             `bson:"password,omitempty" json:"password,omitempty"`
	Salt     string             `bson:"salt,omitempty,omitempty" json:"salt,omitempty"`
	Role     string             `bson:"role,omitempty" json:"role,omitempty"`
}

type Password struct {
	Password    string `bson:"password,omitempty" json:"password,omitempty"`
	Newpassword string `bson:"newpass,omitempty" json:"newpass,omitempty"`
}

type Pengguna struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap  string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	TanggalLahir string             `bson:"tanggallahir,omitempty" json:"tanggallahir,omitempty"`
	JenisKelamin string             `bson:"jeniskelamin,omitempty" json:"jeniskelamin,omitempty"`
	NomorHP      string             `bson:"nomorhp,omitempty" json:"nomorhp,omitempty"`
	Alamat       string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
	Akun         User               `bson:"akun,omitempty" json:"akun,omitempty"`
}

type Admin struct {
	ID   primitive.ObjectID `bson:"_id,omitempty,omitempty" json:"_id,omitempty"`
	Akun User               `bson:"akun,omitempty" json:"akun,omitempty"`
}

type Seller struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	NamaToko    string             `bson:"namatoko,omitempty" json:"namatoko,omitempty"`
	NomorHP     string             `bson:"nomorhp,omitempty" json:"nomorhp,omitempty"`
	Alamat      string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
	Akun        User               `bson:"akun,omitempty" json:"akun,omitempty"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaProduct string             `bson:"productid,omitempty" json:"productid,omitempty"`
	Deskripsi   string             `json:"deskripsi,omitempty" bson:"deskripsi,omitempty"`
	Kategori    string             `json:"kategori,omitempty" bson:"kategori,omitempty"`
	Harga       string             `json:"harga,omitempty" bson:"harga,omitempty"`
}

type Orderproduct struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Pengguna    Pengguna           `bson:"pengguna,omitempty" json:"pengguna,omitempty"`
	Seller      Seller             `bson:"seller,omitempty" json:"seller,omitempty"`
	NamaProduct string             `bson:"namaproduct,omitempty" json:"namaproduct,omitempty"`
	Quantity    string             `bson:"quantity,omitempty" json:"quantity,omitempty"`
	TotalCost   string             `bson:"total_cost,omitempty" json:"total_cost,v"`
	Status      string             `bson:"status,omitempty" json:"status,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Role    string `json:"role,omitempty" bson:"role,omitempty"`
}

type Response struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Payload struct {
	Id   primitive.ObjectID `json:"id"`
	Role string             `json:"role"`
	Exp  time.Time          `json:"exp"`
	Iat  time.Time          `json:"iat"`
	Nbf  time.Time          `json:"nbf"`
}
