package am

import "github.com/gertd/go-pluralize"

var pluralizer = pluralize.NewClient()

func Plural(s string) string {
	return pluralizer.Plural(s)
}

func Singular(s string) string {
	return pluralizer.Singular(s)
}

func IsPlural(s string) bool {
	return pluralizer.IsPlural(s)
}

func IsSingular(s string) bool {
	return pluralizer.IsSingular(s)
}

func AddSingularRule(singular, plural string) {
	pluralizer.AddSingularRule(singular, plural)
}

func AddPluralRule(plural, singular string) {
	pluralizer.AddPluralRule(plural, singular)
}

func AddUncountableRule(word string) {
	pluralizer.AddUncountableRule(word)
}
