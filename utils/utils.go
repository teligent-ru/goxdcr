package utils

import (
	"encoding/binary"
	"errors"
	"expvar"
	"fmt"
	"github.com/couchbase/cbauth"
	"github.com/couchbase/go-couchbase"
	mc "github.com/couchbase/gomemcached"
	mcc "github.com/couchbase/gomemcached/client"
	base "github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/simple_utils"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type BucketBasicStats struct {
	ItemCount int `json:"itemCount"`
}

//Only used by unit test
//TODO: replace with go-couchbase bucket stats API
type CouchBucket struct {
	Name string           `json:"name"`
	Stat BucketBasicStats `json:"basicStats"`
}

var NonExistentBucketError error = errors.New("Bucket doesn't exist")

var logger_utils *log.CommonLogger = log.NewLogger("Utils", log.DefaultLoggerContext)

func loggerForFunc(logger *log.CommonLogger) *log.CommonLogger {
	var l *log.CommonLogger
	if logger != nil {
		l = logger
	} else {
		l = logger_utils
	}
	return l
}

func ValidateSettings(defs base.SettingDefinitions,
	settings map[string]interface{},
	logger *log.CommonLogger) error {
	var l *log.CommonLogger = loggerForFunc(logger)

	l.Debugf("Start validate setting=%v, defs=%v", settings, defs)
	var err *base.SettingsError = nil
	for key, def := range defs {
		val, ok := settings[key]
		if !ok && def.Required {
			if err == nil {
				err = base.NewSettingsError()
			}
			err.Add(key, errors.New("required, but not supplied"))
		} else {
			if val != nil && def.Data_type != reflect.PtrTo(reflect.TypeOf(val)) {
				if err == nil {
					err = base.NewSettingsError()
				}
				err.Add(key, errors.New(fmt.Sprintf("expected type is %v, supplied type is %v",
					def.Data_type, reflect.TypeOf(val))))
			}
		}
	}
	if err != nil {
		l.Infof("setting validation result = %v", *err)
		return *err
	}
	return nil
}

func RecoverPanic(err *error) {
	if r := recover(); r != nil {
		*err = errors.New(fmt.Sprint(r))
	}
}

func LocalPool(localConnectStr string) (couchbase.Pool, error) {
	localURL := fmt.Sprintf("http://%s", localConnectStr)
	client, err := couchbase.ConnectWithAuth(localURL, cbauth.NewAuthHandler(nil))
	if err != nil {
		return couchbase.Pool{}, NewEnhancedError(fmt.Sprintf("Error connecting to couchbase. url=%v", UrlForLog(localURL)), err)
	}
	return client.GetPool("default")
}

// Get bucket in local cluster
func LocalBucket(localConnectStr, bucketName string) (*couchbase.Bucket, error) {
	logger_utils.Debugf("Getting local bucket name=%v\n", bucketName)

	pool, err := LocalPool(localConnectStr)
	if err != nil {
		return nil, err
	}

	bucket, err := pool.GetBucket(bucketName)
	if err != nil {
		return nil, NewEnhancedError(fmt.Sprintf("Error getting bucket, %v, from pool.", bucketName), err)
	}

	logger_utils.Debugf("Got local bucket successfully name=%v\n", bucket.Name)
	return bucket, err
}

func UnwrapError(infos map[string]interface{}) (err error) {
	if infos != nil && len(infos) > 0 {
		err = infos["error"].(error)
	}
	return err
}

// returns an enhanced error with erroe message being "msg + old error message"
func NewEnhancedError(msg string, err error) error {
	return errors.New(msg + "\n err = " + err.Error())
}

// return host address in the form of hostName:port
func GetHostAddr(hostName string, port uint16) string {
	return hostName + base.UrlPortNumberDelimiter + strconv.FormatInt(int64(port), base.ParseIntBase)
}

// extract host name from hostAddr, which is in the form of hostName:port
func GetHostName(hostAddr string) string {
	return strings.Split(hostAddr, base.UrlPortNumberDelimiter)[0]
}

func GetPortNumber(hostAddr string) (uint16, error) {
	port_str := strings.Split(hostAddr, base.UrlPortNumberDelimiter)[1]
	port, err := strconv.ParseUint(port_str, 10, 16)
	if err == nil {
		return uint16(port), nil
	} else {
		return 0, err
	}
}

func GetMapFromExpvarMap(expvarMap *expvar.Map) map[string]interface{} {
	regMap := make(map[string]interface{})

	expvarMap.Do(func(keyValue expvar.KeyValue) {
		valueStr := keyValue.Value.String()
		// first check if valueStr is an integer
		valueInt, err := strconv.Atoi(valueStr)
		if err == nil {
			regMap[keyValue.Key] = valueInt
		} else {
			// then check if valueStr is a float
			valueFloat, err := strconv.ParseFloat(valueStr, 64)
			if err == nil {
				regMap[keyValue.Key] = valueFloat
			} else {
				// should never happen
				logger_utils.Errorf("Invalid value in expvarMap. Only float and integer values are supported")
			}
		}
	})
	return regMap
}

//convert the format returned by go-memcached StatMap - map[string]string to map[uint16]uint64
func ParseHighSeqnoStat(vbnos []uint16, stats_map map[string]string, highseqno_map map[uint16]uint64) error {

	for _, vbno := range vbnos {
		stats_key := fmt.Sprintf(base.VBUCKET_HIGH_SEQNO_STAT_KEY_FORMAT, vbno)
		highseqnostr, ok := stats_map[stats_key]
		if !ok {
			logger_utils.Infof("Can't find high seqno for vbno=%v in stats map. Source topology may have changed.\n", vbno)
			continue
		}
		highseqno, err := strconv.ParseUint(highseqnostr, 10, 64)
		if err != nil {
			return err
		}
		highseqno_map[vbno] = highseqno
	}

	return nil
}

// encode data in a map into a byte array, which can then be used as
// the body part of a http request
// so far only five types are supported: string, int, bool, LogLevel, []byte
// which should be sufficient for all cases at hand
func EncodeMapIntoByteArray(data map[string]interface{}) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	params := make(url.Values)
	for key, val := range data {
		var strVal string
		switch val.(type) {
		case string:
			strVal = val.(string)
		case int:
			strVal = strconv.FormatInt(int64(val.(int)), base.ParseIntBase)
		case bool:
			strVal = strconv.FormatBool(val.(bool))
		case log.LogLevel:
			strVal = val.(log.LogLevel).String()
		case []byte:
			strVal = string(val.([]byte))
		default:
			return nil, simple_utils.IncorrectValueTypeInMapError(key, val, "string/int/bool/LogLevel/[]byte")
		}
		params.Add(key, strVal)
	}

	return []byte(params.Encode()), nil
}

