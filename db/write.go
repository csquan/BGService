package db

import (
	"fmt"
	"github.com/ethereum/BGService/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
)

// 插入
func InsertUser(engine *xorm.Engine, user *types.Users) error {
	rows, err := engine.Table("users").Insert(user)
	if err != nil {
		log.Println(err)
		return err
	}
	if rows == 0 {
		fmt.Println("插入失败")
		return errors.New("insert null")
	}
	fmt.Println("插入成功")
	return nil
}

func InsertInvitation(engine *xorm.Engine, invitation *types.Invitation) error {
	rows, err := engine.Table("invitation").Insert(invitation)
	if err != nil {
		log.Println(err)
		return err
	}
	if rows == 0 {
		fmt.Println("插入失败")
		return errors.New("insert null")
	}
	fmt.Println("插入成功")
	return nil
}

func InsertUserBindInfo(engine *xorm.Engine, UserBindInfo *types.InsertUserBindInfo) error {
	rows, err := engine.Table("userBindInfos").Insert(UserBindInfo)
	if err != nil {
		log.Println(err)
		return err
	}
	if rows == 0 {
		fmt.Println("插入失败")
		return errors.New("insert null")
	}
	fmt.Println("插入成功")
	return nil
}

func DeleteUserBindInfo(engine *xorm.Engine, id int) error {
	user := types.UserBindInfos{ID: id}
	rows, err := engine.Table("userBindInfos").Delete(&user)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if rows == 0 {
		logrus.Error("删除失败")
		return errors.New("delete fail")
	}
	return nil
}

func InsertUserStrategy(engine *xorm.Engine, UserBindInfo *types.UserStrategy) error {
	rows, err := engine.Table("`userStrategy`").Insert(UserBindInfo)
	if err != nil {
		log.Println(err)
		return err
	}
	if rows == 0 {
		fmt.Println("插入失败")
		return errors.New("insert null")
	}
	fmt.Println("插入成功")
	return nil
}

/*
// 删除
func deleteUser(engine *xorm.Engine, name string) {
	user := types.Users{Username: name}
	rows, err := engine.Delete(&user)
	if err != nil {
		log.Println(err)
		return
	}
	if rows == 0 {
		fmt.Println("删除失败")
		return
	}
	fmt.Println("删除成功")
}

// 修改
func UpdateUser(engine *xorm.Engine, user *types.Users) {
	//Update(bean interface{}, condiBeans ...interface{}) bean是需要更新的bean,condiBeans是条件
	update, err := engine.Update(user, types.Users{Id: user.Id})
	if err != nil {
		log.Println(err)
		return
	}
	if update > 0 {
		fmt.Println("更新成功")
		return
	}
	log.Println("更新失败")
}


// 事务
func sessionUserTest(engine *xorm.Engine, user *types.Users) {
	session := engine.NewSession()
	session.Begin()
	_, err := session.Insert(user)
	if err != nil {
		session.Rollback()
		log.Fatal(err)
	}
	user.Username = "mac"
	_, err = session.Update(user, types.Users{Id: user.Id})
	if err != nil {
		session.Rollback()
		log.Fatal(err)
	}
	err = session.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
*/
