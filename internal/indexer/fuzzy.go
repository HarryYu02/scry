package indexer

func levenshteinDist(source, target string) int {
	dpRowNum := len(source) + 1
	dpColNum := len(target) + 1
	dp := make([][]int, dpRowNum)
	for i := range dpRowNum {
		row := make([]int, dpColNum)
		dp[i] = row
	}

	for i := 1; i < dpRowNum; i++ {
		dp[i][0] = i
	}
	for j := 1; j < dpColNum; j++ {
		dp[0][j] = j
	}

	for j := 1; j < dpColNum; j++ {
		for i := 1; i < dpRowNum; i++ {
			cost := 0
			if source[i-1] != target[j-1] {
				cost = 1
			}

			dp[i][j] = min(
				dp[i-1][j] + 1,
				dp[i][j-1] + 1,
				dp[i-1][j-1] + cost,
			)
		}
	}

	return dp[len(source)][len(target)]
}
