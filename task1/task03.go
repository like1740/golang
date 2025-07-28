package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// _ "github.com/go-sql-driver/mysql"
	// "github.com/jmoiron/sqlx"
)

type Student struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Age   int
	Grade string `gorm:"size:50"`
}

type Account struct {
	ID      uint    `gorm:"primaryKey"`
	Balance float64 `gorm:"not null"`
}

type Transaction struct {
	ID            uint    `gorm:"primaryKey"`
	FromAccountID uint    `gorm:"not null"`
	ToAccountID   uint    `gorm:"not null"`
	Amount        float64 `gorm:"not null"`
}

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;unique;not null"`
	Posts     []Post    `gorm:"foreignKey:UserID"`
	PostCount uint      `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Post struct {
	ID            uint      `gorm:"primaryKey"`
	Title         string    `gorm:"size:200;not null"`
	Content       string    `gorm:"type:text;not null"`
	UserID        uint      `gorm:"not null;index"`
	User          User      `gorm:"foreignKey:UserID"`
	Comments      []Comment `gorm:"foreignKey:PostID"`
	CommentCount  uint      `gorm:"default:0"`
	CommentStatus string    `gorm:"size:20;default:'无评论'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	UserID    uint      `gorm:"not null"`
	PostID    uint      `gorm:"not null;index"`
	Post      Post      `gorm:"foreignKey:PostID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// type Employee struct {
// 	ID         int    `db:"id"`
// 	Name       string `db:"name"`
// 	Department string `db:"department"`
// 	Salary     int    `db:"salary"`
// }

// type Book struct {
// 	ID     int     `db:"id"`
// 	Title  string  `db:"title"`
// 	Author string  `db:"author"`
// 	Price  float64 `db:"price"`
// }

func main() {
	/*
		1.1 基本CRUD操作
		假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
	*/
	const (
		username = "root"
		password = "root"
		hostname = "127.0.0.1:3306"
		dbname   = "test"
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	if err := db.AutoMigrate(&Student{}, &Account{}, &Transaction{}); err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}

	// // 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	// if err := db.Create(&Student{Name: "张三", Age: 20, Grade: "三年级"}).Error; err != nil {
	// 	log.Printf("插入记录失败: %v", err)
	// } else {
	// 	fmt.Println("成功插入新记录")
	// }

	// // 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	// var students []Student
	// if err := db.Where("age > ?", 18).Find(&students).Error; err != nil {
	// 	log.Printf("查询失败: %v", err)
	// } else {
	// 	fmt.Println("\n年龄大于18岁的学生:")
	// 	for _, s := range students {
	// 		fmt.Printf("ID: %d, 姓名: %s, 年龄: %d, 年级: %s\n", s.ID, s.Name, s.Age, s.Grade)
	// 	}
	// }

	// // 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	// if err := db.Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级").Error; err != nil {
	// 	log.Printf("更新失败: %v", err)
	// } else {
	// 	fmt.Println("\n成功更新张三的年级")
	// }

	// // 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	// result := db.Where("age < ?", 15).Delete(&Student{})
	// if result.Error != nil {
	// 	log.Printf("删除失败: %v", result.Error)
	// } else {
	// 	fmt.Printf("\n成功删除%d条记录\n", result.RowsAffected)
	// }
	/*
		1.2 基本CRUD操作
		假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
	*/
	// db.Create(&Account{Balance: 1000.0})
	// db.Create(&Account{Balance: 500.0})

	// // 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
	// if err := TransferMoney(db, 1, 2, 100.0); err != nil {
	// 	log.Printf("转账失败: %v", err)
	// } else {
	// 	log.Println("转账成功")
	// }

	// var accountA, accountB Account
	// db.First(&accountA, 1)
	// db.First(&accountB, 2)
	// log.Printf("账户A余额: %.2f", accountA.Balance)
	// log.Printf("账户B余额: %.2f", accountB.Balance)

	/*
		2.1 Sqlx入门: 使用SQL扩展库进行查询
		假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
	*/
	// db, err := sqlx.Connect("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True")
	// if err != nil {
	// 	log.Fatalf("数据库连接失败: %v", err)
	// }
	// defer db.Close()

	// // 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	// techEmployees, err := getEmployeesByDepartment(db, "技术部")
	// if err != nil {
	// 	log.Printf("查询技术部员工失败: %v", err)
	// } else {
	// 	fmt.Println("技术部员工:")
	// 	for _, emp := range techEmployees {
	// 		fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 薪资: %d\n",
	// 			emp.ID, emp.Name, emp.Department, emp.Salary)
	// 	}
	// }

	// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
	// highestPaid, err := getHighestPaidEmployee(db)
	// if err != nil {
	// 	log.Printf("查询最高薪资员工失败: %v", err)
	// } else {
	// 	fmt.Println("\n薪资最高的员工:")
	// 	fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 薪资: %d\n",
	// 		highestPaid.ID, highestPaid.Name, highestPaid.Department, highestPaid.Salary)
	// }

	/*
		2.2 Sqlx入门: 实现类型安全映射
		假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
		要求 ：
			定义一个 Book 结构体，包含与 books 表对应的字段。
			编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
	*/
	// expensiveBooks, err := getBooksByMinPrice(db, 50.0)
	// if err != nil {
	// 	log.Fatalf("查询失败: %v", err)
	// }

	// fmt.Println("价格高于50元的书籍:")
	// for _, book := range expensiveBooks {
	// 	fmt.Printf(
	// 		"ID: %d | 书名: %-20s | 作者: %-10s | 价格: ￥%.2f\n",
	// 		book.ID, book.Title, book.Author, book.Price,
	// 	)
	// }

	/*
		3.1 进阶gorm: 模型定义
		假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
		要求 ：
			使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
			编写Go代码，使用Gorm创建这些模型对应的数据库表。
	*/
	user := User{Name: "赵六", Email: fmt.Sprintf("test-%d@example.com", time.Now().Unix())}
	db.Create(&user)

	post := Post{
		Title:   "GORM使用指南",
		Content: "这是一篇关于GORM的详细教程...",
		UserID:  user.ID,
	}
	db.Create(&post)

	// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	comments := []Comment{
		{Content: "好文章！", UserID: user.ID, PostID: post.ID},
		{Content: "学到了很多", UserID: user.ID, PostID: post.ID},
	}
	db.Create(&comments)

	/*
		3.2 进阶gorm: 关联查询
		基于上述博客系统的模型定义。
	*/
	// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	posts, err := getUserPostsWithComments(db, 1)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Println("\n用户1的文章及评论:")
		for _, post := range posts {
			fmt.Printf("文章: %s (评论数: %d)\n", post.Title, post.CommentCount)
			for _, comment := range post.Comments {
				fmt.Printf("  - 评论: %s\n", comment.Content)
			}
		}
	}

	// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
	mostCommented, err := getMostCommentedPost(db)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("\n评论最多的文章: %s (评论数: %d)\n",
			mostCommented.Title, mostCommented.CommentCount)
	}

	/*
		3.3 进阶gorm: 钩子函数
		继续使用博客系统的模型。
	*/
	// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
	var comment Comment
	db.First(&comment)
	db.Delete(&comment)

	updatedPost, _ := getUserPostsWithComments(db, user.ID)
	fmt.Printf("\n删除评论后的状态: %+v\n", updatedPost[0])

}

// func TransferMoney(db *gorm.DB, fromID, toID uint, amount float64) error {
// 	return db.Transaction(func(tx *gorm.DB) error {
// 		var fromAccount Account
// 		if err := tx.First(&fromAccount, fromID).Error; err != nil {
// 			return fmt.Errorf("找不到转出账户: %w", err)
// 		}

// 		if fromAccount.Balance < amount {
// 			return fmt.Errorf("账户余额不足: 当前余额 %.2f < 转账金额 %.2f",
// 				fromAccount.Balance, amount)
// 		}

// 		if err := tx.Model(&Account{}).Where("id = ?", fromID).
// 			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
// 			return fmt.Errorf("更新转出账户失败: %w", err)
// 		}

// 		if err := tx.Model(&Account{}).Where("id = ?", toID).
// 			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
// 			return fmt.Errorf("更新转入账户失败: %w", err)
// 		}

// 		if err := tx.Create(&Transaction{
// 			FromAccountID: fromID,
// 			ToAccountID:   toID,
// 			Amount:        amount,
// 		}).Error; err != nil {
// 			return fmt.Errorf("创建交易记录失败: %w", err)
// 		}

// 		return nil
// 	})
// }

// func getEmployeesByDepartment(db *sqlx.DB, department string) ([]Employee, error) {
// 	var employees []Employee
// 	query := `SELECT id, name, department, salary FROM employees WHERE department = ?`

// 	err := db.Select(&employees, query, department)
// 	return employees, err
// }

// func getHighestPaidEmployee(db *sqlx.DB) (Employee, error) {
// 	var employee Employee
// 	query := `SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1`

// 	err := db.Get(&employee, query)
// 	return employee, err
// }

// func getBooksByMinPrice(db *sqlx.DB, minPrice float64) ([]Book, error) {
// 	var books []Book

// 	query := `SELECT id, title, author, price FROM books WHERE price > ? ORDER BY price DESC`

// 	err := db.Select(&books, query, minPrice)
// 	if err != nil {
// 		return nil, fmt.Errorf("查询失败: %w", err)
// 	}

// 	return books, nil
// }

func getUserPostsWithComments(db *gorm.DB, userID uint) ([]Post, error) {
	var posts []Post
	err := db.Preload("Comments").Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}

func getMostCommentedPost(db *gorm.DB) (Post, error) {
	var post Post
	err := db.Order("comment_count DESC").First(&post).Error
	return post, err
}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	result := tx.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + ?", 1))
	return result.Error
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var count int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		if err := tx.Model(&Post{}).Where("id = ?", c.PostID).
			Updates(map[string]interface{}{
				"comment_count":  0,
				"comment_status": "无评论",
			}).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Model(&Post{}).Where("id = ?", c.PostID).
			Update("comment_count", count).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Model(&Post{}).Where("id = ?", c.PostID).
		Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
}
