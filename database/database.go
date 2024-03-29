package database

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/iotames/easyim/config"

	"github.com/bwmarrin/snowflake"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iotames/miniutils"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

const (
	WHERE_COMPARE_EQUAL   = "="
	WHERE_COMPARE_LIKE    = "LIKE"
	WHERE_COMPARE_BETWEEN = "BETWEEN"
	WHERE_COMPARE_IN      = "IN"
	WHERE_LINK_AND        = "AND"
	WHERE_LINK_OR         = "OR"
)

var (
	once   sync.Once
	engine *xorm.Engine
	snode  *snowflake.Node
)

func getNodeId() int64 {
	d := config.GetDatabase()
	return int64(d.NodeID)
}

func getEngine() *xorm.Engine {
	if engine != nil {
		return engine
	}
	once.Do(func() {
		engine = newEngine(config.GetDatabase())
	})
	return engine
}

func SetEngine(db config.Database) {
	engine = newEngine(db)
}

func newEngine(db config.Database) *xorm.Engine {
	log.Println("New newEngine Begin")
	var err error
	if db.Driver == config.DRIVER_SQLITE3 {
		engine, err = xorm.NewEngine(db.Driver, config.SQLITE_FILENAME)
	} else {
		engine, err = xorm.NewEngine(db.Driver, db.GetDSN())
	}
	if err != nil {
		panic(err)
	}
	engineInit(engine)
	log.Println("New newEngine End")
	return engine
}

func engineInit(engine *xorm.Engine) {
	log.Println("Init engineInit Begin")
	ormMap := names.GonicMapper{}
	engine.SetMapper(ormMap)
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	engine.DatabaseTZ, _ = time.LoadLocation("Asia/Shanghai")
	engine.SetTableMapper(ormMap)
	engine.SetColumnMapper(ormMap)
	engine.ShowSQL(true)
	log.Println("Init engineInit End")
}

func GetSnowflakeNode() *snowflake.Node {
	if snode == nil {
		node, err := snowflake.NewNode(getNodeId())
		if err != nil {
			logger := miniutils.GetLogger("")
			logger.Error("Error for database.getSnowflakeNode:", err)
			snode = nil
		}
		snode = node
	}
	log.Println("---getSnowflakeNode---", snode)
	return snode
}

type IDitem interface {
	ParseID() snowflake.ID
	GetID() int64
}

type IModel interface {
	GenerateID() int64
	IDitem
}

