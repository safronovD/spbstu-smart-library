package converter

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

type parseFunc func([]gjson.Result)

type jsonParser struct {
	fields map[string]string
}

func NewJsonParser() *jsonParser {
	return &jsonParser{
		fields: make(map[string]string),
	}
}

func (p *jsonParser) String() string {
	str := ""
	for key := range p.fields {
		str += fmt.Sprintf("%v : %v\n", key, p.fields[key])
	}
	return str
}

func (p *jsonParser) Parse(oldJson string) {
	parser := p.createParser()
	bytes := gjson.Get(oldJson, "children").Array()
	for _, value := range bytes {
		temp := value.Map()
		if f, ok := parser[strings.Split(temp["class"].String(), " ")[0]]; ok {
			f(temp["children"].Array())
		}
	}
}

func (p *jsonParser) createParser() map[string]parseFunc {
	return map[string]parseFunc{
		"header": func(results []gjson.Result) {
			p.getAuthor(results)
		},
		"bibliographicDescription": func(results []gjson.Result) {
			p.getTitle(results)
			p.getSpecialization(results)
			p.getLinks(results)
		},
		"subjects": func(results []gjson.Result) {
			p.getKeyWords(results)
		},
		"abstract": func(results []gjson.Result) {
			p.getDescription(results)
		},
		"additional": func(results []gjson.Result) {
			p.getMetaInfo(results)
		},
	}
}
func (p *jsonParser) getAuthor(temp []gjson.Result) {
	p.fields["author"] = temp[0].Map()["children"].Array()[0].Map()["html"].String()
	p.fields["author"] += " "
	p.fields["author"] += temp[0].Map()["children"].Array()[2].Map()["html"].String()
}

func (p *jsonParser) getTitle(temp []gjson.Result) {
	if temp := temp[0].Map()["children"].Array()[0].Map(); temp["class"].String() == "title" {
		p.fields["title"] = temp["children"].Array()[0].Map()["html"].String()
	} else if temp["class"].String() == "publication" {
		p.fields["dateOfPublication"] = temp["children"].Array()[2].Map()["html"].String()
	}
}

func (p *jsonParser) getSpecialization(temp []gjson.Result) {
	if temp := temp[0].Map()["children"].Array()[0].Map(); temp["class"].String() == "title" {
		tempString := ""
		for _, value := range temp["children"].Array() {
			if temp := value.Map(); temp["class"].String() == "otherInfo" {
				tempString += temp["html"].String()
				tempString += ";"
			} else if temp["class"].String() == "parallelTitleProper" {
				p.fields["specializationEn"] = temp["html"].String()
			}
		}
		str := strings.Split(tempString, ";")
		p.fields["worktype"] = str[0]
		p.fields["specializationBaseRu"] = str[1]
		p.fields["specializationExtRu"] = str[2]
	}
}

func (p *jsonParser) getLinks(temp []gjson.Result) {
	for _, list := range temp[0].Map()["children"].Array() {
		if temp := list.Map(); temp["class"].String() == "notes" {
			p.fields["pdfLink"] = temp["children"].Array()[3].Map()["children"].Array()[1].Map()["href"].String()
			for index := range temp["children"].Array() {
				if temp, ok := temp["children"].Array()[index].Map()["children"]; ok {
					if tmp := temp.Array()[0].Map(); tmp["tag"].String() == "a" {
						switch tmp["html"].String() {
						case "Отчет о проверке на объем и корректность внешних заимствований":
							p.fields["structureReport"] = tmp["href"].String()
						case "Отзыв руководителя":
							p.fields["commentLink"] = tmp["href"].String()
						case "Рецензия":
							p.fields["reviewLink"] = tmp["href"].String()
						}
					}
				}
			}
		}
	}
}

func (p *jsonParser) getKeyWords(temp []gjson.Result) {
	p.fields["keyWords"] = ""
	for _, list := range temp {
		if list.Map()["class"].String() == "topicalSubject subject" {
			p.fields["keyWords"] += list.Map()["children"].Array()[0].Map()["html"].String()
			p.fields["keyWords"] += ", "
		} else if list.Map()["class"].String() == "uncontrolledSubject subject" {
			for _, value := range list.Map()["children"].Array() {
				p.fields["keyWords"] += value.Map()["html"].String()
			}
		}
	}
}

func (p *jsonParser) getMetaInfo(temp []gjson.Result) {
	for _, list := range temp {
		switch temp := list.Map(); temp["class"].String() {
		case "personalName relation-727":
			p.fields["teacher"] = temp["children"].Array()[0].Map()["html"].String()
			p.fields["teacher"] += " "
			p.fields["teacher"] += temp["children"].Array()[2].Map()["html"].String()
		case "personalName relation-570":
			p.fields["controlPerson"] = temp["children"].Array()[0].Map()["html"].String()
			p.fields["controlPerson"] += " "
			p.fields["controlPerson"] += temp["children"].Array()[2].Map()["html"].String()
		case "corporateName":
			p.fields["university"] = temp["children"].Array()[0].Map()["html"].String()
			p.fields["institute"] = temp["children"].Array()[2].Map()["html"].String()
		}
	}
}

func (p *jsonParser) getDescription(temp []gjson.Result) {
	p.fields["descriptionRu"] = temp[0].Map()["html"].String()
	p.fields["descriptionEn"] = temp[1].Map()["html"].String()
}
