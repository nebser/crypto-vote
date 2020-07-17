package keyfiles

type KeyFiles struct {
	PrivateKeyFile string
	PublicKeyFile  string
}

type KeyFilesList []KeyFiles

func (kfl KeyFilesList) Find(criteria func(KeyFiles) bool) (KeyFiles, bool) {
	for _, k := range kfl {
		if criteria(k) {
			return k, true
		}
	}
	return KeyFiles{}, false
}

func (kfl KeyFilesList) Filter(criteria func(KeyFiles) bool) (result KeyFilesList) {
	for _, k := range kfl {
		if criteria(k) {
			result = append(result, k)
		}
	}
	return
}
