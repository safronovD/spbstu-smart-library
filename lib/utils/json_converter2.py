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

with open(sys.argv[1]) as json_file:
    json_obj = json.load(json_file)

ops = []

new_json_obj = {}

ops.append(['author', lambda : setValue(new_json_obj, 'author', json_obj["children"][0]["children"][0]["children"][0]["html"] + " " + json_obj["children"][0]["children"][0]["children"][2]["html"])])
ops.append(['title', lambda : setValue(new_json_obj, 'title', json_obj["children"][1]["children"][0]["children"][0]["children"][0]["html"])])
ops.append(['workType', lambda : setValue(new_json_obj, 'workType', json_obj["children"][1]["children"][0]["children"][0]["children"][5]["html"])])
ops.append(['specializationBaseRu', lambda : setValue(new_json_obj, 'specializationBaseRu', json_obj["children"][1]["children"][0]["children"][0]["children"][7]["html"].split(";")[0])])
ops.append(['specializationExtRu', lambda : setValue(new_json_obj, 'specializationExtRu', json_obj["children"][1]["children"][0]["children"][0]["children"][7]["html"].split(";")[1])])
ops.append(['specializationEn', lambda : setValue(new_json_obj, 'specializationEn', json_obj["children"][1]["children"][0]["children"][0]["children"][9]["html"])])
ops.append(['publicationYear', lambda : setValue(new_json_obj, 'publicationYear', json_obj["children"][1]["children"][0]["children"][2]["children"][2]["html"])])
ops.append(['pdfLink', lambda : setValue(new_json_obj, 'pdfLink', json_obj["children"][1]["children"][0]["children"][3]["children"][3]["children"][1]["href"])])
ops.append(['commentLink', lambda : setValue(new_json_obj, 'commentLink', json_obj["children"][1]["children"][0]["children"][3]["children"][5]["children"][0]["href"])])
ops.append(['reviewLink', lambda : setValue(new_json_obj, 'reviewLink', json_obj["children"][1]["children"][0]["children"][3]["children"][6]["children"][0]["href"])])
ops.append(['structureReport', lambda : setValue(new_json_obj, 'structureReport', json_obj["children"][1]["children"][0]["children"][3]["children"][7]["children"][0]["href"])])
ops.append(['descriptionRu', lambda : setValue(new_json_obj, 'descriptionRu', json_obj["children"][2]["children"][0]["html"])])
ops.append(['descriptionEn', lambda : setValue(new_json_obj, 'descriptionEn', json_obj["children"][2]["children"][1]["html"])])
ops.append(['teacher', lambda : setValue(new_json_obj, 'teacher', json_obj["children"][3]["children"][1]["children"][0]["html"] + " " + json_obj["children"][3]["children"][1]["children"][2]["html"])])
ops.append(['controlPerson', lambda : setValue(new_json_obj, 'controlPerson', json_obj["children"][3]["children"][2]["children"][0]["html"] + " " + json_obj["children"][3]["children"][2]["children"][2]["html"])])
ops.append(['institute', lambda : setValue(new_json_obj, 'institute', json_obj["children"][3]["children"][3]["children"][2]["html"])])
ops.append(['univer', lambda : setValue(new_json_obj, 'univer', json_obj["children"][3]["children"][3]["children"][0]["html"])])
ops.append(['transactionDate', lambda : setValue(new_json_obj, 'transactionDate', json_obj["children"][5]["children"][0]["children"][1]["html"])])

for op in ops:
    try:
        op[1]()
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
                if sub["html"][0] in alphabet:
                    new_json_obj["keyWordsRu"].append(sub["html"])
                else:
                    new_json_obj["keyWordsEn"].append(sub["html"])
if new_json_obj['transactionDate'] != "":
    try:
        dates = new_json_obj['transactionDate'].split('.')
        new_json_obj['transactionDate'] = "{}-{}-{}T00:00:00Z".format(dates[2], dates[1], dates[0])
    except:
        pass

with open(sys.argv[2], 'w', encoding='utf-8') as f:
    f.write(json.dumps(new_json_obj, indent=2, ensure_ascii=False))
