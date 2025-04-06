package season

import "fmt"

type ID int

func (i ID) Season() (Season, error) {
	switch i {
	case 1:
		return Winter, nil
	case 2:
		return Spring, nil
	case 3:
		return Summer, nil
	case 4:
		return Autumn, nil
	}
	return "", fmt.Errorf("invalid season %d, season id can only be one of [1,2,3,4]", i)
}
