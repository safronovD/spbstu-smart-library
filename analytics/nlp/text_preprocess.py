from nltk import download
from nltk.corpus import stopwords
from nltk.stem.snowball import SnowballStemmer
from nltk.tokenize import RegexpTokenizer

download('stopwords')
download('punkt')


def stem_sentence(tokens, stemmer):
    output_sentence = []

    for token in tokens:
        output_sentence.append(stemmer.stem(token.lower()))

    return output_sentence


def get_cleared_words_from_text(text):
    tokens = RegexpTokenizer(r'[а-яА-Я]+').tokenize(text)
    words = stem_sentence(tokens, SnowballStemmer("russian", ignore_stopwords=True))
    words_without_stopwords = [word for word in words if not word in stopwords.words()]
    filtered_words = [word for word in words_without_stopwords if len(word) > 1]
    return filtered_words
