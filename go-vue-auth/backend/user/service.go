package user

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service interface {
	Register(name string, email string, password string) (*User, error)
	Login(email string, password string) (*User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Register(
	name string,
	email string,
	password string,
) (*User, error) {

	email = strings.TrimSpace(strings.ToLower(email))
	name = strings.TrimSpace(name)

	existingUser, err := s.repository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	newUser := &User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.repository.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *service) Login(
	email string,
	password string,
) (*User, error) {

	email = strings.TrimSpace(strings.ToLower(email))

	foundUser, err := s.repository.FindByEmail(email)

	if err != nil {
		return nil, err
	}

	if foundUser == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.Password),
		[]byte(password),
	)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return foundUser, nil
}
