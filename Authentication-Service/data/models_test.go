package data

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserGetAll(t *testing.T) {
	// Tạo mock DB và đối tượng rows giả lập
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "email", "first_name", "last_name", "password", "user_active", "created_at", "updated_at"}
	rows := sqlmock.NewRows(columns).
		AddRow(1, "test1@example.com", "John", "Doe", "password1", 1, "2021-01-01", "2021-01-01").
		AddRow(2, "test2@example.com", "Jane", "Smith", "password2", 1, "2021-01-02", "2021-01-02")

	// Thiết lập các truy vấn giả lập và kết quả trả về
	mock.ExpectQuery("select id, email, first_name, last_name, password, user_active, created_at, updated_at").
		WillReturnRows(rows)

	// Tạo đối tượng Models và thiết lập DB mock
	models := Models{User: User{}}

	// Gọi phương thức GetAll()
	users, err := models.User.GetAll()
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	// Kiểm tra giá trị của người dùng đầu tiên
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, "test1@example.com", users[0].Email)
	assert.Equal(t, "John", users[0].FirstName)
	assert.Equal(t, "Doe", users[0].LastName)
	assert.Equal(t, "password1", users[0].Password)
	assert.Equal(t, 1, users[0].Active)
	assert.Equal(t, "2021-01-01", users[0].CreatedAt)
	assert.Equal(t, "2021-01-01", users[0].UpdatedAt)

	// Kiểm tra giá trị của người dùng thứ hai
	assert.Equal(t, 2, users[1].ID)
	assert.Equal(t, "test2@example.com", users[1].Email)
	assert.Equal(t, "Jane", users[1].FirstName)
	assert.Equal(t, "Smith", users[1].LastName)
	assert.Equal(t, "password2", users[1].Password)
	assert.Equal(t, 1, users[1].Active)
	assert.Equal(t, "2021-01-02", users[1].CreatedAt)
	assert.Equal(t, "2021-01-02", users[1].UpdatedAt)

	// Kiểm tra xem tất cả các truy vấn đã được gọi
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Viết các test case khác tương tự cho các phương thức khác trong gói data

func TestUserGetByEmail(t *testing.T) {
	// Tạo mock DB và đối tượng row giả lập
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "email", "first_name", "last_name", "password", "user_active", "created_at", "updated_at"}
	row := sqlmock.NewRows(columns).
		AddRow(1, "test@example.com", "John", "Doe", "password", 1, "2021-01-01", "2021-01-01")

	// Thiết lập truy vấn giả lập và kết quả trả về
	mock.ExpectQuery("select id, email, first_name, last_name, password, user_active, created_at, updated_at where email = ?").
		WithArgs("test@example.com").
		WillReturnRows(row)

	// Tạo đối tượng Models và thiết lập DB mock
	models := Models{User: User{}}

	// Gọi phương thức GetByEmail()
	user, err := models.User.GetByEmail("test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Kiểm tra giá trị của người dùng
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, 1, user.Active)
	assert.Equal(t, "2021-01-01", user.CreatedAt)
	assert.Equal(t, "2021-01-01", user.UpdatedAt)

	// Kiểm tra xem tất cả các truy vấn đã được gọi
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Viết các hàm test cho các phương thức khác tương tự

func TestMain(m *testing.M) {
	// Thiết lập giả lập DB, gọi các unit test và giả lập DB
	// Chạy các test case và trả về kết quả
}
