package user

// UseCase Define the interface with functions that will be used
type UserUseCase interface {
	GetAll() ([]*User, error)
	Get(ID int64) (User, error)
	Store(u *User) error
	Update(u *User) error
	Remove(ID int64) error
}

// UserService define the struct for user service
type UserService struct{}

func (*us *UserService) GetAll() ([]*User, error) {
	return nil, nil
}

func (*us *UserService) Get(ID int64) (User, error) {
	return nil, nil
}

func (s *UserService) Store(b *User) error {
	return nil
}

func (s *UserService) Update(b *User) error {
	return nil
}

func (s *UserService) Remove(ID int64) error {
	return nil
}