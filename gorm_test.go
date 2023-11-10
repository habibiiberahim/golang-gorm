package belajar_golang_gorm

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"strconv"
	"testing"
	"time"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("user:password@(127.0.0.1:3306)/belajar_golang_gorm" +
		"?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		//turn off default transaction
		SkipDefaultTransaction: true,
		//cache prepare statement to memory
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	//Set Connection Pool Database
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db
}

var db = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, db)
}

func TestTruncateTableSample(t *testing.T) {
	err := db.Exec("truncate sample").Error
	assert.Nil(t, err)
}

func TestTruncateTableUserLog(t *testing.T) {
	err := db.Exec("truncate user_logs").Error
	assert.Nil(t, err)
}

func TestTruncateTableTodo(t *testing.T) {
	err := db.Exec("truncate todos").Error
	assert.Nil(t, err)
}
func TestTruncateTableWallet(t *testing.T) {
	err := db.Exec("truncate wallets").Error
	assert.Nil(t, err)
}

func TestTruncateTableUser(t *testing.T) {
	err := db.Exec("truncate users").Error
	assert.Nil(t, err)
}

func TestExecuteSQL(t *testing.T) {
	err := db.Exec("insert into sample (id, name) values (?,?)", "1", "Habibi").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample (id, name) values (?,?)", "2", "Iberahim").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample (id, name) values (?,?)", "3", "Habibi Iberahim").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample (id, name) values (?,?)", "4", "Iberahim Habibi").Error
	assert.Nil(t, err)
}

type Sample struct {
	Id   string
	Name string
}

func TestRawSQL(t *testing.T) {
	var sample Sample
	err := db.Raw("select id, name from sample where id = ?", "1").Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "Habibi", sample.Name)

	var samples []Sample
	err = db.Raw("select id, name from sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(samples))
}

// implement lazy result using rows method for better memory consumtion
func TestSqlRow(t *testing.T) {
	var samples []Sample

	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	for rows.Next() {
		var id string
		var name string
		err := rows.Scan(&id, &name)
		assert.Nil(t, err)

		samples = append(samples, Sample{
			Id:   id,
			Name: name,
		})
	}
	assert.Equal(t, 4, len(samples))
}

func TestScanRow(t *testing.T) {
	var samples []Sample

	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	for rows.Next() {
		err := db.ScanRows(rows, &samples)
		assert.Nil(t, err)
	}
	assert.Equal(t, 4, len(samples))
}

func TestCreateUser(t *testing.T) {
	user := User{
		ID:       "1",
		Password: "rahasia",
		Name: Name{
			FirstName:  "Habibi",
			MiddleName: "Iberahim ",
			LastName:   "S.Kom",
		},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Information: "Ini tidak akan masuk ke database",
	}

	response := db.Create(&user)
	assert.Nil(t, response.Error)
	assert.Equal(t, 1, int(response.RowsAffected))
}

func TestBatchInsert(t *testing.T) {
	var users []User

	for i := 2; i < 10; i++ {
		users = append(users, User{
			ID:       strconv.Itoa(i),
			Password: "rahasia",
			Name:     Name{FirstName: "user " + strconv.Itoa(i)},
		})
	}

	result := db.Create(users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 8, int(result.RowsAffected))
}

func TestTransactionSuccess(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		//
		err := tx.Create(&User{
			ID:       "10",
			Password: "rahasia",
			Name:     Name{FirstName: "user 10"},
		}).Error
		if err != nil {
			return err
		}
		//
		err = tx.Create(&User{
			ID:       "11",
			Password: "rahasia",
			Name:     Name{FirstName: "user 11"},
		}).Error
		if err != nil {
			return err
		}
		//
		err = tx.Create(&User{
			ID:       "12",
			Password: "rahasia",
			Name:     Name{FirstName: "user 12"},
		}).Error
		if err != nil {
			return err
		}

		return nil
	})

	assert.Nil(t, err)
}

func TestTransactionError(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		//
		err := tx.Create(&User{
			ID:       "13",
			Password: "rahasia",
			Name:     Name{FirstName: "user 13"},
		}).Error
		if err != nil {
			return err
		}
		//
		err = tx.Create(&User{
			ID:       "11",
			Password: "rahasia",
			Name:     Name{FirstName: "user 11"},
		}).Error
		if err != nil {
			return err
		}
		return nil
	})

	assert.NotNil(t, err)
}

