package callsign

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestExtractCallSignFromMessage(t *testing.T) {

	de, dx, err := ExtractCallSignFromMessage("CQ DFD         ", false)
	assert.Nil(t, err)
	fmt.Println(de, dx)
}

func TestReq(t *testing.T) {
	var cqDirectionReg = regexp.MustCompile("^[A-Z]{2}$")

	fmt.Println(cqDirectionReg.MatchString("DX"))
	fmt.Println(cqDirectionReg.MatchString("DXCC"))
	fmt.Println(cqDirectionReg.MatchString("EU"))
	fmt.Println(cqDirectionReg.MatchString("1D"))
}
