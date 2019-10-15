package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ApplicationConfig struct {
	Id         int       `orm:"column(id);pk" description:"主键"`
	AppKey     string    `orm:"column(app_key);size(255)" description:"唯一应用标识"`
	AppName    string    `orm:"column(app_name);size(255)" description:"应用名称"`
	Hosts      string    `orm:"column(hosts);size(500);null" description:"部署服务器"`
	Remark     string    `orm:"column(remark);size(500);null" description:"应用描述"`
	Owner      string    `orm:"column(owner);size(50);null" description:"负责人"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);null" description:"创建时间"`
	UpdateTime time.Time `orm:"column(update_time);type(datetime);null" description:"更新时间"`
}

func (t *ApplicationConfig) TableName() string {
	return "application_config"
}

func init() {
	orm.RegisterModel(new(ApplicationConfig))
}

// AddApplicationConfig insert a new ApplicationConfig into database and returns
// last inserted Id on success.
func AddApplicationConfig(m *ApplicationConfig) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetApplicationConfigById retrieves ApplicationConfig by Id. Returns error if
// Id doesn't exist
func GetApplicationConfigById(id int) (v *ApplicationConfig, err error) {
	o := orm.NewOrm()
	v = &ApplicationConfig{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApplicationConfig retrieves all ApplicationConfig matches certain condition. Returns empty list if
// no records exist
func GetAllApplicationConfig(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ApplicationConfig))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []ApplicationConfig
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateApplicationConfig updates ApplicationConfig by Id and returns error if
// the record to be updated doesn't exist
func UpdateApplicationConfigById(m *ApplicationConfig) (err error) {
	o := orm.NewOrm()
	v := ApplicationConfig{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApplicationConfig deletes ApplicationConfig by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApplicationConfig(id int) (err error) {
	o := orm.NewOrm()
	v := ApplicationConfig{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ApplicationConfig{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
