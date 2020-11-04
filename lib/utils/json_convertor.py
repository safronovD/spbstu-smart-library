import json
import sys

# args 
# 1 - input json file
# 2 - output json file

def normalize(json):
    if type(json) == list:
        new_json = []
        for item in json:
            if "class" in item and item["class"] == "punct":
                continue
            if "html" in item:
                new_json.append(item["html"])
                continue
            if "children" in item:
                new_json.append({item["class"] : normalize(item["children"])})
                continue
        return new_json    
    elif type(json) == dict:
        if "html" in json:
            return json["html"]
        if "children" in json:
            return {json["class"] : normalize(json["children"])}
        

with open(sys.argv[1], encoding='utf-8') as f:
    json_str = f.read()

json_obj = json.loads(json_str)
new_json_obj = {json_obj["class"] : normalize(json_obj["children"])}

with open(sys.argv[2], 'w', encoding='utf-8') as f:
    f.write(json.dumps(new_json_obj, indent=2, ensure_ascii=False))