func TestManualTransactionSuccess(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{
		ID:       "13",
		Password: "rahasia",
		Name:     Name{FirstName: "user 13"},
	}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{
		ID:       "14",
		Password: "rahasia",
		Name:     Name{FirstName: "user 14"},
	}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestManualTransactionError(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{
		ID:       "15",
		Password: "rahasia",
		Name:     Name{FirstName: "user 15"},
	}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{
		ID:       "14",
		Password: "rahasia",
		Name:     Name{FirstName: "user 14"},
	}).Error
	assert.NotNil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestQuerySingleObject(t *testing.T) {
	user := User{}
	err := db.First(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)

	user = User{}
	err = db.Last(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "9", user.ID)
}

func TestQuerySingleObjectInlineCondition(t *testing.T) {
	user := User{}
	err := db.Take(&user, "id = ?", "5").Error
	assert.Nil(t, err)
	assert.Equal(t, "5", user.ID)
	assert.Equal(t, "user 5", user.Name.FirstName)
}

func TestQueryAllObjects(t *testing.T) {
	var users []User
	err := db.Find(&users, "id in ?", []string{"1", "2", "3", "4"}).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))
}

func TestQueryCondition(t *testing.T) {
	var users []User
	err := db.Where("first_name like ?", "%user%").
		Where("password = ?", "rahasia").
		Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 13, len(users))
}

func TestOrOperator(t *testing.T) {
	var users []User
	err := db.Where("first_name like ?", "%user%").
		Or("password = ?", "rahasia").
		Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 14, len(users))
}

func TestNotOperator(t *testing.T) {
	var users []User
	err := db.Not("first_name like ?", "%user%").
		Where("password = ?", "rahasia").
		Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

// Always use select field for better performance query
func TestSelectFields(t *testing.T) {
	var users []User
	err := db.Select("id, first_name").Find(&users).Error
	assert.Nil(t, err)

	for _, user := range users {
		assert.NotNil(t, user.ID)
		assert.NotEqual(t, "", user.Name.FirstName)
	}
	assert.Equal(t, 14, len(users))
}

func TestStructCondition(t *testing.T) {
	userCondition := User{
		Name: Name{
			FirstName: "user 5",
		},
		Password: "rahasia",
	}
	var users []User
	err := db.Where(userCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestMapCondition(t *testing.T) {
	mapCondition := map[string]interface{}{
		"middle_name": "",
		"last_name":   "",
	}
	var users []User
	err := db.Where(mapCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 13, len(users))
}

func TestOrderLimitOffset(t *testing.T) {
	var users []User
	err := db.Order("id asc, first_name desc").Limit(5).Offset(5).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
}

type UserResponse struct {
	ID        string
	FirstName string
	LastName  string
}

func TestQueryNonModel(t *testing.T) {
	var users []UserResponse
	err := db.Model(&User{}).Select("id", "first_name", "last_name").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 14, len(users))
}

func TestUpdate(t *testing.T) {
	user := User{}
	err := db.Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	user.Name.FirstName = "Ini"
	user.Name.FirstName = "Baru"
	user.Name.FirstName = "Update"
	user.Password = "newPassword"
	err = db.Save(&user).Error
	assert.Nil(t, err)
}

func TestSelectedColumn(t *testing.T) {
	err := db.Model(User{}).Where("id = ?", "1").Updates(map[string]interface{}{
		"middle_name": "update middle name via map",
		"last_name":   "update last name via map",
	}).Error
	assert.Nil(t, err)

	err = db.Model(User{}).Where("id = ?", "1").Update("password", "updatePassword").Error
	assert.Nil(t, err)

	err = db.Where("id = ?", "1").Updates(User{
		Name: Name{
			FirstName: "Habibi",
			LastName:  "Iberahim",
		},
	}).Error
	assert.Nil(t, err)
}

func TestAutoIncrement(t *testing.T) {
	for i := 0; i < 10; i++ {
		userLog := UserLog{
			UserId: "1",
			Action: "Login",
		}
		err := db.Create(&userLog).Error
		assert.Nil(t, err)
		assert.NotEqual(t, 0, userLog.ID)
	}
}

func TestSaveOrUpdateAutoIncrement(t *testing.T) {
	userLog := UserLog{
		UserId: "1",
		Action: "Test Action",
	}

	err := db.Save(&userLog).Error
	assert.Nil(t, err)
	assert.Equal(t, "1", userLog.UserId)

	userLog.UserId = "2"
	err = db.Save(&userLog).Error
	assert.Nil(t, err)
	assert.Equal(t, "2", userLog.UserId)
}

func TestSaveOrUpdateNonAutoIncrement(t *testing.T) {
	user := User{
		ID:       "21",
		Password: "rahasia",
		Name: Name{
			FirstName: "Habibi Iberahim",
		},
	}

	err := db.Save(&user).Error
	assert.Nil(t, err)

	user.Name.FirstName = "user 21 updated"
	err = db.Save(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, user.Name.FirstName, "user 21 updated")
}

func TestConflict(t *testing.T) {
	user := User{
		ID:       "25",
		Password: "rahasia",
		Name: Name{
			FirstName: "Habibi Iberahim",
		},
	}

	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ? ", "25").Error
	assert.Nil(t, err)

	err = db.Delete(&user).Error
	assert.Nil(t, err)

	err = db.Delete(&user, "id = ?", "21").Error
	assert.Nil(t, err)

	err = db.Where("id = ? ", "10").Delete(&User{}).Error
	assert.Nil(t, err)
}

func TestSoftDelete(t *testing.T) {
	todo := Todo{
		UserId:      "1",
		Title:       "Todo 1",
		Description: "Action Todo 1",
	}

	err := db.Create(&todo).Error
	assert.Nil(t, err)

	err = db.Delete(&todo).Error
	assert.Nil(t, err)
	assert.NotNil(t, todo.DeletedAt)

	var todos []Todo
	err = db.Find(&todos).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(todos))
}

func TestUnscoped(t *testing.T) {
	var todo Todo
	err := db.Unscoped().First(&todo, "id = ? ", "1").Error
	assert.Nil(t, err)

	err = db.Unscoped().Delete(&todo).Error
	assert.Nil(t, err)

	var todos []Todo
	err = db.Unscoped().Find(&todos).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(todos))
}

