package converter

import (
	"strings"
)

type fillFunc func(data *jsonParser)

type jsonFiller struct {
	newJson *prettyJson
}

func NewJsonFiller() *jsonFiller {
	return &jsonFiller{
		newJson: NewPrettyJson(),
	}
}

func (f *jsonFiller) Fill(schema string, data *jsonParser) *prettyJson {
	filler := f.createFiller()
	// TODO: Recreate schema logic
	for _, attribute := range strings.Split(schema, ", ") {
		if fill, ok := filler[attribute]; ok {
			fill(data)
		}
	}
	return f.newJson
}

func (f *jsonFiller) createFiller() map[string]fillFunc {
	return map[string]fillFunc{
		"author":         func(data *jsonParser) { f.fillAuthor(data) },
		"worktype":       func(data *jsonParser) { f.fillWorkType(data) },
		"title":          func(data *jsonParser) { f.fillTitle(data) },
		"specialization": func(data *jsonParser) { f.fillSpecialization(data) },
		"description":    func(data *jsonParser) { f.fillDescription(data) },
		"personalities":  func(data *jsonParser) { f.fillPersonalities(data) },
		"university":     func(data *jsonParser) { f.fillUniversity(data) },
		"keyWords":       func(data *jsonParser) { f.fillKeyWords(data) },
		"links":          func(data *jsonParser) { f.fillLinks(data) },
	}
}

func (f *jsonFiller) fillAuthor(data *jsonParser) {
	if item, ok := data.fields["author"]; ok {
		f.newJson.Author = item
	}
}

func (f *jsonFiller) fillWorkType(data *jsonParser) {
	if item, ok := data.fields["worktype"]; ok {
		f.newJson.Worktype = item
	}
}

func (f *jsonFiller) fillTitle(data *jsonParser) {
	if item, ok := data.fields["title"]; ok {
		f.newJson.Title = item
	}
}

func (f *jsonFiller) fillSpecialization(data *jsonParser) {
	if item, ok := data.fields["specializationBaseRu"]; ok {
		f.newJson.SpecializationBaseRu = item
	}
	if item, ok := data.fields["specializationExtRu"]; ok {
		f.newJson.SpecializationExtRU = item
	}
	if item, ok := data.fields["specializationEn"]; ok {
		f.newJson.SpecializationEn = item
	}
}

func (f *jsonFiller) fillDescription(data *jsonParser) {
	if item, ok := data.fields["descriptionRu"]; ok {
		f.newJson.DescriptionRu = item
	}
	if item, ok := data.fields["descriptionEn"]; ok {
		f.newJson.DescriptionEn = item
	}
}

func (f *jsonFiller) fillPersonalities(data *jsonParser) {
	if item, ok := data.fields["teacher"]; ok {
		f.newJson.Teacher = item
	}
	if item, ok := data.fields["controlPerson"]; ok {
		f.newJson.ControlPerson = item
	}
}

func (f *jsonFiller) fillUniversity(data *jsonParser) {
	if item, ok := data.fields["university"]; ok {
		f.newJson.University = item
	}
	if item, ok := data.fields["institute"]; ok {
		f.newJson.Institute = item
	}
}

func (f *jsonFiller) fillKeyWords(data *jsonParser) {
	if item, ok := data.fields["keyWords"]; ok {
		f.newJson.KeyWords = strings.Split(item, ", ")
	}
}

func (f *jsonFiller) fillLinks(data *jsonParser) {
	if item, ok := data.fields["commentLink"]; ok {
		f.newJson.Links["commentLink"] = item
	}
	if item, ok := data.fields["pdfLink"]; ok {
		f.newJson.Links["pdfLink"] = item
	}
	if item, ok := data.fields["structureReport"]; ok {
		f.newJson.Links["structureReport"] = item
	}
}
