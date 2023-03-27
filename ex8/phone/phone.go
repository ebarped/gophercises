package phone

import "unicode"

type Phone struct {
	Id     int
	Number string
}

func New(id int, number string) Phone {
	return Phone{
		Id:     id,
		Number: number,
	}
}

// Remove removes a Phone from a Phone slice if the number matches
func Remove(s []Phone, p Phone) []Phone {
	for i, v := range s {
		if v.Number == p.Number {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// NormalizeNumber formats a phone number to match the `##########“ pattern
func NormalizeNumber(p Phone) Phone {
	var result string

	for _, c := range p.Number {
		if unicode.IsNumber(c) {
			result += string(c)
		}
	}

	return Phone{Id: p.Id, Number: result}
}