func TestLock(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Take(&user, "id = ?", "1").Error
		if err != nil {
			return err
		}

		user.Name.FirstName = "Iberahim"
		user.Name.MiddleName = "Habibi"
		err = tx.Save(&user).Error
		return err
	})
	assert.Nil(t, err)
}

func TestCreateWallet(t *testing.T) {
	wallet := Wallet{
		ID:      "1",
		UserId:  "1",
		Balance: 1000000,
	}

	err := db.Create(&wallet).Error
	assert.Nil(t, err)
}

func TestRetrieveReleation(t *testing.T) {
	var user User
	//Double Query for get relation
	err := db.Model(&User{}).Preload("Wallet").Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	assert.Equal(t, user.ID, "1")
	assert.Equal(t, user.Wallet.UserId, "1")
}

func TestRetrieveRelationJoin(t *testing.T) {
	var user User
	err := db.Model(&User{}).Joins("Wallet").Take(&user, "users.id = ?", "1").Error
	assert.Nil(t, err)

	assert.Equal(t, user.ID, "1")
	assert.Equal(t, user.Wallet.UserId, "1")
}

func TestAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       "50",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 50",
		},
		Information: "",
		Wallet: Wallet{
			ID:      "50",
			UserId:  "50",
			Balance: 1000000,
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestSkipAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       "51",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 51",
		},
		Information: "",
		Wallet: Wallet{
			ID:      "51",
			UserId:  "51",
			Balance: 1000000,
		},
	}

	err := db.Omit(clause.Associations).Create(&user).Error
	assert.Nil(t, err)
}

func TestUserAndAddress(t *testing.T) {
	user := User{
		ID:       "52",
		Password: "rahasia",
		Name: Name{
			FirstName: "Habibi",
		},
		Wallet: Wallet{
			ID:      "52",
			UserId:  "52",
			Balance: 1000000,
		},
		Addresses: []Address{
			{
				UserId:  "52",
				Address: "Banjarmasin",
			},
			{
				UserId:  "52",
				Address: "Banjarbaru",
			},
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestPreloadJoinOneToMany(t *testing.T) {
	var users []User
	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)
}

func TestTakePreloadJoinOneToMany(t *testing.T) {
	var user User
	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").
		Take(&user, "users.id = ?", "52").Error
	assert.Nil(t, err)
}

func TestBelongsTo(t *testing.T) {
	fmt.Println("Preload")
	var addresses []Address
	err := db.Model(&Address{}).Preload("User").Find(&addresses).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(addresses))

	fmt.Println("Joins")
	addresses = []Address{}
	err = db.Model(&Address{}).Joins("User").Find(&addresses).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(addresses))
}

// split table with realtion one to one for better performance query (ex: User split with Wallet)
func TestBelongsToOneToOne(t *testing.T) {
	fmt.Println("Preload")
	var wallets []Wallet
	err := db.Model(&Wallet{}).Preload("User").Find(&wallets).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(wallets))

	fmt.Println("Joins")
	wallets = []Wallet{}
	err = db.Model(&Wallet{}).Joins("User").Find(&wallets).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(wallets))
}

func TestCreateManyToMany(t *testing.T) {
	product := Product{
		ID:    "P001",
		Name:  "Contoh Product",
		Price: 100000,
	}
	err := db.Create(&product).Error
	assert.Nil(t, err)

	err = db.Table("user_like_product").Create(map[string]interface{}{
		"user_id":    "1",
		"product_id": "P001",
	}).Error
	assert.Nil(t, err)

	err = db.Table("user_like_product").Create(map[string]interface{}{
		"user_id":    "2",
		"product_id": "P001",
	}).Error
	assert.Nil(t, err)
}

