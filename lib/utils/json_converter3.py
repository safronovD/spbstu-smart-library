import json
import sys

# args 
# 1 - input json file
# 2 - output json file
        
def setValue(dictV, key, value):
    dictV[key] = value

def findSubjects(dictV):
    for json in dictV["children"]:
        if json["class"] == "subjects":
            return json["children"]
    return []

alphabet=set('абвгдеёжзийклмнопрстуфхцчшщъыьэюя')

def getValue(dictV, tags, skip=0):
    if len(tags) == 1:
        for item in dictV:
            if item["class"] == tags[0] and skip > 0:
                skip = skip - 1
                continue
            if item["class"] == tags[0]:
                try:
                    return item["html"]
                except:
                    return item["children"]
        return ""

    if type(dictV) == list:
        for item in dictV:
            if item["class"] == tags[0]:
                return getValue(item["children"], tags[1:], skip)

    if type(dictV) == dict:
        if dictV["class"] == tags[0]:
            return getValue(dictV["children"], tags[1:], skip)
        else:
            return ""


with open(sys.argv[1]) as json_file:
    json_obj = json.load(json_file)

{
    "author": "Колясов Олег Аркадьевич",
    "title": "Техносферная безопасность при реконструкции линейной части магистрального газопровода",
    "workType": "выпускная квалификационная работа бакалавра",
    "specializationBaseRu": "20.03.01 - Техносферная безопасность ",
    "specializationExtRu": " 20.03.01_07 - Техносферная безопасность (общий профиль)",
    "specializationEn": " Technosphere safety in the reconstruction of the linear part of the gas pipeline",
    "publicationYear": "2019",
    "pdfLink": "http://elib.spbstu.ru/dl/3/2019/vr/vr19-3567.pdf",
    "commentLink": "http://elib.spbstu.ru/dl/3/2019/vr/rev/vr19-3567-o.pdf",
    "reviewLink": "http://elib.spbstu.ru/dl/3/2019/vr/rev/vr19-3567-a.pdf",
    "structureReport": "",
    "descriptionRu": "В настоящей выпускной квалификационной работе рассмотрена тема обеспечения техносферной безопасности при реконструкции линейной части магистрального газопровода.\nЦелью выполнения работы является разработка необходимых технологий и мероприятий, обеспечивающих охрану окружающей среды, промышленную безопасность и охрану труда, при проведении реконструкции участка магистрального газопровода.",
    "descriptionEn": "In this final qualifying work is overviewed the development process of technosphere safety of reconstruction of the Gas Pipeline.\nThe purpose of the work is to develop the necessary technologies and measures to ensure environmental protection, industrial safety and labor protection during the reconstruction of the main gas pipeline section. ",
    "teacher": "Ефремов Сергей Владимирович",
    "controlPerson": "Ефремов Сергей Владимирович",
    "institute": "Инженерно-строительный институт",
    "univer": "Санкт-Петербургский политехнический университет Петра Великого",
    "transactionDate": "2019-11-01T00:00:00Z",
    "topicalSubjects": [],
    "keyWordsRu": [
        "безопасность",
        "газ",
        "трубопровод",
        "техносфера",
        "окружающая среда",
        "охрана труда"
    ],
    "keyWordsEn": [
        "safety",
        "gas",
        "pipeline",
        "technosphere",
        "environment",
        "labor protection"
    ]
}
ops = []

new_json_obj = {}