type BaseModel struct {
	// TODO 分布式ID 雪花算法 https://www.itqiankun.com/article/1565747019
	ID        int64     `xorm:"pk unique"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func (b *BaseModel) GenerateID() int64 {
	if b.ID == 0 {
		id := GetSnowflakeNode().Generate().Int64()
		if id == 0 {
			panic("Error: getSnowflakeNode().Generate().Int64() == 0")
		}
		b.ID = id
	}
	return b.ID
}

func (b BaseModel) ParseID() snowflake.ID {
	return snowflake.ParseInt64(b.ID)
}

func (b BaseModel) GetID() int64 {
	return b.ID
}

func (b BaseModel) ToMap(m IModel) map[string]interface{} {
	typeof := reflect.TypeOf(m).Elem()
	typevalue := reflect.ValueOf(m).Elem()
	fieldLen := typeof.NumField()
	fieldsMap := make(map[string]interface{}, fieldLen+2)
	for i := 0; i < fieldLen; i++ {
		field := typeof.Field(i)
		fvalue := typevalue.Field(i)
		value := fvalue.Interface()
		if field.Name == "BaseModel" {
			for j := 0; j < field.Type.NumField(); j++ {
				fieldj := field.Type.Field(j)
				fvaluej := fvalue.Field(j)
				valuej := fvaluej.Interface()
				if fieldj.Name == "ID" {
					valuej = fmt.Sprintf("%d", valuej)
				}
				fieldsMap[fieldj.Name] = valuej
			}
		} else {
			if strings.Contains(field.Name, "ID") {
				value = fmt.Sprintf("%d", value)
			}
			fieldsMap[field.Name] = value
		}
	}
	return fieldsMap
}

func TableColToObj(t string) string {
	tmp := (names.GonicMapper{}).Table2Obj(t)
	replaceMap := map[string]string{"Id": "ID"}
	for k, v := range replaceMap {
		keyIndex := strings.Index(tmp, k)
		lastIndex := len(tmp) - 2 // 搜索词在末尾
		if k == "Id" && lastIndex == keyIndex {
			tmp = strings.ReplaceAll(tmp, k, v)
		}
	}
	return tmp
}
func ObjToTableCol(o string) string {
	return (names.GonicMapper{}).Obj2Table(o)
}

func CreateTables() {
	err := getEngine().CreateTables(getAllTables()...)
	if err != nil {
		panic(fmt.Errorf("error for database.CreateTables:%v", err))
	}
}

func SyncTables() {
	err := getEngine().Sync(getAllTables()...)
	if err != nil {
		panic(fmt.Errorf("error for database.SyncTables:%v", err))
	}
}

func GetModel(m IModel) (bool, error) {
	b, err := getEngine().Get(m)
	if err != nil {
		logger := miniutils.GetLogger("")
		logger.Error("Error for database.GetModel:", err)
	}
	return b, err
}
func Query(sqlOrArgs ...interface{}) (resultsSlice []map[string][]byte, err error) {
	return getEngine().Query(sqlOrArgs...)
}
func Exec(sqlOrArgs ...interface{}) (int64, error) {
	result, err := getEngine().Exec(sqlOrArgs...)
	rowsNum, _ := result.RowsAffected()
	return rowsNum, err
}

// GetModelWhere 添加复杂条件. 参数 m IModel 各属性必须为零值，否则查询条件会冲突
// GetModelWhere(new(User), "name = ? AND age = ?", "Tom", 19)
func GetModelWhere(m IModel, query interface{}, args ...interface{}) (bool, error) {
	b, err := getEngine().Where(query, args...).Get(m)
	if err != nil {
		logger := miniutils.GetLogger("")
		logger.Error("Error for database.GetModel:", err)
	}
	return b, err
}

// 转化map为Like条件
func GetWhereLikeArgs(filter map[string]string) (query interface{}, args []interface{}) {
	q := ""
	i := 0
	for k, v := range filter {
		if strings.TrimSpace(v) == "" {
			continue
		}
		args = append(args, `%`+v+`%`)
		field := ObjToTableCol(k)
		qOne := fmt.Sprintf("`%s` LIKE ?", field)
		if i > 0 {
			q += fmt.Sprintf(" AND %s", qOne)
		} else {
			q += qOne
		}
		i++
	}
	query = q
	return
}

func GetWhereOne(field, compare string, v interface{}) string {
	field = ObjToTableCol(field)
	result := fmt.Sprintf(`%s %s`, field, compare)
	switch v.(type) {
	case string:
		if strings.TrimSpace(v.(string)) != "" {
			val := v.(string)
			if compare == WHERE_COMPARE_LIKE {
				val = `'%` + val + `%'`
			}
			result += " " + val
		}
	case []string:
		if len(v.([]string)) > 0 {
			vals := v.([]string)
			if compare == WHERE_COMPARE_IN {
				val := "("
				for i, inv := range vals {
					if i < (len(vals) - 1) {
						val += fmt.Sprintf(`'%s',`, inv)
					} else {
						val += fmt.Sprintf(`'%s'`, inv)
					}
				}
				val += ")"
				result += " " + val
			}
			if compare == WHERE_COMPARE_BETWEEN {
				result += fmt.Sprintf(" %s AND %s", vals[0], vals[1])
			}
		}
	}
	return result
}

// GetAll 获取多条记录
// users := make([]Userinfo, 0)
// GetAll(&users, 50, 3, "age > ? or name = ?", 30, "xlw")
//
// GetAll(&users, 50, 3, map[string]interface{}{"Name": "jinzhu", "Age": 0})
func GetAll(rows interface{}, limit, page int, query interface{}, args ...interface{}) error {
	start := (page - 1) * limit
	err := getEngine().Where(query, args...).Limit(limit, start).Find(rows)
	if err != nil {
		logger := miniutils.GetLogger("")
		logger.Error("Error for database.GetAll:", err)
	}
	return err
}

func CreateModel(m IModel) (int64, error) {
	m.GenerateID()
	return getEngine().Insert(m)
}

// UpdateModel 更新数据
// dt参数指定更新的字段，字段名用数据库中的字段名，不用go结构体字段名，如 ExecutedAt -> executed_at
func UpdateModel(m IModel, dt map[string]interface{}) (int64, error) {
	modelID := m.GetID() // m.ParseID().Int64()
	if dt == nil {
		return getEngine().ID(modelID).Update(m)
	}
	return getEngine().Table(m).ID(modelID).Update(dt)
}

func DeleteModel(m IModel) (int64, error) {
	return getEngine().Delete(m)
}

func BatchDelete(m IModel, codes []string) (int64, error) {
	return getEngine().In("ID", codes).Delete(m)
}

func BatchUpdate(m IModel, codes []string) (int64, error) {
	return getEngine().In("ID", codes).Update(m)
}

func NewSession() *xorm.Session {
	return getEngine().NewSession()
}
