package models

import "free_cms/common/models"

type BooksPreg struct {
	models.BooksPreg
}

func NewBooksPreg()(booksPreg *BooksPreg){
	return &BooksPreg{}
}

func (m *BooksPreg)FindById(id int)(booksPreg BooksPreg,err error){
	err = models.Db.Where("id=?",id).First(&booksPreg).Error
	return
}