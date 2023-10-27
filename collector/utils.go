package collector

func SliceDifference[T comparable](slice1, slice2 []T) []T {
	result := []T{}
	set := make(map[T]struct{})

	for _, v := range slice2 {
		set[v] = struct{}{}
	}

	for _, v := range slice1 {
		if _, found := set[v]; !found {
			result = append(result, v)
		}
	}

	return result
}

// Utility function to compute percentile
func Percentile(sortedValues []float64, p int) float64 {
	index := (float64(p) / 100) * float64(len(sortedValues))
	if int(index) == len(sortedValues) {
		return sortedValues[len(sortedValues)-1]
	}
	return sortedValues[int(index)]
}

// Utility function to compute sum
func Sum(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum
}

// Utility function to compute max
func Max(values []float64) float64 {
	max := values[0]
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

// Utility function to compute min
func Min(values []float64) float64 {
	min := values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}

// Utility function to compute avg
func Avg(values []float64) float64 {
	return Sum(values) / float64(len(values))
}
