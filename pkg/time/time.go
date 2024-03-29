package time

import "time"

// WrapperTime は現在時刻を提供するインターフェースです。
type WrapperTime interface {
	Now() time.Time
}

// CustomTime はTimeProviderインターフェースを実装します。
type CustomTime struct{}

// Now は現在の時刻を返します。
func (c CustomTime) Now() time.Time {
	return time.Now()
}
