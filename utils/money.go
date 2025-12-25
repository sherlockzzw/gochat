package utils

import (
	"fmt"
	"strconv"
)

// Money 金额工具函数
// 所有金额在数据库中以"分"为单位存储(int64)
// 1元 = 100分

// YuanToFen 元转分
// 例如: 1.23元 -> 123分
// 注意: 由于浮点数精度问题，建议使用字符串输入
func YuanToFen(yuan float64) int64 {
	return int64(yuan * 100)
}

// YuanToFenFromString 从字符串元转分(推荐使用，避免精度问题)
// 例如: "1.23" -> 123分, "100" -> 10000分
func YuanToFenFromString(yuanStr string) (int64, error) {
	yuan, err := strconv.ParseFloat(yuanStr, 64)
	if err != nil {
		return 0, err
	}
	return YuanToFen(yuan), nil
}

// FenToYuan 分转元
// 例如: 123分 -> 1.23元
func FenToYuan(fen int64) float64 {
	return float64(fen) / 100.0
}

// FormatMoney 格式化金额显示(分转元，保留2位小数)
// 例如: 123分 -> "1.23", 10000分 -> "100.00"
func FormatMoney(fen int64) string {
	yuan := FenToYuan(fen)
	return fmt.Sprintf("%.2f", yuan)
}

// FormatMoneyYuan 格式化金额显示(分转元，保留2位小数，带"元"单位)
// 例如: 123分 -> "1.23元"
func FormatMoneyYuan(fen int64) string {
	return FormatMoney(fen) + "元"
}