func UrlForLog(urlStr string) string {
	result, err := url.Parse(urlStr)
	if err == nil {
		if result.User != nil {
			result.User = url.UserPassword(result.User.Username(), "xxxx")
		}
		return result.String()
	} else {
		return urlStr
	}
}

func GetMatchedKeys(expression string, keys []string) (map[string][][]int, error) {
	logger_utils.Infof("GetMatchedKeys expression=%v, expression in bytes=%v\n", expression, []byte(expression))
	if !utf8.ValidString(expression) {
		return nil, errors.New("expression is not valid utf8")
	}
	for _, key := range keys {
		logger_utils.Infof("key=%v, key_bytes=%v\n", key, []byte(key))
		if !utf8.ValidString(key) {
			return nil, errors.New("key is not valid utf8")
		}
	}

	regExp, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}

	matchesMap := make(map[string][][]int)

	for _, key := range keys {
		var matches [][]int
		if RegexpMatch(regExp, []byte(key)) {
			matches = regExp.FindAllStringIndex(key, -1)
		} else {
			matches = make([][]int, 0)
		}
		logger_utils.Debugf("key=%v, matches with byte index=%v\n", key, matches)
		convertedMatches, err := convertByteIndexToRuneIndex(key, matches)
		if err != nil {
			return nil, err
		}
		matchesMap[key] = convertedMatches
	}

	return matchesMap, nil
}

func RegexpMatch(regExp *regexp.Regexp, key []byte) bool {
	return regExp.Match(key)
}

