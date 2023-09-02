package sort

func insertionSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		for j := i; j > 0 && arr[j-1] > arr[j]; j-- {
			arr[j-1], arr[j] = arr[j], arr[j-1]
		}
	}
}

func quickSort(arr []int) {
	if len(arr) <= 1 {
		return
	}

	pivot := arr[0]
	i, j := 1, len(arr)-1
	for i <= j {
		switch {
		case arr[i] < pivot:
			i++
		case arr[j] >= pivot:
			j--
		default:
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	if j > 0 {
		arr[0], arr[j] = arr[j], arr[0]
	}
	if j >= 2 {
		quickSort(arr[:j])
	}
	if j <= len(arr)-3 {
		quickSort(arr[j+1:])
	}
}
