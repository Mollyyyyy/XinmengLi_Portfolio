from pyspark.ml.recommendation import ALS
from pyspark.mllib.evaluation import RankingMetrics
from pyspark.ml.evaluation import RegressionEvaluator
import pyspark.sql.functions as F
import numpy as np
import itertools as it
from pyspark.sql.functions import expr
from pyspark.sql import SparkSession
from pyspark.ml.recommendation import ALS, ALSModel
def als(spark):
    train = spark.read.parquet('goodread_train_1.parquet')
    #train.createOrReplaceTempView('train')
    #train = spark.sql('SELECT * FROM train')

    val = spark.read.parquet('goodread_validation_1.parquet')
    #val.createOrReplaceTempView('val')
    #val = spark.sql('SELECT * FROM val')

    #test = spark.read.parquet('goodread_test_1.parquet')
    #test.createOrReplaceTempView('test')
    #test = spark.sql('SELECT * FROM test ORDER BY index ASC')

    rank  = [10, 20, 50, 100]
    lam = [0.01, 0.1, 1]
    choice = it.product(rank, lam)

    ## Pick out users from validation set
    ## Change val to test here to get the ranking metrics on the test dataset.
    user_id_val = val.select('user_id').distinct()
    true_label_val = val.select('user_id', 'book_id').groupBy('user_id').agg(expr('collect_list(book_id) as true_item'))

    for i in choice:
        als = ALS(rank = i[0], regParam=i[1], userCol="user_id", itemCol="book_id", ratingCol="rating", maxIter=20,nonnegative=True, coldStartStrategy="drop")
        model = als.fit(train)
        print('Finish Training for {}'.format(i))
        model.save("model_1_iter20_2_"+str(i[0])+'_'+str(i[1]))
        #model = ALSModel.load('model_1_iter20_2_200_0.01')
        # Make top 500 recommendations for users in validation test
        res = model.recommendForUserSubset(user_id_val,500)
        #res.createOrReplaceTempView('res')
        pred_label = res.select('user_id','recommendations.book_id')
        pred_true_rdd = pred_label.join(F.broadcast(true_label_val), 'user_id', 'inner').rdd.map(lambda row: (row[1], row[2]))

        print('Start Evaluating for {}'.format(i))
        metrics = RankingMetrics(pred_true_rdd)
        map_ = metrics.meanAveragePrecision
        ndcg = metrics.ndcgAt(500)
        mpa = metrics.precisionAt(500)
        print(i, 'map score: ', map_, 'ndcg score: ', ndcg, 'precision at 500: ', mpa)
    
    return 0

'''
memory = "5g"

spark = (SparkSession.builder
            .appName('als')
            .master('yarn')
            .config('spark.executor.memory', memory)
            .config('spark.driver.memory', memory)
            .config('spark.executor.memoryOverhead', '4096')
            .config("spark.sql.broadcastTimeout", "36000")
            .config("spark.storage.memoryFraction","0")
            .config("spark.memory.offHeap.enabled","true")
            .config("spark.memory.offHeap.size","16g")
            .getOrCreate())
'''