// given a matches map, convert the indices from byte index to rune index
func convertByteIndexToRuneIndex(key string, matches [][]int) ([][]int, error) {
	convertedMatches := make([][]int, 0)
	if len(key) == 0 || len(matches) == 0 {
		return matches, nil
	}

	// parse key and build a byte index to rune index map
	indexMap := make(map[int]int)
	byteIndex := 0
	runeIndex := 0
	keyBytes := []byte(key)
	keyLen := len(key)
	for {
		indexMap[byteIndex] = runeIndex
		if byteIndex < keyLen {
			_, runeLen := utf8.DecodeRune(keyBytes[byteIndex:])
			byteIndex += runeLen
			runeIndex++
		} else {
			break
		}
	}

	logger_utils.Debugf("key=%v, indexMap=%v\n", key, indexMap)

	var ok bool
	for _, match := range matches {
		convertedMatch := make([]int, 2)
		convertedMatch[0], ok = indexMap[match[0]]
		if !ok {
			// should not happen
			errMsg := InvalidRuneIndexErrorMessage(key, match[0])
			logger_utils.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
		convertedMatch[1], ok = indexMap[match[1]]
		if !ok {
			// should not happen
			errMsg := InvalidRuneIndexErrorMessage(key, match[1])
			logger_utils.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
		convertedMatches = append(convertedMatches, convertedMatch)
	}

	return convertedMatches, nil
}

func InvalidRuneIndexErrorMessage(key string, index int) string {
	return fmt.Sprintf("byte index, %v, in match for key, %v, is not a starting index for a rune", index, key)
}

func LocalBucketUUID(local_connStr string, bucketName string) (string, error) {
	local_default_pool, err := LocalPool(local_connStr)
	if err != nil {
		return "", err
	}
	bucket, ok := local_default_pool.BucketMap[bucketName]
	if !ok {
		return "", NonExistentBucketError
	}
	return bucket.UUID, nil
}

func ReplicationStatusNotFoundError(topic string) error {
	return fmt.Errorf("Cannot find replication status for topic %v", topic)
}

func BucketNotFoundError(bucketName string) error {
	return fmt.Errorf("Bucket `%v` not found.", bucketName)
}

func GetMemcachedConnection(serverAddr, bucketName string, userAgent string, logger *log.CommonLogger) (*mcc.Client, error) {
	logger.Infof("GetMemcachedConnection serverAddr=%v, bucketName=%v\n", serverAddr, bucketName)
	if serverAddr == "" {
		panic("serverAddr is empty")
	}
	username, password, err := cbauth.GetMemcachedServiceAuth(serverAddr)
	logger.Debugf("memcached auth: username=%v, password=%v, err=%v\n", username, password, err)
	if err != nil {
		return nil, err
	}

	conn, err := base.NewConn(serverAddr, username, password)
	if err != nil {
		return nil, err
	}

	_, err = conn.SelectBucket(bucketName)
	if err != nil {
		return nil, err
	}

	SendHELO(conn, userAgent, base.HELOTimeout, base.HELOTimeout, logger)

	return conn, nil
}

// send helo with specified user agent string to memcached
func SendHELO(client *mcc.Client, userAgent string, readTimeout, writeTimeout time.Duration, logger *log.CommonLogger) {
	helo := ComposeHELORequest(userAgent)

	conn := client.Hijack()
	conn.(net.Conn).SetWriteDeadline(time.Now().Add(writeTimeout))
	_, err := conn.Write(helo.Bytes())
	if err != nil {
		logger.Warnf("Error sending HELO command. userAgent=%v, err=%v.", userAgent, err)
		return
	}

	conn.(net.Conn).SetReadDeadline(time.Now().Add(readTimeout))
	response, err := client.Receive()
	if err != nil {
		logger.Warnf("Received error response from HELO command. userAgent=%v, err=%v.", userAgent, err)
	} else if response.Status != mc.SUCCESS {
		logger.Warnf("Received unexpected response from HELO command. userAgent=%v, response status=%v.", userAgent, response.Status)
	} else {
		logger.Infof("Successfully sent HELO command with userAgent=%v", userAgent)
	}
}

// compose a HELO command with specified user agent string
func ComposeHELORequest(userAgent string) *mc.MCRequest {
	value := make([]byte, 2)
	// tcp nodelay
	binary.BigEndian.PutUint16(value[0:2], 0x03)
	return &mc.MCRequest{
		Key:    []byte(userAgent),
		Opcode: mc.HELLO,
		Body:   value,
	}
}

func GetIntSettingFromSettings(settings map[string]interface{}, settingName string) (int, error) {
	settingObj := GetSettingFromSettings(settings, settingName)
	if settingObj == nil {
		return -1, nil
	}

	setting, ok := settingObj.(int)
	if !ok {
		return -1, fmt.Errorf("Setting %v is of wrong type", settingName)
	}

	return setting, nil
}

func GetStringSettingFromSettings(settings map[string]interface{}, settingName string) (string, error) {
	settingObj := GetSettingFromSettings(settings, settingName)
	if settingObj == nil {
		return "", nil
	}

	setting, ok := settingObj.(string)
	if !ok {
		return "", fmt.Errorf("Setting %v is of wrong type", settingName)
	}

	return setting, nil
}

func GetSettingFromSettings(settings map[string]interface{}, settingName string) interface{} {
	if settings == nil {
		return nil
	}

	setting, ok := settings[settingName]
	if !ok {
		return nil
	}

	return setting
}

func GetMemcachedClient(serverAddr, bucketName string, kv_mem_clients map[string]*mcc.Client, userAgent string, logger *log.CommonLogger) (*mcc.Client, error) {
	client, ok := kv_mem_clients[serverAddr]
	if ok {
		return client, nil
	} else {
		if bucketName == "" {
			panic("unexpected empty bucketName")
		}

		var client, err = GetMemcachedConnection(serverAddr, bucketName, userAgent, logger)
		if err == nil {
			kv_mem_clients[serverAddr] = client
			return client, nil
		} else {
			return nil, err
		}
	}
}

func GetServerVBucketsMap(connStr, bucketName string, bucketInfo map[string]interface{}) (map[string][]uint16, error) {
	vbucketServerMapObj, ok := bucketInfo[base.VBucketServerMapKey]
	if !ok {
		return nil, fmt.Errorf("Error getting vbucket server map from bucket info. connStr=%v, bucketName=%v, bucketInfo=%v\n", connStr, bucketName, bucketInfo)
	}
	vbucketServerMap, ok := vbucketServerMapObj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Vbucket server map is of wrong type. connStr=%v, bucketName=%v, vbucketServerMap=%v\n", connStr, bucketName, vbucketServerMapObj)
	}

	// get server list
	serverListObj, ok := vbucketServerMap[base.ServerListKey]
	if !ok {
		return nil, fmt.Errorf("Error getting server list from vbucket server map. connStr=%v, bucketName=%v, vbucketServerMap=%v\n", connStr, bucketName, vbucketServerMap)
	}
	serverList, ok := serverListObj.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Server list is of wrong type. connStr=%v, bucketName=%v, serverList=%v\n", connStr, bucketName, serverListObj)
	}

	servers := make([]string, len(serverList))
	for index, serverName := range serverList {
		serverNameStr, ok := serverName.(string)
		if !ok {
			return nil, fmt.Errorf("Server name is of wrong type. connStr=%v, bucketName=%v, serverName=%v\n", connStr, bucketName, serverName)
		}
		servers[index] = serverNameStr
	}

	// get vbucket "map"
	vbucketMapObj, ok := vbucketServerMap[base.VBucketMapKey]
	if !ok {
		return nil, fmt.Errorf("Error getting vbucket map from vbucket server map. connStr=%v, bucketName=%v, vbucketServerMap=%v\n", connStr, bucketName, vbucketServerMap)
	}
	vbucketMap, ok := vbucketMapObj.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Vbucket map is of wrong type. connStr=%v, bucketName=%v, vbucketMap=%v\n", connStr, bucketName, vbucketMapObj)
	}

	serverVBMap := make(map[string][]uint16)

	for vbno, indexListObj := range vbucketMap {
		indexList, ok := indexListObj.([]interface{})
		if !ok {
			return nil, fmt.Errorf("Index list is of wrong type. connStr=%v, bucketName=%v, indexList=%v\n", connStr, bucketName, indexListObj)
		}
		if len(indexList) == 0 {
			return nil, fmt.Errorf("Index list is empty. connStr=%v, bucketName=%v, vbno=%v\n", connStr, bucketName, vbno)
		}
		indexFloat, ok := indexList[0].(float64)
		if !ok {
			return nil, fmt.Errorf("Master index is of wrong type. connStr=%v, bucketName=%v, index=%v\n", connStr, bucketName, indexList[0])
		}
		indexInt := int(indexFloat)
		if indexInt < 0 || indexInt >= len(servers) {
			return nil, fmt.Errorf("Master index is out of range. connStr=%v, bucketName=%v, index=%v\n", connStr, bucketName, indexInt)
		}

		server := servers[indexInt]
		var vbList []uint16
		vbList, ok = serverVBMap[server]
		if !ok {
			vbList = make([]uint16, 0)
		}
		vbList = append(vbList, uint16(vbno))
		serverVBMap[server] = vbList
	}
	return serverVBMap, nil
}

// get conflict resolution type setting from bucket info
// default to seqno if not found
func GetConflictResolutionTypeFromBucketInfo(bucketName string, bucketInfo map[string]interface{}) (string, error) {
	conflictResolutionType := base.ConflictResolutionType_Seqno
	conflictResolutionTypeObj, ok := bucketInfo[base.ConflictResolutionTypeKey]
	if ok {
		conflictResolutionType, ok = conflictResolutionTypeObj.(string)
		if !ok {
			return "", fmt.Errorf("ConflictResolutionType on bucket %v is of wrong type.", bucketName)
		}
	}
	return conflictResolutionType, nil
}
