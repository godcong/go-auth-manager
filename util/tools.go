package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/godcong/wego/util"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"io"
	"io/ioutil"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

/*RandomKind RandomKind */
type RandomKind int

/*random kinds */
const (
	RandomNum      RandomKind = iota // 纯数字
	RandomLower                      // 小写字母
	RandomUpper                      // 大写字母
	RandomLowerNum                   // 数字、小写字母
	RandomUpperNum                   // 数字、大写字母
	RandomAll                        // 数字、大小写字母
)

/*RandomString defines */
var (
	RandomString = map[RandomKind]string{
		RandomNum:      "0123456789",
		RandomLower:    "abcdefghijklmnopqrstuvwxyz",
		RandomUpper:    "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		RandomLowerNum: "0123456789abcdefghijklmnopqrstuvwxyz",
		RandomUpperNum: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		RandomAll:      "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	}
)

//GenerateRandomString2 随机字符串
func GenerateRandomString2(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

//GenerateRandomString 随机字符串
func GenerateRandomString(size int, kind ...RandomKind) string {
	bytes := RandomString[RandomAll]
	if kind != nil {
		if k, b := RandomString[kind[0]]; b == true {
			bytes = k
		}
	}
	var result []byte
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

// UnmarshalJSON ...
func UnmarshalJSON(reader io.Reader, v interface{}) error {
	bytes, err := ioutil.ReadAll(reader)
	log.Info(string(bytes))
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	err = jsoniter.Unmarshal(bytes, v)
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON ...
func MarshalJSON(v interface{}) ([]byte, error) {
	bytes, err := jsoniter.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

// DecryptJWT ...
func DecryptJWT(key []byte, token string) (string, error) {
	cl := jwt.Claims{}
	signed, err := jwt.ParseSigned(token)
	if err != nil {
		return "", err
	}

	err = signed.Claims(key, &cl)
	if err != nil {
		return "", err
	}

	expected := jwt.Expected{
		Issuer: "godcong",
		Time:   time.Now(),
	}

	err = cl.Validate(expected)
	if err != nil {
		return "", err
	}

	return cl.Subject, nil
}

// EncryptJWT ...
func EncryptJWT(key []byte, sub []byte, expiry ...time.Duration) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", nil
	}
	if expiry == nil {
		expiry = []time.Duration{time.Hour * 14 * 24}
	}

	cl := jwt.Claims{
		Subject:   string(sub),
		Issuer:    "godcong",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Expiry:    jwt.NewNumericDate(time.Now().Add(expiry[0])),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        GenerateRandomString(16),
	}

	raw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	return raw, err
}

// SHA256 ...
func SHA256(v, key, salt string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(v))
	m.Write([]byte("."))
	m.Write([]byte(salt))
	return strings.ToUpper(fmt.Sprintf("%x", m.Sum(nil)))
}

// StructureName ...
func StructureName(s interface{}) string {
	return reflect.Indirect(reflect.ValueOf(s)).Type().Name()
}

// GenSpreadSign ...
func GenSpreadSign() string {
	return util.GenMD5(uuid.New().String())
}
