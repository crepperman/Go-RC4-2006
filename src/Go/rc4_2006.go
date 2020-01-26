package main


import (
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func RC4(key string, plaintext string) []byte {
	k := []rune(key)
	dst := make([]rune, len(plaintext))

	var S [256]int
	KSA(k, &S)
	b := PRGA(plaintext, S, dst)
	bs_UTF8BE, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), b)

	return bs_UTF8BE
}
func KSA(key []rune, S *[256]int) {
	i := 0
	for ; i < 256; i++ {
		S[i] = i
	}

	j := 0
	for i = 0; i < 256; i++ {
		ss := string(key)[i%len(key)]
		j = (j + S[i] + int(ss)) % 256
		S[i], S[j] = S[j], S[i]
	}
}

func PRGA(src1 string, S [256]int, dst []rune) []byte {
	i, j, KK := 0, 0, 0
	ret := []byte{}
	for _, vb := range src1 {
		i = (i + 1) % 256
		j = (j + S[i]) % 256
		S[i], S[j] = S[j], S[i]
		KK = S[(S[i]+S[j])%256]
		bs_UTF16BE1, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder(),[]byte(string(rune(vb))))

		xor, errr := XOR(bs_UTF16BE1, []byte{0, byte(KK)})
		if errr != nil {
			xor, errr := XOR([]byte{bs_UTF16BE1[0], bs_UTF16BE1[1]}, []byte{0, byte(KK)})
			if errr != nil {
				fmt.Println("ERR:", errr)
			} else {
				ret = append(ret, xor[:]...)
			}
			i = (i + 1) % 256
			j = (j + S[i]) % 256
			S[i], S[j] = S[j], S[i]
			KK = S[(S[i]+S[j])%256]
			_xor, errr := XOR([]byte{bs_UTF16BE1[2], bs_UTF16BE1[3]}, []byte{0, byte(KK)})
			if errr != nil {
				fmt.Println("ERR:", errr)
			} else {
				ret = append(ret, _xor[:]...)
			}
		} else {
			ret = append(ret, xor[:]...)
		}
	}
	return ret
}

func XOR(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}
	buf := make([]byte, len(a))
	for i, _ := range a {
		buf[i] = a[i] ^ b[i]
	}
	return buf, nil
}
