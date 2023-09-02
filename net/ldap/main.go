package main

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"log"
)

const (
	addr = "backend.123sou.cn:389"
	user = "uid=java1,dc=devopsman,dc=cn"
)

func main() {
	l, err := ldap.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.Bind("cn=admin,dc=devopsman,dc=cn", "admin123")
	if err != nil {
		log.Fatal(err)
	}
	delUser(l, user)
	password := addUser(l, user)
	checkUser(user, password)
}

func checkUser(user, password string) bool {
	l, err := ldap.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	_, err = l.SimpleBind(&ldap.SimpleBindRequest{
		Username: user,
		Password: password,
	})
	if err != nil {
		fmt.Println("用户验证失败：", err)
		return false
	}
	return true
}

func delUser(l *ldap.Conn, user string) {
	err := l.Del(&ldap.DelRequest{
		DN:       user,
		Controls: nil,
	})
	if err != nil {
		fmt.Println("删除用户错误：", err)
		return
	}
}

func addUser(l *ldap.Conn, user string) string {
	//创建新用户
	addResponse := ldap.NewAddRequest(user, []ldap.Control{})
	addResponse.Attribute("cn", []string{"java1"})
	addResponse.Attribute("sn", []string{"java1"})
	addResponse.Attribute("uid", []string{"java1"})
	addResponse.Attribute("homeDirectory", []string{"/home/java1"})
	addResponse.Attribute("loginShell", []string{"java1"})
	addResponse.Attribute("gidNumber", []string{"0"})
	addResponse.Attribute("uidNumber", []string{"8001"})
	addResponse.Attribute("objectClass", []string{"shadowAccount", "posixAccount", "top", "inetOrgPerson"})
	err := l.Add(addResponse)
	if err != nil {
		fmt.Println("创建用户失败: ", err)
		return ""
	}

	//随机给用户生成密码，并将新密码输出
	passwordModifyRequest2 := ldap.NewPasswordModifyRequest("uid=java1,dc=devopsman,dc=cn", "", "")
	passwordModifyResponse2, err := l.PasswordModify(passwordModifyRequest2)
	if err != nil {
		fmt.Println("修改密码失败：", err)
		return ""
	}
	generatedPassword := passwordModifyResponse2.GeneratedPassword
	fmt.Println("生成的密码: ", generatedPassword)
	return generatedPassword
}
