package transaction

import (
	"context"
	"fmt"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/MentalMentos/techFin/pkg/helpers"
	"github.com/MentalMentos/techFin/pkg/logger"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"time"
)

type TransactionRepo struct {
	db          db.Client
	redisClient redis.IRedis
	logger      logger.Logger
}

func NewTransactionRepo(db db.Client, redisClient redis.IRedis, logger logger.Logger) *TransactionRepo {
	return &TransactionRepo{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
	}
}

// CreateTransaction создает новую транзакцию и кэширует её в Redis
func (r *TransactionRepo) CreateTransaction(ctx context.Context, tx pgx.Tx, userID int, amount float64, targetUserID *int) error {
	// Вставка новой транзакции в базу данных
	_, err := tx.Exec(ctx, "INSERT INTO transactions (user_id, amount, target_user_id, status) VALUES ($1, $2, $3, 'completed');",
		userID, amount, targetUserID)
	if err != nil {
		// Логирование ошибки при создании транзакции в базе данных и откат транзакции
		r.logger.Info(helpers.RepoPrefix, helpers.RepoCreateTransactionError)
		return errors.Wrap(err, helpers.RepoCreateTransactionError)
	}

	// Создание ключа для транзакции в Redis
	transactionKey := fmt.Sprintf("transaction:%d:%d", userID, time.Now().Unix())
	transactionData := map[string]interface{}{
		"user_id":        userID,
		"amount":         amount,
		"target_user_id": targetUserID,
		"status":         "completed",
	}

	// Кэширование данных транзакции в Redis с временем жизни 24 часа
	err = r.redisClient.SetObject(ctx, transactionKey, transactionData, 24*time.Hour)
	if err != nil {
		// Логирование ошибки при сохранении транзакции в кэш
		r.logger.Info(helpers.RepoPrefix, helpers.RepoCacheTransactionError)
		return errors.Wrap(err, helpers.RepoCacheTransactionError)
	}

	// Логирование успешного создания транзакции и её кэширования
	r.logger.Info(helpers.RepoPrefix, "Transaction created successfully and cached")
	return nil
}

// GetLastTransactions извлекает последние транзакции пользователя, сначала проверяя кэш
func (r *TransactionRepo) GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	// Создание ключа для последней транзакции в Redis
	cacheKey := fmt.Sprintf("last_transactions:%d", userID)
	var cachedTransactions []models.Transaction
	// Попытка получить данные из кэша
	err := r.redisClient.GetObject(ctx, cacheKey, &cachedTransactions)
	if err == nil && len(cachedTransactions) > 0 {
		// Если данные найдены в кэше, возвращаем их
		r.logger.Info(helpers.RepoPrefix, "Transactions retrieved from Redis cache")
		return cachedTransactions, nil
	}

	// Если транзакции не найдены в кэше, извлекаем из базы данных
	rows, err := r.db.DB().QueryContext(ctx, db.Query{
		Name:     "get_last_transactions",
		QueryRaw: "SELECT id, amount, target_user_id, status, created_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT 10",
	}, userID)
	if err != nil {
		// Логирование ошибки при запросе данных из базы данных
		r.logger.Info(helpers.RepoPrefix, helpers.RepoGetTransactionsError)
		return nil, errors.Wrap(err, helpers.RepoGetTransactionsError)
	}
	defer rows.Close()

	// Чтение и сканирование результатов из базы данных
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Amount, &t.TargetUserID, &t.Status, &t.CreatedAt); err != nil {
			// Логирование ошибки при сканировании строки
			r.logger.Info(helpers.RepoPrefix, helpers.RepoScanTransactionError)
			return nil, errors.Wrap(err, helpers.RepoScanTransactionError)
		}
		transactions = append(transactions, t)
	}

	// Проверка на ошибку при итерации по строкам
	if err := rows.Err(); err != nil {
		// Логирование ошибки при итерации по строкам
		r.logger.Info(helpers.RepoPrefix, helpers.RepoIterateTransactionsError)
		return nil, errors.Wrap(err, helpers.RepoIterateTransactionsError)
	}

	// Кэширование извлечённых транзакций с временем жизни 15 минут
	err = r.redisClient.SetObject(ctx, cacheKey, transactions, 15*time.Minute)
	if err != nil {
		// Логирование ошибки при сохранении транзакций в кэш
		r.logger.Info(helpers.RepoPrefix, helpers.RepoCacheTransactionError)
		return nil, errors.Wrap(err, helpers.RepoCacheTransactionError)
	}

	// Логирование успешного получения транзакций и их кэширования
	r.logger.Info(helpers.RepoPrefix, "Transactions retrieved from database and cached")
	return transactions, nil
}
