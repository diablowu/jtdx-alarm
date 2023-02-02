package city

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizePrefix(t *testing.T) {

	given := [][]string{
		[]string{"=KL7OG(4)[8];", "=KL7OG"},
		[]string{"KL7OG(4)[8];", "KL7OG"},
		[]string{"=KL7OG(4)[8]", "=KL7OG"},
		[]string{"=KL7OG(4)[8]<344234/dfasdf>~adsfasd~;", "=KL7OG"},
	}

	for _, s := range given {
		assert.Equal(t, s[1], normalizePrefix(s[0]))
	}
}

func Test_removePairTag(t *testing.T) {
	given := [][]string{
		[]string{"(", ")", "=KL7OG(4)[8];", "=KL7OG[8];"},
		[]string{"[", "]", "KL7OG(4)[8];", "KL7OG(4);"},
		[]string{"<", ">", "=KL7OG(4)[8]<344234/dfasdf>~adsfasd~;", "=KL7OG(4)[8]~adsfasd~;"},
		[]string{"~", "~", "=KL7OG(4)[8]<344234/dfasdf>~adsfasd~;", "=KL7OG(4)[8]<344234/dfasdf>;"},
	}

	for _, s := range given {
		assert.Equal(t, s[3], removePairTag(s[2], s[0], s[1]))
	}
}

func Test_loadFromCTYData(t *testing.T) {
	err := LoadFromCTYData("testdata/cty.dat")
	assert.Nil(t, err)
	fmt.Println(len(PrefixData))
	count := 0
	for k, v := range PrefixData {
		count = count + 1
		fmt.Printf("%s => %s \n", k, v)
		if count >= 20 {
			break
		}

	}
}

func Test_finddxcc(t *testing.T) {
	err := LoadFromCTYData("testdata/cty.dat")
	assert.Nil(t, err)
	fmt.Println(len(PrefixData))

	given := [][]string{
		[]string{"BI1NIZ", "BY"},
		[]string{"BI1NIZ/P", "BY"},
		[]string{"KB3UWW", "KH6"},
		[]string{"JS6SCO", "JA"},
		[]string{"BU2FF", "BV"},
	}
	for _, g := range given {
		dxcc := FindDXCC(g[0])
		assert.NotNil(t, dxcc)
		assert.Equal(t, g[1], dxcc.DXCCName)
	}
}
