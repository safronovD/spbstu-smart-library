import json
import sys

# args 
# 1 - input json file
# 2 - output json file

def normalize(json, class_type):
    if type(json) == list:
        new_json = []
        for item in json:
            if "class" in item and item["class"] == "punct":
                continue
            if "html" in item and class_type != "notes":   #ES can't handle list with string and dicts
                new_json.append(item["html"])
                continue
            if "children" in item:
                new_json.append({item["class"] : normalize(item["children"], item["class"])})
                continue
        return new_json    
    elif type(json) == dict:
        if "html" in json:
            return json["html"]
        if "children" in json:
            return {json["class"] : normalize(json["children"], json["class"])}
        

with open(sys.argv[1]) as json_file:
    json_obj = json.load(json_file)

new_json_obj = {json_obj["class"] : normalize(json_obj["children"], json_obj["class"])}

with open(sys.argv[2], 'w', encoding='utf-8') as f:
    f.write(json.dumps(new_json_obj, indent=2, ensure_ascii=False))
