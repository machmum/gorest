package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/machmum/gorest/utl/model/postgresql"
	"strings"
)

func main() {
	var psn = `postgres://localhost:5432/postgres`
	db, err := gorm.Open("postgres", psn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SingularTable(true)

	dbInsert := `insert into public.oauth_channel values (1,3,'1.0','admin','admin',now(),null,null);
	insert into public.profile values (1,'user','user','$2a$10$8XSmzoUQdMz1GCPHoen33OBGJsaZsJnV9Yuxrgh62b.uHMG5Qq0sG','username','description','user@mail.com','user',1,now(),null,null,null,null);
	insert into public.product values (1, 1, 'product ck1-a', 'this is description', '1000', '100', 1, '1', now(), null, null, null, null);
	insert into public.product_image values (1,1,'https://tinyjpg.com/images/social/website.jpg',1,1,now(),null,null,null,null);`
	queries := strings.Split(dbInsert, ";")

	createSchema(db, &gorestdb.OauthChannel{}, &gorestdb.Profile{}, &gorestdb.Product{}, &gorestdb.ProductImage{})

	for _, v := range queries[0 : len(queries)-1] {
		if err := db.Exec(v).Error; err != nil {
			logrus.Fatal(err)
		}
	}

	logrus.Fatal("done")
}

func createSchema(db *gorm.DB, models ...interface{}) {
	for _, model := range models {
		if exist := db.HasTable(model); !exist {
			if err := db.CreateTable(model).Error; err != nil {
				logrus.Fatal(err)
			}
		}
	}
}