func TestPreloadManyToMany(t *testing.T) {
	var product Product
	err := db.Preload("LikedByUsers").First(&product, "id = ?", "P001").Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(product.LikedByUsers))
}

func TestPreloadManyToManyUser(t *testing.T) {
	var user User
	err := db.Preload("LikeProducts").Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(user.LikeProducts))
}

func TestAssociationFind(t *testing.T) {
	var product Product
	err := db.First(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	var users []User
	err = db.Model(&product).Where("users.first_name LIKE ?", "User%").Association("LikedByUsers").Find(&users)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestAssociationAppend(t *testing.T) {
	var user User
	err := db.First(&user, "id = ? ", "52").Error
	assert.Nil(t, err)

	var product Product
	err = db.First(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Append(&user)
	assert.Nil(t, err)
}

func TestAssociationReplace(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Take(&user, "id = ?", "1").Error
		assert.Nil(t, err)

		wallet := Wallet{
			ID:      "01",
			UserId:  user.ID,
			Balance: 100000,
		}

		err = tx.Model(&user).Association("Wallet").Replace(&wallet)
		return err
	})
	assert.NotNil(t, err)

}

func TestAssociationDelete(t *testing.T) {
	var user User
	err := db.First(&user, "id = ? ", "3").Error
	assert.Nil(t, err)

	var product Product
	err = db.First(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Delete(&user)
	assert.Nil(t, err)
}

func TestAssociationClear(t *testing.T) {
	var product Product
	err := db.First(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Clear()
	assert.Nil(t, err)
}

func TestPreloadingWithCondition(t *testing.T) {
	var user User
	err := db.Preload("Wallet", "balance > ?", 1000000).First(&user, "id = ?", "1").Error
	assert.Nil(t, err)
}

func TestNestedPreloading(t *testing.T) {
	var wallet Wallet
	err := db.Preload("User.Addresses").Take(&wallet, "id = ?", "52").Error
	assert.Nil(t, err)
}

func TestPreloadAll(t *testing.T) {
	var user User
	err := db.Preload(clause.Associations).Take(&user, "id = ?", "52").Error
	assert.Nil(t, err)
}

func TestJoinQuery(t *testing.T) {
	var users []User
	//inner joins
	err := db.Joins("join wallets on wallets.user_id = users.id").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))

	users = []User{}
	//left joins
	err = db.Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 17, len(users))
}

func TestJoinCondition(t *testing.T) {
	var users []User
	err := db.Joins("join wallets on wallets.user_id = users.id AND wallets.balance > ?", 500000).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))

	users = []User{}
	//Alias using name field
	err = db.Joins("Wallet").Where("Wallet.balance > ?", 500000).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))
}

func TestCount(t *testing.T) {
	var users = []User{}
	var count int64
	//Alias using name field
	err := db.Model(User{}).Joins("Wallet").Where("Wallet.balance > ?", 500000).
		Find(&users).
		Count(&count).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(2), count)
}

type AggregationResult struct {
	TotalBalance int64
	MaxBalance   int64
	MinBalance   int64
	AvgBalance   float64
}

func TestAggregation(t *testing.T) {
	var result AggregationResult
	err := db.Model(Wallet{}).Select("sum(balance) as total_balance, min(balance) as min_balance, " +
		"max(balance) as max_balance, avg(balance) as avg_balance").Take(&result).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(4000000), result.TotalBalance)
	assert.Equal(t, int64(1000000), result.MinBalance)
	assert.Equal(t, int64(3000000), result.MaxBalance)
	assert.Equal(t, float64(2000000), result.AvgBalance)
}

func TestGroupByHaving(t *testing.T) {
	var results []AggregationResult
	err := db.Model(Wallet{}).Select("sum(balance) as total_balance, min(balance) as min_balance, "+
		"max(balance) as max_balance, avg(balance) as avg_balance").
		Joins("User").Group("User.id").Having("sum(balance) > ?", 3000000).
		Find(&results).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(results))
}

func TestContext(t *testing.T) {
	ctx := context.Background()

	var users []User
	err := db.WithContext(ctx).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 17, len(users))
}
func BrokeWalletBalance(db *gorm.DB) *gorm.DB {
	return db.Where("balance = ?", 0)
}

func RichWalletBalance(db *gorm.DB) *gorm.DB {
	return db.Where("balance > ?", 1000000)
}

func TestScopes(t *testing.T) {
	var wallets []Wallet
	err := db.Scopes(BrokeWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(wallets))

	wallets = []Wallet{}
	err = db.Scopes(RichWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(wallets))
}

func TestMigrator(t *testing.T) {
	err := db.Migrator().AutoMigrate(GuestBook{})
	assert.Nil(t, err)
}

func TestHook(t *testing.T) {
	user := User{
		Password: "rahahsia",
		Name: Name{
			FirstName: "User Hook",
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
	assert.NotEqual(t, "", user.ID)
}
