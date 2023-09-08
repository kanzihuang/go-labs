package main

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"log"
)

const (
	addr = "backend.123sou.cn:389"
	dc   = "dc=eryajf,dc=net"
)

func main() {
	l, err := ldap.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	admin := fmt.Sprintf("cn=%s,%s", "admin", dc)
	err = l.Bind(admin, "123456")
	if err != nil {
		log.Fatal("admin登录失败：", err)
	}

	user := "p"
	password := "123456"
	addUser(l, dc, user, password)
	getAllUsers(l, dc)
	getUser(l, dc, user)
	checkUser(dc, user, password)
	delUser(l, dc, user)
}

func getUser(l *ldap.Conn, dc, uid string) string {
	searchRequest := ldap.NewSearchRequest(dc,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldap.EscapeFilter(uid)),
		[]string{"dn", "uid", "cn", "mail", "telephone"},
		nil)

	sr, err := l.Search(searchRequest)
	if err != nil {
		fmt.Println("获取用户出错：", err)
	}

	if len(sr.Entries) != 1 {
		fmt.Println("用户不存在或返回条目过多")
		return ""
	}
	for _, entry := range sr.Entries {
		email := entry.GetAttributeValue("mail")
		fmt.Printf("email: %s\n", email)
	}
	return sr.Entries[0].GetAttributeValue("mail")
}

func getAllUsers(l *ldap.Conn, dc string) []string {
	searchRequest := ldap.NewSearchRequest(dc,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson))"),
		[]string{"uid", "cn", "mail", "telephone"},
		nil)

	sr, err := l.Search(searchRequest)
	if err != nil {
		fmt.Println("获取用户出错", err)
	}

	for _, entry := range sr.Entries {
		email := entry.GetAttributeValue("mail")
		fmt.Printf("email: %s\n", email)
	}
	//return sr.Entries[0].GetAttributeValue("mail")
	return []string{}
}

func checkUser(dc, uid, password string) bool {
	l, err := ldap.Dial("tcp", addr)
	if err != nil {
		log.Fatal("LDAP连接错误：", err)
	}
	defer l.Close()

	if _, err := l.SimpleBind(&ldap.SimpleBindRequest{
		Username: fmt.Sprintf("uid=%s,ou=people,%s", uid, dc),
		Password: password,
	}); err != nil {
		fmt.Println("用户验证失败：", err)
		return false
	}
	return true
}

func delUser(l *ldap.Conn, dc, uid string) {
	err := l.Del(&ldap.DelRequest{
		DN:       fmt.Sprintf("uid=%s,%s", uid, dc),
		Controls: nil,
	})
	if err != nil {
		fmt.Println("删除用户错误：", err)
		return
	}
}

func addUser(l *ldap.Conn, uc, uid, password string) {
	//创建新用户
	user := fmt.Sprintf("uid=%s,%s", uid, dc)
	addResponse := ldap.NewAddRequest(user, []ldap.Control{})
	addResponse.Attribute("cn", []string{uid})
	addResponse.Attribute("sn", []string{uid})
	addResponse.Attribute("uid", []string{uid})
	addResponse.Attribute("homeDirectory", []string{fmt.Sprintf("/home/%s", uid)})
	addResponse.Attribute("loginShell", []string{uid})
	addResponse.Attribute("gidNumber", []string{"0"})
	addResponse.Attribute("uidNumber", []string{"8001"})
	addResponse.Attribute("mail", []string{uid + "@ldap.com"})
	addResponse.Attribute("objectClass", []string{"shadowAccount", "posixAccount", "top", "inetOrgPerson"})
	err := l.Add(addResponse)
	if err != nil {
		fmt.Println("创建用户失败: ", err)
		return
	}

	//随机给用户生成密码，并将新密码输出
	passwordModifyRequest2 := ldap.NewPasswordModifyRequest(fmt.Sprintf("uid=%s,%s", uid, uc), "", password)
	if passwordModifyResponse2, err := l.PasswordModify(passwordModifyRequest2); err != nil {
		fmt.Println("修改密码失败：", err)
		return
	} else {
		generatedPassword := passwordModifyResponse2.GeneratedPassword
		fmt.Println("生成的密码: ", generatedPassword)
	}
	return
}
