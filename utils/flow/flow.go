package flow

import (
	"strconv"

	"gin-auth-mongo/databases"
	"gin-auth-mongo/utils/consts"
	"gin-auth-mongo/utils/datetime"
)

// check if there is a blocked record for the IP in Redis
// true: blocked
// false: not blocked
func CheckIPBlocked(ip string) bool {

	blockedKey := consts.FLOW_LIMIT_BLOCKED_KEY + ip
	isBlocked, err := databases.RedisGet(blockedKey)

	// blocked
	if err != nil || isBlocked == "true" {
		return true
	}

	return false
}

// set blocked record for the IP in Redis
func SetIPBlocked(ip string) {

	blockedKey := consts.FLOW_LIMIT_BLOCKED_KEY + ip

	// block the IP for FLOW_LIMIT_BLOCKED minutes
	databases.RedisSet(blockedKey, "true", consts.FLOW_LIMIT_BLOCKED, datetime.MINUTES)
}

// get request count for the IP from Redis
func GetIPRequestCount(ip string) int {

	counterKey := consts.FLOW_LIMIT_COUNTER_KEY + ip
	count, err := databases.RedisGet(counterKey)
	if err != nil {
		return 0
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		return 0
	}

	return countInt
}

// set counter expiry
func SetIPRequestCountExpiry(ip string) {

	counterKey := consts.FLOW_LIMIT_COUNTER_KEY + ip
	databases.RedisExpire(counterKey, consts.FLOW_LIMIT_PERIOD, datetime.SECONDS)
}

// increase request count for the IP in Redis (1 time)
func IncreaseIPRequestCountV3(ip string) (int, error) {

	counterKey := consts.FLOW_LIMIT_COUNTER_KEY + ip
	count, err := databases.RedisIncr(counterKey)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