ops.append(["author", lambda : getValue(json_obj, ["bibliographicRecord", "header", "personalName relation-070", "entry"]) + " " + getValue(json_obj, ["bibliographicRecord", "header", "personalName relation-070", "expansionOfInitials"])])
ops.append(["title", lambda : getValue(json_obj, ["bibliographicRecord", "bibliographicDescription monograph multilevel hLevel_0", "general", "title", "titleProper"])])
ops.append(["workType", lambda : getValue(json_obj, ["bibliographicRecord", "bibliographicDescription monograph multilevel hLevel_0", "general", "title", "otherInfo"])])
ops.append(["specializationBaseRu", lambda : getValue(json_obj, ["bibliographicRecord", "bibliographicDescription monograph multilevel hLevel_0", "general", "title", "otherInfo"], 1).split(';')[0]])
ops.append(["specializationExtRu", lambda : getValue(json_obj, ["bibliographicRecord", "bibliographicDescription monograph multilevel hLevel_0", "general", "title", "otherInfo"], 1).split(';')[1]])
ops.append(["specializationEn", lambda : getValue(json_obj, ["bibliographicRecord", "bibliographicDescription monograph multilevel hLevel_0", "general", "title", "parallelTitleProper"])])
ops.append(["publicationYear", lambda : getValue(json_obj, ["bibliographicRecord", "bibliographicDescription monograph multilevel hLevel_0", "general", "publication", "dateOfPublication"])])
ops.append(["pdfLink", lambda : getValue(json_obj, ["bibliographicRecord", "bibliographicDescription monograph multilevel hLevel_0", "general", "notes", "note"], 3)[1]["href"] ])
ops.append(["commentLink", lambda : getValue(json_obj, ["bibliographicRecord", "bibliographicDescription monograph multilevel hLevel_0", "general", "notes", "note"], 5)[0]["href"] ])
ops.append(["reviewLink", lambda : getValue(json_obj, ["bibliographicRecord", "bibliographicDescription monograph multilevel hLevel_0", "general", "notes", "note"], 6)[0]["href"] ])
ops.append(["structureReport", lambda : getValue(json_obj, ["bibliographicRecord", "bibliographicDescription monograph multilevel hLevel_0", "general", "notes", "note"], 7)[0]["href"] ])
ops.append(["descriptionRu", lambda : getValue(json_obj, ["bibliographicRecord", "abstract endsWithFullStop", "abstract endsWithFullStop"])])
ops.append(["descriptionEn", lambda : getValue(json_obj, ["bibliographicRecord", "abstract endsWithFullStop", "abstract endsWithFullStop"], 1)])
ops.append(["teacher", lambda : getValue(json_obj, ["bibliographicRecord", "additional", "personalName relation-727", "entry"]) + " " + getValue(json_obj, ["bibliographicRecord", "additional", "personalName relation-727", "expansionOfInitials"])])
ops.append(["controlPerson", lambda : getValue(json_obj, ["bibliographicRecord", "additional"])[2]["children"][0]["html"] + " " + getValue(json_obj, ["bibliographicRecord", "additional"])[2]["children"][2]["html"]])
ops.append(["institute", lambda : getValue(json_obj, ["bibliographicRecord", "additional", "corporateName", "entry"])])
ops.append(["univer", lambda : getValue(json_obj, ["bibliographicRecord", "additional", "corporateName", "subdivision"])])
ops.append(["transactionDate", lambda : getValue(json_obj, ["bibliographicRecord", "originatingSource", "originalCataloguingAgency", "dateOfTransaction"])])

for op in ops:
    try:
        new_json_obj[op[0]] = op[1]()        
    except:
        new_json_obj[op[0]] = ""


new_json_obj["topicalSubjects"] = []
new_json_obj["keyWordsRu"] = []
new_json_obj["keyWordsEn"] = []

for subject in findSubjects(json_obj):
    if subject["class"] == "topicalSubject subject":
        new_json_obj["topicalSubjects"].append(subject["children"][0]["html"])
        continue

    if subject["class"] == "uncontrolledSubject subject":
        for sub in subject["children"]:
            if sub["class"] == "subjectTerm":
                if sub["html"][0].lower() in alphabet:
                    new_json_obj["keyWordsRu"].append(sub["html"])
                else:
                    new_json_obj["keyWordsEn"].append(sub["html"])

 
if new_json_obj['transactionDate'] != "":
    try:
        dates = new_json_obj['transactionDate'].split('.')
        new_json_obj['transactionDate'] = "{}-{}-{}T00:00:00Z".format(dates[2], dates[1], dates[0])
    except:
        pass
else:
    #Wrong date
    new_json_obj['transactionDate'] = "{}-{}-{}T00:00:00Z".format("1970", "01", "01")

try:
    if new_json_obj["descriptionRu"][0].lower() not in alphabet:
        tmp = new_json_obj["descriptionRu"]
        new_json_obj["descriptionRu"] = new_json_obj["descriptionEn"]
        new_json_obj["descriptionEn"] = tmp
except:
    pass

for key in new_json_obj:
    if new_json_obj[key] == "" or new_json_obj[key] == None:
        new_json_obj[key] = "unknown"

with open(sys.argv[2], 'w', encoding='utf-8') as f:
    f.write(json.dumps(new_json_obj, indent=2, ensure_ascii=False))
