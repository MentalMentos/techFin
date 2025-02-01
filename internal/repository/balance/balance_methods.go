package balance

import (
	"context"
	"fmt"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/pkg/helpers"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

// GetBalance извлекает баланс пользователя, сначала проверяя Redis, а затем базу данных
func (r *BalanceRepository) GetBalance(ctx context.Context, userID int) (float64, error) {
	// Проверка наличия баланса в Redis
	balanceStr, err := r.redisClient.Get(ctx, fmt.Sprintf("balance:%d", userID))
	if err == nil {
		// Попытка преобразовать значение в число с плавающей точкой
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err == nil {
			// Если баланс найден и успешно преобразован, возвращаем его
			r.logger.Info(helpers.RepoPrefix, "Balance retrieved from Redis cache")
			return balance, nil
		}
		// Логирование ошибки при парсинге значения из Redis
		r.logger.Info(helpers.RepoPrefix, helpers.RepoRedisParseError)
	}

	// Если баланс не найден в кэше, извлекаем его из базы данных
	var balance float64
	err = r.db.DB().QueryRowContext(ctx, db.Query{
		Name:     "get_balance",
		QueryRaw: "SELECT balance FROM users WHERE id=$1;",
	}, userID).Scan(&balance)
	if err != nil {
		// Логирование ошибки при извлечении баланса из базы данных
		r.logger.Info(helpers.RepoPrefix, helpers.RepoGetBalanceError)
		return 0, errors.Wrap(err, helpers.RepoGetBalanceError)
	}

	// Кэширование баланса в Redis с временем жизни 15 минут
	err = r.redisClient.Set(ctx, fmt.Sprintf("balance:%d", userID), fmt.Sprintf("%f", balance), 15*time.Minute)
	if err != nil {
		// Логирование ошибки при кэшировании баланса в Redis
		r.logger.Info(helpers.RepoPrefix, helpers.RepoCacheBalanceError)
	}

	// Логирование успешного извлечения баланса из базы данных и кэширования
	r.logger.Info(helpers.RepoPrefix, "Balance retrieved from database and cached")
	return balance, nil
}

// UpdateBalance обновляет баланс пользователя и кэширует новое значение
func (r *BalanceRepository) UpdateBalance(ctx context.Context, tx pgx.Tx, userID int, amount float64) (float64, error) {
	var updatedBalance float64

	// Обновление баланса в базе данных
	_, err := tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2;", amount, userID)
	if err != nil {
		// Логирование ошибки при обновлении баланса
		r.logger.Info(helpers.RepoPrefix, helpers.RepoUpdateBalanceError)
		return 0, errors.Wrap(err, helpers.RepoUpdateBalanceError)
	}

	// Получение обновленного баланса
	err = tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&updatedBalance)
	if err != nil {
		// Логирование ошибки при извлечении обновленного баланса
		r.logger.Info(helpers.RepoPrefix, helpers.RepoFetchBalanceError)
		return 0, errors.Wrap(err, helpers.RepoFetchBalanceError)
	}

	// Кэширование обновленного баланса в Redis с временем жизни 15 минут
	err = r.redisClient.Set(ctx, fmt.Sprintf("balance:%d", userID), fmt.Sprintf("%f", updatedBalance), 15*time.Minute)
	if err != nil {
		// Логирование ошибки при кэшировании обновленного баланса
		r.logger.Info(helpers.RepoPrefix, helpers.RepoCacheBalanceError)
	}

	// Логирование успешного обновления и кэширования баланса
	r.logger.Info(helpers.RepoPrefix, "Balance updated successfully and cached")
	return updatedBalance, nil
}
