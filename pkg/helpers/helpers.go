package helpers

const (
	ServicePrefix             = "[ SERVICE ]"
	ServiceGetBalanceError    = "FAILED_TO_GET_BALANCE"
	ServiceUpdateBalanceError = "FAILED_TO_UPDATE_BALANCE"
	ServiceTransferError      = "FAILED_TO_TRANSFER_FUNDS"
	ServiceTransactionError   = "FAILED_TO_GET_TRANSACTIONS"
	ServiceInsufficientFunds  = "INSUFFICIENT_FUNDS"
	ServiceInvalidAmount      = "INVALID_TRANSFER_AMOUNT"

	RepoPrefix                   = "[ REPO ]"
	RepoGetBalanceError          = "FAILED_TO_GET_BALANCE_FROM_DB"
	RepoCacheBalanceError        = "FAILED_TO_CACHE_BALANCE_IN_REDIS"
	RepoUpdateBalanceError       = "FAILED_TO_UPDATE_BALANCE"
	RepoFetchBalanceError        = "FAILED_TO_FETCH_UPDATED_BALANCE"
	RepoRedisParseError          = "FAILED_TO_PARSE_CACHED_BALANCE"
	RepoCreateTransactionError   = "FAILED_TO_CREATE_TRANSACTION"
	RepoCacheTransactionError    = "FAILED_TO_CACHE_TRANSACTION"
	RepoGetTransactionsError     = "FAILED_TO_GET_TRANSACTIONS"
	RepoScanTransactionError     = "FAILED_TO_SCAN_TRANSACTION"
	RepoIterateTransactionsError = "FAILED_TO_ITERATE_TRANSACTIONS"

	HandlerPrefix = "[ HANDL ]"

	AppPrefix = " [ APP ] "

	PgPrefix            = " [ POSTGRES ] "
	ReconnectDB         = "RECONNECTING TO DATABASE..."        // ReconnectDB contains reconnect db message
	DisconnectDB        = "DISCONNECTED FROM DATABASE"         // DisconnectDB contains disconnect db message
	PgConnectFailed     = "FAILED TO CONNECT TO DATABASE"      // PgConnectFailed contains error message for failed to connect to database
	PgConnectSuccess    = "SUCCESSFULLY CONNECTED TO POSTGRES" // PgConnectSuccess contains success message for successfully connected to database
	PgTransactionFailed = "FAILED TO FETCH TRANSACTION"        // PgTransactionFailed contains error message for failed to fetch transaction
	PgMigrateFailed     = "FAILED TO MIGRATE DATABASE"         // PgMigrateFailed contains error message for failed to migrate database
	NoRowsAffected      = "NO ROWS AFFECTED"                   // NoRowsAffected contains error message for no rows affected
	FailedToRollback    = "FAILED TO ROLLBACK"                 // FailedToRollback contains error message for failed to rollback
	FailedToClose       = "FAILED TO CLOSE"                    // FailedToClose contains error message for failed to close
)
