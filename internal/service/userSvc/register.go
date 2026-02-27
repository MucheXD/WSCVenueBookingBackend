package userSvc

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
)

func RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {
	// 校验前端传入的数据是否合法，避免不必要的数据库查询和潜在的安全风险
	// 这里的校验规则可以根据实际需求进行调整
	if !isValidUsername(user.Username) {
		return nil, fmt.Errorf("%w: username invalid", ErrRegisterInfoInvalid)
	}
	if !isValidSchoolID(user.SchoolID) {
		return nil, fmt.Errorf("%w: school ID invalid", ErrRegisterInfoInvalid)
	}
	if !isValidPasswordHash(user.PasswordHash) {
		return nil, fmt.Errorf("%w: password hash invalid", ErrRegisterInfoInvalid)
	}
	if !isValidPasswordSalt(user.PasswordSalt) {
		return nil, fmt.Errorf("%w: password salt invalid", ErrRegisterInfoInvalid)
	}

	// 用户名去重
	usrnameExist, err := repository.IsUsernameExists(user.Username)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrCheckUsernameExistsInDB, err)
	}
	if usrnameExist {
		return nil, ErrUsernameAlreadyExists
	}
	// 创建用户
	user.UID = generateNewUserID()
	user.RegisterTime = time.Now().UTC()
	user.PhoneNumber = ""
	user.PermMap = systemPermission.Map(systemPermission.RegisterDefault)
	user.PermVAGID = 0
	err = repository.CreateNewUser(user)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrCreateUserInDB, err)
	}
	return user, nil
}

// 生成新的用户ID，确保唯一性，具有一定的时间排序特性，同时包含随机部分以增加复杂度和安全性
func generateNewUserID() string {
	// UID 生成方案：偏移的毫秒UNIX时间戳(14位->9字符) + 随机字符串(3字符)
	// TTTTTTTTTRRR
	const timestampOffset = 1145141919810               // perfect offset fit 13 digits for unix, surely random xd
	const alphabet = "UXQG4Y7261EATRJ5VNL3BD9PMZKSH8FC" // 32 chars, without 0 O I W

	tsPart := convert10To32Inverse(
		(time.Now().UnixMilli()-timestampOffset)%0x1FFFFFFFFFFF, alphabet) // 0x1FFFFFFFFFFF = 32^9 - 1, 9 chars for timestamp part
	if len(tsPart) < 9 { // 时间部分可能长度不足，补齐即可
		padding := strings.Repeat("U", 9-len(tsPart))
		tsPart = padding + tsPart
	}
	randPart := utils.GenerateRandomString(3, alphabet)
	ostr := tsPart + randPart
	// TTTTTTTTTRRR -> TTTTTTRTRTRT
	return ostr[0:6] + ostr[9:10] + ostr[6:7] + ostr[10:11] + ostr[7:8] + ostr[11:12] + ostr[8:9]
}

// 使用指定字符表将整数转换为逆序32进制字符串，用于UID生成
func convert10To32Inverse(n int64, alphabet string) string {
	if n == 0 {
		return string(alphabet[0])
	}
	if n < 0 {
		n = -n
	}
	var result []byte
	// 辗转相除法
	for n > 0 {
		remainder := n % 32
		result = append(result, alphabet[remainder])
		n = n / 32
	}
	return string(result)
}

// 预编译正则表达式以提高性能
// ^[a-zA-Z0-9_]+$ 表示只允许字母、数字和下划线
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// IsValidUsername 校验用户名是否合法
// 规则：4-20个字符，仅限英文字母、数字和下划线，不含特殊符号
func isValidUsername(username string) bool {
	length := len(username)
	if length < 4 || length > 20 {
		return false
	}
	return usernameRegex.MatchString(username)
}

// 这里可以根据实际需求定义学校ID的格式规则
// 必须是数字，长度为12位
var schoolIDRegex = regexp.MustCompile(`^\d{12}$`)

func isValidSchoolID(schoolID string) bool {
	return schoolIDRegex.MatchString(schoolID)
}

var passwordHashRegex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
var passwordSaltRegex = regexp.MustCompile(`^[a-zA-Z0-9]{16}$`)

func isValidPasswordHash(passwordHash string) bool {
	return passwordHashRegex.MatchString(passwordHash)
}
func isValidPasswordSalt(passwordSalt string) bool {
	return passwordSaltRegex.MatchString(passwordSalt)
}
