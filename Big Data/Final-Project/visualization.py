import numpy as np
import matplotlib.pyplot as plt

from matplotlib.ticker import NullFormatter
from sklearn import manifold, datasets
from time import time
from pyspark.ml.recommendation import ALS, ALSModel
import umap
import seaborn as sns
from sklearn.preprocessing import LabelEncoder

'''
this file is used to run line by line in spark directly. 
'''

#train = spark.read.parquet('goodread_train_5.parquet')
#val = spark.read.parquet('goodread_validation_5.parquet') 
#als = ALS(rank = 10, regParam=0.01, userCol="user_id", itemCol="book_id", ratingCol="rating", nonnegative=True, coldStartStrategy="drop")
#model = als.fit(train)

# load model from als.py saved ALS file
model2 = ALSModel.load('als_5_model')

genre = spark.read.csv('book_genre.csv', schema = 'index INT, book_id INT, genres STRING')
genre.createOrReplaceTempView('genre')
X = model2.itemFactors
X.createOrReplaceTempView('X')
joined = spark.sql('select genre.genres, X.features from genre inner join X on genre.book_id = X.id')
#count = spark.sql('select count(*) from joined')
#count.show() = 155667

feature = joined.select('features').collect()
feature = np.array([row['features'] for row in feature])
genre = joined.select('genres').collect()
genre = np.array([row['genres'] for row in genre])

LE = LabelEncoder()
g = LE.fit_transform(genre)

file = open("genre_number", "wb")
np.save(file, g)
file.close

# UMAP
reducer = umap.UMAP()
embedding = reducer.fit_transform(feature)
file = open("umap", "wb")
np.save(file, embedding)
file.close


# t-SNE
X_embedded = manifold.TSNE(n_components=2).fit_transform(feature)
X_embedded.shape
file = open("tsne", "wb")
np.save(file, X_embedded)
file.close

# plots are in jupyter notebook
