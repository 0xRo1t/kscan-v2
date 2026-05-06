package hydra

func DefaultMqttList() *AuthList {
	a := NewAuthList()
	a.Username = []string{
		"admin",
		"mqtt",
		"user",
		"test",
		"guest",
		"root",
		"",
	}
	a.Password = []string{
		"",
		"admin",
		"mqtt",
		"123456",
		"password",
		"public",
		"12345",
		"1234",
		"123",
		"test",
		"guest",
		"root",
		"zaq1@WSX",
		"qweasdzxc",
		"Passw0rd",
		"password123",
		"1q2w3e4r",
		"1qaz2wsx",
		"abc123",
		"123qwe",
		"0000",
		"1234567",
		"12345678",
		"123456789",
	}
	a.Special = []Auth{
		NewSpecialAuth("admin", "admin"),
		NewSpecialAuth("admin", "public"),
		NewSpecialAuth("admin", ""),
		NewSpecialAuth("", ""),
	}
	return a
}
