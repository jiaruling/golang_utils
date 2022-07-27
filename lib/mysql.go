package lib

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

// https://gorm.io/zh_CN/docs/index.html
// https://www.topgoer.com/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%93%8D%E4%BD%9C/go%E6%93%8D%E4%BD%9Cmysql/mysql%E4%BD%BF%E7%94%A8.html

var mysqlGormMap map[string]*gorm.DB
var mysqlXMap map[string]*sqlx.DB

type MySQL struct {
	Name            string
	User            string
	Password        string
	Ip              string
	Port            int
	DB              string
	Parameter       string        // default charset=utf8mb4&parseTime=True&loc=Local
	MaxIdleConns    int           // default 5
	MaxOpenConns    int           // default 100
	ConnMaxLifetime time.Duration // default time.Hour
}

type MysqlGorm struct {
	MySQL
}

func NewMysqlGorm(user, password, ip string, port int, db string) *MysqlGorm {
	if mysqlGormMap == nil {
		mysqlGormMap = make(map[string]*gorm.DB)
	}

	return &MysqlGorm{
		MySQL{
			User:     user,
			Password: password,
			Ip:       ip,
			Port:     port,
			DB:       db,
		},
	}
}

func (m *MysqlGorm) InitMysqlGorm() {
	if mysqlGormMap == nil {
		mysqlGormMap = make(map[string]*gorm.DB)
	}
	if m.Name == "" {
		m.Name = "default"
	}
	if _, ok := mysqlGormMap[m.Name]; ok {
		return
	}
	if m.Parameter == "" {
		m.Parameter = "charset=utf8mb4&parseTime=True&loc=Local"
	}
	DNS := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?%s", m.User, m.Password, m.Ip, m.Port, m.DB, m.Parameter)
	conn, err := gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		log.Fatalln("open mysql failed,", err.Error())
		return
	}
	// 设置数据库连接池
	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatalln("设置数据库连接池失败,", err.Error())
		return
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	if m.MaxIdleConns == 0 {
		m.MaxIdleConns = 5
	}
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	if m.MaxOpenConns == 0 {
		m.MaxOpenConns = 100
	}
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	if m.ConnMaxLifetime == 0 {
		m.ConnMaxLifetime = time.Hour
	}
	sqlDB.SetConnMaxLifetime(m.ConnMaxLifetime)
	mysqlGormMap[m.Name] = conn
}

func GetMysqlGorm(name ...string) *gorm.DB {
	if mysqlGormMap == nil {
		return nil
	}
	var n string
	if len(name) > 0 {
		n = name[0]
	} else {
		n = "default"
	}
	if _, ok := mysqlGormMap[n]; ok {
		return mysqlGormMap[n]
	} else {
		return nil
	}
}

func DestroyMysqlGorm(name ...string) {
	if mysqlGormMap == nil {
		return
	}
	var n string
	if len(name) > 0 {
		n = name[0]
	} else {
		n = "default"
	}
	if _, ok := mysqlGormMap[n]; ok {
		delete(mysqlGormMap, n)
	}
}

func DestroyMysqlGormAll() {
	if mysqlGormMap == nil {
		return
	}
	for k, _ := range mysqlGormMap {
		delete(mysqlGormMap, k)
	}
}

// https://github.com/jmoiron/sqlx

type MysqlX struct {
	MySQL
}

func NewMysqlX(user, password, ip string, port int, db string) *MysqlX {
	if mysqlXMap == nil {
		mysqlXMap = make(map[string]*sqlx.DB)
	}
	return &MysqlX{
		MySQL{
			User:     user,
			Password: password,
			Ip:       ip,
			Port:     port,
			DB:       db,
		},
	}
}

func (m *MysqlX) InitMysqlX() {
	if mysqlXMap == nil {
		mysqlXMap = make(map[string]*sqlx.DB)
	}
	if m.Name == "" {
		m.Name = "default"
	}
	if _, ok := mysqlXMap[m.Name]; ok {
		return
	}
	if m.Parameter == "" {
		m.Parameter = "charset=utf8mb4&parseTime=True&loc=Local"
	}
	DNS := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?%s", m.User, m.Password, m.Ip, m.Port, m.DB, m.Parameter)
	conn, err := sqlx.Open("mysql", DNS)
	if err != nil {
		log.Fatalln("open mysql failed,", err.Error())
		return
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	if m.MaxIdleConns == 0 {
		m.MaxIdleConns = 5
	}
	conn.SetMaxIdleConns(m.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	if m.MaxOpenConns == 0 {
		m.MaxOpenConns = 100
	}
	conn.SetMaxOpenConns(m.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	if m.ConnMaxLifetime == 0 {
		m.ConnMaxLifetime = time.Hour
	}
	conn.SetConnMaxLifetime(m.ConnMaxLifetime)
	mysqlXMap[m.Name] = conn
}

func GetMysqlX(name ...string) *sqlx.DB {
	if mysqlXMap == nil {
		return nil
	}
	var n string
	if len(name) > 0 {
		n = name[0]
	} else {
		n = "default"
	}
	if _, ok := mysqlXMap[n]; ok {
		return mysqlXMap[n]
	} else {
		return nil
	}
}

func DestroyMysqlX(name ...string) {
	if mysqlXMap == nil {
		return
	}
	var n string
	if len(name) > 0 {
		n = name[0]
	} else {
		n = "default"
	}
	if _, ok := mysqlXMap[n]; ok {
		_ = mysqlXMap[n].Close()
		delete(mysqlXMap, n)
	}
}

func DestroyMysqlXAll() {
	if mysqlXMap == nil {
		return
	}
	for k, r := range mysqlXMap {
		_ = r.Close()
		delete(mysqlXMap, k)
	}
}
