package repository_test

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"CHIPMUNK-T0T/ito_web_app/internal/functional"
	"CHIPMUNK-T0T/ito_web_app/internal/repository"
	"CHIPMUNK-T0T/ito_web_app/test/mock"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepository(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		db, mock, err := mock.NewMockDB()
		if err != nil {
			t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
		}

		repo := repository.NewUserRepository(db)
		user, _ := domain.NewUser("testuser", "password")

		mock.ExpectBegin()
		expectedSQL := regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`password`) VALUES (?,?,?,?,?)")
		mock.ExpectExec(expectedSQL).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Username(), user.Password()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.Create(user)
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未実行のクエリが存在します: %v", err)
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		db, mock, err := mock.NewMockDB()
		if err != nil {
			t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
		}

		repo := repository.NewUserRepository(db)

		hashedPassword, _ := functional.Encrypt("password")
		rows := sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "testuser", string(hashedPassword), time.Now(), time.Now(), nil)

		expectedSQL := regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")
		mock.ExpectQuery(expectedSQL).
			WithArgs(1, 1).
			WillReturnRows(rows)

		user, err := repo.FindByID(1)
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}
		if user.Username() != "testuser" {
			t.Errorf("期待値: testuser, 実際: %s", user.Username())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未実行のクエリが存在します: %v", err)
		}
	})

	t.Run("FindByUserNameAndPassword", func(t *testing.T) {
		db, mock, err := mock.NewMockDB()
		if err != nil {
			t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
		}

		repo := repository.NewUserRepository(db)

		hashedPassword, _ := functional.Encrypt("password")
		rows := sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "testuser", string(hashedPassword), time.Now(), time.Now(), nil)

		expectedSQL := regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")
		mock.ExpectQuery(expectedSQL).
			WithArgs("testuser", 1).
			WillReturnRows(rows)

		user, err := repo.FindByUserNameAndPassword("testuser", "password")
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}
		if user.Username() != "testuser" {
			t.Errorf("期待値: testuser, 実際: %s", user.Username())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未実行のクエリが存在します: %v", err)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		db, mock, err := mock.NewMockDB()
		if err != nil {
			t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
		}

		repo := repository.NewUserRepository(db)

		rows := sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "user1", "pass1", time.Now(), time.Now(), nil).
			AddRow(2, "user2", "pass2", time.Now(), time.Now(), nil)

		expectedSQL := regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")
		mock.ExpectQuery(expectedSQL).
			WillReturnRows(rows)

		users, err := repo.FindAll()
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}
		if len(users) != 2 {
			t.Errorf("期待値: 2, 実際: %d", len(users))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未実行のクエリが存在します: %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		db, mock, err := mock.NewMockDB()
		if err != nil {
			t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
		}

		repo := repository.NewUserRepository(db)
		user := domain.NewUserWithID(1, "updated", "newpassword")

		// 更新前のユーザー存在確認のSELECTクエリ
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, "testuser", "oldpassword", time.Now(), time.Now(), nil))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`username`=?,`password`=? WHERE `users`.`deleted_at` IS NULL AND `id` = ?")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Username(), user.Password(), user.ID()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.Update(user)
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未実行のクエリが存在します: %v", err)
		}
	})

	// t.Run("Delete", func(t *testing.T) {
	// 	db, mock, err := mock.NewMockDB()
	// 	if err != nil {
	// 		t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
	// 	}

	// 	repo := repository.NewUserRepository(db)

	// 	// 削除前のユーザー存在確認のSELECTクエリ
	// 	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
	// 		WithArgs(1, 1).
	// 		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at", "deleted_at"}).
	// 			AddRow(1, "testuser", "password", time.Now(), time.Now(), nil))

	// 	// 削除用のクエリ
	// 	mock.ExpectBegin()
	// 	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).
	// 		WithArgs(sqlmock.AnyArg(), 1).
	// 		WillReturnResult(sqlmock.NewResult(1, 1))
	// 	mock.ExpectCommit()

	// 	err = repo.Delete(1)
	// 	if err != nil {
	// 		t.Errorf("予期しないエラー: %v", err)
	// 	}

	// 	// 削除確認用のクエリ（トランザクション期待値なし）
	// 	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
	// 		WithArgs(1, 1).
	// 		WillReturnError(sql.ErrNoRows)

	// 	_, err = repo.FindByID(1)
	// 	if err == nil {
	// 		t.Error("削除されたユーザーが見つかりました")
	// 	}

	// 	if err := mock.ExpectationsWereMet(); err != nil {
	// 		t.Errorf("未実行のクエリが存在します: %v", err)
	// 	}
	// })
}