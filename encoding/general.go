package encoding

// valueLength caclulates how large the necessary value needs to be to fit in the appropriate
// packet length
func valueLength(value uint32) int {
	/* length of enumerated is variable, as per 20.2.11 */
	if value < 0x100 {
		return size8
	} else if value < 0x10000 {
		return size16
	} else if value < 0x1000000 {
		return size24
	}
	return size32
}
