package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type BooksPreg struct {
	Model
	Name       string    `json:"name"                           form:"name"`
	ListAuthor string `json:"list_author"                       form:"list_author"`
	ListDescribe string `json:"list_describe"                   form:"list_describe"`
	ListMsgImg string `json:"list_msg_img"                      form:"list_msg_img"`
	ListA string `json:"list_a"                                 form:"list_a"`
	ContentText string `json:"content_text"                     form:"content_text"`
	ListMsgLastTime string `json:"list_msg_last_time"           form:"list_msg_last_time"`
	ListAuthorFilter string `json:"list_author_filter"          form:"list_author_filter"`
	ListDescribeFilter string `json:"list_describe_filter"      form:"list_describe_filter"`
	ContentTextFilter string `json:"content_text_filter"        form:"content_text_filter"`

	CreatedAtText int64 `json:"create_at_text" gorm:"-"`
}


func NewBooksPreg() (books *BooksPreg) {
	return &BooksPreg{}
}

func (m *BooksPreg) AfterFind(scope *gorm.Scope) (err error) {
	m.CreatedAtText = m.CreatedAt.Unix()
	return
}

func (m *BooksPreg) Pagination(offset, limit int, key string) (res []BooksPreg, count int) {
	query := Db
	if key != "" {
		query = query.Where("name like ?", "%"+key+"%")
	}
	query.Offset(offset).Limit(limit).Order("id desc").Find(&res)
	query.Model(BooksPreg{}).Count(&count)
	return
}

func (m *BooksPreg) Create() (err error, newAttr *BooksPreg) {
	err = Db.Create(m).Error
	newAttr = m
	return
}

func (m *BooksPreg) Update() (err error, newAttr BooksPreg) {
	if m.Id > 0 {
		err = Db.Where("id=?", m.Id).Save(m).Error
	} else {
		err = errors.New("id参数错误")
	}
	newAttr = *m
	return
}

func (m *BooksPreg) Delete() (err error) {
	if m.Id > 0 {
		err = Db.Delete(m).Error
	} else {
		err = errors.New("id参数错误")
	}
	return
}

func (m *BooksPreg)DelBatch(ids []int)(err error){
	if len(ids)>0{
		err = Db.Where("id in (?)",ids).Delete(m).Error
	}else {
		err = errors.New("id参数错误")
	}
	return
}

func (m *BooksPreg)FindById(id int)(booksPreg BooksPreg,err error){
	err = Db.Where("id=?",id).First(&booksPreg).Error
	return
}