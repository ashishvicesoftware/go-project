package database

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
)

type UserProvider struct {
	db     *gorm.DB
	parent DB
}

func (a UserProvider) Create(user User) (User, error) {
	fmt.Println(user, "user")
	if user.ID == "" {
		return User{}, &invalidRequest{"create user", "id must be included"}
	}
	if db := a.db.Create(&user); db.Error != nil {
		return User{}, db.Error
	}
	return user, nil
}

func (a UserProvider) Get(id string) (User, error) {

	// fetched_ID, err := strconv.Unquote(id)

	users := User{ID: id}

	// fmt.Println(err, "error")

	if db := a.db.Where(User{ID: id}).Find(&users); db.Error != nil {
		if IsNotFound(db.Error) {
			return User{}, &recordNotFound{"get contact", 0}
		}
		return User{}, db.Error
	}
	return users, nil
}

func (a UserProvider) GetAll() ([]User, error) {
	user := make([]User, 0)
	if db := a.db.Order("id").Find(&user); db.Error != nil {
		return nil, db.Error
	}
	return user, nil
}

func (a UserProvider) GetToUpdate(id string) (User, error) {
	fmt.Println(id, "project detail")
	users := User{ID: id}
	if db := a.db.Where(User{ID: id}).Take(&users); db.Error != nil {
		if IsNotFound(db.Error) {
			return User{}, &recordNotFound{"get contact", 0}
		}
		fmt.Println(users, "GetProject")
		return User{}, db.Error
	}
	return users, nil
}

func (a UserProvider) GetUserByEmail(email string) (User, error) {

	userData := User{Email: email}
	if db := a.db.Where(User{Email: email}).Take(&userData); db.Error != nil {
		if IsNotFound(db.Error) {
			return User{}, &recordNotFound{"email not found", 0}
		}
		return User{}, db.Error
	}
	return userData, nil
}

func (a UserProvider) Update(user User) (User, error) {
	fmt.Println(user.ID, "user.ID")
	existing, err := a.GetToUpdate(user.ID)
	if err != nil {
		if IsNotFound(err) {
			return User{}, &recordNotFound{"update note", 0}
		}
		return User{}, err

	}

	if user.Name == "" {
		user.Name = existing.Name
	}

	if user.Email == "" {
		user.Email = existing.Email
	}

	if user.ImageUrl == "" {
		user.ImageUrl = existing.ImageUrl
	}

	user.IsActive = existing.IsActive

	if user.CreatedAt.IsZero() {
		user.CreatedAt = existing.CreatedAt
	}

	if db := a.db.Save(&user); db.Error != nil {
		return User{}, db.Error
	}
	return user, nil
}

func (a UserProvider) Delete(id string) error {

	fetched_ID, err := strconv.Unquote(id)

	fmt.Println(err, "error")

	db := a.db.Delete(&User{ID: fetched_ID})
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return &recordNotFound{"delete project", 0}
	}
	return nil
}

func (a UserProvider) UpdateUserDetails(user User) (User, error) {
	existing, err := a.GetToUpdate(user.ID)
	if err != nil {
		if IsNotFound(err) {
			return User{}, &recordNotFound{"update Users", 0}
		}
		return User{}, err

	}

	if user.ImageUrl == "" {
		user.ImageUrl = existing.ImageUrl
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = existing.CreatedAt
	}

	// if user.doIntro == false {
	// 	user.doIntro = existing.doIntro
	// }
	if user.IsActive == false {
		user.IsActive = existing.IsActive
	}

	if db := a.db.Save(&user); db.Error != nil {
		return User{}, db.Error
	}
	return user, nil
}
