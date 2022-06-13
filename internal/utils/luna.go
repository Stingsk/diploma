package utils

const (
	lunaMultiplier = 10
)

func Valid(number int64) bool {
	return (number%lunaMultiplier+checksum(number/lunaMultiplier))%lunaMultiplier == 0
}

func checksum(number int64) int64 {
	var luna int64

	for i := 0; number > 0; i++ {
		cur := number % lunaMultiplier

		if i%2 == 0 {
			cur *= 2
			if cur > lunaMultiplier-1 {
				cur = cur%lunaMultiplier + cur/lunaMultiplier
			}
		}

		luna += cur
		number /= lunaMultiplier
	}

	return luna % lunaMultiplier
}
