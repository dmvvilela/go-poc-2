package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

type Contact struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;" json:"name"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Contact) Prepare() {
	c.ID = 0
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.Email = html.EscapeString(strings.TrimSpace(c.Email))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Contact) Validate() error {
	if c.Name == "" {
		return errors.New("Required Name")
	}
	if c.Email == "" {
		return errors.New("Required Email")
	}
	if err := checkmail.ValidateFormat(c.Email); err != nil {
		return errors.New("Invalid Email")
	}
	return nil
}

func (c *Contact) SaveContact(db *gorm.DB) (*Contact, error) {
	var err error
	err = db.Debug().Create(&c).Error
	if err != nil {
		return &Contact{}, err
	}
	return c, nil
}

func (c *Contact) FindAllContacts(db *gorm.DB) (*[]Contact, error) {
	var err error
	contacts := []Contact{}
	err = db.Debug().Model(&Contact{}).Limit(100).Find(&contacts).Error
	if err != nil {
		return &[]Contact{}, err
	}
	return &contacts, err
}

func (c *Contact) FindContactByID(db *gorm.DB, cid uint64) (*Contact, error) {
	var err error
	err = db.Debug().Model(Contact{}).Where("id = ?", cid).Take(&c).Error
	if err != nil {
		return &Contact{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Contact{}, errors.New("Contact Not Found")
	}
	return c, err
}

func (c *Contact) UpdateContact(db *gorm.DB, cid uint32) (*Contact, error) {
	db = db.Debug().Model(&Contact{}).Where("id = ?", cid).Take(&Contact{}).UpdateColumns(
		map[string]interface{}{
			"name":       c.Name,
			"email":      c.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Contact{}, db.Error
	}

	// This is the display the updated user
	err := db.Debug().Model(&Contact{}).Where("id = ?", cid).Take(&c).Error
	if err != nil {
		return &Contact{}, err
	}
	return c, nil
}

func (c *Contact) DeleteContact(db *gorm.DB, cid uint64) (int64, error) {
	db = db.Debug().Model(&Contact{}).Where("id = ?", cid).Take(&Contact{}).Delete(&Contact{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
