package servicework

import (
	"context"
	"errors"
	"testing"
	"time"

	"a1337b04rd/internal/domain/entity"
)

// --- Мок UserRepository ---

type mockUserRepo struct {
	createErr error
	created   *entity.User
	getErr    error
	userByID  *entity.User
	updateErr error
	updated   *entity.User
	deleteErr error
	deletedID int
}

func (m *mockUserRepo) CreateUser(ctx context.Context, u *entity.User) (*entity.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	// копируем для проверки
	*m.created = *u
	return u, nil
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.userByID, nil
}

func (m *mockUserRepo) UpdateUser(ctx context.Context, u *entity.User) (*entity.User, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	// копируем
	*m.updated = *u
	return u, nil
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id int) error {
	m.deletedID = id
	return m.deleteErr
}

// --- Тесты userService ---

func TestCreateUser_Success(t *testing.T) {
	stored := &entity.User{}
	mr := &mockUserRepo{created: stored}
	svc := NewUserService(mr)

	input := entity.User{CharacterName: "Alice"}
	out, err := svc.CreateUser(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateUser вернул ошибку: %v", err)
	}
	// repo получил правильные данные
	if stored.CharacterName != "Alice" {
		t.Errorf("ожидали CharacterName 'Alice', получили '%s'", stored.CharacterName)
	}
	// CreatedAt установлен
	if out.CreatedAt.IsZero() {
		t.Error("ожидали, что CreatedAt установлен, получили zero")
	}
}

func TestCreateUser_MissingName(t *testing.T) {
	mr := &mockUserRepo{}
	svc := NewUserService(mr)

	_, err := svc.CreateUser(context.Background(), entity.User{})
	if err == nil {
		t.Fatal("ожидали ошибку, получили nil")
	}
	expected := "character name is required"
	if err.Error() != expected {
		t.Errorf("ожидали '%s', получили '%v'", expected, err)
	}
}

func TestCreateUser_RepoError(t *testing.T) {
	stored := &entity.User{}
	mr := &mockUserRepo{created: stored, createErr: errors.New("db fail")}
	svc := NewUserService(mr)

	_, err := svc.CreateUser(context.Background(), entity.User{CharacterName: "Bob"})
	if err == nil {
		t.Fatal("ожидали ошибку, получили nil")
	}
	if err.Error() != "db fail" {
		t.Errorf("ожидали 'db fail', получили '%v'", err)
	}
}

func TestGetUserByID_Success(t *testing.T) {
	now := time.Now()
	user := &entity.User{ID: 5, CharacterName: "Carol", CreatedAt: now}
	mr := &mockUserRepo{userByID: user}
	svc := NewUserService(mr)

	out, err := svc.GetUserByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetUserByID вернул ошибку: %v", err)
	}
	if out.ID != 5 || out.CharacterName != "Carol" {
		t.Errorf("ожидали %+v, получили %+v", user, out)
	}
}

func TestGetUserByID_Error(t *testing.T) {
	mr := &mockUserRepo{getErr: errors.New("not found")}
	svc := NewUserService(mr)

	_, err := svc.GetUserByID(context.Background(), 1)
	if err == nil {
		t.Fatal("ожидали ошибку, получили nil")
	}
	if err.Error() != "not found" {
		t.Errorf("ожидали 'not found', получили '%v'", err)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	updated := &entity.User{}
	mr := &mockUserRepo{updated: updated}
	svc := NewUserService(mr)

	input := entity.User{ID: 7, CharacterName: "Dave"}
	out, err := svc.UpdateUser(context.Background(), input)
	if err != nil {
		t.Fatalf("UpdateUser вернул ошибку: %v", err)
	}
	// repo получил правильные данные
	if updated.ID != 7 || updated.CharacterName != "Dave" {
		t.Errorf("репозиторий получил %+v", updated)
	}
	// сервис вернул те же данные
	if out.ID != 7 || out.CharacterName != "Dave" {
		t.Errorf("ожидали %+v, получили %+v", input, out)
	}
}

func TestUpdateUser_Error(t *testing.T) {
	updated := &entity.User{}
	mr := &mockUserRepo{updated: updated, updateErr: errors.New("upd fail")}
	svc := NewUserService(mr)

	_, err := svc.UpdateUser(context.Background(), entity.User{ID: 8})
	if err == nil {
		t.Fatal("ожидали ошибку, получили nil")
	}
	if err.Error() != "upd fail" {
		t.Errorf("ожидали 'upd fail', получили '%v'", err)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	mr := &mockUserRepo{}
	svc := NewUserService(mr)

	if err := svc.DeleteUser(context.Background(), 9); err != nil {
		t.Errorf("DeleteUser вернул ошибку: %v", err)
	}
	if mr.deletedID != 9 {
		t.Errorf("ожидали удаление ID=9, получили %d", mr.deletedID)
	}
}

func TestDeleteUser_Error(t *testing.T) {
	mr := &mockUserRepo{deleteErr: errors.New("del fail")}
	svc := NewUserService(mr)

	err := svc.DeleteUser(context.Background(), 10)
	if err == nil {
		t.Fatal("ожидали ошибку, получили nil")
	}
	if err.Error() != "del fail" {
		t.Errorf("ожидали 'del fail', получили '%v'", err)
	}
}
