package goethx

// Suggested values for the public eth blockchain

// PublicBlockCountValid represents after how many blocks
// a tx is considered valid
const PublicBlockCountValid = 3

// PublicPollingInterval represents the number of seconds between each tx check
const PublicPollingInterval = 2000

// PublicPollingTimeOut represents the number of seconds
// after we stop looking at tx
const PublicPollingTimeOut = 7200000
