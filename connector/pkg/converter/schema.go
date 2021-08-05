package converter

type prettyJson struct {
	Author               string            `json:"author,omitempty"`
	Worktype             string            `json:"work_type,omitempty"`
	Title                string            `json:"title,omitempty"`
	SpecializationBaseRu string            `json:"specialization_base_ru,omitempty"`
	SpecializationExtRU  string            `json:"specialization_ext_ru,omitempty"`
	SpecializationEn     string            `json:"specialization_en,omitempty"`
	DescriptionRu        string            `json:"description_ru,omitempty"`
	DescriptionEn        string            `json:"description_en,omitempty"`
	Teacher              string            `json:"teacher,omitempty"`
	ControlPerson        string            `json:"control_person,omitempty"`
	University           string            `json:"university,omitempty"`
	Institute            string            `json:"institute,omitempty"`
	KeyWords             []string          `json:"key_words,omitempty"`
	Links                map[string]string `json:"links,omitempty"`
}

func NewPrettyJson() *prettyJson {
	return &prettyJson{Links: make(map[string]string)}
}
