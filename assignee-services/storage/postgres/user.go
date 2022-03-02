package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	// "time"
	pb "two_services/assignee-services/genproto"

	"github.com/go-redis/redis"
	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"golang.org/x/crypto/bcrypt"

	// "github.com/google/uuid"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

var (
	isnullphoto sql.NullString
	isnullbio   sql.NullString
	isupdated   sql.NullString
	newaddress  []byte
	newphone    []byte
)

func (r *UserRepo) GetEmail(req *pb.Email) (string, error) {
	query := `SELECT password FROM assignees WHERE email=$1`
	email := ""
	err := r.db.QueryRow(query, req.Email).Scan(
		&email,
	)
	if err != nil {
		return "", err
	}
	return email, err
}

func (r *UserRepo) Regist(user *pb.CreateUserReqWithCode) (*pb.Mess, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	json, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = client.Set(user.Username, json, time.Second*3600).Err()
	if err != nil {
		return nil, err
	}
	return &pb.Mess{Res: "Ok!"}, nil
}
func (r *UserRepo) Verfy(get *pb.Check) (*pb.CreateUserReq, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	user := pb.CreateUserReqWithCode{}
	res, err := client.Get(get.Username).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		return nil, err
	}
	if user.Username == get.Username && user.VerfCode == get.Code {
		newUser := pb.CreateUserReq{
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Username:     user.Username,
			ProfilePhoto: user.ProfilePhoto,
			Bio:          user.Bio,
			Email:        user.Email,
			Gender:       user.Gender,
			Address:      user.Address,
			Phone:        user.Phone,
			Password:     user.Password,
		}
		return &newUser, nil

	} else {
		return nil, errors.New("In Verify")
	}
}
func (r *UserRepo) Login(user *pb.EmailWithPassword) (*pb.UserRes, error) {
	query := `SELECT 
					id,
					first_name,
					last_name,
					username,
					profile_photo,
					bio,
					email,
					gender,
					address,
					phone_num,
					created_at,
					accses_token,
					refresh_token,
					password	
					FROM assignees
			WHERE email=$1`
	newVar := pb.UserRes{}
	pass := ""
	err := r.db.QueryRow(query, user.Email).Scan(
		&newVar.Id,
		&newVar.FirstName,
		&newVar.LastName,
		&newVar.Username,
		&isnullphoto,
		&isnullbio,
		&newVar.Email,
		&newVar.Gender,
		&newaddress,
		&newphone,
		&newVar.CreatedAt,
		&newVar.AccesToken,
		&newVar.RefreshToken,
		&pass,
	)
	if isnullbio.Valid {
		newVar.Bio = isnullbio.String
	}
	if isnullphoto.Valid {
		newVar.ProfilePhoto = isnullphoto.String
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(newaddress, &newVar.Address)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(newphone, &newVar.PhoneNum)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(user.Password))
	if err != nil {
		return nil, err
	}
	return &newVar, nil
}
func (r *UserRepo) Create(user *pb.CreateUserReqWithCode) (*pb.UserRes, error) {
	query := `INSERT INTO assignees (
		id,
		first_name, 
		last_name, 
		username, 
		profile_photo, 
		bio, 
		email, 
		gender, 
		address, 
		phone_num,
		created_at,
		accses_token,
		refresh_token,
		password) 
					VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) 
					RETURNING 
					id,
					first_name,
					last_name,
					username,
					profile_photo,
					bio,
					email,
					gender,
					address,
					phone_num,
					created_at,
					accses_token,
					refresh_token`
	newUser := pb.UserRes{}

	newaddress, err := json.Marshal(user.Address)
	if err != nil {
		return nil, err
	}
	newphone, err = json.Marshal(user.Phone)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(query,
		user.Id,
		user.FirstName,
		user.LastName,
		user.Username,
		user.ProfilePhoto,
		user.Bio,
		user.Email,
		user.Gender,
		newaddress,
		newphone,
		time.Now().Format(time.RFC3339),
		user.AccsesToken,
		user.RefreshToken,
		user.Password,
	).Scan(
		&newUser.Id,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Username,
		&isnullphoto,
		&isnullbio,
		&newUser.Email,
		&newUser.Gender,
		&newaddress,
		&newphone,
		&newUser.CreatedAt,
		&newUser.AccesToken,
		&newUser.RefreshToken,
	)
	if isnullbio.Valid {
		newUser.Bio = isnullbio.String
	}
	if isnullphoto.Valid {
		newUser.ProfilePhoto = isnullphoto.String
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(newaddress, &newUser.Address)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(newphone, &newUser.PhoneNum)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &newUser, nil
}

func (r *UserRepo) Get(user *pb.GetOrDeleteUser) (*pb.UserRes, error) {

	query := `SELECT id,first_name,last_name,username,profile_photo,bio,email,gender,address,phone_num,created_at,updated_at FROM assignees WHERE id=$1 AND deleted_at IS NULL`
	newUser := pb.UserRes{}
	err := r.db.QueryRow(query, user.Id).Scan(
		&newUser.Id,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Username,
		&isnullphoto,
		&isnullbio,
		&newUser.Email,
		&newUser.Gender,
		&newaddress,
		&newphone,
		&newUser.CreatedAt,
		&isupdated,
	)
	if err != nil {
		return nil, err
	}
	if isupdated.Valid {
		newUser.UpdatedAt = isupdated.String
	}
	if isnullbio.Valid {
		newUser.Bio = isnullbio.String
	}
	if isnullphoto.Valid {
		newUser.ProfilePhoto = isnullphoto.String
	}
	err = json.Unmarshal(newaddress, &newUser.Address)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(newphone, &newUser.PhoneNum)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r *UserRepo) Update(user *pb.UpdateUserReq) (*pb.UserRes, error) {
	query := `UPDATE assignees SET first_name=$1,last_name=$2,username=$3,profile_photo=$4,bio=$5,email=$6,gender=$7,address=$8,phone_num=$9,updated_at=$10 WHERE id=$11 RETURNING  id,first_name,last_name,username,profile_photo,bio,email,gender,address,phone_num,created_at,updated_at`
	newUser := pb.UserRes{}
	newaddress, err := json.Marshal(user.Address)
	if err != nil {
		return nil, err
	}
	newphone, err = json.Marshal(user.Phone)
	if err != nil {
		return nil, err
	}
	listUsers, err := r.List(&pb.Empty{})
	if err == nil {
		return nil, err
	}
	for _, val := range listUsers.Users {
		if val.Username == user.Username && val.Email == user.Email {
			return nil, err
		}
	}

	err = r.db.QueryRow(query, user.FirstName, user.LastName, user.Username, user.ProfilePhoto, user.Bio, user.Email, user.Gender, newaddress, newphone, time.Now().String(), user.Id).Scan(
		&newUser.Id,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Username,
		&isnullphoto,
		&isnullbio,
		&newUser.Email,
		&newUser.Gender,
		&newaddress,
		&newphone,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if isnullbio.Valid {
		newUser.Bio = isnullbio.String
	}
	if isnullphoto.Valid {
		newUser.ProfilePhoto = isnullphoto.String
	}
	err = json.Unmarshal(newaddress, &newUser.Address)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(newphone, &newUser.PhoneNum)
	if err != nil {
		return nil, err
	}
	return &newUser, nil
}

func (r *UserRepo) Delete(user *pb.GetOrDeleteUser) (*pb.ErrOrStatus, error) {
	query := `UPDATE assignees SET deleted_at=$1 WHERE id=$2`
	_, err := r.db.Exec(query, time.Now().Format(time.RFC3339), user.Id)
	if err != nil {
		return nil, err
	}

	return &pb.ErrOrStatus{
		Message: "Ok!",
	}, nil
}

func (r *UserRepo) List(user *pb.Empty) (*pb.UsersList, error) {
	query := `SELECT id,first_name,last_name,username,profile_photo,bio,email,gender,address,phone_num,created_at,updated_at FROM assignees WHERE deleted_at IS NULL`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	newUsers := pb.UsersList{}
	for rows.Next() {
		newUser := pb.UserRes{}
		err := rows.Scan(
			&newUser.Id,
			&newUser.FirstName,
			&newUser.LastName,
			&newUser.Username,
			&isnullphoto,
			&isnullbio,
			&newUser.Email,
			&newUser.Gender,
			&newaddress,
			&newphone,
			&newUser.CreatedAt,
			&isupdated,
		)
		if err != nil {
			return nil, err
		}
		if isupdated.Valid {
			newUser.UpdatedAt = isupdated.String
		}
		if isnullbio.Valid {
			newUser.Bio = isnullbio.String
		}
		if isnullphoto.Valid {
			newUser.ProfilePhoto = isnullphoto.String
		}
		err = json.Unmarshal(newaddress, &newUser.Address)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(newphone, &newUser.PhoneNum)
		if err != nil {
			return nil, err
		}
		newUsers.Users = append(newUsers.Users, &newUser)
	}
	return &newUsers, nil
}

func (r *UserRepo) ListEmail(all_id *pb.Ids) (*pb.Emails, error) {
	emails := pb.Emails{}
	for _, j := range all_id.Id {
		email := pb.Email{}
		query := `SELECT email FROM assignees WHERE id=$1`
		err := r.db.QueryRow(query, j.Id).Scan(
			&email.Email,
		)
		if err != nil {
			return nil, err
		}
		emails.Email = append(emails.Email, &email)
	}
	return &emails, nil
}

func (r *UserRepo) ChekUser(user *pb.EmailWithUsername) (*pb.Bool, bool) {
	username := `SELECT COUNT(*) FROM assignees WHERE username=$1 AND deleted_at IS NULL`
	count_username := 0
	err := r.db.QueryRow(username, user.Username).Scan(
		&count_username,
	)
	if err != nil {
		return &pb.Bool{Chekfild: true}, true
	}

	email := `SELECT COUNT(*) FROM assignees WHERE email=$1 AND deleted_at IS NULL`
	count_email := 0
	err = r.db.QueryRow(email, user.Email).Scan(
		&count_email,
	)
	if err != nil {
		return &pb.Bool{Chekfild: true}, true
	}
	if count_email == 0 && count_username == 0 {
		return &pb.Bool{Chekfild: false}, false
	} else {
		return &pb.Bool{Chekfild: true}, true
	}
}
func (r *UserRepo) UpdateToken(req *pb.TokensReq) (*pb.Tokens, error) {
	tokens := pb.Tokens{}
	query := `UPDATE assignees SET refresh_token = $1, accses_token = $2 WHERE id=$3 RETURNING accses_token,refresh_token`
	err := r.db.QueryRow(query, req.RefreshToken, req.AccessToken, req.Id).Scan(
		&tokens.AccessToken,
		&tokens.RefreshToken,
	)
	if err != nil {
		return nil, err
	}
	return &tokens, nil
}

func (r *UserRepo) Filtr(req *pb.FiltrReq) (*pb.UsersList, error) {
	sql := sqlbuilder.NewSelectBuilder()
	sql.Select("id", "first_name", "last_name", "username", "profile_photo", "bio", "email", "gender", "address", "phone_num", "created_at", "accses_token", "refresh_token").From("assignees").Where(sql.IsNull("deleted_at"))

	if v,ok := req.Filtr["username"]; ok {
		v = v + "%"
		sql.Where(sql.Like("username",v))
	}
	if v,ok := req.Filtr["email"]; ok {
		v = "%" + v + "%"
		sql.Where(sql.Like("email",v))
	}
	query, args := sql.BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := r.db.Query(query,args...)
	if err != nil {
		fmt.Println(err)
		fmt.Println(query,args)
		return nil, err
	}
	users := pb.UsersList{}
	for rows.Next() {
		newUser := pb.UserRes{}
		err = rows.Scan(
			&newUser.Id,
			&newUser.FirstName,
			&newUser.LastName,
			&newUser.Username,
			&isnullphoto,
			&isnullbio,
			&newUser.Email,
			&newUser.Gender,
			&newaddress,
			&newphone,
			&newUser.CreatedAt,
			&newUser.AccesToken,
			&newUser.RefreshToken,
		)
		if isnullbio.Valid {
			newUser.Bio = isnullbio.String
		}
		if isnullphoto.Valid {
			newUser.ProfilePhoto = isnullphoto.String
		}
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(newaddress, &newUser.Address)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(newphone, &newUser.PhoneNum)
		if err != nil {
			return nil, err
		}
		users.Users = append(users.Users, &newUser)
	}
	return &users, nil
}
