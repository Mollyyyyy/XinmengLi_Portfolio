from pyspark.sql.types import IntegerType
from pyspark.ml.feature import StringIndexer
from pyspark.sql.functions import monotonically_increasing_id
import random
import time

def create (spark, file_path):
    start = time.time()
    full_set = spark.read.csv(file_path, header=True, schema='user_id INT, book_id INT, is_read INT, rating INT, is_reviewed INT')
    full_set = full_set.filter(rating!=0)
    full_set.createOrReplaceTempView('full_set')
    a = spark.sql('Select user_id, count(book_id) FROM full_set GROUP BY user_id HAVING count(book_id)>10 SORT BY user_id')
    user = a.select('user_id').collect()
    user_list = [int(row['user_id']) for row in user]
    random.shuffle(user_list)
    user_num = len(user_list)
    temp_train, temp_val, temp_test = user_list[:int(0.6*user_num)], user_list[int(0.6*user_num):int(0.8*user_num)], user_list[int(0.8*user_num):]

    temp_train = spark.createDataFrame(temp_train,IntegerType())
    temp_train.createOrReplaceTempView('temp_train')
    train_60 = spark.sql('select user_id, book_id, rating from full_set inner join temp_train on temp_train.value = full_set.user_id')

    temp_val = spark.createDataFrame(temp_val,IntegerType())
    temp_val.createOrReplaceTempView('temp_val')
    val_full = spark.sql('select user_id, book_id, rating from full_set inner join temp_val on temp_val.value = full_set.user_id')
    val_full.createOrReplaceTempView('val_full')

    temp_test = spark.createDataFrame(temp_test,IntegerType())
    temp_test.createOrReplaceTempView('temp_test')
    test_full = spark.sql('select user_id, book_id, rating from full_set inner join temp_test on temp_test.value = full_set.user_id')
    test_full.createOrReplaceTempView('test_full')
    
    val_ind = spark.sql('SELECT ROW_NUMBER() OVER(PARTITION BY user_id ORDER BY book_id ) AS index,user_id,rating,book_id FROM val_full')
    test_ind = spark.sql('SELECT ROW_NUMBER() OVER(PARTITION BY user_id ORDER BY book_id ) AS index,user_id,rating,book_id FROM test_full')
    val_ind.createOrReplaceTempView('val_ind')
    test_ind.createOrReplaceTempView('test_ind')

    val_half = spark.sql('select user_id, book_id, rating from val_ind where index%2=0')
    val_add = spark.sql('select user_id, book_id, rating from val_ind where index%2=1')
    test_half = spark.sql('select user_id, book_id, rating from test_ind where index%2=0')
    test_add = spark.sql('select user_id, book_id, rating from test_ind where index%2=1')

    val = val_half.withColumn("index", monotonically_increasing_id())
    test = test_half.withColumn("index", monotonically_increasing_id())
    val.write.parquet("goodread_validation.parquet")
    test.write.parquet("goodread_test.parquet")

    union1 = train_60.unionAll(val_add)
    union2 = union1.unionAll(test_add)
    train = union2.withColumn("index", monotonically_increasing_id())
    train.write.parquet('goodread_train.parquet')
    
    end = time.time()
    print(end-start)

    return val_half
