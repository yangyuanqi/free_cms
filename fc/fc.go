package fc

import (
	"flag"
	"fmt"
	"free_cms/models"
	"github.com/astaxie/beego"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var TplModel = `package models

type Config struct {
	Model
	$attr$
}


func NewConfig() (config *Config) {
	return &Config{}
}

func (m *Config) Pagination(offset, limit int, key string) (res []Config, count int) {
	query := models.Db
	if key != "" {
		query = query.Where("name like ?", "%"+key+"%")
	}
	query.Offset(offset).Limit(limit).Order("id desc").Find(&res)
	query.Model(Config{}).Count(&count)
	return
}

func (m *Config) Create() (err error, newAttr *Config) {
	err = models.Db.Create(m).Error
	newAttr = m
	return
}

func (m *Config) Update() (err error, newAttr Config) {
	if m.Id > 0 {
		err = models.Db.Where("id=?", m.Id).Save(m).Error
	} else {
		err = errors.New("id参数错误")
	}
	newAttr = *m
	return
}

func (m *Config) Delete() (err error) {
	if m.Id > 0 {
		err = models.Db.Delete(m).Error
	} else {
		err = errors.New("id参数错误")
	}
	return
}

func (m *Config) DelBath(ids []int) (err error) {
	if len(ids) > 0 {
		err = models.Db.Where("id in (?)", ids).Delete(m).Error
	} else {
		err = errors.New("id参数错误")
	}
	return
}

func (m *Config) FindById(id int) (config Config, err error) {
	err = Db.Where("id=?", id).First(&config).Error
	return
}

`

type TabelAttr struct {
	Field   string `gorm:"column:Field"`
	Type    string `gorm:"column:Type"`
	Null    string `gorm:"column:Null"`
	Key     string `gorm:"column:Key"`
	Default string `gorm:"column:Default"`
	Extra   string `gorm:"column:Extra"`
}



func main() {
	tableName := flag.String("t", "", "表名")
	modelPath := flag.String("m", "", "模型地址")
	controllerPath := flag.String("c", "", "控制器地址")
	viewPath := flag.String("v", "", "视图地址")
	flag.Parse()

	tablePrefix := beego.AppConfig.String("tablePrefix")
	var tableAttr []TabelAttr
	models.Db.Raw("desc " + tablePrefix + *tableName).Scan(&tableAttr)

	//获取数据
	if *modelPath != "" {
		if err := os.MkdirAll(path.Clean(*modelPath), 777); err != nil {
			log.Println("模型文件创建失败")
		}
		modelData := createModel(tableAttr, *tableName, *modelPath, TplModel)
		if err := ioutil.WriteFile(path.Join(*modelPath, "models", fmt.Sprintf("%s.go", *tableName)), []byte(modelData), os.ModeType); err != nil {
			log.Println(err)
		}
	}

	var controllerData string
	if *controllerPath != "" {
		controllerData = createModel(tableAttr, *tableName, *modelPath, "")

		if err := os.MkdirAll(*controllerPath, 777); err != nil {
			log.Println("控制器文件创建失败")
		}

		if err := ioutil.WriteFile(path.Join(*modelPath, fmt.Sprintf("%s.go", *tableName)), []byte(controllerData), os.ModeType); err != nil {
			log.Println(err)
		}
	}

	var viewData string
	if *viewPath != "" {
		viewData = createModel(tableAttr, *tableName, *modelPath, "")

		if err := os.MkdirAll(*viewPath, 777); err != nil {
			log.Println("视图文件创建失败")
		}

		if err := ioutil.WriteFile(path.Join(*modelPath, fmt.Sprintf("%s.go", *tableName)), []byte(viewData), os.ModeType); err != nil {
			log.Println(err)
		}
	}

}

func createModel(tableAttr []TabelAttr, tableName, modelPath string, tpl string) (tplModel string) {
	//字段名称
	var attr string
	for _, v := range tableAttr {
		fieldName := Hump(v.Field, "max")
		//类型
		typeName := getType(v.Type)
		f1 := fieldName + strings.Repeat(" ", 20-len(fieldName))
		f2 := typeName + strings.Repeat(" ", 10-len(typeName))
		attr += f1 + f2 + fmt.Sprintf("`json:\"%s\"%sform:\"%s\"%sgorm:\"default:'%s'\"`\n	",
			v.Field,
			strings.Repeat(" ", 15-len(v.Field)),
			v.Field,
			strings.Repeat(" ", 15-len(v.Field)),
			v.Default,
		)
	}

	//替换字符
	tpl = strings.Replace(tpl, "package models", "package "+path.Clean(modelPath), -1)
	tpl = strings.Replace(tpl, "Config", Hump(tableName, "max"), -1)
	tpl = strings.Replace(tpl, "config", Hump(tableName, "min"), -1)
	tplModel = strings.Replace(tpl, "$attr$", attr, -1)
	return
}

//字段类型转换
func getType(typeName string) (str string) {
	if strings.Index(typeName, "bigint") >= 0 || strings.Index(typeName, "int") >= 0 || strings.Index(typeName, "tinyint") >= 0 {
		return "int"
	}
	if strings.Index(typeName, "varchar") >= 0 || strings.Index(typeName, "char") >= 0 {
		return "string"
	}
	if strings.Index(typeName, "datetime") >= 0 || strings.Index(typeName, "time") >= 0 {
		return "time.Time"
	}
	return "string"
}

//字段转换大驼峰，小驼峰
func Hump(v, t string) (new string) {
	field := strings.Split(v, "_")
	if t == "min" {
		for k, v := range field {
			if k > 0 {
				field[k] = strings.Title(v)
			}
		}
	}
	if t == "max" {
		for k, v := range field {
			field[k] = strings.Title(v)
		}
	}
	new = strings.Join(field, "")
	return
}