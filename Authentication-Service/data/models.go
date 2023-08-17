package data

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeOut = 3 * time.Second

var db *sql.DB

type Models struct {
	User User
}

// User is the structure which holds one user from the database.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"user_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New is function used to create an instance of the data package.It return the type
// Model, which embed all the types we want to be available to our application

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User: User{},
	}
}

func (user *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)

	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.Active, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`
	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()
	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at where id = $1`
	row := db.QueryRowContext(ctx, query, id)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.Active, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()
	query := `update users set first_name = $1, last_name = $2, password = $3, user_active = $4, updated_at = $5 where id = $6`

	_, err := db.ExecContext(ctx, query, u.FirstName, u.LastName, u.Password, u.Active, time.Now(), u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()
	query := `delete from users where id = $1`
	_, err := db.ExecContext(ctx, query, u.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete by id
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()
	query := `delete from users where id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// inders a new user into the database , and return Id of the newly created user instead of the user

func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()
	var newId int
	query := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)`
	err := db.QueryRowContext(ctx, query, user.Email, user.FirstName, user.LastName, user.Password, user.Active, time.Now(), time.Now()).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

// ResetPassword is the method we will use to change a user's password
func (u *User) ResetPassword(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()
	hashPassword, err := bcrypt.GenerateFromPassword([]byte("password"), 12)
	query := `update users set password = $1 where id = $2`
	_, err = db.ExecContext(ctx, query, hashPassword, id)
	if err != nil {
		return err
	}
	return nil
}

// MachePassword is the method compares a password with the one in the database
func (u *User) MachePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
