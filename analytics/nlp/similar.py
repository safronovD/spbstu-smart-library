import nltk
from nltk.corpus import stopwords
import PyPDF2
import pdfminer.high_level
import os
import re

nltk.download('stopwords')
nltk.download('punkt')


def convertPdfToText(path):
    f = open(path, 'rb')
    text = pdfminer.high_level.extract_text(f)
    f.close()
    return text

def getClearedWordsFromPdf(path):
    text = convertPdfToText(path)
    words = list(filter(lambda word: word not in stopwords.words("russian") and len(word) > 1, nltk.tokenize.RegexpTokenizer(r'[а-яА-Я]+').tokenize(text)))
    print(words)

def compareTwoPdf(path1, path2):
    text1 = convertPdfToText(path1)
    text2 = convertPdfToText(path2)

def main():
    path = "../../connector/output/EBOOKS/pdfs/"
    cur_path = os.getcwd()
    
    os.chdir(path)

    pdfs = list(filter(lambda s: re.match(".*\.pdf$", s), os.listdir(".")))
    first = pdfs[0]
    getClearedWordsFromPdf(first)
    for pdf in pdfs[1:]:
        #compareTwoPdf(first, pdf)
        pass

    os.chdir(cur_path)
if __name__ == "__main__":
    main()