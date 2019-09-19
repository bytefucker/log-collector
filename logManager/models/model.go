package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	etcd "go.etcd.io/etcd/clientv3"
)

// mysql配置信息
type MysqlInfo struct {
	Address  string
	Port     string
	User     string
	Password string
	DbName   string
}

var (
	// mysql客户端
	mysqlClient *sqlx.DB
	// etcd客户端
	etcdClient *etcd.Client
)

// 初始化mysql
func InitMySql(dbInfo MysqlInfo) (err error) {
	dbDevInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbInfo.User, dbInfo.Password, dbInfo.Address, dbInfo.Port, dbInfo.DbName)
	database, err := sqlx.Open("mysql", dbDevInfo)
	if err != nil {
		return
	}
	mysqlClient = database
	return
}

// 初始化etcd
func InitEtcd(etcdAddr []string) (err error) {
	etcdclient, err := etcd.New(etcd.Config{
		Endpoints:   etcdAddr,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	etcdClient = etcdclient
	return
}

/// 获取当前时间
func NowTime() (now_time string) {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 获取etcd中指定key的值
func GetEtcdKeyData(key string, defVal interface{}) (data []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcdClient.Get(ctx, key)
	cancel()
	if err != nil {
		err = fmt.Errorf("get etcd key: [%s] failed, err: %s", err)
		return
	}

	for _, ev := range resp.Kvs {
		if string(ev.Key) == key {
			data = ev.Value
			return
		}
	}
	if err != nil || resp.Count == 0 {
		logs.Warn("not found key: [%s] form etcd, init null", key)
		data, err = json.Marshal(&defVal)
		if err != nil {
			logs.Warn("not found key: [%s] form etcd, init null failed", key)
			return
		}
		err = nil
	}
	return
}

// 设置数据到etcd
func SetEtcdKeyData(key, data string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = etcdClient.Put(ctx, key, data)
	cancel()

	if err != nil {
		err = fmt.Errorf("set data to etcd failed, key:[%s], data:[%s]. err: %s", key, data, err)
		return
	}
	return
}
